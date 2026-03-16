package cron

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func newTestCronService(t *testing.T) *Service {
	t.Helper()

	store, err := NewSQLiteCronStore(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatalf("NewSQLiteCronStore() error = %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	return NewService(store, staticClock{now: time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)})
}

func TestService_CreateJob_ComputesNextRunForRecurring(t *testing.T) {
	t.Parallel()

	service := newTestCronService(t)

	jobID, err := service.CreateJob(context.Background(), CronJob{
		OwnerUserID:  "user-1",
		TargetChatID: 100,
		ScheduleType: "cron",
		CronExpr:     "*/5 * * * *",
		Prompt:       "Resumo diario",
		Active:       true,
	})
	if err != nil {
		t.Fatalf("CreateJob() error = %v", err)
	}
	if jobID == "" {
		t.Fatalf("expected non-empty job id")
	}

	job, err := service.store.GetJob(context.Background(), jobID)
	if err != nil {
		t.Fatalf("GetJob() error = %v", err)
	}
	if job.NextRunAt == nil {
		t.Fatalf("expected computed next run")
	}
}

func TestService_PauseResumeDeleteJob(t *testing.T) {
	t.Parallel()

	service := newTestCronService(t)
	ctx := context.Background()

	jobID, err := service.CreateJob(ctx, CronJob{
		OwnerUserID:  "user-1",
		TargetChatID: 100,
		ScheduleType: "once",
		RunAt:        timePtr(time.Date(2026, 3, 13, 9, 0, 0, 0, time.UTC)),
		Prompt:       "Lembrete",
		Active:       true,
	})
	if err != nil {
		t.Fatalf("CreateJob() error = %v", err)
	}

	if err := service.PauseJob(ctx, jobID); err != nil {
		t.Fatalf("PauseJob() error = %v", err)
	}
	job, _ := service.store.GetJob(ctx, jobID)
	if job.Active {
		t.Fatalf("expected paused job")
	}

	if err := service.ResumeJob(ctx, jobID); err != nil {
		t.Fatalf("ResumeJob() error = %v", err)
	}
	job, _ = service.store.GetJob(ctx, jobID)
	if !job.Active {
		t.Fatalf("expected resumed job")
	}

	if err := service.DeleteJob(ctx, jobID); err != nil {
		t.Fatalf("DeleteJob() error = %v", err)
	}
	job, err = service.store.GetJob(ctx, jobID)
	if err != nil {
		t.Fatalf("GetJob() after delete error = %v", err)
	}
	if job != nil {
		t.Fatalf("expected deleted job to be absent")
	}
}

func timePtr(v time.Time) *time.Time {
	return &v
}
