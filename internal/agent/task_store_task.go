package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func (s *SQLiteTaskStore) createTask(ctx context.Context, task TeamTask, dependsOn []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin create task tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	status, blockReason, err := s.resolveCreateTaskStateTx(ctx, tx, task, dependsOn)
	if err != nil {
		return err
	}
	if err := s.insertTaskTx(ctx, tx, task, status, blockReason, dependsOn); err != nil {
		return err
	}
	if err := s.recordTaskCreatedEventsTx(ctx, tx, task, status, blockReason); err != nil {
		return err
	}
	if err := s.reopenRecoveryDependentsTx(ctx, tx, task); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit create task tx: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) claimNextTask(ctx context.Context, teamID, agentName string) (*TeamTask, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin claim tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.requeueExpiredRunningTasksTx(ctx, tx, teamID); err != nil {
		return nil, err
	}
	if err := s.renewWorkerLeaseTx(ctx, tx, teamID, agentName); err != nil {
		return nil, err
	}
	if paused, err := s.isTeamUnavailableTx(ctx, tx, teamID); err != nil {
		return nil, err
	} else if paused {
		return nil, nil
	}

	task, err := s.selectNextClaimableTaskTx(ctx, tx, teamID, agentName)
	if err != nil || task == nil {
		return task, err
	}

	now := time.Now().UTC()
	claimed, err := s.markTaskClaimedTx(ctx, tx, task.ID, agentName, now)
	if err != nil || !claimed {
		return nil, err
	}
	if err := s.insertTaskEventTx(ctx, tx, TaskEvent{
		TeamID:    task.TeamID,
		TaskID:    &task.ID,
		AgentName: agentName,
		EventType: "task_claimed",
		Payload:   task.Title,
	}); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit claim tx: %w", err)
	}

	task.AssignedAgent = &agentName
	task.Status = TaskRunning
	task.StartedAt = &now
	return task, nil
}

func (s *SQLiteTaskStore) isTeamUnavailableTx(ctx context.Context, tx *sql.Tx, teamID string) (bool, error) {
	row := tx.QueryRowContext(ctx, `SELECT status FROM teams WHERE id = ?`, teamID)
	var status string
	if err := row.Scan(&status); err != nil {
		return false, fmt.Errorf("get team status for claim: %w", err)
	}
	return status != TeamStatusActive, nil
}

func (s *SQLiteTaskStore) resolveCreateTaskStateTx(ctx context.Context, tx *sql.Tx, task TeamTask, dependsOn []string) (TaskStatus, string, error) {
	if err := s.validateTaskDependenciesTx(ctx, tx, task, dependsOn); err != nil {
		return "", "", err
	}

	status := task.Status
	blockReason := task.ErrorMessage
	if len(dependsOn) > 0 {
		blocked, reason, err := s.evaluateDependencyStateTx(ctx, tx, dependsOn)
		if err != nil {
			return "", "", err
		}
		if blocked {
			status = TaskBlocked
			blockReason = reason
		}
	}

	return status, blockReason, nil
}

func (s *SQLiteTaskStore) insertTaskTx(ctx context.Context, tx *sql.Tx, task TeamTask, status TaskStatus, blockReason string, dependsOn []string) error {
	allowedToolsJSON, err := json.Marshal(task.AllowedTools)
	if err != nil {
		return fmt.Errorf("marshal task allowed tools: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO tasks (id, team_id, run_id, parent_task_id, title, prompt, working_dir, allowed_tools, assigned_agent, status, result_summary, error_message)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.TeamID, task.RunID, task.ParentTaskID, task.Title, task.Prompt, task.Workdir, string(allowedToolsJSON), task.AssignedAgent, status, task.ResultSummary, blockReason,
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	for _, dep := range dependsOn {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO task_dependencies (task_id, depends_on_task_id) VALUES (?, ?)`,
			task.ID, dep,
		); err != nil {
			return fmt.Errorf("insert task dependency: %w", err)
		}
	}
	return nil
}

func (s *SQLiteTaskStore) recordTaskCreatedEventsTx(ctx context.Context, tx *sql.Tx, task TeamTask, status TaskStatus, blockReason string) error {
	if err := s.insertTaskEventTx(ctx, tx, TaskEvent{
		TeamID:    task.TeamID,
		TaskID:    &task.ID,
		AgentName: MasterAgentName,
		EventType: "task_created",
		Payload:   task.Title,
	}); err != nil {
		return err
	}
	if status != TaskBlocked {
		return nil
	}
	return s.insertTaskEventTx(ctx, tx, TaskEvent{
		TeamID:    task.TeamID,
		TaskID:    &task.ID,
		AgentName: MasterAgentName,
		EventType: "task_blocked",
		Payload:   blockReason,
	})
}

func (s *SQLiteTaskStore) reopenRecoveryDependentsTx(ctx context.Context, tx *sql.Tx, task TeamTask) error {
	if task.ParentTaskID == nil || *task.ParentTaskID == "" {
		return nil
	}
	return s.reopenDependentsForRecoveryTx(ctx, tx, task.TeamID, *task.ParentTaskID, task.ID)
}

func (s *SQLiteTaskStore) selectNextClaimableTaskTx(ctx context.Context, tx *sql.Tx, teamID, agentName string) (*TeamTask, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT t.id, t.team_id, t.run_id, t.parent_task_id, t.title, t.prompt, t.working_dir, t.allowed_tools, t.assigned_agent, t.status, t.result_summary, t.error_message, t.created_at, t.started_at, t.finished_at
		FROM tasks t
		WHERE t.team_id = ?
		  AND t.status = ?
		  AND (t.assigned_agent IS NULL OR t.assigned_agent = ?)
		  AND NOT EXISTS (
			SELECT 1
			FROM task_dependencies d
			JOIN tasks dep ON dep.id = d.depends_on_task_id
			WHERE d.task_id = t.id
			  AND dep.status != ?
		  )
		ORDER BY t.created_at ASC, t.id ASC
		LIMIT 1
	`, teamID, TaskPending, agentName, TaskCompleted)

	task, err := scanTask(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select next runnable task: %w", err)
	}
	return task, nil
}

func (s *SQLiteTaskStore) markTaskClaimedTx(ctx context.Context, tx *sql.Tx, taskID, agentName string, now time.Time) (bool, error) {
	res, err := tx.ExecContext(ctx,
		`UPDATE tasks SET assigned_agent = ?, status = ?, started_at = ? WHERE id = ? AND (assigned_agent IS NULL OR assigned_agent = ?) AND status = ?`,
		agentName, TaskRunning, now, taskID, agentName, TaskPending,
	)
	if err != nil {
		return false, fmt.Errorf("claim task update: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("claim task rows affected: %w", err)
	}
	return rowsAffected == 1, nil
}
