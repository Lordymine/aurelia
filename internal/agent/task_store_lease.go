package agent

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (s *SQLiteTaskStore) requeueRunningTasks(ctx context.Context, teamID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin requeue running tasks tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	rows, err := tx.QueryContext(ctx, `
		SELECT id
		FROM tasks
		WHERE team_id = ? AND status = ? AND assigned_agent IS NOT NULL AND assigned_agent != ''
	`, teamID, TaskRunning)
	if err != nil {
		return fmt.Errorf("query running tasks for requeue: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var taskIDs []string
	for rows.Next() {
		var taskID string
		if err := rows.Scan(&taskID); err != nil {
			return fmt.Errorf("scan running task for requeue: %w", err)
		}
		taskIDs = append(taskIDs, taskID)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("running task rows: %w", err)
	}

	for _, taskID := range taskIDs {
		if _, err := tx.ExecContext(ctx,
			`UPDATE tasks SET status = ?, error_message = ?, started_at = NULL WHERE id = ? AND team_id = ? AND status = ?`,
			TaskPending, "requeued after process restart", taskID, teamID, TaskRunning,
		); err != nil {
			return fmt.Errorf("requeue running task: %w", err)
		}
		if err := s.insertTaskEventTx(ctx, tx, TaskEvent{
			TeamID:    teamID,
			TaskID:    &taskID,
			AgentName: MasterAgentName,
			EventType: "task_requeued",
			Payload:   "requeued after process restart",
		}); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit requeue running tasks tx: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) renewWorkerLeaseTx(ctx context.Context, tx *sql.Tx, teamID, agentName string) error {
	now := time.Now().UTC()
	expiresAt := now.Add(s.leaseDuration)
	res, err := tx.ExecContext(ctx, `
		UPDATE team_members
		SET status = ?, last_heartbeat_at = ?, lease_expires_at = ?
		WHERE team_id = ? AND agent_name = ?
	`, "active", now, expiresAt, teamID, agentName)
	if err != nil {
		return fmt.Errorf("renew worker lease: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("renew worker lease rows affected: %w", err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("worker %s is not registered in team %s", agentName, teamID)
	}
	return nil
}

func (s *SQLiteTaskStore) heartbeatWorker(ctx context.Context, teamID, agentName string) error {
	now := time.Now().UTC()
	expiresAt := now.Add(s.leaseDuration)
	_, err := s.db.ExecContext(ctx, `
		UPDATE team_members
		SET status = ?, last_heartbeat_at = ?, lease_expires_at = ?
		WHERE team_id = ? AND agent_name = ?
	`, "active", now, expiresAt, teamID, agentName)
	if err != nil {
		return fmt.Errorf("heartbeat worker lease: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) markWorkerIdleTx(ctx context.Context, tx *sql.Tx, teamID, agentName string) error {
	now := time.Now().UTC()
	if _, err := tx.ExecContext(ctx, `
		UPDATE team_members
		SET status = ?, last_heartbeat_at = ?, lease_expires_at = NULL
		WHERE team_id = ? AND agent_name = ?
	`, "idle", now, teamID, agentName); err != nil {
		return fmt.Errorf("mark worker idle: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) requeueExpiredRunningTasksTx(ctx context.Context, tx *sql.Tx, teamID string) error {
	now := time.Now().UTC()
	rows, err := tx.QueryContext(ctx, `
		SELECT t.id, t.assigned_agent
		FROM tasks t
		JOIN team_members m ON m.team_id = t.team_id AND m.agent_name = t.assigned_agent
		WHERE t.team_id = ?
		  AND t.status = ?
		  AND t.assigned_agent IS NOT NULL
		  AND m.lease_expires_at IS NOT NULL
		  AND m.lease_expires_at <= ?
	`, teamID, TaskRunning, now)
	if err != nil {
		return fmt.Errorf("query expired running tasks: %w", err)
	}
	defer func() { _ = rows.Close() }()

	type expiredTask struct {
		taskID    string
		agentName string
	}
	var expired []expiredTask
	for rows.Next() {
		var item expiredTask
		if err := rows.Scan(&item.taskID, &item.agentName); err != nil {
			return fmt.Errorf("scan expired running task: %w", err)
		}
		expired = append(expired, item)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("expired running task rows: %w", err)
	}

	for _, item := range expired {
		if _, err := tx.ExecContext(ctx, `
			UPDATE tasks
			SET status = ?, assigned_agent = NULL, error_message = ?, started_at = NULL
			WHERE id = ? AND team_id = ? AND status = ?
		`, TaskPending, fmt.Sprintf("requeued after expired lease from %s", item.agentName), item.taskID, teamID, TaskRunning); err != nil {
			return fmt.Errorf("requeue expired running task: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE team_members
			SET status = ?, lease_expires_at = NULL
			WHERE team_id = ? AND agent_name = ?
		`, "stale", teamID, item.agentName); err != nil {
			return fmt.Errorf("mark stale worker: %w", err)
		}
		if err := s.insertTaskEventTx(ctx, tx, TaskEvent{
			TeamID:    teamID,
			TaskID:    &item.taskID,
			AgentName: MasterAgentName,
			EventType: "task_requeued",
			Payload:   fmt.Sprintf("requeued after expired lease from %s", item.agentName),
		}); err != nil {
			return err
		}
	}

	return nil
}
