package server

import (
	"log"
	"sync"
	"time"
)

const (
	defaultMatchmakingInterval = 30 * time.Second
	matchGroupSize             = 4
)

// Matchmaker periodically inspects QueueManager and creates 4-player matches.
type Matchmaker struct {
	queue    *QueueManager
	sessions *SessionManager
	interval time.Duration
	ticker   *time.Ticker
	done     chan struct{}
	wg       sync.WaitGroup
	mu       sync.Mutex
	running  bool
}

// NewMatchmaker initializes a new Matchmaker instance.
func NewMatchmaker(queue *QueueManager, sessions *SessionManager, interval time.Duration) *Matchmaker {
	if interval <= 0 {
		interval = defaultMatchmakingInterval
	}
	return &Matchmaker{
		queue:    queue,
		sessions: sessions,
		interval: interval,
		done:     make(chan struct{}),
	}
}

// Start launches the background matchmaking daemon.
func (m *Matchmaker) Start() {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return
	}
	m.running = true
	m.ticker = time.NewTicker(m.interval)
	m.mu.Unlock()

	m.wg.Add(1)
	go m.run()
	log.Printf("[Matchmaker] Started matchmaking daemon with %v interval", m.interval)
}

// Stop gracefully terminates the matchmaking daemon.
func (m *Matchmaker) Stop() {
	m.mu.Lock()
	if !m.running {
		m.mu.Unlock()
		return
	}
	m.running = false
	close(m.done)
	if m.ticker != nil {
		m.ticker.Stop()
	}
	m.mu.Unlock()

	m.wg.Wait()
	log.Println("[Matchmaker] Stopped matchmaking daemon")
}

// RunOnce runs an immediate matchmaking pass (useful for manual triggers and unit tests).
func (m *Matchmaker) RunOnce() int {
	matchesCreated := 0
	for {
		group, ok := m.queue.DequeueGroup(matchGroupSize)
		if !ok {
			// Fewer than 4 players remaining -> leave them in queue for next pass
			break
		}

		// Filter out any players who disconnected while waiting in queue
		validGroup := make([]*NetPlayer, 0, matchGroupSize)
		for _, p := range group {
			p.mu.Lock()
			connected := p.IsConnected
			p.mu.Unlock()
			if connected {
				validGroup = append(validGroup, p)
			}
		}

		// If some players disconnected right before dequeue, re-enqueue remaining valid players at front
		if len(validGroup) < matchGroupSize {
			for _, p := range validGroup {
				m.queue.Enqueue(p)
			}
			break
		}

		// Create and spawn session
		session := m.sessions.CreateGame(validGroup)
		go session.Run()
		matchesCreated++
		log.Printf("[Matchmaker] Created game session %s with 4 players", session.ID)
	}
	return matchesCreated
}

func (m *Matchmaker) run() {
	defer m.wg.Done()

	for {
		select {
		case <-m.ticker.C:
			m.RunOnce()
		case <-m.done:
			return
		}
	}
}
