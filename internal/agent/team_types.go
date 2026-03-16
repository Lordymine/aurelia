package agent

import (
	"context"
	"time"
)

type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskBlocked   TaskStatus = "blocked"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
	TaskCancelled TaskStatus = "cancelled"
)

type TeamTask struct {
	ID            string
	TeamID        string
	RunID         string
	ParentTaskID  *string
	Title         string
	Prompt        string
	Workdir       string
	AllowedTools  []string
	AssignedAgent *string
	Status        TaskStatus
	ResultSummary string
	ErrorMessage  string
	CreatedAt     time.Time
	StartedAt     *time.Time
	FinishedAt    *time.Time
}

type MailMessage struct {
	ID         string
	TeamID     string
	FromAgent  string
	ToAgent    string
	TaskID     *string
	Kind       string
	Body       string
	CreatedAt  time.Time
	ConsumedAt *time.Time
}

type TaskEvent struct {
	ID        int64
	TeamID    string
	TaskID    *string
	AgentName string
	EventType string
	Payload   string
	CreatedAt time.Time
}

type TeamStatusSnapshot struct {
	TeamKey    string
	TeamID     string
	TeamStatus string
	Pending    int
	Running    int
	Blocked    int
	Completed  int
	Failed     int
	Cancelled  int
	TotalTasks int
}

type TeamManager interface {
	CreateTeam(ctx context.Context, teamKey, userID, leadAgent string) (string, error)
	GetTeamIDByKey(ctx context.Context, teamKey string) (string, error)
	RegisterTeammate(ctx context.Context, teamID, agentName, roleDescription string) error
	CreateTask(ctx context.Context, task TeamTask, dependsOn []string) error
	GetTask(ctx context.Context, teamID, taskID string) (*TeamTask, error)
	ListTasks(ctx context.Context, teamID string) ([]TeamTask, error)
	ClaimNextTask(ctx context.Context, teamID, agentName string) (*TeamTask, error)
	HeartbeatWorker(ctx context.Context, teamID, agentName string) error
	CompleteTask(ctx context.Context, teamID, taskID, agentName, result string) error
	FailTask(ctx context.Context, teamID, taskID, agentName, reason string) error
	GetTeamStatus(ctx context.Context, teamID string) (string, error)
	SetTeamStatus(ctx context.Context, teamID, status string) error
	CancelActiveTasks(ctx context.Context, teamID, reason string) error
	PostMessage(ctx context.Context, msg MailMessage) error
	PullMessages(ctx context.Context, teamID, agentName string, limit int) ([]MailMessage, error)
	ListEvents(ctx context.Context, teamID string, limit int) ([]TaskEvent, error)
}
