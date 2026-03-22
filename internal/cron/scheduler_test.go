package cron

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

type fakeCronStore struct {
	dueJobs         []CronJob
	jobs            map[string]CronJob
	executions      []CronExecution
	listDueJobsErr  error
	updateJobErr    error
	recordExecErr   error
	listDueJobsSeen []time.Time
}

func (f *fakeCronStore) CreateJob(ctx context.Context, job CronJob) error {
	if f.jobs == nil {
		f.jobs = make(map[string]CronJob)
	}
	f.jobs[job.ID] = job
	return nil
}

func (f *fakeCronStore) UpdateJob(ctx context.Context, job CronJob) error {
	if f.updateJobErr != nil {
		return f.updateJobErr
	}
	if f.jobs == nil {
		f.jobs = make(map[string]CronJob)
	}
	f.jobs[job.ID] = job
	return nil
}

func (f *fakeCronStore) DeleteJob(ctx context.Context, jobID string) error {
	delete(f.jobs, jobID)
	return nil
}

func (f *fakeCronStore) GetJob(ctx context.Context, jobID string) (*CronJob, error) {
	job, ok := f.jobs[jobID]
	if !ok {
		return nil, nil
	}
	copy := job
	return &copy, nil
}

func (f *fakeCronStore) ListJobsByChat(ctx context.Context, chatID int64) ([]CronJob, error) {
	var result []CronJob
	for _, job := range f.jobs {
		if job.TargetChatID == chatID {
			result = append(result, job)
		}
	}
	return result, nil
}

func (f *fakeCronStore) ListDueJobs(ctx context.Context, now time.Time, limit int) ([]CronJob, error) {
	f.listDueJobsSeen = append(f.listDueJobsSeen, now)
	if f.listDueJobsErr != nil {
		return nil, f.listDueJobsErr
	}
	return f.dueJobs, nil
}

func (f *fakeCronStore) RecordExecution(ctx context.Context, exec CronExecution) error {
	if f.recordExecErr != nil {
		return f.recordExecErr
	}
	f.executions = append(f.executions, exec)
	return nil
}

func (f *fakeCronStore) ResolveJobID(ctx context.Context, prefix string) (string, error) {
	for id := range f.jobs {
		if len(id) >= len(prefix) && id[:len(prefix)] == prefix {
			return id, nil
		}
	}
	return "", fmt.Errorf("cron job %s not found", prefix)
}

func (f *fakeCronStore) ListExecutionsByJob(ctx context.Context, jobID string) ([]CronExecution, error) {
	var result []CronExecution
	for _, exec := range f.executions {
		if exec.JobID == jobID {
			result = append(result, exec)
		}
	}
	return result, nil
}

type fakeCronRuntime struct {
	results map[string]string
	errors  map[string]error
	seen    []string
}

func (f *fakeCronRuntime) ExecuteJob(ctx context.Context, job CronJob) (*ExecutionResult, error) {
	f.seen = append(f.seen, job.ID)
	if err := f.errors[job.ID]; err != nil {
		return nil, err
	}
	return &ExecutionResult{Output: f.results[job.ID]}, nil
}

type staticClock struct {
	now time.Time
}

func (c staticClock) Now() time.Time {
	return c.now
}

func TestComputeNextRun_EveryFiveMinutes(t *testing.T) {
	base := time.Date(2026, 3, 22, 10, 0, 0, 0, time.UTC)
	next, err := computeNextRun("*/5 * * * *", base)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 3, 22, 10, 5, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("got %v, want %v", next, want)
	}
}

func TestComputeNextRun_DailyAt9AM(t *testing.T) {
	base := time.Date(2026, 3, 22, 10, 0, 0, 0, time.UTC)
	next, err := computeNextRun("0 9 * * *", base)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 3, 23, 9, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("got %v, want %v", next, want)
	}
}

func TestComputeNextRun_WeekdaysOnly(t *testing.T) {
	// Friday 2026-03-27
	base := time.Date(2026, 3, 27, 18, 0, 0, 0, time.UTC)
	next, err := computeNextRun("0 9 * * 1-5", base)
	if err != nil {
		t.Fatal(err)
	}
	// Should skip Sat+Sun → Monday 2026-03-30
	if next.Weekday() != time.Monday {
		t.Errorf("expected Monday, got %v (%v)", next.Weekday(), next)
	}
}

