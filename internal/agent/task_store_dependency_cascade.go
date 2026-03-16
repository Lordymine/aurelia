package agent

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (s *SQLiteTaskStore) cancelDependentsTx(ctx context.Context, tx *sql.Tx, teamID, taskID, reason string) error {
	dependentIDs, err := s.listCancellableDependentsTx(ctx, tx, teamID, taskID)
	if err != nil {
		return err
	}

	for _, dependentID := range dependentIDs {
		if err := s.cancelDependentCascadeTx(ctx, tx, teamID, dependentID, reason); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteTaskStore) listCancellableDependentsTx(ctx context.Context, tx *sql.Tx, teamID, dependencyTaskID string) ([]string, error) {
	return s.listDependentsByStatusTx(ctx, tx, teamID, dependencyTaskID, TaskPending, TaskBlocked)
}

func (s *SQLiteTaskStore) cancelDependentCascadeTx(ctx context.Context, tx *sql.Tx, teamID, dependentID, reason string) error {
	if err := s.markDependentCancelledTx(ctx, tx, teamID, dependentID, reason); err != nil {
		return err
	}
	if err := s.insertTaskEventTx(ctx, tx, TaskEvent{
		TeamID:    teamID,
		TaskID:    &dependentID,
		AgentName: MasterAgentName,
		EventType: "task_cancelled",
		Payload:   reason,
	}); err != nil {
		return err
	}

	nextReason := fmt.Sprintf("cancelled because dependency %s was cancelled", dependentID)
	return s.cancelDependentsTx(ctx, tx, teamID, dependentID, nextReason)
}

func (s *SQLiteTaskStore) markDependentCancelledTx(ctx context.Context, tx *sql.Tx, teamID, dependentID, reason string) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE tasks SET status = ?, error_message = ?, finished_at = ? WHERE id = ? AND team_id = ? AND status IN (?, ?)`,
		TaskCancelled, reason, time.Now().UTC(), dependentID, teamID, TaskPending, TaskBlocked,
	)
	if err != nil {
		return fmt.Errorf("cancel dependent task: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) reopenDependentsForRecoveryTx(ctx context.Context, tx *sql.Tx, teamID, dependencyTaskID, recoveryTaskID string) error {
	parentStatus, err := s.getTaskStatusTx(ctx, tx, dependencyTaskID)
	if err != nil {
		return err
	}
	if parentStatus != TaskFailed && parentStatus != TaskCancelled {
		return nil
	}

	dependentIDs, err := s.listCancelledDependentsTx(ctx, tx, teamID, dependencyTaskID)
	if err != nil {
		return err
	}

	for _, dependentID := range dependentIDs {
		if err := s.reopenDependentForRecoveryTx(ctx, tx, teamID, dependencyTaskID, recoveryTaskID, dependentID); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteTaskStore) listCancelledDependentsTx(ctx context.Context, tx *sql.Tx, teamID, dependencyTaskID string) ([]string, error) {
	return s.listDependentsByStatusTx(ctx, tx, teamID, dependencyTaskID, TaskCancelled)
}

func (s *SQLiteTaskStore) reopenDependentForRecoveryTx(
	ctx context.Context,
	tx *sql.Tx,
	teamID, dependencyTaskID, recoveryTaskID, dependentID string,
) error {
	reason := fmt.Sprintf("blocked while dependency %s is in recovery via %s", dependencyTaskID, recoveryTaskID)
	if err := s.markDependentBlockedForRecoveryTx(ctx, tx, teamID, dependentID, reason); err != nil {
		return err
	}
	if err := s.insertRecoveryEventsTx(ctx, tx, teamID, dependentID, reason); err != nil {
		return err
	}
	return s.reopenDependentsForRecoveryTx(ctx, tx, teamID, dependentID, recoveryTaskID)
}

func (s *SQLiteTaskStore) markDependentBlockedForRecoveryTx(
	ctx context.Context,
	tx *sql.Tx,
	teamID, dependentID, reason string,
) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE tasks SET status = ?, error_message = ?, finished_at = NULL WHERE id = ? AND team_id = ? AND status = ?`,
		TaskBlocked, reason, dependentID, teamID, TaskCancelled,
	)
	if err != nil {
		return fmt.Errorf("reopen dependent task: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) insertRecoveryEventsTx(ctx context.Context, tx *sql.Tx, teamID, dependentID, reason string) error {
	if err := s.insertTaskEventTx(ctx, tx, TaskEvent{
		TeamID:    teamID,
		TaskID:    &dependentID,
		AgentName: MasterAgentName,
		EventType: "task_reopened",
		Payload:   reason,
	}); err != nil {
		return err
	}

	return s.insertTaskEventTx(ctx, tx, TaskEvent{
		TeamID:    teamID,
		TaskID:    &dependentID,
		AgentName: MasterAgentName,
		EventType: "task_blocked",
		Payload:   reason,
	})
}
