package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

func (s *SQLiteTaskStore) getTask(ctx context.Context, teamID, taskID string) (*TeamTask, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, team_id, run_id, parent_task_id, title, prompt, working_dir, allowed_tools, assigned_agent, status, result_summary, error_message, created_at, started_at, finished_at
		FROM tasks
		WHERE team_id = ? AND id = ?
	`, teamID, taskID)

	task, err := scanTask(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get task: %w", err)
	}
	return task, nil
}

func (s *SQLiteTaskStore) listTasks(ctx context.Context, teamID string) ([]TeamTask, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, team_id, run_id, parent_task_id, title, prompt, working_dir, allowed_tools, assigned_agent, status, result_summary, error_message, created_at, started_at, finished_at
		FROM tasks
		WHERE team_id = ?
		ORDER BY created_at ASC, id ASC
	`, teamID)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var tasks []TeamTask
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan list task row: %w", err)
		}
		tasks = append(tasks, *task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list tasks rows: %w", err)
	}
	return tasks, nil
}

func scanTask(scanner interface{ Scan(dest ...any) error }) (*TeamTask, error) {
	var task TeamTask
	var parentID, assignedAgent sql.NullString
	var allowedToolsJSON string
	var startedAt, finishedAt sql.NullTime
	if err := scanner.Scan(
		&task.ID,
		&task.TeamID,
		&task.RunID,
		&parentID,
		&task.Title,
		&task.Prompt,
		&task.Workdir,
		&allowedToolsJSON,
		&assignedAgent,
		&task.Status,
		&task.ResultSummary,
		&task.ErrorMessage,
		&task.CreatedAt,
		&startedAt,
		&finishedAt,
	); err != nil {
		return nil, err
	}
	if parentID.Valid {
		task.ParentTaskID = &parentID.String
	}
	if assignedAgent.Valid {
		task.AssignedAgent = &assignedAgent.String
	}
	if allowedToolsJSON != "" {
		if err := json.Unmarshal([]byte(allowedToolsJSON), &task.AllowedTools); err != nil {
			return nil, fmt.Errorf("unmarshal task allowed tools: %w", err)
		}
	}
	if startedAt.Valid {
		ts := startedAt.Time
		task.StartedAt = &ts
	}
	if finishedAt.Valid {
		ts := finishedAt.Time
		task.FinishedAt = &ts
	}
	return &task, nil
}
