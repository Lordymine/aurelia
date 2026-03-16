package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
	cronpkg "github.com/kocar/aurelia/internal/cron"
)

type fakeScheduleService struct {
	created []cronpkg.CronJob
	err     error
}

func (f *fakeScheduleService) CreateJob(ctx context.Context, job cronpkg.CronJob) (string, error) {
	f.created = append(f.created, job)
	if f.err != nil {
		return "", f.err
	}
	if job.ID == "" {
		job.ID = "generated-job-id"
	}
	return job.ID, nil
}

func TestCreateScheduleTool_Execute_CreatesRecurringJobUsingChatContext(t *testing.T) {
	t.Parallel()

	service := &fakeScheduleService{}
	tool := NewCreateScheduleTool(service)

	ctx := agent.WithTeamContext(context.Background(), "12345", "user-1")
	result, err := tool.Execute(ctx, map[string]interface{}{
		"schedule_type": "cron",
		"cron_expr":     "0 8 * * 1-5",
		"prompt":        "Me mande o resumo diario",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(service.created) != 1 {
		t.Fatalf("expected one created job, got %d", len(service.created))
	}
	job := service.created[0]
	if job.ScheduleType != "cron" {
		t.Fatalf("unexpected schedule type: %q", job.ScheduleType)
	}
	if job.TargetChatID != 12345 {
		t.Fatalf("expected target chat id 12345, got %d", job.TargetChatID)
	}
	if job.OwnerUserID != "user-1" {
		t.Fatalf("expected owner user id user-1, got %q", job.OwnerUserID)
	}
	if job.CronExpr != "0 8 * * 1-5" || job.Prompt != "Me mande o resumo diario" {
		t.Fatalf("unexpected created job payload: %#v", job)
	}
	if !strings.Contains(result, job.ID) {
		t.Fatalf("expected created job id in result, got %q", result)
	}
}

func TestCreateScheduleTool_Execute_CreatesOnceJob(t *testing.T) {
	t.Parallel()

	service := &fakeScheduleService{}
	tool := NewCreateScheduleTool(service)

	ctx := agent.WithTeamContext(context.Background(), "99", "user-2")
	_, err := tool.Execute(ctx, map[string]interface{}{
		"schedule_type": "once",
		"run_at":        "2026-03-12T09:00:00-03:00",
		"prompt":        "Lembre-me da daily",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(service.created) != 1 {
		t.Fatalf("expected one created job, got %d", len(service.created))
	}
	job := service.created[0]
	if job.ScheduleType != "once" {
		t.Fatalf("unexpected schedule type: %q", job.ScheduleType)
	}
	if job.RunAt == nil {
		t.Fatalf("expected parsed run_at in job")
	}
}

func TestCreateScheduleTool_Execute_AllowsExplicitTargetChatID(t *testing.T) {
	t.Parallel()

	service := &fakeScheduleService{}
	tool := NewCreateScheduleTool(service)

	ctx := agent.WithTeamContext(context.Background(), "123", "user-3")
	_, err := tool.Execute(ctx, map[string]interface{}{
		"schedule_type":  "cron",
		"cron_expr":      "0 9 * * *",
		"prompt":         "Mandar status",
		"target_chat_id": float64(777),
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if service.created[0].TargetChatID != 777 {
		t.Fatalf("expected explicit target chat id 777, got %d", service.created[0].TargetChatID)
	}
}

func TestCreateScheduleTool_Execute_ErrorsWithoutTeamContext(t *testing.T) {
	t.Parallel()

	tool := NewCreateScheduleTool(&fakeScheduleService{})

	_, err := tool.Execute(context.Background(), map[string]interface{}{
		"schedule_type": "cron",
		"cron_expr":     "0 8 * * *",
		"prompt":        "Teste",
	})
	if err == nil {
		t.Fatalf("expected missing context error")
	}
}


