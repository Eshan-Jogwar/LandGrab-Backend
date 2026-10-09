package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LandGrab/server"
	"github.com/gorilla/websocket"
)

func TestQueueAndMatchmakerLogic(t *testing.T) {
	qm := server.NewQueueManager()
	sm := server.NewSessionManager()
	mm := server.NewMatchmaker(qm, sm, 100*time.Millisecond)

	// Enqueue 7 mock players
	players := make([]*server.NetPlayer, 7)
	for i := 0; i < 7; i++ {
		players[i] = &server.NetPlayer{
			TempID:      string(rune('A' + i)),
			Send:        make(chan []byte, 10),
			IsConnected: true,
		}
		qm.Enqueue(players[i])
	}

	if qm.Len() != 7 {
		t.Fatalf("expected queue length 7, got %d", qm.Len())
	}

	// Run matchmaking pass
	created := mm.RunOnce()
	if created != 1 {
		t.Errorf("expected 1 game created, got %d", created)
	}

	// 7 - 4 = exactly 3 players remaining in queue
	if qm.Len() != 3 {
		t.Errorf("expected 3 players left in queue, got %d", qm.Len())
	}

	if sm.Count() != 1 {
		t.Errorf("expected 1 active session in SessionManager, got %d", sm.Count())
	}

	// Add 1 more player -> total 4 in queue
	extraPlayer := &server.NetPlayer{
		TempID:      "H",
		Send:        make(chan []byte, 10),
		IsConnected: true,
	}
	qm.Enqueue(extraPlayer)

	created = mm.RunOnce()
	if created != 1 {
		t.Errorf("expected 2nd game created, got %d", created)
	}

	// Queue should now be empty
	if qm.Len() != 0 {
		t.Errorf("expected queue to be empty, got %d", qm.Len())
	}

	if sm.Count() != 2 {
		t.Errorf("expected 2 active sessions, got %d", sm.Count())
	}
}

func TestServerHealthAndStatusEndpoints(t *testing.T) {
	s := server.NewServer(":0", 1*time.Second)

	// Health check
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	s.HealthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"status": "ok"`) {
		t.Errorf("unexpected health body: %s", w.Body.String())
	}
}

func TestWebSocketConnectionAndMatchmaking(t *testing.T) {
	s := server.NewServer(":0", 100*time.Millisecond)
	ts := httptest.NewServer(s.Handler())
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	// Connect 4 clients
	conns := make([]*websocket.Conn, 4)
	for i := 0; i < 4; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("failed to dial websocket client %d: %v", i, err)
		}
		defer conn.Close()
		conns[i] = conn

		// Receive MsgConnected
		var msg server.ServerMessage
		err = conn.ReadJSON(&msg)
		if err != nil {
			t.Fatalf("failed to read connected message: %v", err)
		}
		if msg.Type != server.MsgConnected {
			t.Errorf("expected MsgConnected, got %s", msg.Type)
		}

		// Receive MsgQueued
		err = conn.ReadJSON(&msg)
		if err != nil {
			t.Fatalf("failed to read queued message: %v", err)
		}
		if msg.Type != server.MsgQueued {
			t.Errorf("expected MsgQueued, got %s", msg.Type)
		}
	}

	// Trigger matchmaker pass
	s.Matchmaker().RunOnce()

	// All 4 clients should receive MsgGameStart followed by MsgGameState
	for i := 0; i < 4; i++ {
		var msg server.ServerMessage
		_ = conns[i].SetReadDeadline(time.Now().Add(2 * time.Second))

		err := conns[i].ReadJSON(&msg)
		if err != nil {
			t.Fatalf("client %d failed reading game start: %v", i, err)
		}
		if msg.Type != server.MsgGameStart {
			t.Errorf("client %d expected game_start, got %s", i, msg.Type)
		}

		err = conns[i].ReadJSON(&msg)
		if err != nil {
			t.Fatalf("client %d failed reading game state: %v", i, err)
		}
		if msg.Type != server.MsgGameState {
			t.Errorf("client %d expected game_state, got %s", i, msg.Type)
		}
	}

	// Client 1 (Player 1) is active: should receive MsgYourTurn
	var turnMsg server.ServerMessage
	err := conns[0].ReadJSON(&turnMsg)
	if err != nil {
		t.Fatalf("failed to read your_turn for player 1: %v", err)
	}
	if turnMsg.Type != server.MsgYourTurn {
		t.Errorf("expected MsgYourTurn, got %s", turnMsg.Type)
	}

	// Player 1 sends roll_dice
	rollAction := server.ClientMessage{
		Type: server.ActionRollDice,
	}
	err = conns[0].WriteJSON(rollAction)
	if err != nil {
		t.Fatalf("failed writing roll_dice: %v", err)
	}

	// Client 0 should receive MsgActionResult
	var actionRes server.ServerMessage
	_ = conns[0].SetReadDeadline(time.Now().Add(2 * time.Second))
	err = conns[0].ReadJSON(&actionRes)
	if err != nil {
		t.Fatalf("failed reading action result: %v", err)
	}
	if actionRes.Type != server.MsgActionResult {
		t.Errorf("expected MsgActionResult, got %s", actionRes.Type)
	}

	// All clients receive updated MsgGameState reflecting movement
	var stateMsg server.ServerMessage
	err = conns[0].ReadJSON(&stateMsg)
	if err != nil {
		t.Fatalf("failed reading updated state: %v", err)
	}
	if stateMsg.Type != server.MsgGameState {
		t.Errorf("expected updated MsgGameState, got %s", stateMsg.Type)
	}
}
