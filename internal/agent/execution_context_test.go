package agent

import (
	"context"
	"testing"
)

func TestRunContext_RoundTrip(t *testing.T) {
	t.Parallel()

	ctx := WithRunContext(context.Background(), "run-123")

	runID, ok := RunContextFromContext(ctx)
	if !ok {
		t.Fatalf("expected run context to be available")
	}
	if runID != "run-123" {
		t.Fatalf("unexpected run id: %q", runID)
	}
}

func TestRunContext_Missing(t *testing.T) {
	t.Parallel()

	if runID, ok := RunContextFromContext(context.Background()); ok || runID != "" {
		t.Fatalf("expected missing run context, got runID=%q ok=%v", runID, ok)
	}
}

func TestWorkdirContext_RoundTrip(t *testing.T) {
	t.Parallel()

	ctx := WithWorkdirContext(context.Background(), `C:\repo-alvo`)

	workdir, ok := WorkdirFromContext(ctx)
	if !ok {
		t.Fatalf("expected workdir context to be available")
	}
	if workdir != `C:\repo-alvo` {
		t.Fatalf("unexpected workdir: %q", workdir)
	}
}

func TestWorkdirContext_Missing(t *testing.T) {
	t.Parallel()

	if workdir, ok := WorkdirFromContext(context.Background()); ok || workdir != "" {
		t.Fatalf("expected missing workdir context, got workdir=%q ok=%v", workdir, ok)
	}
}
