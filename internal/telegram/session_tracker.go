package telegram

import (
	"fmt"
	"sync"
)

// sessionUsage tracks accumulated token usage and cost for a chat session.
type sessionUsage struct {
	InputTokens  int
	OutputTokens int
	CostUSD      float64
	NumTurns     int
}

// TotalTokens returns the total tokens consumed in this session.
func (u sessionUsage) TotalTokens() int {
	return u.InputTokens + u.OutputTokens
}

// String returns a human-readable summary of the session usage.
func (u sessionUsage) String() string {
	return fmt.Sprintf("Tokens: %d (in: %d, out: %d) | Turns: %d | Cost: $%.4f",
		u.TotalTokens(), u.InputTokens, u.OutputTokens, u.NumTurns, u.CostUSD)
}

// sessionTracker accumulates token usage per chat for auto-reset decisions.
type sessionTracker struct {
	mu    sync.RWMutex
	usage map[int64]*sessionUsage
}

func newSessionTracker() *sessionTracker {
	return &sessionTracker{
		usage: make(map[int64]*sessionUsage),
	}
}

// Add accumulates token usage for a chat. Returns the updated total tokens.
func (t *sessionTracker) Add(chatID int64, inputTokens, outputTokens, numTurns int, costUSD float64) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	u, ok := t.usage[chatID]
	if !ok {
		u = &sessionUsage{}
		t.usage[chatID] = u
	}
	u.InputTokens += inputTokens
	u.OutputTokens += outputTokens
	u.NumTurns += numTurns
	u.CostUSD += costUSD
	return u.TotalTokens()
}

// Get returns the current usage for a chat.
func (t *sessionTracker) Get(chatID int64) sessionUsage {
	t.mu.RLock()
	defer t.mu.RUnlock()
	u := t.usage[chatID]
	if u == nil {
		return sessionUsage{}
	}
	return *u
}

// Clear resets usage tracking for a chat.
func (t *sessionTracker) Clear(chatID int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.usage, chatID)
}
