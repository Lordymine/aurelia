package agent

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (s *SQLiteTaskStore) createTeam(ctx context.Context, id, teamKey, userID, leadAgent string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO teams (id, team_key, user_id, lead_agent, status) VALUES (?, ?, ?, ?, ?)`,
		id, teamKey, userID, leadAgent, TeamStatusActive,
	)
	if err != nil {
		return fmt.Errorf("insert team: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) getTeamIDByKey(ctx context.Context, teamKey string) (string, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id FROM teams WHERE team_key = ?`, teamKey)
	var teamID string
	if err := row.Scan(&teamID); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("get team id by key: %w", err)
	}
	return teamID, nil
}

func (s *SQLiteTaskStore) getTeamStatus(ctx context.Context, teamID string) (string, error) {
	row := s.db.QueryRowContext(ctx, `SELECT status FROM teams WHERE id = ?`, teamID)
	var status string
	if err := row.Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("get team status: %w", err)
	}
	return status, nil
}

func (s *SQLiteTaskStore) setTeamStatus(ctx context.Context, teamID, status string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE teams SET status = ? WHERE id = ?`, status, teamID)
	if err != nil {
		return fmt.Errorf("set team status: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) listTeams(ctx context.Context) ([]struct {
	TeamID  string
	TeamKey string
	UserID  string
}, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, team_key, user_id FROM teams ORDER BY created_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var teams []struct {
		TeamID  string
		TeamKey string
		UserID  string
	}
	for rows.Next() {
		var item struct {
			TeamID  string
			TeamKey string
			UserID  string
		}
		if err := rows.Scan(&item.TeamID, &item.TeamKey, &item.UserID); err != nil {
			return nil, fmt.Errorf("scan team row: %w", err)
		}
		teams = append(teams, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list teams rows: %w", err)
	}
	return teams, nil
}

func (s *SQLiteTaskStore) cancelActiveTasks(ctx context.Context, teamID, reason string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin cancel active tasks tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	rows, err := tx.QueryContext(ctx, `
		SELECT id
		FROM tasks
		WHERE team_id = ?
		  AND status IN (?, ?, ?)
	`, teamID, TaskPending, TaskBlocked, TaskRunning)
	if err != nil {
		return fmt.Errorf("query cancellable tasks: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var taskIDs []string
	for rows.Next() {
		var taskID string
		if err := rows.Scan(&taskID); err != nil {
			return fmt.Errorf("scan cancellable task: %w", err)
		}
		taskIDs = append(taskIDs, taskID)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("cancellable task rows: %w", err)
	}

	now := sql.NullTime{Time: time.Now().UTC(), Valid: true}
	for _, taskID := range taskIDs {
		if _, err := tx.ExecContext(ctx, `
			UPDATE tasks
			SET status = ?, error_message = ?, finished_at = ?, assigned_agent = CASE WHEN status = ? THEN NULL ELSE assigned_agent END
			WHERE id = ? AND team_id = ? AND status IN (?, ?, ?)
		`, TaskCancelled, reason, now, TaskRunning, taskID, teamID, TaskPending, TaskBlocked, TaskRunning); err != nil {
			return fmt.Errorf("cancel active task: %w", err)
		}
		if err := s.insertTaskEventTx(ctx, tx, TaskEvent{
			TeamID:    teamID,
			TaskID:    &taskID,
			AgentName: MasterAgentName,
			EventType: "task_cancelled",
			Payload:   reason,
		}); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit cancel active tasks tx: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) registerTeammate(ctx context.Context, teamID, memberID, agentName, roleDescription string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO team_members (id, team_id, agent_name, role_description, status, last_heartbeat_at, lease_expires_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		memberID, teamID, agentName, roleDescription, "idle", nil, nil,
	)
	if err != nil {
		return fmt.Errorf("insert teammate: %w", err)
	}
	return nil
}
