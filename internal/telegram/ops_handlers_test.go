package telegram

import (
	"context"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/observability"
)

type fakeOpsService struct {
	recent   []observability.Operation
	failures []observability.Operation
	err      error
}

func (f *fakeOpsService) ListRecentOperations(ctx context.Context, limit int) ([]observability.Operation, error) {
	return f.recent, f.err
}

func (f *fakeOpsService) ListFailedOperations(ctx context.Context, limit int) ([]observability.Operation, error) {
	return f.failures, f.err
}

func TestOpsCommandHandler_HandleText(t *testing.T) {
	t.Parallel()

	handler := NewOpsCommandHandler(&fakeOpsService{
		failures: []observability.Operation{{
			Component:  "agent.tool",
			Operation:  "run_command",
			Status:     "error",
			DurationMS: 50,
			RunID:      "run-1",
			Summary:    "command failed",
		}},
		recent: []observability.Operation{{
			Component:  "agent.loop",
			Operation:  "llm_generate",
			Status:     "ok",
			DurationMS: 120,
			RunID:      "run-2",
			Summary:    "tool_calls=0",
		}},
	})

	reply, err := handler.HandleText(context.Background(), "/ops")
	if err != nil {
		t.Fatalf("HandleText() error = %v", err)
	}
	if !strings.Contains(reply, "Falhas recentes:") || !strings.Contains(reply, "Eventos recentes:") {
		t.Fatalf("unexpected reply %q", reply)
	}
	if !strings.Contains(reply, "agent.tool/run_command") {
		t.Fatalf("expected failure line in reply %q", reply)
	}
}
