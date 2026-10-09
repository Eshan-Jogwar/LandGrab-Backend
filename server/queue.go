package server

import (
	"sync"
)

// QueueManager provides a thread-safe waiting lobby for matchmaking.
type QueueManager struct {
	mu      sync.Mutex
	players []*NetPlayer
}

// NewQueueManager initializes an empty QueueManager.
func NewQueueManager() *QueueManager {
	return &QueueManager{
		players: make([]*NetPlayer, 0),
	}
}

// Enqueue adds a player to the waiting queue and broadcasts their queued status.
func (qm *QueueManager) Enqueue(p *NetPlayer) int {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	// Check if already in queue
	for i, existing := range qm.players {
		if existing.TempID == p.TempID {
			return i + 1
		}
	}

	qm.players = append(qm.players, p)
	pos := len(qm.players)

	p.SendJSON(ServerMessage{
		Type: MsgQueued,
		Payload: QueuedPayload{
			TempID:       p.TempID,
			Position:     pos,
			TotalQueued:  len(qm.players),
			RequiredSize: 4,
		},
	})

	return pos
}

// DequeueGroup atomically extracts a group of the given size.
// If fewer players than size are present, it leaves the queue untouched and returns false.
func (qm *QueueManager) DequeueGroup(size int) ([]*NetPlayer, bool) {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if len(qm.players) < size {
		return nil, false
	}

	group := make([]*NetPlayer, size)
	copy(group, qm.players[:size])
	qm.players = qm.players[size:]

	return group, true
}

// Remove removes a player by their TempID (e.g. if their socket disconnected while waiting).
func (qm *QueueManager) Remove(tempID string) bool {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	for i, p := range qm.players {
		if p.TempID == tempID {
			qm.players = append(qm.players[:i], qm.players[i+1:]...)
			return true
		}
	}
	return false
}

// Len returns the current count of players waiting in the queue.
func (qm *QueueManager) Len() int {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	return len(qm.players)
}

// Clear flushes all waiting players (e.g. during server shutdown).
func (qm *QueueManager) Clear() {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.players = make([]*NetPlayer, 0)
}
