package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for dev/game clients
		return true
	},
}

// Server provides the complete networking layer for LandGrab.
type Server struct {
	Addr       string
	queue      *QueueManager
	sessions   *SessionManager
	matchmaker *Matchmaker
	httpServer *http.Server

	playerCounter uint64
	startTime     time.Time
}

// NewServer initializes all server components.
func NewServer(addr string, matchmakingInterval time.Duration) *Server {
	queue := NewQueueManager()
	sessions := NewSessionManager()
	matchmaker := NewMatchmaker(queue, sessions, matchmakingInterval)

	s := &Server{
		Addr:       addr,
		queue:      queue,
		sessions:   sessions,
		matchmaker: matchmaker,
		startTime:  time.Now(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/status", s.handleStatus)
	mux.HandleFunc("/games/", s.handleGetGame)

	// Apply middlewares: CORS -> Logging -> Recovery -> Handler
	handler := RecoveryMiddleware(LoggingMiddleware(CORSMiddleware(mux)))

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

// Handler returns the HTTP handler for testing.
func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
}

// Matchmaker returns the server's matchmaker instance.
func (s *Server) Matchmaker() *Matchmaker {
	return s.matchmaker
}

// HealthHandler exposes the health check for testing.
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	s.handleHealth(w, r)
}

// Start launches the matchmaker daemon and HTTP listener.
func (s *Server) Start() error {
	s.matchmaker.Start()
	log.Printf("[Server] LandGrab server listening on %s", s.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully terminates the server and background matchmaker.
func (s *Server) Shutdown(ctx context.Context) error {
	s.matchmaker.Stop()
	return s.httpServer.Shutdown(ctx)
}

// handleWebSocket handles WebSocket handshake, assigns temporary ID, and enqueues player.
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[Server] WebSocket upgrade error: %v", err)
		return
	}

	// Generate sequential temporary ID
	idNum := atomic.AddUint64(&s.playerCounter, 1)
	tempID := fmt.Sprintf("guest-%d", idNum)

	netPlayer := NewNetPlayer(tempID, conn)

	// Check if this is a reconnection request
	reconnectID := r.URL.Query().Get("reconnect_id")
	gameID := r.URL.Query().Get("game_id")

	if reconnectID != "" && gameID != "" {
		if game, ok := s.sessions.GetGame(gameID); ok {
			if existingPlayer, ok := game.TempIDIndex[reconnectID]; ok {
				existingPlayer.mu.Lock()
				existingPlayer.Conn = conn
				existingPlayer.IsConnected = true
				existingPlayer.mu.Unlock()

				go existingPlayer.WritePump()
				go existingPlayer.ReadPump()

				existingPlayer.SendJSON(ServerMessage{
					Type: MsgConnected,
					Payload: ConnectedPayload{
						TempID: reconnectID,
					},
				})
				game.broadcastGameState()
				return
			}
		}
	}

	// New player connection: Send welcome message and enqueue
	netPlayer.SendJSON(ServerMessage{
		Type: MsgConnected,
		Payload: ConnectedPayload{
			TempID: tempID,
		},
	})

	// When disconnected before match, remove from queue
	netPlayer.SetDisconnectHandler(func(id string) {
		s.queue.Remove(id)
	})

	s.queue.Enqueue(netPlayer)

	go netPlayer.WritePump()
	go netPlayer.ReadPump()
}

// handleHealth returns basic 200 OK health check.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

// handleStatus returns server metrics, queue size, and active game counts.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := map[string]interface{}{
		"status":         "healthy",
		"uptime_seconds": int(time.Since(s.startTime).Seconds()),
		"queue": map[string]interface{}{
			"players_waiting": s.queue.Len(),
		},
		"active_sessions": map[string]interface{}{
			"total_games": s.sessions.Count(),
		},
	}

	_ = json.NewEncoder(w).Encode(stats)
}

// handleGetGame returns an active game's current state snapshot via HTTP GET.
func (s *Server) handleGetGame(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Path[len("/games/"):]
	if gameID == "" {
		http.Error(w, `{"error": "Game ID required"}`, http.StatusBadRequest)
		return
	}

	game, ok := s.sessions.GetGame(gameID)
	if !ok {
		http.Error(w, `{"error": "Game not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(game.CreateSnapshot())
}
