package agent

import (
	"context"
	"testing"
)

func TestToolRegistryFilterDefinitions_NilMeansAll(t *testing.T) {
	t.Parallel()

	registry := NewToolRegistry()
	registry.Register(Tool{Name: "read_file"}, func(ctx context.Context, args map[string]interface{}) (string, error) { return "", nil })
	registry.Register(Tool{Name: "run_command"}, func(ctx context.Context, args map[string]interface{}) (string, error) { return "", nil })

	defs := registry.FilterDefinitions(nil)
	if len(defs) != 2 {
		t.Fatalf("expected all tools, got %d", len(defs))
	}
}

func TestToolRegistryFilterDefinitions_EmptySliceMeansNone(t *testing.T) {
	t.Parallel()

	registry := NewToolRegistry()
	registry.Register(Tool{Name: "read_file"}, func(ctx context.Context, args map[string]interface{}) (string, error) { return "", nil })

	defs := registry.FilterDefinitions([]string{})
	if len(defs) != 0 {
		t.Fatalf("expected no tools, got %d", len(defs))
	}
}
