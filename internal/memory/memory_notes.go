package memory

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// AddNote stores a compact long-term memory note.
func (m *MemoryManager) AddNote(ctx context.Context, note Note) error {
	note.ConversationID = strings.TrimSpace(note.ConversationID)
	note.Topic = strings.TrimSpace(note.Topic)
	note.Summary = strings.TrimSpace(strings.ReplaceAll(note.Summary, "\x00", ""))
	note.Kind = strings.TrimSpace(note.Kind)
	note.Source = strings.TrimSpace(note.Source)

	if note.ConversationID == "" || note.Topic == "" || note.Summary == "" {
		return fmt.Errorf("conversation_id, topic and summary are required")
	}
	if note.Kind == "" {
		note.Kind = "note"
	}
	if note.Source == "" {
		note.Source = "unknown"
	}

	existing, ok, err := m.findLatestNoteByTopicKind(ctx, note.ConversationID, note.Topic, note.Kind)
	if err != nil {
		return err
	}
	if ok {
		return m.consolidateExistingNote(ctx, existing, note)
	}

	query := `
		INSERT OR IGNORE INTO memory_notes (conversation_id, topic, summary, kind, importance, source)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = m.db.ExecContext(ctx, query, note.ConversationID, note.Topic, note.Summary, note.Kind, note.Importance, note.Source)
	if err != nil {
		return fmt.Errorf("failed to insert note: %w", err)
	}
	return nil
}

func (m *MemoryManager) consolidateExistingNote(ctx context.Context, existing, note Note) error {
	if strings.TrimSpace(existing.Summary) == note.Summary {
		return nil
	}

	mergedSummary := mergeNoteSummaries(existing.Summary, note.Summary)
	importance := existing.Importance
	if note.Importance > importance {
		importance = note.Importance
	}

	updateQuery := `
		UPDATE memory_notes
		SET summary = ?, importance = ?, source = ?, created_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.db.ExecContext(ctx, updateQuery, mergedSummary, importance, note.Source, existing.ID)
	if err != nil {
		return fmt.Errorf("failed to consolidate note: %w", err)
	}
	return nil
}

func (m *MemoryManager) findLatestNoteByTopicKind(ctx context.Context, conversationID, topic, kind string) (Note, bool, error) {
	query := `
		SELECT id, conversation_id, topic, summary, kind, importance, source, created_at
		FROM memory_notes
		WHERE conversation_id = ? AND topic = ? AND kind = ?
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`

	var note Note
	err := m.db.QueryRowContext(ctx, query, conversationID, topic, kind).
		Scan(&note.ID, &note.ConversationID, &note.Topic, &note.Summary, &note.Kind, &note.Importance, &note.Source, &note.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Note{}, false, nil
		}
		return Note{}, false, fmt.Errorf("failed to query latest note: %w", err)
	}
	return note, true, nil
}

// ListRecentNotes retrieves the newest notes for a conversation.
func (m *MemoryManager) ListRecentNotes(ctx context.Context, conversationID string, limit int) ([]Note, error) {
	if limit <= 0 {
		limit = 5
	}

	query := `
		SELECT id, conversation_id, topic, summary, kind, importance, source, created_at
		FROM memory_notes
		WHERE conversation_id = ?
		ORDER BY created_at DESC, id DESC
		LIMIT ?
	`
	rows, err := m.db.QueryContext(ctx, query, strings.TrimSpace(conversationID), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.ConversationID, &note.Topic, &note.Summary, &note.Kind, &note.Importance, &note.Source, &note.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan note row: %w", err)
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	reverseNotes(notes)
	return notes, nil
}

func reverseNotes(notes []Note) {
	for i, j := 0, len(notes)-1; i < j; i, j = i+1, j-1 {
		notes[i], notes[j] = notes[j], notes[i]
	}
}

func mergeNoteSummaries(existing, incoming string) string {
	existing = strings.TrimSpace(existing)
	incoming = strings.TrimSpace(incoming)

	switch {
	case existing == "":
		return incoming
	case incoming == "":
		return existing
	case strings.Contains(existing, incoming):
		return existing
	case strings.Contains(incoming, existing):
		return incoming
	default:
		return existing + " | " + incoming
	}
}
