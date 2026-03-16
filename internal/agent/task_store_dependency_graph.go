package agent

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (s *SQLiteTaskStore) resolveRootTaskIDTx(ctx context.Context, tx *sql.Tx, taskID string) (string, error) {
	currentID := taskID
	seen := map[string]bool{}

	for {
		if seen[currentID] {
			return "", fmt.Errorf("parent task cycle detected at %s", currentID)
		}
		seen[currentID] = true

		var parentID sql.NullString
		if err := tx.QueryRowContext(ctx, `SELECT parent_task_id FROM tasks WHERE id = ?`, currentID).Scan(&parentID); err != nil {
			if err == sql.ErrNoRows {
				return "", fmt.Errorf("parent task %s does not exist", currentID)
			}
			return "", fmt.Errorf("resolve root task id: %w", err)
		}
		if !parentID.Valid || strings.TrimSpace(parentID.String) == "" {
			return currentID, nil
		}
		currentID = parentID.String
	}
}

func (s *SQLiteTaskStore) hasDependencyPathTx(ctx context.Context, tx *sql.Tx, fromTaskID, targetTaskID string) (bool, error) {
	if fromTaskID == targetTaskID {
		return true, nil
	}

	queue := []string{fromTaskID}
	visited := map[string]bool{}

	for len(queue) > 0 {
		currentID := queue[0]
		queue = queue[1:]
		if visited[currentID] {
			continue
		}
		visited[currentID] = true

		rows, err := tx.QueryContext(ctx, `SELECT depends_on_task_id FROM task_dependencies WHERE task_id = ?`, currentID)
		if err != nil {
			return false, fmt.Errorf("query dependency path: %w", err)
		}

		var nextIDs []string
		for rows.Next() {
			var nextID string
			if err := rows.Scan(&nextID); err != nil {
				_ = rows.Close()
				return false, fmt.Errorf("scan dependency path: %w", err)
			}
			if nextID == targetTaskID {
				_ = rows.Close()
				return true, nil
			}
			nextIDs = append(nextIDs, nextID)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return false, fmt.Errorf("dependency path rows: %w", err)
		}
		_ = rows.Close()

		queue = append(queue, nextIDs...)
	}

	return false, nil
}

func (s *SQLiteTaskStore) listChildTasksTx(ctx context.Context, tx *sql.Tx, parentTaskID string) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id FROM tasks WHERE parent_task_id = ? ORDER BY created_at ASC, id ASC`, parentTaskID)
	if err != nil {
		return nil, fmt.Errorf("query child tasks: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var childIDs []string
	for rows.Next() {
		var childID string
		if err := rows.Scan(&childID); err != nil {
			return nil, fmt.Errorf("scan child task: %w", err)
		}
		childIDs = append(childIDs, childID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("child task rows: %w", err)
	}
	return childIDs, nil
}
