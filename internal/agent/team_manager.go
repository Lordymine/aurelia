package agent

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type DefaultTeamManager struct {
	store *SQLiteTaskStore
}

func NewTeamManager(store *SQLiteTaskStore) (TeamManager, error) {
	if store == nil {
		return nil, fmt.Errorf("task store is required")
	}
	return &DefaultTeamManager{store: store}, nil
}

func (m *DefaultTeamManager) CreateTeam(ctx context.Context, teamKey, userID, leadAgent string) (string, error) {
	teamID := uuid.NewString()
	if err := m.store.createTeam(ctx, teamID, teamKey, userID, leadAgent); err != nil {
		return "", err
	}
	return teamID, nil
}

func (m *DefaultTeamManager) GetTeamIDByKey(ctx context.Context, teamKey string) (string, error) {
	return m.store.getTeamIDByKey(ctx, teamKey)
}

func (m *DefaultTeamManager) RegisterTeammate(ctx context.Context, teamID, agentName, roleDescription string) error {
	return m.store.registerTeammate(ctx, teamID, uuid.NewString(), agentName, roleDescription)
}

func (m *DefaultTeamManager) CreateTask(ctx context.Context, task TeamTask, dependsOn []string) error {
	if task.Status == "" {
		task.Status = TaskPending
	}
	return m.store.createTask(ctx, task, dependsOn)
}

func (m *DefaultTeamManager) GetTask(ctx context.Context, teamID, taskID string) (*TeamTask, error) {
	return m.store.getTask(ctx, teamID, taskID)
}

func (m *DefaultTeamManager) ListTasks(ctx context.Context, teamID string) ([]TeamTask, error) {
	return m.store.listTasks(ctx, teamID)
}

func (m *DefaultTeamManager) ClaimNextTask(ctx context.Context, teamID, agentName string) (*TeamTask, error) {
	return m.store.claimNextTask(ctx, teamID, agentName)
}

func (m *DefaultTeamManager) HeartbeatWorker(ctx context.Context, teamID, agentName string) error {
	return m.store.heartbeatWorker(ctx, teamID, agentName)
}

func (m *DefaultTeamManager) CompleteTask(ctx context.Context, teamID, taskID, agentName, result string) error {
	return m.store.completeTask(ctx, teamID, taskID, agentName, result)
}

func (m *DefaultTeamManager) FailTask(ctx context.Context, teamID, taskID, agentName, reason string) error {
	return m.store.failTask(ctx, teamID, taskID, agentName, reason)
}

func (m *DefaultTeamManager) GetTeamStatus(ctx context.Context, teamID string) (string, error) {
	return m.store.getTeamStatus(ctx, teamID)
}

func (m *DefaultTeamManager) SetTeamStatus(ctx context.Context, teamID, status string) error {
	return m.store.setTeamStatus(ctx, teamID, status)
}

func (m *DefaultTeamManager) CancelActiveTasks(ctx context.Context, teamID, reason string) error {
	return m.store.cancelActiveTasks(ctx, teamID, reason)
}

func (m *DefaultTeamManager) PostMessage(ctx context.Context, msg MailMessage) error {
	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}
	return m.store.postMessage(ctx, msg)
}

func (m *DefaultTeamManager) PullMessages(ctx context.Context, teamID, agentName string, limit int) ([]MailMessage, error) {
	if limit <= 0 {
		limit = 20
	}
	return m.store.pullMessages(ctx, teamID, agentName, limit)
}

func (m *DefaultTeamManager) ListEvents(ctx context.Context, teamID string, limit int) ([]TaskEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	return m.store.listEvents(ctx, teamID, limit)
}
