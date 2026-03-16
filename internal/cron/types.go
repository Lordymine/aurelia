package cron

import (
	"context"
	"time"

	"github.com/kocar/aurelia/internal/agent"
)

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

type CronExecution struct {
	ID            string
	JobID         string
	StartedAt     time.Time
	FinishedAt    *time.Time
	Status        string
	OutputSummary string
	ErrorMessage  string
}

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

type AgentExecutor interface {
	Execute(ctx context.Context, systemPrompt string, history []agent.Message, allowedTools []string) ([]agent.Message, string, error)
}

type Runtime interface {
	ExecuteJob(ctx context.Context, job CronJob) (string, error)
}

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now().UTC()
}

type SchedulerConfig struct {
	PollInterval time.Duration
}


