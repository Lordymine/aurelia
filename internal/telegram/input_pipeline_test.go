package telegram

import (
	"strings"
	"testing"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/bridge"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/session"
)

func TestBuildSystemPrompt_WithoutDependencies(t *testing.T) {
	t.Parallel()

	bc := &BotController{
		config: &config.AppConfig{
			Providers:        map[string]config.ProviderConfig{},
		},
	}

	prompt, err := bc.buildSystemPrompt("hello", nil, 0, 0)
	if err != nil {
		t.Fatalf("buildSystemPrompt() error = %v", err)
	}
	if !strings.Contains(prompt, "## Scheduling Tasks") {
		t.Fatalf("expected cron instructions in prompt, got %q", prompt)
	}
}

func TestBuildSystemPrompt_WithAgent(t *testing.T) {
	t.Parallel()

	bc := &BotController{
		config: &config.AppConfig{
			Providers:        map[string]config.ProviderConfig{},
		},
	}

	agent := &agents.Agent{
		Name:   "coder",
		Prompt: "You are a coding assistant.",
	}

	prompt, err := bc.buildSystemPrompt("write some code", agent, 0, 0)
	if err != nil {
		t.Fatalf("buildSystemPrompt() error = %v", err)
	}
	if !strings.Contains(prompt, "You are a coding assistant.") {
		t.Fatalf("expected agent prompt in system prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "# Agent Instructions") {
		t.Fatalf("expected agent header in system prompt, got %q", prompt)
	}
}

func TestBuildSystemPrompt_AgentWithEmptyPrompt(t *testing.T) {
	t.Parallel()

	bc := &BotController{
		config: &config.AppConfig{
			Providers:        map[string]config.ProviderConfig{},
		},
	}

	agent := &agents.Agent{
		Name: "empty",
	}

	prompt, err := bc.buildSystemPrompt("hello", agent, 0, 0)
	if err != nil {
		t.Fatalf("buildSystemPrompt() error = %v", err)
	}
	// With empty agent, prompt should only contain cron instructions
	if !strings.Contains(prompt, "## Scheduling Tasks") {
		t.Fatalf("expected cron instructions in prompt, got %q", prompt)
	}
	if strings.Contains(prompt, "# Agent Instructions") {
		t.Fatalf("expected no agent instructions with empty agent, got %q", prompt)
	}
}

// newTestBotController creates a minimal BotController for testing bridge event
// processing. The bot field is nil — only use for tests that don't trigger Send.
func newTestBotController() *BotController {
	return &BotController{
		config: &config.AppConfig{
			MaxSessionTokens: 100000,
			Providers:        map[string]config.ProviderConfig{},
		},
		sessions: session.NewStore(),
		tracker:  session.NewTracker(),
	}
}

// noopProgress returns a progressReporter that does nothing (no bot required).
func noopProgress() *progressReporter {
	return &progressReporter{}
}

func TestProcessBridgeEventsAsync_ProcessDeath_NoTerminalEvent(t *testing.T) {
	bc := newTestBotController()
	chat := &telebot.Chat{ID: 1}

	ch := make(chan bridge.Event, 2)
	ch <- bridge.Event{Type: "system", SessionID: "sess-1"}
	ch <- bridge.Event{Type: "assistant", Text: "partial response"}
	close(ch) // simulate process death — no terminal event

	outcome := bc.processBridgeEventsAsync(chat, ch, noopProgress(), "test", 1)
	if outcome != outcomeProcessDeath {
		t.Fatalf("expected outcomeProcessDeath, got %d", outcome)
	}

	// Session should still have been set from the system event
	sid := bc.sessions.Get(1)
	if sid != "sess-1" {
		t.Fatalf("expected session sess-1, got %q", sid)
	}
}

func TestProcessBridgeEventsAsync_EmptyChannelIsDeath(t *testing.T) {
	bc := newTestBotController()
	chat := &telebot.Chat{ID: 1}

	ch := make(chan bridge.Event)
	close(ch) // immediate close — no events at all

	outcome := bc.processBridgeEventsAsync(chat, ch, noopProgress(), "test", 1)
	if outcome != outcomeProcessDeath {
		t.Fatalf("expected outcomeProcessDeath, got %d", outcome)
	}
}

func TestProcessBridgeEventsAsync_SessionSetFromSystemEvent(t *testing.T) {
	bc := newTestBotController()
	chat := &telebot.Chat{ID: 42}

	ch := make(chan bridge.Event, 1)
	ch <- bridge.Event{Type: "system", SessionID: "sess-xyz"}
	close(ch) // death after system event

	bc.processBridgeEventsAsync(chat, ch, noopProgress(), "test", 1)

	sid, active := bc.sessions.GetWithState(42)
	if sid != "sess-xyz" {
		t.Fatalf("expected session sess-xyz, got %q", sid)
	}
	if !active {
		t.Fatal("session should be active after Set")
	}
}

func TestBridgeRecovery_DeactivateAllPreservesIDs(t *testing.T) {
	sessions := session.NewStore()
	sessions.Set(1, "sess-a")
	sessions.Set(2, "sess-b")

	sessions.DeactivateAll()

	// Sessions should be cold but IDs preserved
	sid, active := sessions.GetWithState(1)
	if active {
		t.Fatal("session 1 should be inactive")
	}
	if sid != "sess-a" {
		t.Fatalf("session 1 ID should be preserved, got %q", sid)
	}

	sid, active = sessions.GetWithState(2)
	if active {
		t.Fatal("session 2 should be inactive")
	}
	if sid != "sess-b" {
		t.Fatalf("session 2 ID should be preserved, got %q", sid)
	}

	// Get still works
	if id := sessions.Get(1); id != "sess-a" {
		t.Fatalf("Get(1) = %q, want sess-a", id)
	}
}

// --- P3: Backoff tests ---

func TestBridgeFailureTracker_RecordAndCooldown(t *testing.T) {
	var tracker bridgeFailureTracker

	// First two failures: not in cooldown
	tracker.record()
	if tracker.inCooldown() {
		t.Fatal("should not be in cooldown after 1 failure")
	}
	tracker.record()
	if tracker.inCooldown() {
		t.Fatal("should not be in cooldown after 2 failures")
	}

	// Third failure: enters cooldown
	inCooldown := tracker.record()
	if !inCooldown {
		t.Fatal("record should return true after 3 failures")
	}
	if !tracker.inCooldown() {
		t.Fatal("should be in cooldown after 3 failures")
	}
}

func TestBridgeFailureTracker_ResetClearsCooldown(t *testing.T) {
	var tracker bridgeFailureTracker

	tracker.record()
	tracker.record()
	tracker.record()

	if !tracker.inCooldown() {
		t.Fatal("should be in cooldown after 3 failures")
	}

	tracker.reset()

	if tracker.inCooldown() {
		t.Fatal("should not be in cooldown after reset")
	}
}

func TestBridgeFailureTracker_OldFailuresExpire(t *testing.T) {
	var tracker bridgeFailureTracker

	// Manually add old failures outside the window
	tracker.mu.Lock()
	old := time.Now().Add(-2 * time.Minute)
	tracker.failures = []time.Time{old, old}
	tracker.mu.Unlock()

	// New failure + old ones outside window = only 1 recent failure
	tracker.record()

	if tracker.inCooldown() {
		t.Fatal("should not be in cooldown — old failures should have expired")
	}
}

func TestBridgeFailureTracker_EmptyNotInCooldown(t *testing.T) {
	var tracker bridgeFailureTracker
	if tracker.inCooldown() {
		t.Fatal("empty tracker should not be in cooldown")
	}
}