func TestComputeNextRun_InvalidExpr(t *testing.T) {
	_, err := computeNextRun("not a cron", time.Now())
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestScheduler_RunDueJobs_ExecutesDueRecurringJobAndSchedulesNextRun(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	nextExpected := time.Date(2026, 3, 12, 10, 5, 0, 0, time.UTC)

	store := &fakeCronStore{
		jobs: map[string]CronJob{},
		dueJobs: []CronJob{
			{
				ID:           "job-recurring",
				OwnerUserID:  "user-1",
				TargetChatID: 100,
				ScheduleType: "cron",
				CronExpr:     "*/5 * * * *",
				Prompt:       "check status",
				Active:       true,
				LastStatus:   "idle",
			},
		},
	}
	runtime := &fakeCronRuntime{
		results: map[string]string{"job-recurring": "all good"},
	}

	scheduler, err := NewScheduler(store, runtime, staticClock{now: now}, SchedulerConfig{PollInterval: time.Minute})
	if err != nil {
		t.Fatalf("NewScheduler() error = %v", err)
	}

	processed, err := scheduler.RunDueJobs(context.Background())
	if err != nil {
		t.Fatalf("RunDueJobs() error = %v", err)
	}
	if processed != 1 {
		t.Fatalf("expected 1 processed job, got %d", processed)
	}
	if len(runtime.seen) != 1 || runtime.seen[0] != "job-recurring" {
		t.Fatalf("unexpected executed jobs: %#v", runtime.seen)
	}

	job := store.jobs["job-recurring"]
	if job.LastStatus != "success" {
		t.Fatalf("expected job status success, got %q", job.LastStatus)
	}
	if job.NextRunAt == nil || !job.NextRunAt.Equal(nextExpected) {
		t.Fatalf("unexpected next run at: %#v", job.NextRunAt)
	}
	if len(store.executions) != 1 || store.executions[0].Status != "success" {
		t.Fatalf("expected one successful execution, got %#v", store.executions)
	}
}

func TestScheduler_RunDueJobs_DeactivatesOnceJobAfterSuccess(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	runAt := now.Add(-time.Minute)

	store := &fakeCronStore{
		jobs: map[string]CronJob{},
		dueJobs: []CronJob{
			{
				ID:           "job-once",
				OwnerUserID:  "user-1",
				TargetChatID: 100,
				ScheduleType: "once",
				RunAt:        &runAt,
				Prompt:       "remind me",
				Active:       true,
				NextRunAt:    &runAt,
				LastStatus:   "idle",
			},
		},
	}
	runtime := &fakeCronRuntime{
		results: map[string]string{"job-once": "done once"},
	}

	scheduler, err := NewScheduler(store, runtime, staticClock{now: now}, SchedulerConfig{PollInterval: time.Minute})
	if err != nil {
		t.Fatalf("NewScheduler() error = %v", err)
	}

	processed, err := scheduler.RunDueJobs(context.Background())
	if err != nil {
		t.Fatalf("RunDueJobs() error = %v", err)
	}
	if processed != 1 {
		t.Fatalf("expected 1 processed job, got %d", processed)
	}

	job := store.jobs["job-once"]
	if job.Active {
		t.Fatalf("expected once job to be deactivated after success")
	}
	if job.NextRunAt != nil {
		t.Fatalf("expected once job next run to be cleared, got %#v", job.NextRunAt)
	}
}

func TestScheduler_RunDueJobs_RecordsFailureAndKeepsRecurringJobActive(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)

	store := &fakeCronStore{
		jobs: map[string]CronJob{},
		dueJobs: []CronJob{
			{
				ID:           "job-failing",
				OwnerUserID:  "user-1",
				TargetChatID: 100,
				ScheduleType: "cron",
				CronExpr:     "*/5 * * * *",
				Prompt:       "this fails",
				Active:       true,
				LastStatus:   "idle",
			},
		},
	}
	runtime := &fakeCronRuntime{
		errors: map[string]error{"job-failing": errors.New("runtime failed")},
	}

	scheduler, err := NewScheduler(store, runtime, staticClock{now: now}, SchedulerConfig{PollInterval: time.Minute})
	if err != nil {
		t.Fatalf("NewScheduler() error = %v", err)
	}

	processed, err := scheduler.RunDueJobs(context.Background())
	if err != nil {
		t.Fatalf("RunDueJobs() error = %v", err)
	}
	if processed != 1 {
		t.Fatalf("expected 1 processed job, got %d", processed)
	}

	job := store.jobs["job-failing"]
	if job.LastStatus != "failed" {
		t.Fatalf("expected failed status, got %q", job.LastStatus)
	}
	if job.LastError == "" {
		t.Fatalf("expected last error to be persisted")
	}
	if !job.Active {
		t.Fatalf("expected recurring failing job to remain active")
	}
	if len(store.executions) != 1 || store.executions[0].Status != "failed" {
		t.Fatalf("expected one failed execution, got %#v", store.executions)
	}
}

func TestScheduler_RunDueJobs_PropagatesStoreError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("db unavailable")
	store := &fakeCronStore{
		jobs:           map[string]CronJob{},
		listDueJobsErr: expectedErr,
	}
	runtime := &fakeCronRuntime{}

	scheduler, err := NewScheduler(store, runtime, staticClock{now: time.Now().UTC()}, SchedulerConfig{PollInterval: time.Minute})
	if err != nil {
		t.Fatalf("NewScheduler() error = %v", err)
	}

	_, err = scheduler.RunDueJobs(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected store error %v, got %v", expectedErr, err)
	}
}
