package cron

import "fmt"

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
	return nil
}
