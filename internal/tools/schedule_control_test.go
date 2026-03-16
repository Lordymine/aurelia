package tools

import (
	"context"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
)

type fakeScheduleController struct {
	paused  []string
	resumed []string
	deleted []string
	err     error
}

func (f *fakeScheduleController) PauseJob(ctx context.Context, jobID string) error {
	if f.err != nil {
		return f.err
	}
	f.paused = append(f.paused, jobID)
	return nil
}

func (f *fakeScheduleController) ResumeJob(ctx context.Context, jobID string) error {
	if f.err != nil {
		return f.err
	}
	f.resumed = append(f.resumed, jobID)
	return nil
}

func (f *fakeScheduleController) DeleteJob(ctx context.Context, jobID string) error {
	if f.err != nil {
		return f.err
	}
	f.deleted = append(f.deleted, jobID)
	return nil
}

func TestPauseScheduleTool_Execute(t *testing.T) {
	t.Parallel()

	controller := &fakeScheduleController{}
	tool := NewPauseScheduleTool(controller)

	ctx := agent.WithTeamContext(context.Background(), "123", "user-1")
	result, err := tool.Execute(ctx, map[string]interface{}{"job_id": "job-a"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(controller.paused) != 1 || controller.paused[0] != "job-a" {
		t.Fatalf("unexpected paused calls: %#v", controller.paused)
	}
	if result == "" {
		t.Fatalf("expected non-empty result")
	}
}

func TestResumeScheduleTool_Execute(t *testing.T) {
	t.Parallel()

	controller := &fakeScheduleController{}
	tool := NewResumeScheduleTool(controller)

	ctx := agent.WithTeamContext(context.Background(), "123", "user-1")
	result, err := tool.Execute(ctx, map[string]interface{}{"job_id": "job-a"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(controller.resumed) != 1 || controller.resumed[0] != "job-a" {
		t.Fatalf("unexpected resumed calls: %#v", controller.resumed)
	}
	if result == "" {
		t.Fatalf("expected non-empty result")
	}
}

func TestDeleteScheduleTool_Execute(t *testing.T) {
	t.Parallel()

	controller := &fakeScheduleController{}
	tool := NewDeleteScheduleTool(controller)

	ctx := agent.WithTeamContext(context.Background(), "123", "user-1")
	result, err := tool.Execute(ctx, map[string]interface{}{"job_id": "job-a"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(controller.deleted) != 1 || controller.deleted[0] != "job-a" {
		t.Fatalf("unexpected deleted calls: %#v", controller.deleted)
	}
	if result == "" {
		t.Fatalf("expected non-empty result")
	}
}


