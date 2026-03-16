package agent

import (
	"fmt"
	"strings"
)

func (s *SQLiteTaskStore) initialize() error {
	query := `
	CREATE TABLE IF NOT EXISTS teams (
		id TEXT PRIMARY KEY,
		team_key TEXT NOT NULL UNIQUE,
		user_id TEXT NOT NULL,
		lead_agent TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS team_members (
		id TEXT PRIMARY KEY,
		team_id TEXT NOT NULL,
		agent_name TEXT NOT NULL,
		role_description TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		team_id TEXT NOT NULL,
		run_id TEXT NOT NULL DEFAULT '',
		parent_task_id TEXT,
		title TEXT NOT NULL,
		prompt TEXT NOT NULL,
		working_dir TEXT NOT NULL DEFAULT '',
		allowed_tools TEXT NOT NULL DEFAULT '[]',
		assigned_agent TEXT,
		status TEXT NOT NULL,
		result_summary TEXT NOT NULL DEFAULT '',
		error_message TEXT NOT NULL DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		started_at DATETIME,
		finished_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS task_dependencies (
		task_id TEXT NOT NULL,
		depends_on_task_id TEXT NOT NULL,
		PRIMARY KEY (task_id, depends_on_task_id)
	);

	CREATE TABLE IF NOT EXISTS mail_messages (
		id TEXT PRIMARY KEY,
		team_id TEXT NOT NULL,
		from_agent TEXT NOT NULL,
		to_agent TEXT NOT NULL,
		task_id TEXT,
		kind TEXT NOT NULL,
		body TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		consumed_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS task_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		team_id TEXT NOT NULL,
		task_id TEXT,
		agent_name TEXT NOT NULL,
		event_type TEXT NOT NULL,
		payload TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("initialize task store schema: %w", err)
	}
	if _, err := s.db.Exec(`ALTER TABLE tasks ADD COLUMN run_id TEXT NOT NULL DEFAULT ''`); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
		return fmt.Errorf("migrate tasks.run_id: %w", err)
	}
	if _, err := s.db.Exec(`ALTER TABLE tasks ADD COLUMN working_dir TEXT NOT NULL DEFAULT ''`); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
		return fmt.Errorf("migrate tasks.working_dir: %w", err)
	}
	if _, err := s.db.Exec(`ALTER TABLE tasks ADD COLUMN allowed_tools TEXT NOT NULL DEFAULT '[]'`); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
		return fmt.Errorf("migrate tasks.allowed_tools: %w", err)
	}
	if _, err := s.db.Exec(`ALTER TABLE team_members ADD COLUMN last_heartbeat_at DATETIME`); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
		return fmt.Errorf("migrate team_members.last_heartbeat_at: %w", err)
	}
	if _, err := s.db.Exec(`ALTER TABLE team_members ADD COLUMN lease_expires_at DATETIME`); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
		return fmt.Errorf("migrate team_members.lease_expires_at: %w", err)
	}
	return nil
}
