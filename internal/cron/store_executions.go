package cron

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *SQLiteCronStore) RecordExecution(ctx context.Context, exec CronExecution) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO cron_executions (id, job_id, started_at, finished_at, status, output_summary, error_message)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, exec.ID, exec.JobID, exec.StartedAt, exec.FinishedAt, exec.Status, exec.OutputSummary, exec.ErrorMessage)
	if err != nil {
		return fmt.Errorf("insert cron execution: %w", err)
	}
	return nil
}

func (s *SQLiteCronStore) ListExecutionsByJob(ctx context.Context, jobID string) ([]CronExecution, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, job_id, started_at, finished_at, status, output_summary, error_message
		FROM cron_executions
		WHERE job_id = ?
		ORDER BY started_at ASC, id ASC
	`, jobID)
	if err != nil {
		return nil, fmt.Errorf("list cron executions: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var executions []CronExecution
	for rows.Next() {
		var exec CronExecution
		var finishedAt sql.NullTime
		if err := rows.Scan(&exec.ID, &exec.JobID, &exec.StartedAt, &finishedAt, &exec.Status, &exec.OutputSummary, &exec.ErrorMessage); err != nil {
			return nil, fmt.Errorf("scan cron execution row: %w", err)
		}
		if finishedAt.Valid {
			ts := finishedAt.Time
			exec.FinishedAt = &ts
		}
		executions = append(executions, exec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cron execution rows: %w", err)
	}
	return executions, nil
}
