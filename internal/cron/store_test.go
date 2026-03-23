package cron

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func newTestCronStore(t *testing.T) *SQLiteCronStore {
	t.Helper()

	store, err := NewSQLiteCronStore(filepath.Join(t.TempDir(), "cron.db"))
	if err != nil {
		t.Fatalf("NewSQLiteCronStore() error = %v", err)
	}

	t.Cleanup(func() {
		_ = store.Close()
	})

	return store
}

func TestSQLiteCronStore_CreateJobAndGetJob(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newTestCronStore(t)

	nextRunAt := time.Date(2026, 3, 12, 8, 0, 0, 0, time.UTC)
	job := CronJob{
		ID:           "job-1",
		OwnerUserID:  "user-1",
		TargetChatID: 12345,
		ScheduleType: "cron",
		CronExpr:     "0 8 * * 1-5",
		Prompt:       "Me mande o resumo da manha",
		Active:       true,
		NextRunAt:    &nextRunAt,
		LastStatus:   "idle",
	}

	if err := store.CreateJob(ctx, job); err != nil {
		t.Fatalf("CreateJob() error = %v", err)
	}

	got, err := store.GetJob(ctx, "job-1")
	if err != nil {
		t.Fatalf("GetJob() error = %v", err)
	}
	if got == nil {
		t.Fatalf("expected persisted job, got nil")
	}
	if got.Prompt != job.Prompt {
		t.Fatalf("unexpected prompt: got %q want %q", got.Prompt, job.Prompt)
	}
	if got.NextRunAt == nil || !got.NextRunAt.Equal(nextRunAt) {
		t.Fatalf("unexpected next run at: %#v", got.NextRunAt)
	}
}

func TestSQLiteCronStore_ListJobsByChat_ReturnsOnlyChatJobs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newTestCronStore(t)

	firstRun := time.Date(2026, 3, 12, 8, 0, 0, 0, time.UTC)
	for _, job := range []CronJob{
		{
			ID:           "job-chat-a",
			OwnerUserID:  "user-1",
			TargetChatID: 100,
			ScheduleType: "cron",
			CronExpr:     "0 8 * * *",
			Prompt:       "job a",
			Active:       true,
			NextRunAt:    &firstRun,
			LastStatus:   "idle",
		},
		{
			ID:           "job-chat-b",
			OwnerUserID:  "user-2",
			TargetChatID: 200,
			ScheduleType: "cron",
			CronExpr:     "0 9 * * *",
			Prompt:       "job b",
			Active:       true,
			NextRunAt:    &firstRun,
			LastStatus:   "idle",
		},
	} {
		if err := store.CreateJob(ctx, job); err != nil {
			t.Fatalf("CreateJob(%s) error = %v", job.ID, err)
		}
	}

	jobs, err := store.ListJobsByChat(ctx, 100)
	if err != nil {
		t.Fatalf("ListJobsByChat() error = %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job for chat 100, got %d", len(jobs))
	}
	if jobs[0].ID != "job-chat-a" {
		t.Fatalf("unexpected job id: %q", jobs[0].ID)
	}
}

func TestSQLiteCronStore_ListDueJobs_ReturnsOnlyActiveDueJobsOrderedByNextRun(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newTestCronStore(t)

	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	pastA := now.Add(-10 * time.Minute)
	pastB := now.Add(-2 * time.Minute)
	future := now.Add(20 * time.Minute)

	jobs := []CronJob{
		{
			ID:           "job-due-oldest",
			OwnerUserID:  "user-1",
			TargetChatID: 100,
			ScheduleType: "cron",
			CronExpr:     "*/5 * * * *",
			Prompt:       "oldest due",
			Active:       true,
			NextRunAt:    &pastA,
			LastStatus:   "idle",
		},
		{
			ID:           "job-due-newest",
			OwnerUserID:  "user-1",
			TargetChatID: 100,
			ScheduleType: "cron",
			CronExpr:     "*/5 * * * *",
			Prompt:       "newest due",
			Active:       true,
			NextRunAt:    &pastB,
			LastStatus:   "idle",
		},
		{
			ID:           "job-future",
			OwnerUserID:  "user-1",
			TargetChatID: 100,
			ScheduleType: "cron",
			CronExpr:     "*/5 * * * *",
			Prompt:       "future",
			Active:       true,
			NextRunAt:    &future,
			LastStatus:   "idle",
		},
		{
			ID:           "job-inactive",
			OwnerUserID:  "user-1",
			TargetChatID: 100,
			ScheduleType: "cron",
			CronExpr:     "*/5 * * * *",
			Prompt:       "inactive",
			Active:       false,
			NextRunAt:    &pastA,
			LastStatus:   "idle",
		},
	}

	for _, job := range jobs {
		if err := store.CreateJob(ctx, job); err != nil {
			t.Fatalf("CreateJob(%s) error = %v", job.ID, err)
		}
	}

	dueJobs, err := store.ListDueJobs(ctx, now, 10)
	if err != nil {
		t.Fatalf("ListDueJobs() error = %v", err)
	}
	if len(dueJobs) != 2 {
		t.Fatalf("expected 2 due jobs, got %d", len(dueJobs))
	}
	if dueJobs[0].ID != "job-due-oldest" || dueJobs[1].ID != "job-due-newest" {
		t.Fatalf("unexpected due order: %#v", dueJobs)
	}
}

