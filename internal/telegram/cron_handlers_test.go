package telegram

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/cron"
)

type fakeCronCommandService struct {
	addRecurringCalls []struct {
		userID string
		chatID int64
		expr   string
		prompt string
	}
	addOnceCalls []struct {
		userID    string
		chatID    int64
		timestamp string
		prompt    string
	}
	listCalls []struct {
		chatID int64
	}
	pauseCalls  []string
	resumeCalls []string
	deleteCalls []string

	jobs    []cron.CronJob
	addErr  error
	listErr error
}

func (f *fakeCronCommandService) AddRecurringJob(ctx context.Context, userID string, chatID int64, expr, prompt string) (string, error) {
	f.addRecurringCalls = append(f.addRecurringCalls, struct {
		userID string
		chatID int64
		expr   string
		prompt string
	}{userID: userID, chatID: chatID, expr: expr, prompt: prompt})
	if f.addErr != nil {
		return "", f.addErr
	}
	return "job-recurring-1", nil
}

func (f *fakeCronCommandService) AddOnceJob(ctx context.Context, userID string, chatID int64, timestamp, prompt string) (string, error) {
	f.addOnceCalls = append(f.addOnceCalls, struct {
		userID    string
		chatID    int64
		timestamp string
		prompt    string
	}{userID: userID, chatID: chatID, timestamp: timestamp, prompt: prompt})
	if f.addErr != nil {
		return "", f.addErr
	}
	return "job-once-1", nil
}

func (f *fakeCronCommandService) ListJobs(ctx context.Context, chatID int64) ([]cron.CronJob, error) {
	f.listCalls = append(f.listCalls, struct{ chatID int64 }{chatID: chatID})
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.jobs, nil
}

func (f *fakeCronCommandService) PauseJob(ctx context.Context, jobID string) error {
	f.pauseCalls = append(f.pauseCalls, jobID)
	return nil
}

func (f *fakeCronCommandService) ResumeJob(ctx context.Context, jobID string) error {
	f.resumeCalls = append(f.resumeCalls, jobID)
	return nil
}

func (f *fakeCronCommandService) DeleteJob(ctx context.Context, jobID string) error {
	f.deleteCalls = append(f.deleteCalls, jobID)
	return nil
}

func TestCronCommandHandler_HandleText_AddRecurring(t *testing.T) {
	t.Parallel()

	service := &fakeCronCommandService{}
	handler := NewCronCommandHandler(service)

	reply, err := handler.HandleText(context.Background(), "user-1", 12345, `/cron add "0 8 * * 1-5" "Me mande o resumo da manha"`)
	if err != nil {
		t.Fatalf("HandleText() error = %v", err)
	}
	if len(service.addRecurringCalls) != 1 {
		t.Fatalf("expected one recurring add call, got %d", len(service.addRecurringCalls))
	}
	call := service.addRecurringCalls[0]
	if call.expr != "0 8 * * 1-5" || call.prompt != "Me mande o resumo da manha" {
		t.Fatalf("unexpected recurring add args: %#v", call)
	}
	if !strings.Contains(reply, "job-recurring-1") {
		t.Fatalf("expected reply to include created job id, got %q", reply)
	}
}

func TestCronCommandHandler_HandleText_AddOnce(t *testing.T) {
	t.Parallel()

	service := &fakeCronCommandService{}
	handler := NewCronCommandHandler(service)

	reply, err := handler.HandleText(context.Background(), "user-1", 12345, `/cron once "2026-03-12T09:00:00-03:00" "Lembre-me da daily"`)
	if err != nil {
		t.Fatalf("HandleText() error = %v", err)
	}
	if len(service.addOnceCalls) != 1 {
		t.Fatalf("expected one once add call, got %d", len(service.addOnceCalls))
	}
	call := service.addOnceCalls[0]
	if call.timestamp != "2026-03-12T09:00:00-03:00" || call.prompt != "Lembre-me da daily" {
		t.Fatalf("unexpected once add args: %#v", call)
	}
	if !strings.Contains(reply, "job-once-1") {
		t.Fatalf("expected reply to include created once job id, got %q", reply)
	}
}

func TestCronCommandHandler_HandleText_ListJobs(t *testing.T) {
	t.Parallel()

	service := &fakeCronCommandService{
		jobs: []cron.CronJob{
			{ID: "job-a", ScheduleType: "cron", CronExpr: "0 8 * * *", Prompt: "job a", Active: true, LastStatus: "idle"},
			{ID: "job-b", ScheduleType: "once", Prompt: "job b", Active: false, LastStatus: "success"},
		},
	}
	handler := NewCronCommandHandler(service)

	reply, err := handler.HandleText(context.Background(), "user-1", 12345, `/cron list`)
	if err != nil {
		t.Fatalf("HandleText() error = %v", err)
	}
	if len(service.listCalls) != 1 {
		t.Fatalf("expected one list call, got %d", len(service.listCalls))
	}
	if !strings.Contains(reply, "job-a") || !strings.Contains(reply, "job-b") {
		t.Fatalf("expected listed jobs in reply, got %q", reply)
	}
}

func TestCronCommandHandler_HandleText_PauseResumeDelete(t *testing.T) {
	t.Parallel()

	service := &fakeCronCommandService{}
	handler := NewCronCommandHandler(service)

	if _, err := handler.HandleText(context.Background(), "user-1", 12345, `/cron pause job-a`); err != nil {
		t.Fatalf("pause HandleText() error = %v", err)
	}
	if _, err := handler.HandleText(context.Background(), "user-1", 12345, `/cron resume job-a`); err != nil {
		t.Fatalf("resume HandleText() error = %v", err)
	}
	if _, err := handler.HandleText(context.Background(), "user-1", 12345, `/cron del job-a`); err != nil {
		t.Fatalf("delete HandleText() error = %v", err)
	}

	if len(service.pauseCalls) != 1 || service.pauseCalls[0] != "job-a" {
		t.Fatalf("unexpected pause calls: %#v", service.pauseCalls)
	}
	if len(service.resumeCalls) != 1 || service.resumeCalls[0] != "job-a" {
		t.Fatalf("unexpected resume calls: %#v", service.resumeCalls)
	}
	if len(service.deleteCalls) != 1 || service.deleteCalls[0] != "job-a" {
		t.Fatalf("unexpected delete calls: %#v", service.deleteCalls)
	}
}

func TestCronCommandHandler_HandleText_ReturnsUsageForInvalidCommand(t *testing.T) {
	t.Parallel()

	handler := NewCronCommandHandler(&fakeCronCommandService{})

	reply, err := handler.HandleText(context.Background(), "user-1", 12345, `/cron unknown`)
	if err != nil {
		t.Fatalf("HandleText() error = %v", err)
	}
	if !strings.Contains(strings.ToLower(reply), "uso") {
		t.Fatalf("expected usage reply, got %q", reply)
	}
}

func TestCronCommandHandler_HandleText_PropagatesServiceError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("invalid cron expression")
	service := &fakeCronCommandService{addErr: expectedErr}
	handler := NewCronCommandHandler(service)

	_, err := handler.HandleText(context.Background(), "user-1", 12345, `/cron add "bad expr" "prompt"`)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected service error %v, got %v", expectedErr, err)
	}
}


