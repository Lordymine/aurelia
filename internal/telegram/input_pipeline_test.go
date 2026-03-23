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