func TestSQLiteCronStore_UpdateJob_PersistsPauseAndNextRun(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newTestCronStore(t)

	firstRun := time.Date(2026, 3, 12, 8, 0, 0, 0, time.UTC)
	job := CronJob{
		ID:           "job-update",
		OwnerUserID:  "user-1",
		TargetChatID: 100,
		ScheduleType: "cron",
		CronExpr:     "0 8 * * *",
		Prompt:       "job update",
		Active:       true,
		NextRunAt:    &firstRun,
		LastStatus:   "idle",
	}
	if err := store.CreateJob(ctx, job); err != nil {
		t.Fatalf("CreateJob() error = %v", err)
	}

	secondRun := time.Date(2026, 3, 13, 8, 0, 0, 0, time.UTC)
	job.Active = false
	job.NextRunAt = &secondRun
	job.LastStatus = "failed"
	job.LastError = "last failure"

	if err := store.UpdateJob(ctx, job); err != nil {
		t.Fatalf("UpdateJob() error = %v", err)
	}

	got, err := store.GetJob(ctx, "job-update")
	if err != nil {
		t.Fatalf("GetJob() after update error = %v", err)
	}
	if got == nil {
		t.Fatalf("expected updated job")
	}
	if got.Active {
		t.Fatalf("expected job to be paused")
	}
	if got.NextRunAt == nil || !got.NextRunAt.Equal(secondRun) {
		t.Fatalf("unexpected next run after update: %#v", got.NextRunAt)
	}
	if got.LastError != "last failure" {
		t.Fatalf("unexpected last error: %q", got.LastError)
	}
}

func TestSQLiteCronStore_RecordExecutionAndListExecutionsByJob(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := newTestCronStore(t)

	runAt := time.Date(2026, 3, 12, 8, 0, 0, 0, time.UTC)
	job := CronJob{
		ID:           "job-exec",
		OwnerUserID:  "user-1",
		TargetChatID: 100,
		ScheduleType: "once",
		RunAt:        &runAt,
		Prompt:       "run once",
		Active:       true,
		NextRunAt:    &runAt,
		LastStatus:   "idle",
	}
	if err := store.CreateJob(ctx, job); err != nil {
		t.Fatalf("CreateJob() error = %v", err)
	}

	finishedAt := runAt.Add(5 * time.Second)
	exec := CronExecution{
		ID:            "exec-1",
		JobID:         "job-exec",
		StartedAt:     runAt,
		FinishedAt:    &finishedAt,
		Status:        "success",
		OutputSummary: "job finished",
	}
	if err := store.RecordExecution(ctx, exec); err != nil {
		t.Fatalf("RecordExecution() error = %v", err)
	}

	executions, err := store.ListExecutionsByJob(ctx, "job-exec")
	if err != nil {
		t.Fatalf("ListExecutionsByJob() error = %v", err)
	}
	if len(executions) != 1 {
		t.Fatalf("expected 1 execution, got %d", len(executions))
	}
	if executions[0].Status != "success" {
		t.Fatalf("unexpected execution status: %q", executions[0].Status)
	}
	if executions[0].OutputSummary != "job finished" {
		t.Fatalf("unexpected execution summary: %q", executions[0].OutputSummary)
	}
}

func TestDueJobsIndex(t *testing.T) {
	store, err := NewSQLiteCronStore(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	var count int
	err = store.db.QueryRow(
		"SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_cron_jobs_due'",
	).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("expected idx_cron_jobs_due index to exist")
	}
}
