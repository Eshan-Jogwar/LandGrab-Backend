package server

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// SessionManager tracks and provides CRUD operations for active NetGame sessions.
type SessionManager struct {
	mu          sync.RWMutex
	sessions    map[string]*NetGame
	gameCounter uint64
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*NetGame),
	}
}

// CreateGame initializes a new 4-player game session and registers it.
func (sm *SessionManager) CreateGame(players []*NetPlayer) *NetGame {
	if len(players) != 4 {
		panic(fmt.Sprintf("CreateGame requires exactly 4 players, got %d", len(players)))
	}

	gameNum := atomic.AddUint64(&sm.gameCounter, 1)
	gameID := fmt.Sprintf("game-%d", gameNum)

	var fourPlayers [4]*NetPlayer
	copy(fourPlayers[:], players[:4])

	ng := NewNetGame(gameID, fourPlayers)

	// Bind session cleanup callback
	ng.onFinish = func(id string) {
		sm.DeleteGame(id)
	}

	sm.mu.Lock()
	sm.sessions[gameID] = ng
	sm.mu.Unlock()

	return ng
}

// GetGame retrieves an active NetGame session by ID.
func (sm *SessionManager) GetGame(id string) (*NetGame, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	ng, ok := sm.sessions[id]
	return ng, ok
}

// DeleteGame removes a game session from the registry.
func (sm *SessionManager) DeleteGame(id string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, ok := sm.sessions[id]; ok {
		delete(sm.sessions, id)
		return true
	}
	return false
}

// Count returns the number of currently active sessions.
func (sm *SessionManager) Count() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.sessions)
}

// ListGameIDs returns all active session IDs.
func (sm *SessionManager) ListGameIDs() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	ids := make([]string, 0, len(sm.sessions))
	for id := range sm.sessions {
		ids = append(ids, id)
	}
	return ids
}
