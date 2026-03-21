package cron

import (
	"context"
	"time"
)

// CronJob represents a scheduled job.
type CronJob struct {
	ID           string
	OwnerUserID  string
	TargetChatID int64
	ScheduleType string
	CronExpr     string
	RunAt        *time.Time
	Prompt       string
	Active       bool
	LastRunAt    *time.Time
	NextRunAt    *time.Time
	LastStatus   string
	LastError    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CronExecution records the result of a job run.
type CronExecution struct {
	ID            string
	JobID         string
	StartedAt     time.Time
	FinishedAt    *time.Time
	Status        string
	OutputSummary string
	ErrorMessage  string
}

// Store persists cron jobs and executions.
type Store interface {
	CreateJob(ctx context.Context, job CronJob) error
	UpdateJob(ctx context.Context, job CronJob) error
	DeleteJob(ctx context.Context, jobID string) error
	GetJob(ctx context.Context, jobID string) (*CronJob, error)
	ListJobsByChat(ctx context.Context, chatID int64) ([]CronJob, error)
	ListDueJobs(ctx context.Context, now time.Time, limit int) ([]CronJob, error)
	RecordExecution(ctx context.Context, exec CronExecution) error
	ListExecutionsByJob(ctx context.Context, jobID string) ([]CronExecution, error)
}

// BridgeExecutor is the interface for executing a prompt via the Claude Code bridge.
type BridgeExecutor interface {
	Execute(ctx context.Context, systemPrompt string, userPrompt string) (string, error)
}

// Runtime executes a cron job and returns its output.
type Runtime interface {
	ExecuteJob(ctx context.Context, job CronJob) (string, error)
}

// Clock abstracts time for testing.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now().UTC()
}

// SchedulerConfig configures the cron scheduler.
type SchedulerConfig struct {
	PollInterval time.Duration
}
