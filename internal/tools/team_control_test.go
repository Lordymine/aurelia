package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
)

type fakeTeamController struct {
	paused    bool
	resumed   bool
	cancelled bool
	reason    string
	snapshot  agent.TeamStatusSnapshot
}

func (f *fakeTeamController) Pause(ctx context.Context, teamKey string) error {
	f.paused = true
	return nil
}

func (f *fakeTeamController) Resume(ctx context.Context, teamKey string) error {
	f.resumed = true
	return nil
}

func (f *fakeTeamController) Cancel(ctx context.Context, teamKey, reason string) error {
	f.cancelled = true
	f.reason = reason
	return nil
}

func (f *fakeTeamController) BuildStatusSnapshot(ctx context.Context, teamKey string) (agent.TeamStatusSnapshot, error) {
	return f.snapshot, nil
}

func TestPauseTeamTool_Execute(t *testing.T) {
	controller := &fakeTeamController{}
	tool := NewPauseTeamTool(controller)
	ctx := agent.WithTeamContext(context.Background(), "team-key", "user-1")

	got, err := tool.Execute(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !controller.paused {
		t.Fatal("expected pause to be invoked")
	}
	if !strings.Contains(strings.ToLower(got), "pausei") {
		t.Fatalf("unexpected message: %q", got)
	}
}

func TestCancelTeamTool_Execute_ForwardsReason(t *testing.T) {
	controller := &fakeTeamController{}
	tool := NewCancelTeamTool(controller)
	ctx := agent.WithTeamContext(context.Background(), "team-key", "user-1")

	_, err := tool.Execute(ctx, map[string]interface{}{"reason": "usuario pediu para parar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !controller.cancelled {
		t.Fatal("expected cancel to be invoked")
	}
	if controller.reason != "usuario pediu para parar" {
		t.Fatalf("unexpected reason: %q", controller.reason)
	}
}

func TestTeamStatusTool_Execute_RendersSnapshot(t *testing.T) {
	controller := &fakeTeamController{
		snapshot: agent.TeamStatusSnapshot{
			TeamStatus: "paused",
			Pending:    2,
			Running:    1,
			Completed:  3,
			TotalTasks: 6,
		},
	}
	tool := NewTeamStatusTool(controller)
	ctx := agent.WithTeamContext(context.Background(), "team-key", "user-1")

	got, err := tool.Execute(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "status=paused") || !strings.Contains(got, "pendentes=2") {
		t.Fatalf("unexpected status output: %q", got)
	}
}


