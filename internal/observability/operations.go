package observability

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Operation struct {
	ID         int64
	RunID      string
	TeamID     string
	TaskID     string
	AgentName  string
	Component  string
	Operation  string
	Status     string
	DurationMS int64
	Summary    string
	CreatedAt  time.Time
}

type Recorder interface {
	RecordOperation(ctx context.Context, operation Operation) error
	ListRecentOperations(ctx context.Context, limit int) ([]Operation, error)
	ListFailedOperations(ctx context.Context, limit int) ([]Operation, error)
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open observability database: %w", err)
	}
	if err := initialize(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *SQLiteStore) RecordOperation(ctx context.Context, operation Operation) error {
	if s == nil || s.db == nil {
		return nil
	}

	operation.Component = strings.TrimSpace(operation.Component)
	operation.Operation = strings.TrimSpace(operation.Operation)
	operation.Status = strings.TrimSpace(operation.Status)
	operation.Summary = compact(operation.Summary, 400)
	if operation.Component == "" || operation.Operation == "" {
		return fmt.Errorf("component and operation are required")
	}
	if operation.Status == "" {
		operation.Status = "ok"
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO operational_events (
			run_id, team_id, task_id, agent_name, component, operation, status, duration_ms, summary
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		operation.RunID,
		operation.TeamID,
		operation.TaskID,
		operation.AgentName,
		operation.Component,
		operation.Operation,
		operation.Status,
		operation.DurationMS,
		operation.Summary,
	)
	if err != nil {
		return fmt.Errorf("record operational event: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ListRecentOperations(ctx context.Context, limit int) ([]Operation, error) {
	return s.listOperations(ctx, limit, false)
}

func (s *SQLiteStore) ListFailedOperations(ctx context.Context, limit int) ([]Operation, error) {
	return s.listOperations(ctx, limit, true)
}

func (s *SQLiteStore) listOperations(ctx context.Context, limit int, failuresOnly bool) ([]Operation, error) {
	if s == nil || s.db == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, run_id, team_id, task_id, agent_name, component, operation, status, duration_ms, summary, created_at
		FROM operational_events
	`
	var args []any
	if failuresOnly {
		query += ` WHERE status != 'ok'`
	}
	query += ` ORDER BY created_at DESC, id DESC LIMIT ?`
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list operational events: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var operations []Operation
	for rows.Next() {
		var operation Operation
		if err := rows.Scan(
			&operation.ID,
			&operation.RunID,
			&operation.TeamID,
			&operation.TaskID,
			&operation.AgentName,
			&operation.Component,
			&operation.Operation,
			&operation.Status,
			&operation.DurationMS,
			&operation.Summary,
			&operation.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan operational event: %w", err)
		}
		operations = append(operations, operation)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("operational events rows: %w", err)
	}
	return operations, nil
}

func Observe(ctx context.Context, recorder Recorder, operation Operation) {
	if recorder == nil {
		return
	}
	if err := recorder.RecordOperation(ctx, operation); err != nil {
		log.Printf("component=%q level=%q msg=%q error=%q", "observability", "error", "failed to persist operational event", err.Error())
	}
}

func Log(level, component, message string, fields map[string]string) {
	parts := []string{
		"component=" + quote(component),
		"level=" + quote(level),
		"msg=" + quote(message),
	}

	keys := make([]string, 0, len(fields))
	for key, value := range fields {
		if strings.TrimSpace(value) == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		parts = append(parts, key+"="+quote(fields[key]))
	}
	log.Print(strings.Join(parts, " "))
}

func MergeFields(base map[string]string, extras map[string]string) map[string]string {
	merged := make(map[string]string, len(base)+len(extras))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range extras {
		if strings.TrimSpace(value) == "" {
			continue
		}
		merged[key] = value
	}
	return merged
}

func initialize(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS operational_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			run_id TEXT NOT NULL DEFAULT '',
			team_id TEXT NOT NULL DEFAULT '',
			task_id TEXT NOT NULL DEFAULT '',
			agent_name TEXT NOT NULL DEFAULT '',
			component TEXT NOT NULL,
			operation TEXT NOT NULL,
			status TEXT NOT NULL,
			duration_ms INTEGER NOT NULL DEFAULT 0,
			summary TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_operational_events_created_at
			ON operational_events (created_at, id);
		CREATE INDEX IF NOT EXISTS idx_operational_events_status
			ON operational_events (status, created_at, id);
	`)
	if err != nil {
		return fmt.Errorf("initialize observability schema: %w", err)
	}
	return nil
}

func compact(text string, maxRunes int) string {
	text = strings.TrimSpace(strings.Join(strings.Fields(text), " "))
	if text == "" {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= maxRunes {
		return text
	}
	return string(runes[:maxRunes]) + "..."
}

func quote(value string) string {
	return fmt.Sprintf("%q", strings.TrimSpace(value))
}
