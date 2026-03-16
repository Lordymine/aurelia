package memory

import (
	"context"
	"fmt"
	"strings"
)

// AddArchiveEntry appends a raw long-term archive entry.
func (m *MemoryManager) AddArchiveEntry(ctx context.Context, entry ArchiveEntry) error {
	entry.ConversationID = strings.TrimSpace(entry.ConversationID)
	entry.SessionID = strings.TrimSpace(entry.SessionID)
	entry.Role = strings.TrimSpace(entry.Role)
	entry.Content = strings.TrimSpace(strings.ReplaceAll(entry.Content, "\x00", ""))
	entry.MessageType = strings.TrimSpace(entry.MessageType)

	if entry.ConversationID == "" || entry.SessionID == "" || entry.Role == "" || entry.Content == "" {
		return fmt.Errorf("conversation_id, session_id, role and content are required")
	}
	if entry.MessageType == "" {
		entry.MessageType = "chat"
	}

	query := `
		INSERT INTO memory_archive (conversation_id, session_id, role, content, message_type)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := m.db.ExecContext(ctx, query, entry.ConversationID, entry.SessionID, entry.Role, entry.Content, entry.MessageType)
	if err != nil {
		return fmt.Errorf("failed to insert archive entry: %w", err)
	}
	return nil
}

// ListArchiveEntries retrieves raw archive entries for a conversation.
func (m *MemoryManager) ListArchiveEntries(ctx context.Context, conversationID string, limit int) ([]ArchiveEntry, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT id, conversation_id, session_id, role, content, message_type, created_at
		FROM memory_archive
		WHERE conversation_id = ?
		ORDER BY created_at DESC, id DESC
		LIMIT ?
	`
	rows, err := m.db.QueryContext(ctx, query, strings.TrimSpace(conversationID), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query archive: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var entries []ArchiveEntry
	for rows.Next() {
		var entry ArchiveEntry
		if err := rows.Scan(&entry.ID, &entry.ConversationID, &entry.SessionID, &entry.Role, &entry.Content, &entry.MessageType, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan archive row: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	reverseArchiveEntries(entries)
	return entries, nil
}

func reverseArchiveEntries(entries []ArchiveEntry) {
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
}
