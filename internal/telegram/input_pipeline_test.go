package telegram

import (
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/config"
)

func TestBuildSystemPrompt_WithoutDependencies(t *testing.T) {
	t.Parallel()

	bc := &BotController{
		config: &config.AppConfig{
			MemoryWindowSize: 5,
			Providers:        map[string]config.ProviderConfig{},
		},
	}

	prompt, err := bc.buildSystemPrompt("hello", nil)
	if err != nil {
		t.Fatalf("buildSystemPrompt() error = %v", err)
	}
	if prompt != "" {
		t.Fatalf("expected empty prompt with nil deps, got %q", prompt)
	}
}

func TestBuildSystemPrompt_WithAgent(t *testing.T) {
	t.Parallel()

	bc := &BotController{
		config: &config.AppConfig{
			MemoryWindowSize: 5,
			Providers:        map[string]config.ProviderConfig{},
		},
	}

	agent := &agents.Agent{
		Name:   "coder",
		Prompt: "You are a coding assistant.",
	}

	prompt, err := bc.buildSystemPrompt("write some code", agent)
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
			MemoryWindowSize: 5,
			Providers:        map[string]config.ProviderConfig{},
		},
	}

	agent := &agents.Agent{
		Name: "empty",
	}

	prompt, err := bc.buildSystemPrompt("hello", agent)
	if err != nil {
		t.Fatalf("buildSystemPrompt() error = %v", err)
	}
	if prompt != "" {
		t.Fatalf("expected empty prompt with empty agent, got %q", prompt)
	}
}

func TestTruncate(t *testing.T) {
	t.Parallel()

	short := "hello"
	if got := truncate(short, 10); got != short {
		t.Fatalf("truncate(%q, 10) = %q, want %q", short, got, short)
	}

	long := strings.Repeat("a", 100)
	got := truncate(long, 10)
	if len([]rune(got)) != 13 { // 10 + "..."
		t.Fatalf("truncate(100 chars, 10) should produce 13 runes, got %d", len([]rune(got)))
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("truncate should end with ..., got %q", got)
	}

	exact := "abcde"
	if got := truncate(exact, 5); got != exact {
		t.Fatalf("truncate(%q, 5) = %q, want %q", exact, got, exact)
	}
}
