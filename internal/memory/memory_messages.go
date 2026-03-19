package memory

import (
	"context"
	"fmt"
	"strings"
)

// AddMessage inserts a new message into the history
func (m *MemoryManager) AddMessage(ctx context.Context, conversationID, role, content string) error {
	cleanContent := strings.ReplaceAll(content, "\x00", "")
	query := `INSERT INTO messages (conversation_id, role, content) VALUES (?, ?, ?)`
	_, err := m.db.ExecContext(ctx, query, conversationID, role, cleanContent)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}
	return nil
}

// GetRecentMessages retrieves the most recent N messages for a conversation
func (m *MemoryManager) GetRecentMessages(ctx context.Context, conversationID string) ([]Message, error) {
	query := `
		SELECT id, conversation_id, role, content, created_at 
		FROM messages 
		WHERE conversation_id = ? 
		ORDER BY created_at DESC, id DESC 
		LIMIT ?
	`
	rows, err := m.db.QueryContext(ctx, query, conversationID, m.memoryWindowSize)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message row: %w", err)
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	reverseMessages(messages)
	return messages, nil
}

// ListMessages retrieves up to limit recent messages for a conversation in ascending order.
// When limit is zero or negative, it returns all messages for the conversation.
func (m *MemoryManager) ListMessages(ctx context.Context, conversationID string, limit int) ([]Message, error) {
	baseQuery := `
		SELECT id, conversation_id, role, content, created_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY created_at DESC, id DESC
	`

	var (
		rowsErr error
		rows    interface {
			Next() bool
			Scan(dest ...any) error
			Err() error
			Close() error
		}
	)

	if limit > 0 {
		query := baseQuery + "\nLIMIT ?"
		rows, rowsErr = m.db.QueryContext(ctx, query, conversationID, limit)
	} else {
		rows, rowsErr = m.db.QueryContext(ctx, baseQuery, conversationID)
	}
	if rowsErr != nil {
		return nil, fmt.Errorf("failed to list messages: %w", rowsErr)
	}
	defer func() { _ = rows.Close() }()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message row: %w", err)
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	reverseMessages(messages)
	return messages, nil
}

// TrimMessages keeps only the newest keep messages for a conversation.
func (m *MemoryManager) TrimMessages(ctx context.Context, conversationID string, keep int) error {
	if keep < 0 {
		keep = 0
	}

	query := `
		DELETE FROM messages
		WHERE conversation_id = ?
		  AND id NOT IN (
			SELECT id
			FROM messages
			WHERE conversation_id = ?
			ORDER BY created_at DESC, id DESC
			LIMIT ?
		  )
	`
	if _, err := m.db.ExecContext(ctx, query, conversationID, conversationID, keep); err != nil {
		return fmt.Errorf("failed to trim messages: %w", err)
	}
	return nil
}

// EnsureConversation makes sure the conversation entry exists
func (m *MemoryManager) EnsureConversation(ctx context.Context, conversationID string, userID int64, provider string) error {
	query := `INSERT INTO conversations (id, user_id, provider) VALUES (?, ?, ?) ON CONFLICT(id) DO UPDATE SET provider = ?`
	_, err := m.db.ExecContext(ctx, query, conversationID, userID, provider, provider)
	if err != nil {
		return fmt.Errorf("failed to upsert conversation: %w", err)
	}
	return nil
}

func reverseMessages(messages []Message) {
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
}
