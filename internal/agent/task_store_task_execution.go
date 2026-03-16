package agent

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (s *SQLiteTaskStore) completeTask(ctx context.Context, teamID, taskID, agentName, result string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin complete task tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.finishTaskTx(ctx, tx, teamID, taskID, agentName, TaskCompleted, result); err != nil {
		return err
	}
	if err := s.unblockDependentsTx(ctx, tx, teamID, taskID); err != nil {
		return err
	}
	if err := s.markWorkerIdleTx(ctx, tx, teamID, agentName); err != nil {
		return err
	}
	if parentTaskID, ok, err := s.getParentTaskIDTx(ctx, tx, taskID); err != nil {
		return err
	} else if ok {
		if err := s.unblockDependentsTx(ctx, tx, teamID, parentTaskID); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit complete task tx: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) failTask(ctx context.Context, teamID, taskID, agentName, reason string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin fail task tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.finishTaskTx(ctx, tx, teamID, taskID, agentName, TaskFailed, reason); err != nil {
		return err
	}
	if err := s.cancelDependentsTx(ctx, tx, teamID, taskID, fmt.Sprintf("cancelled because dependency %s failed: %s", taskID, reason)); err != nil {
		return err
	}
	if err := s.markWorkerIdleTx(ctx, tx, teamID, agentName); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail task tx: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) finishTaskTx(ctx context.Context, tx *sql.Tx, teamID, taskID, agentName string, status TaskStatus, payload string) error {
	now := time.Now().UTC()
	column := "result_summary"
	eventType := "task_completed"
	if status == TaskFailed {
		column = "error_message"
		eventType = "task_failed"
	}

	res, err := tx.ExecContext(ctx,
		fmt.Sprintf(`UPDATE tasks SET status = ?, %s = ?, finished_at = ? WHERE id = ? AND team_id = ? AND assigned_agent = ? AND status = ?`, column),
		status, payload, now, taskID, teamID, agentName, TaskRunning,
	)
	if err != nil {
		return fmt.Errorf("finish task update: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("finish task rows affected: %w", err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("task %s was not running for agent %s", taskID, agentName)
	}
	return s.insertTaskEventTx(ctx, tx, TaskEvent{
		TeamID:    teamID,
		TaskID:    &taskID,
		AgentName: agentName,
		EventType: eventType,
		Payload:   payload,
	})
}
