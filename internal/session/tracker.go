package session

import (
	"fmt"
	"sync"
)

const estimatedTokensPerTurn = 3000

// Usage tracks accumulated token usage and cost for a chat session.
type Usage struct {
	InputTokens  int
	OutputTokens int
	CostUSD      float64
	NumTurns     int
}

func (u Usage) TotalTokens() int {
	return u.InputTokens + u.OutputTokens
}

func (u Usage) String() string {
	return fmt.Sprintf("Tokens: %d (in: %d, out: %d) | Turns: %d | Cost: $%.4f",
		u.TotalTokens(), u.InputTokens, u.OutputTokens, u.NumTurns, u.CostUSD)
}

// Tracker accumulates token usage per chat for auto-reset decisions.
type Tracker struct {
	mu    sync.RWMutex
	usage map[int64]*Usage
}

func NewTracker() *Tracker {
	return &Tracker{
		usage: make(map[int64]*Usage),
	}
}

// Add accumulates token usage for a chat. Returns the updated total tokens.
func (t *Tracker) Add(chatID int64, inputTokens, outputTokens, numTurns int, costUSD float64) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	u, ok := t.usage[chatID]
	if !ok {
		u = &Usage{}
		t.usage[chatID] = u
	}
	u.InputTokens += inputTokens
	u.OutputTokens += outputTokens
	u.NumTurns += numTurns
	u.CostUSD += costUSD
	return u.TotalTokens()
}

// RecordUsage estimates tokens from turns, tracks them, and returns true
// if the session should be reset (threshold exceeded).
func (t *Tracker) RecordUsage(chatID int64, numTurns int, costUSD float64, maxTokens int) bool {
	estimatedTokens := numTurns * estimatedTokensPerTurn
	totalTokens := t.Add(chatID, estimatedTokens, 0, numTurns, costUSD)
	return maxTokens > 0 && totalTokens >= maxTokens
}

// Get returns the current usage for a chat.
func (t *Tracker) Get(chatID int64) Usage {
	t.mu.RLock()
	defer t.mu.RUnlock()
	u := t.usage[chatID]
	if u == nil {
		return Usage{}
	}
	return *u
}

// Clear resets usage tracking for a chat.
func (t *Tracker) Clear(chatID int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.usage, chatID)
}
