package agent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaskExecutor interface {
	ExecuteTask(ctx context.Context, task TeamTask) (string, error)
}

type WorkerRuntime struct {
	agentName string
	manager   TeamManager
	executor  TaskExecutor
}

func NewWorkerRuntime(agentName string, manager TeamManager, executor TaskExecutor) *WorkerRuntime {
	return &WorkerRuntime{
		agentName: agentName,
		manager:   manager,
		executor:  executor,
	}
}

func (w *WorkerRuntime) RunOnce(ctx context.Context, teamID string) (bool, error) {
	if w.manager == nil {
		return false, fmt.Errorf("worker manager is required")
	}
	if w.executor == nil {
		return false, fmt.Errorf("worker executor is required")
	}

	task, err := w.manager.ClaimNextTask(ctx, teamID, w.agentName)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, nil
	}

	taskCtx := WithTaskContext(ctx, teamID, task.ID)
	if task.RunID != "" {
		taskCtx = WithRunContext(taskCtx, task.RunID)
	}
	if task.Workdir != "" {
		taskCtx = WithWorkdirContext(taskCtx, task.Workdir)
	}
	stopHeartbeat := w.startHeartbeat(taskCtx, teamID)
	defer stopHeartbeat()

	result, err := w.executor.ExecuteTask(taskCtx, *task)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return true, nil
		}
		failErr := w.manager.FailTask(ctx, teamID, task.ID, w.agentName, err.Error())
		if failErr != nil {
			if currentTask, getErr := w.manager.GetTask(ctx, teamID, task.ID); getErr == nil && currentTask != nil && currentTask.Status == TaskCancelled {
				return true, nil
			}
			return true, fmt.Errorf("task execution failed: %v; failTask: %w", err, failErr)
		}
		taskID := task.ID
		_ = w.manager.PostMessage(ctx, MailMessage{
			ID:        uuid.NewString(),
			TeamID:    teamID,
			FromAgent: w.agentName,
			ToAgent:   MasterAgentName,
			TaskID:    &taskID,
			Kind:      "blocker",
			Body:      err.Error(),
		})
		return true, nil
	}

	if err := w.manager.CompleteTask(ctx, teamID, task.ID, w.agentName, result); err != nil {
		return true, err
	}

	taskID := task.ID
	if err := w.manager.PostMessage(ctx, MailMessage{
		ID:        uuid.NewString(),
		TeamID:    teamID,
		FromAgent: w.agentName,
		ToAgent:   MasterAgentName,
		TaskID:    &taskID,
		Kind:      "result",
		Body:      result,
	}); err != nil {
		return true, err
	}

	return true, nil
}

func (w *WorkerRuntime) startHeartbeat(ctx context.Context, teamID string) func() {
	if w.manager == nil {
		return func() {}
	}

	stop := make(chan struct{})
	ticker := time.NewTicker(25 * time.Millisecond)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-stop:
				return
			case <-ticker.C:
				_ = w.manager.HeartbeatWorker(ctx, teamID, w.agentName)
			}
		}
	}()

	return func() {
		close(stop)
	}
}

func (w *WorkerRuntime) RunUntilIdle(ctx context.Context, teamID string) (int, error) {
	processedCount := 0

	for {
		processed, err := w.RunOnce(ctx, teamID)
		if err != nil {
			return processedCount, err
		}
		if !processed {
			return processedCount, nil
		}
		processedCount++
	}
}
