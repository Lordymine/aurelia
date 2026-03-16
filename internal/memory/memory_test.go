package memory

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestDB(t *testing.T) *MemoryManager {
	// Create a temporary file for the database since walking WAL might be tricky in pure :memory:
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	mm, err := NewMemoryManager(dbPath, 5)
	if err != nil {
		t.Fatalf("failed to create memory manager: %v", err)
	}
	t.Cleanup(func() {
		_ = mm.Close()
		_ = os.RemoveAll(tempDir)
	})
	return mm
}

func TestEnsureConversation(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	err := mm.EnsureConversation(ctx, "conv-1", 123, "provider-a")
	if err != nil {
		t.Errorf("EnsureConversation failed: %v", err)
	}

	// Upsert with new provider should work
	err = mm.EnsureConversation(ctx, "conv-1", 123, "provider-b")
	if err != nil {
		t.Errorf("EnsureConversation upsert failed: %v", err)
	}
}

func TestAddAndGetMessages(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	_ = mm.EnsureConversation(ctx, "conv-1", 123, "provider-a")

	// Add 6 messages
	for i := 1; i <= 6; i++ {
		err := mm.AddMessage(ctx, "conv-1", "user", "msg "+string(rune('0'+i)))
		if err != nil {
			t.Fatalf("failed to add message: %v", err)
		}
	}

	msgs, err := mm.GetRecentMessages(ctx, "conv-1")
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}

	// We set window size to 5, so we should get the last 5 messages
	if len(msgs) != 5 {
		t.Errorf("expected 5 messages, got %d", len(msgs))
	}

	// Should be in chronological order (oldest first among the 5 most recent -> msg 2 to 6)
	if msgs[0].Content != "msg 2" {
		t.Errorf("expected first message to be 'msg 2', got '%s'", msgs[0].Content)
	}
	if msgs[4].Content != "msg 6" {
		t.Errorf("expected last message to be 'msg 6', got '%s'", msgs[4].Content)
	}
}

func TestNullBytesStripping(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	_ = mm.EnsureConversation(ctx, "conv-1", 123, "provider-a")

	dirtyMsg := "hello" + string([]byte{0}) + "world"
	err := mm.AddMessage(ctx, "conv-1", "user", dirtyMsg)
	if err != nil {
		t.Fatalf("failed to add message with null bytes: %v", err)
	}

	msgs, err := mm.GetRecentMessages(ctx, "conv-1")
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}

	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}

	if strings.Contains(msgs[0].Content, "\x00") {
		t.Errorf("null byte was not stripped from content")
	}
	if msgs[0].Content != "helloworld" {
		t.Errorf("expected 'helloworld', got '%s'", msgs[0].Content)
	}
}

func TestGetRecentMessages_IsolatesConversations(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	_ = mm.EnsureConversation(ctx, "conv-a", 123, "provider-a")
	_ = mm.EnsureConversation(ctx, "conv-b", 456, "provider-a")

	if err := mm.AddMessage(ctx, "conv-a", "user", "from-a"); err != nil {
		t.Fatalf("AddMessage(conv-a) error = %v", err)
	}
	if err := mm.AddMessage(ctx, "conv-b", "user", "from-b"); err != nil {
		t.Fatalf("AddMessage(conv-b) error = %v", err)
	}

	msgsA, err := mm.GetRecentMessages(ctx, "conv-a")
	if err != nil {
		t.Fatalf("GetRecentMessages(conv-a) error = %v", err)
	}
	if len(msgsA) != 1 || msgsA[0].Content != "from-a" {
		t.Fatalf("expected conv-a isolation, got %#v", msgsA)
	}

	msgsB, err := mm.GetRecentMessages(ctx, "conv-b")
	if err != nil {
		t.Fatalf("GetRecentMessages(conv-b) error = %v", err)
	}
	if len(msgsB) != 1 || msgsB[0].Content != "from-b" {
		t.Fatalf("expected conv-b isolation, got %#v", msgsB)
	}
}
