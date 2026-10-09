package server

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/LandGrab/gamecore/state"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// ActionHandler is a callback to handle player actions received over WebSocket.
type ActionHandler func(action PlayerAction)

// DisconnectHandler is called when a player's connection drops.
type DisconnectHandler func(tempID string)

// NetPlayer bundles a network connection and temporary ID with a gamecore player instance.
type NetPlayer struct {
	TempID      string
	Player      *state.Player
	Conn        *websocket.Conn
	Send        chan []byte
	GameID      string
	IsConnected bool
	mu          sync.Mutex

	onAction     ActionHandler
	onDisconnect DisconnectHandler
}

// NewNetPlayer initializes a new NetPlayer with dedicated send buffer.
func NewNetPlayer(tempID string, conn *websocket.Conn) *NetPlayer {
	return &NetPlayer{
		TempID:      tempID,
		Conn:        conn,
		Send:        make(chan []byte, 256),
		IsConnected: true,
	}
}

// SetActionHandler binds an action receiver callback (used when entering a NetGame).
func (np *NetPlayer) SetActionHandler(handler ActionHandler) {
	np.mu.Lock()
	defer np.mu.Unlock()
	np.onAction = handler
}

// SetDisconnectHandler binds a cleanup callback.
func (np *NetPlayer) SetDisconnectHandler(handler DisconnectHandler) {
	np.mu.Lock()
	defer np.mu.Unlock()
	np.onDisconnect = handler
}

// SendJSON marshals and enqueues a ServerMessage to the player's write pump.
func (np *NetPlayer) SendJSON(msg ServerMessage) bool {
	np.mu.Lock()
	if !np.IsConnected {
		np.mu.Unlock()
		return false
	}
	np.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[NetPlayer %s] Error marshaling message: %v", np.TempID, err)
		return false
	}

	select {
	case np.Send <- data:
		return true
	default:
		log.Printf("[NetPlayer %s] Outbound buffer full, dropping message", np.TempID)
		return false
	}
}

// ReadPump runs in a dedicated goroutine reading incoming frames from the WebSocket.
func (np *NetPlayer) ReadPump() {
	defer func() {
		np.Close()
		np.mu.Lock()
		onDisc := np.onDisconnect
		np.mu.Unlock()
		if onDisc != nil {
			onDisc(np.TempID)
		}
	}()

	np.Conn.SetReadLimit(maxMessageSize)
	_ = np.Conn.SetReadDeadline(time.Now().Add(pongWait))
	np.Conn.SetPongHandler(func(string) error {
		_ = np.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := np.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[NetPlayer %s] WebSocket error: %v", np.TempID, err)
			}
			break
		}

		var clientMsg ClientMessage
		if err := json.Unmarshal(message, &clientMsg); err != nil {
			np.SendJSON(ServerMessage{
				Type:  MsgError,
				Error: "Malformed JSON payload",
			})
			continue
		}

		np.mu.Lock()
		handler := np.onAction
		np.mu.Unlock()

		if handler != nil {
			handler(PlayerAction{
				SenderTempID: np.TempID,
				Message:      clientMsg,
			})
		}
	}
}

// WritePump runs in a dedicated goroutine writing messages and heartbeats to the WebSocket.
func (np *NetPlayer) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		np.Close()
	}()

	for {
		select {
		case message, ok := <-np.Send:
			np.mu.Lock()
			if !np.IsConnected {
				np.mu.Unlock()
				return
			}
			_ = np.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = np.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				np.mu.Unlock()
				return
			}
			err := np.Conn.WriteMessage(websocket.TextMessage, message)
			np.mu.Unlock()
			if err != nil {
				return
			}

		case <-ticker.C:
			np.mu.Lock()
			if !np.IsConnected {
				np.mu.Unlock()
				return
			}
			_ = np.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := np.Conn.WriteMessage(websocket.PingMessage, nil)
			np.mu.Unlock()
			if err != nil {
				return
			}
		}
	}
}

// Close gracefully closes the socket connection and marks the player disconnected.
func (np *NetPlayer) Close() {
	np.mu.Lock()
	defer np.mu.Unlock()
	if np.IsConnected {
		np.IsConnected = false
		_ = np.Conn.Close()
	}
}
