package agent

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (s *SQLiteTaskStore) evaluateDependencyStateTx(ctx context.Context, tx *sql.Tx, dependencyIDs []string) (bool, string, error) {
	if len(dependencyIDs) == 0 {
		return false, "", nil
	}

	for _, depID := range dependencyIDs {
		status, err := s.resolveDependencyStatusTx(ctx, tx, depID)
		if err != nil {
			if err == sql.ErrNoRows {
				return true, fmt.Sprintf("blocked by missing dependency %s", depID), nil
			}
			return false, "", err
		}
		switch status {
		case TaskCompleted:
			continue
		case TaskFailed, TaskCancelled:
			return true, fmt.Sprintf("blocked by failed dependency %s", depID), nil
		default:
			return true, fmt.Sprintf("blocked by dependency %s", depID), nil
		}
	}

	return false, "", nil
}

func (s *SQLiteTaskStore) validateTaskDependenciesTx(ctx context.Context, tx *sql.Tx, task TeamTask, dependencyIDs []string) error {
	protectedIDs := map[string]struct{}{}
	if strings.TrimSpace(task.ID) != "" {
		protectedIDs[task.ID] = struct{}{}
	}
	if task.ParentTaskID != nil && strings.TrimSpace(*task.ParentTaskID) != "" {
		rootID, err := s.resolveRootTaskIDTx(ctx, tx, *task.ParentTaskID)
		if err != nil {
			return err
		}
		protectedIDs[rootID] = struct{}{}
	}

	for _, depID := range dependencyIDs {
		depID = strings.TrimSpace(depID)
		if depID == "" {
			return fmt.Errorf("dependency id cannot be empty")
		}
		if _, ok := protectedIDs[depID]; ok {
			return fmt.Errorf("dependency cycle detected for task %s", task.ID)
		}

		exists, err := s.taskExistsInTeamTx(ctx, tx, task.TeamID, depID)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("dependency %s does not exist in team %s", depID, task.TeamID)
		}

		for protectedID := range protectedIDs {
			reaches, err := s.hasDependencyPathTx(ctx, tx, depID, protectedID)
			if err != nil {
				return err
			}
			if reaches {
				return fmt.Errorf("dependency cycle detected between %s and %s", depID, protectedID)
			}
		}
	}

	return nil
}

func (s *SQLiteTaskStore) areAllDependenciesCompletedTx(ctx context.Context, tx *sql.Tx, taskID string) (bool, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT d.depends_on_task_id
		FROM task_dependencies d
		WHERE d.task_id = ?
	`, taskID)
	if err != nil {
		return false, fmt.Errorf("query dependency completion: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var depID string
		if err := rows.Scan(&depID); err != nil {
			return false, fmt.Errorf("scan dependency completion: %w", err)
		}
		status, err := s.resolveDependencyStatusTx(ctx, tx, depID)
		if err != nil {
			return false, err
		}
		if status != TaskCompleted {
			return false, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("dependency completion rows: %w", err)
	}
	return true, nil
}

func (s *SQLiteTaskStore) getTaskStatusTx(ctx context.Context, tx *sql.Tx, taskID string) (TaskStatus, error) {
	var status TaskStatus
	if err := tx.QueryRowContext(ctx, `SELECT status FROM tasks WHERE id = ?`, taskID).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", fmt.Errorf("read task status: %w", err)
	}
	return status, nil
}

func (s *SQLiteTaskStore) taskExistsInTeamTx(ctx context.Context, tx *sql.Tx, teamID, taskID string) (bool, error) {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM tasks WHERE team_id = ? AND id = ?`, teamID, taskID).Scan(&count); err != nil {
		return false, fmt.Errorf("check task existence: %w", err)
	}
	return count > 0, nil
}
