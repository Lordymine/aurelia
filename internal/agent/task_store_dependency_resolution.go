package agent

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *SQLiteTaskStore) resolveDependencyStatusTx(ctx context.Context, tx *sql.Tx, taskID string) (TaskStatus, error) {
	status, err := s.getTaskStatusTx(ctx, tx, taskID)
	if err != nil {
		return "", err
	}
	if status == TaskCompleted {
		return TaskCompleted, nil
	}

	children, err := s.listChildTasksTx(ctx, tx, taskID)
	if err != nil {
		return "", err
	}

	return s.resolveRecoveryStatusTx(ctx, tx, status, children)
}

func (s *SQLiteTaskStore) resolveRecoveryStatusTx(
	ctx context.Context,
	tx *sql.Tx,
	parentStatus TaskStatus,
	childIDs []string,
) (TaskStatus, error) {
	hasCompletedRecovery := false
	hasActiveRecovery := false

	for _, childID := range childIDs {
		childStatus, err := s.resolveDependencyStatusTx(ctx, tx, childID)
		if err != nil {
			return "", err
		}
		switch childStatus {
		case TaskCompleted:
			hasCompletedRecovery = true
		case TaskPending, TaskRunning, TaskBlocked:
			hasActiveRecovery = true
		}
	}

	switch {
	case hasCompletedRecovery:
		return TaskCompleted, nil
	case hasActiveRecovery:
		return TaskBlocked, nil
	default:
		return parentStatus, nil
	}
}

func (s *SQLiteTaskStore) getParentTaskIDTx(ctx context.Context, tx *sql.Tx, taskID string) (string, bool, error) {
	var parentID sql.NullString
	if err := tx.QueryRowContext(ctx, `SELECT parent_task_id FROM tasks WHERE id = ?`, taskID).Scan(&parentID); err != nil {
		return "", false, fmt.Errorf("read parent task id: %w", err)
	}
	if !parentID.Valid || parentID.String == "" {
		return "", false, nil
	}
	return parentID.String, true, nil
}
