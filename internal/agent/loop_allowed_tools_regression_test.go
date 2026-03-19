package agent

import (
	"context"
	"testing"
)

type captureToolsProvider struct {
	toolsSeen []string
}

func (c *captureToolsProvider) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	c.toolsSeen = c.toolsSeen[:0]
	for _, tool := range tools {
		c.toolsSeen = append(c.toolsSeen, tool.Name)
	}
	return &ModelResponse{Content: "ok"}, nil
}

func TestLoopRun_EmptyAllowedToolsDoesNotExpandToFullRegistry(t *testing.T) {
	t.Parallel()

	provider := &captureToolsProvider{}
	registry := NewToolRegistry()
	registry.Register(Tool{Name: "read_file"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "run_command"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})

	loop := NewLoop(provider, registry, 1)
	ctx := WithDynamicToolAccess(context.Background(), true)
	if _, _, err := loop.Run(ctx, "base prompt", nil, []string{}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if len(provider.toolsSeen) != 1 || provider.toolsSeen[0] != requestToolAccessToolName {
		t.Fatalf("expected only request_tool_access, got %v", provider.toolsSeen)
	}
}
