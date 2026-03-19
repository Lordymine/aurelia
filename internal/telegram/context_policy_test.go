package telegram

import (
	"context"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/memory"
)

func TestDecideContextAction(t *testing.T) {
	t.Parallel()

	if got := memory.DecideContextAction(1000, 200); got != memory.ContextActionNone {
		t.Fatalf("expected none, got %s", got)
	}
	if got := memory.DecideContextAction(1000, 750); got != memory.ContextActionSummarize {
		t.Fatalf("expected summarize, got %s", got)
	}
	if got := memory.DecideContextAction(1000, 900); got != memory.ContextActionRotate {
		t.Fatalf("expected rotate, got %s", got)
	}
}

func TestManageConversationContext_SummarizesAndTrims(t *testing.T) {
	t.Parallel()

	mem := setupTelegramMemory(t, 20)
	ctx := context.Background()
	if err := mem.EnsureConversation(ctx, "42", 42, "openai"); err != nil {
		t.Fatalf("EnsureConversation() error = %v", err)
	}

	payload := strings.Repeat("x", 1000)
	for i := 0; i < 18; i++ {
		role := "user"
		if i%2 == 1 {
			role = "assistant"
		}
		if err := mem.AddMessage(ctx, "42", role, payload); err != nil {
			t.Fatalf("AddMessage(%d) error = %v", i, err)
		}
	}

	bc := &BotController{
		config: &config.AppConfig{
			LLMProvider:      "kimi",
			LLMModel:         "moonshot-v1-8k",
			MemoryWindowSize: 20,
		},
		memory:        mem,
		contextPolicy: memory.NewContextPolicy(mem),
	}

	if err := bc.manageConversationContext(ctx, "42"); err != nil {
		t.Fatalf("manageConversationContext() error = %v", err)
	}

	messages, err := mem.ListMessages(ctx, "42", 0)
	if err != nil {
		t.Fatalf("ListMessages() error = %v", err)
	}
	if len(messages) != 10 {
		t.Fatalf("expected 10 messages after summarize trim, got %d", len(messages))
	}

	summary, ok, err := mem.GetLatestNote(ctx, "42", memory.ContextSummaryTopic, memory.ContextSummaryKind)
	if err != nil {
		t.Fatalf("GetLatestNote() error = %v", err)
	}
	if !ok {
		t.Fatal("expected context summary note")
	}
	if !strings.Contains(summary.Summary, "usuario:") && !strings.Contains(summary.Summary, "assistente:") {
		t.Fatalf("unexpected summary %q", summary.Summary)
	}
}

func TestSummaryPrefixUsesLatestContextSummary(t *testing.T) {
	t.Parallel()

	mem := setupTelegramMemory(t, 5)
	ctx := context.Background()
	if err := mem.AddNote(ctx, memory.Note{
		ConversationID: "42",
		Topic:          memory.ContextSummaryTopic,
		Kind:           memory.ContextSummaryKind,
		Summary:        "decisoes e contexto anteriores",
		Source:         "test",
	}); err != nil {
		t.Fatalf("AddNote() error = %v", err)
	}

	bc := &BotController{memory: mem, contextPolicy: memory.NewContextPolicy(mem)}
	prefix, err := bc.summaryPrefix(ctx, "42")
	if err != nil {
		t.Fatalf("summaryPrefix() error = %v", err)
	}
	if prefix == nil {
		t.Fatal("expected summary prefix")
	}
	if !strings.Contains(prefix.Content, "decisoes e contexto anteriores") {
		t.Fatalf("unexpected prefix %#v", prefix)
	}
}

func setupTelegramMemory(t *testing.T, window int) *memory.MemoryManager {
	t.Helper()

	tempDir := t.TempDir()
	mem, err := memory.NewMemoryManager(tempDir+"\\telegram.db", window)
	if err != nil {
		t.Fatalf("NewMemoryManager() error = %v", err)
	}
	t.Cleanup(func() { _ = mem.Close() })
	return mem
}
