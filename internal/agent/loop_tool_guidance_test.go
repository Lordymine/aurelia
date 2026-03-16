package agent

import (
	"context"
	"strings"
	"testing"
)

type capturePromptProvider struct {
	systemPrompt string
}

func (c *capturePromptProvider) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	c.systemPrompt = systemPrompt
	return &ModelResponse{Content: "ok"}, nil
}

func TestLoop_Run_AppendsToolUsageGuidanceForLocalExecution(t *testing.T) {
	t.Parallel()

	provider := &capturePromptProvider{}
	registry := NewToolRegistry()
	registry.Register(Tool{Name: "run_command", Description: "runs commands"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "read_file", Description: "reads files"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "write_file", Description: "writes files"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "list_dir", Description: "lists dirs"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "create_schedule", Description: "creates schedules"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "list_schedules", Description: "lists schedules"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "pause_schedule", Description: "pauses schedules"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})

	loop := NewLoop(provider, registry, 3)

	_, _, err := loop.Run(context.Background(), "base prompt", nil, []string{"run_command", "read_file", "write_file", "list_dir", "create_schedule", "list_schedules", "pause_schedule"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if !strings.Contains(provider.systemPrompt, "Se o usuario pedir para rodar") {
		t.Fatalf("expected system prompt to include execution guidance, got %q", provider.systemPrompt)
	}
	if !strings.Contains(provider.systemPrompt, "`workdir`") {
		t.Fatalf("expected system prompt to mention workdir, got %q", provider.systemPrompt)
	}
	if !strings.Contains(provider.systemPrompt, "Nao diga que o ambiente esta bloqueado") {
		t.Fatalf("expected system prompt to forbid inventing environment restrictions, got %q", provider.systemPrompt)
	}
	if !strings.Contains(provider.systemPrompt, "create_schedule") {
		t.Fatalf("expected system prompt to mention scheduling tools, got %q", provider.systemPrompt)
	}
	if !strings.Contains(provider.systemPrompt, "# RUNTIME CAPABILITIES") {
		t.Fatalf("expected runtime capabilities block, got %q", provider.systemPrompt)
	}
	if !strings.Contains(provider.systemPrompt, "- run_command") {
		t.Fatalf("expected runtime capabilities to list exact tools, got %q", provider.systemPrompt)
	}
}
