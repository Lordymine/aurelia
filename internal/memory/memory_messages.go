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
