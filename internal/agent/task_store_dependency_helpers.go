package agent

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (s *SQLiteTaskStore) listDependentsByStatusTx(
	ctx context.Context,
	tx *sql.Tx,
	teamID, dependencyTaskID string,
	statuses ...TaskStatus,
) ([]string, error) {
	if len(statuses) == 0 {
		return nil, nil
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(statuses)), ",")
	args := make([]any, 0, len(statuses)+2)
	args = append(args, teamID, dependencyTaskID)
	for _, status := range statuses {
		args = append(args, status)
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT t.id
		FROM tasks t
		JOIN task_dependencies d ON d.task_id = t.id
		WHERE t.team_id = ?
		  AND d.depends_on_task_id = ?
		  AND t.status IN (%s)
	`, placeholders)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query dependents by status: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var dependentIDs []string
	for rows.Next() {
		var dependentID string
		if err := rows.Scan(&dependentID); err != nil {
			return nil, fmt.Errorf("scan dependent task id: %w", err)
		}
		dependentIDs = append(dependentIDs, dependentID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dependent task rows: %w", err)
	}

	return dependentIDs, nil
}

func (s *SQLiteTaskStore) updateTaskStatusTx(
	ctx context.Context,
	tx *sql.Tx,
	teamID, taskID string,
	fromStatus, toStatus TaskStatus,
	errorMessage string,
) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE tasks SET status = ?, error_message = ? WHERE id = ? AND team_id = ? AND status = ?`,
		toStatus, errorMessage, taskID, teamID, fromStatus,
	)
	if err != nil {
		return fmt.Errorf("update task status: %w", err)
	}
	return nil
}
