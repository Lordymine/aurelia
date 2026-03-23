package cron

import (
	"fmt"
	"strings"
)

func (s *SQLiteCronStore) initialize() error {
	query := `
	CREATE TABLE IF NOT EXISTS cron_jobs (
		id TEXT PRIMARY KEY,
		owner_user_id TEXT NOT NULL,
		target_chat_id INTEGER NOT NULL,
		schedule_type TEXT NOT NULL,
		cron_expr TEXT NOT NULL DEFAULT '',
		run_at DATETIME,
		prompt TEXT NOT NULL,
		active INTEGER NOT NULL DEFAULT 1,
		last_run_at DATETIME,
		next_run_at DATETIME,
		last_status TEXT NOT NULL DEFAULT 'idle',
		last_error TEXT NOT NULL DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS cron_executions (
		id TEXT PRIMARY KEY,
		job_id TEXT NOT NULL,
		started_at DATETIME NOT NULL,
		finished_at DATETIME,
		status TEXT NOT NULL,
		output_summary TEXT NOT NULL DEFAULT '',
		error_message TEXT NOT NULL DEFAULT ''
	);
	`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("initialize cron schema: %w", err)
	}

	// Migration: add columns if missing (safe to re-run)
	for _, col := range []string{
		"ALTER TABLE cron_executions ADD COLUMN session_id TEXT DEFAULT ''",
		"ALTER TABLE cron_executions ADD COLUMN cost_usd REAL DEFAULT 0",
		"ALTER TABLE cron_executions ADD COLUMN tokens_used INTEGER DEFAULT 0",
	} {
		_, err := s.db.Exec(col)
		if err != nil && !strings.Contains(err.Error(), "duplicate column") {
			return fmt.Errorf("migration: %w", err)
		}
	}

	// Index for ListDueJobs query
	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_cron_jobs_due ON cron_jobs(active, next_run_at)`)
	if err != nil {
		return fmt.Errorf("create due jobs index: %w", err)
	}

	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_cron_jobs_chat ON cron_jobs(target_chat_id)`)
	if err != nil {
		return fmt.Errorf("create chat jobs index: %w", err)
	}

	return nil
}
