package cron

import (
	"context"
	"database/sql"
	"time"

	"github.com/kocar/aurelia/internal/bridge"
)

// CronJob represents a scheduled job.
type CronJob struct {
	ID           string
	OwnerUserID  string
	TargetChatID int64
	AgentName    string // agent from registry to execute this job
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
	SessionID     string
	CostUSD       float64
	TokensUsed    int
}

// ExecutionResult holds the output and metadata from a job execution.
type ExecutionResult struct {
	Output    string
	SessionID string
	CostUSD   float64
	NumTurns  int
}

// Store persists cron jobs and executions.
type Store interface {
	CreateJob(ctx context.Context, job CronJob) error
	UpdateJob(ctx context.Context, job CronJob) error
	DeleteJob(ctx context.Context, jobID string) error
	GetJob(ctx context.Context, jobID string) (*CronJob, error)
	ResolveJobID(ctx context.Context, prefix string) (string, error)
	ListJobsByChat(ctx context.Context, chatID int64) ([]CronJob, error)
	ListDueJobs(ctx context.Context, now time.Time, limit int) ([]CronJob, error)
	RecordExecution(ctx context.Context, exec CronExecution) error
	RecordExecutionTx(ctx context.Context, tx *sql.Tx, exec CronExecution) error
	UpdateJobTx(ctx context.Context, tx *sql.Tx, job CronJob) error
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
	ListExecutionsByJob(ctx context.Context, jobID string) ([]CronExecution, error)
}

// BridgeExecutor is the interface for executing a request via the Claude Code bridge.
type BridgeExecutor interface {
	Execute(ctx context.Context, req bridge.Request) (*bridge.Event, error)
}

// Runtime executes a cron job and returns its result.
type Runtime interface {
	ExecuteJob(ctx context.Context, job CronJob) (*ExecutionResult, error)
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
