package cron

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (s *SQLiteCronStore) CreateJob(ctx context.Context, job CronJob) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO cron_jobs (
			id, owner_user_id, target_chat_id, schedule_type, cron_expr, run_at, prompt, active,
			last_run_at, next_run_at, last_status, last_error
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		job.ID,
		job.OwnerUserID,
		job.TargetChatID,
		job.ScheduleType,
		job.CronExpr,
		job.RunAt,
		job.Prompt,
		boolToInt(job.Active),
		job.LastRunAt,
		job.NextRunAt,
		job.LastStatus,
		job.LastError,
	)
	if err != nil {
		return fmt.Errorf("insert cron job: %w", err)
	}
	return nil
}

// UpdateJobTx updates a cron job within an existing transaction.
func (s *SQLiteCronStore) UpdateJobTx(ctx context.Context, tx *sql.Tx, job CronJob) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE cron_jobs
		SET owner_user_id = ?, target_chat_id = ?, schedule_type = ?, cron_expr = ?, run_at = ?, prompt = ?, active = ?,
			last_run_at = ?, next_run_at = ?, last_status = ?, last_error = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		job.OwnerUserID,
		job.TargetChatID,
		job.ScheduleType,
		job.CronExpr,
		job.RunAt,
		job.Prompt,
		boolToInt(job.Active),
		job.LastRunAt,
		job.NextRunAt,
		job.LastStatus,
		job.LastError,
		job.ID,
	)
	if err != nil {
		return fmt.Errorf("update cron job: %w", err)
	}
	return nil
}

func (s *SQLiteCronStore) UpdateJob(ctx context.Context, job CronJob) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE cron_jobs
		SET owner_user_id = ?, target_chat_id = ?, schedule_type = ?, cron_expr = ?, run_at = ?, prompt = ?, active = ?,
			last_run_at = ?, next_run_at = ?, last_status = ?, last_error = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		job.OwnerUserID,
		job.TargetChatID,
		job.ScheduleType,
		job.CronExpr,
		job.RunAt,
		job.Prompt,
		boolToInt(job.Active),
		job.LastRunAt,
		job.NextRunAt,
		job.LastStatus,
		job.LastError,
		job.ID,
	)
	if err != nil {
		return fmt.Errorf("update cron job: %w", err)
	}
	return nil
}

func (s *SQLiteCronStore) DeleteJob(ctx context.Context, jobID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM cron_jobs WHERE id = ?`, jobID)
	if err != nil {
		return fmt.Errorf("delete cron job: %w", err)
	}
	return nil
}

// ResolveJobID resolves a short (prefix) ID to the full UUID.
// Returns the full ID if exactly one match is found.
func (s *SQLiteCronStore) ResolveJobID(ctx context.Context, prefix string) (string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id FROM cron_jobs WHERE id LIKE ?`, prefix+"%")
	if err != nil {
		return "", fmt.Errorf("resolve job id: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var matches []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return "", fmt.Errorf("scan job id: %w", err)
		}
		matches = append(matches, id)
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("resolve job id rows: %w", err)
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("cron job %s not found", prefix)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("ambiguous job id prefix %q matches %d jobs", prefix, len(matches))
	}
}

func (s *SQLiteCronStore) GetJob(ctx context.Context, jobID string) (*CronJob, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, owner_user_id, target_chat_id, schedule_type, cron_expr, run_at, prompt, active,
		       last_run_at, next_run_at, last_status, last_error, created_at, updated_at
		FROM cron_jobs WHERE id = ?
	`, jobID)
	job, err := scanCronJob(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get cron job: %w", err)
	}
	return job, nil
}

func (s *SQLiteCronStore) ListJobsByChat(ctx context.Context, chatID int64) ([]CronJob, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, owner_user_id, target_chat_id, schedule_type, cron_expr, run_at, prompt, active,
		       last_run_at, next_run_at, last_status, last_error, created_at, updated_at
		FROM cron_jobs
		WHERE target_chat_id = ?
		ORDER BY created_at ASC, id ASC
	`, chatID)
	if err != nil {
		return nil, fmt.Errorf("list cron jobs by chat: %w", err)
	}
	defer func() { _ = rows.Close() }()

	return scanCronJobs(rows)
}

func (s *SQLiteCronStore) ListDueJobs(ctx context.Context, now time.Time, limit int) ([]CronJob, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, owner_user_id, target_chat_id, schedule_type, cron_expr, run_at, prompt, active,
		       last_run_at, next_run_at, last_status, last_error, created_at, updated_at
		FROM cron_jobs
		WHERE active = 1 AND next_run_at IS NOT NULL AND next_run_at <= ?
		ORDER BY next_run_at ASC, id ASC
		LIMIT ?
	`, now.UTC(), limit)
	if err != nil {
		return nil, fmt.Errorf("list due cron jobs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	return scanCronJobs(rows)
}
