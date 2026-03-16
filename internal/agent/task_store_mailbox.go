package agent

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (s *SQLiteTaskStore) postMessage(ctx context.Context, msg MailMessage) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO mail_messages (id, team_id, from_agent, to_agent, task_id, kind, body, consumed_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.TeamID, msg.FromAgent, msg.ToAgent, msg.TaskID, msg.Kind, msg.Body, nil,
	)
	if err != nil {
		return fmt.Errorf("insert mail message: %w", err)
	}
	return nil
}

func (s *SQLiteTaskStore) pullMessages(ctx context.Context, teamID, agentName string, limit int) ([]MailMessage, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin pull messages tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	rows, err := tx.QueryContext(ctx, `
		SELECT id, team_id, from_agent, to_agent, task_id, kind, body, created_at
		FROM mail_messages
		WHERE team_id = ? AND to_agent = ? AND consumed_at IS NULL
		ORDER BY rowid ASC
		LIMIT ?
	`, teamID, agentName, limit)
	if err != nil {
		return nil, fmt.Errorf("query mailbox: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var messages []MailMessage
	var ids []string
	for rows.Next() {
		var msg MailMessage
		var taskID sql.NullString
		if err := rows.Scan(&msg.ID, &msg.TeamID, &msg.FromAgent, &msg.ToAgent, &taskID, &msg.Kind, &msg.Body, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan mailbox row: %w", err)
		}
		if taskID.Valid {
			msg.TaskID = &taskID.String
		}
		messages = append(messages, msg)
		ids = append(ids, msg.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("mailbox rows: %w", err)
	}

	if len(ids) == 0 {
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit empty mailbox tx: %w", err)
		}
		return nil, nil
	}

	consumedAt := time.Now().UTC()
	for i := range ids {
		if _, err := tx.ExecContext(ctx, `UPDATE mail_messages SET consumed_at = ? WHERE id = ?`, consumedAt, ids[i]); err != nil {
			return nil, fmt.Errorf("mark mailbox message consumed: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit pull messages tx: %w", err)
	}

	for i := range messages {
		messages[i].ConsumedAt = &consumedAt
	}
	return messages, nil
}

func (s *SQLiteTaskStore) listEvents(ctx context.Context, teamID string, limit int) ([]TaskEvent, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, team_id, task_id, agent_name, event_type, payload, created_at
		FROM task_events
		WHERE team_id = ?
		ORDER BY id ASC
		LIMIT ?
	`, teamID, limit)
	if err != nil {
		return nil, fmt.Errorf("list task events: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var events []TaskEvent
	for rows.Next() {
		var event TaskEvent
		var taskID sql.NullString
		if err := rows.Scan(&event.ID, &event.TeamID, &taskID, &event.AgentName, &event.EventType, &event.Payload, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan task event row: %w", err)
		}
		if taskID.Valid {
			event.TaskID = &taskID.String
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("task event rows: %w", err)
	}
	return events, nil
}

func (s *SQLiteTaskStore) insertTaskEventTx(ctx context.Context, tx *sql.Tx, event TaskEvent) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO task_events (team_id, task_id, agent_name, event_type, payload)
		 VALUES (?, ?, ?, ?, ?)`,
		event.TeamID, event.TaskID, event.AgentName, event.EventType, event.Payload,
	)
	if err != nil {
		return fmt.Errorf("insert task event: %w", err)
	}
	return nil
}
