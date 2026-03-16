package memory

import (
	"context"
	"strings"
	"testing"
)

func TestUpsertAndGetFact(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	err := mm.UpsertFact(ctx, Fact{
		Scope:    "user",
		EntityID: "123",
		Key:      "user.name",
		Value:    "Rafael",
		Source:   "persona",
	})
	if err != nil {
		t.Fatalf("UpsertFact() error = %v", err)
	}

	fact, ok, err := mm.GetFact(ctx, "user", "123", "user.name")
	if err != nil {
		t.Fatalf("GetFact() error = %v", err)
	}
	if !ok {
		t.Fatal("expected fact to exist")
	}
	if fact.Value != "Rafael" {
		t.Fatalf("expected value Rafael, got %q", fact.Value)
	}
}

func TestUpsertFact_UpdatesExistingValue(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	initial := Fact{
		Scope:    "agent",
		EntityID: "default",
		Key:      "agent.role",
		Value:    "Team Lead",
		Source:   "persona",
	}
	updated := Fact{
		Scope:    "agent",
		EntityID: "default",
		Key:      "agent.role",
		Value:    "Chief Architect",
		Source:   "memory",
	}

	if err := mm.UpsertFact(ctx, initial); err != nil {
		t.Fatalf("UpsertFact(initial) error = %v", err)
	}
	if err := mm.UpsertFact(ctx, updated); err != nil {
		t.Fatalf("UpsertFact(updated) error = %v", err)
	}

	fact, ok, err := mm.GetFact(ctx, "agent", "default", "agent.role")
	if err != nil {
		t.Fatalf("GetFact() error = %v", err)
	}
	if !ok {
		t.Fatal("expected fact to exist")
	}
	if fact.Value != "Chief Architect" {
		t.Fatalf("expected updated value, got %q", fact.Value)
	}
	if fact.Source != "memory" {
		t.Fatalf("expected updated source, got %q", fact.Source)
	}
}

func TestGetFact_MissingFact(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	_, ok, err := mm.GetFact(ctx, "user", "999", "user.name")
	if err != nil {
		t.Fatalf("GetFact() error = %v", err)
	}
	if ok {
		t.Fatal("expected missing fact")
	}
}

func TestAddAndListNotes(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	err := mm.AddNote(ctx, Note{
		ConversationID: "conv-1",
		Topic:          "architecture",
		Summary:        "Decidido manter SQLite com facts e notes.",
		Kind:           "decision",
		Importance:     8,
		Source:         "conversation",
	})
	if err != nil {
		t.Fatalf("AddNote() error = %v", err)
	}

	notes, err := mm.ListRecentNotes(ctx, "conv-1", 5)
	if err != nil {
		t.Fatalf("ListRecentNotes() error = %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
	if notes[0].Summary != "Decidido manter SQLite com facts e notes." {
		t.Fatalf("unexpected note summary %q", notes[0].Summary)
	}
}

func TestAddNote_DeduplicatesSameNote(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	note := Note{
		ConversationID: "conv-1",
		Topic:          "architecture",
		Summary:        "Decidido manter SQLite com facts e notes.",
		Kind:           "decision",
		Importance:     8,
		Source:         "conversation",
	}

	if err := mm.AddNote(ctx, note); err != nil {
		t.Fatalf("AddNote() first error = %v", err)
	}
	if err := mm.AddNote(ctx, note); err != nil {
		t.Fatalf("AddNote() second error = %v", err)
	}

	notes, err := mm.ListRecentNotes(ctx, "conv-1", 10)
	if err != nil {
		t.Fatalf("ListRecentNotes() error = %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected deduplicated notes length 1, got %d", len(notes))
	}
}

func TestAddNote_ConsolidatesByTopicAndKind(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	first := Note{
		ConversationID: "conv-1",
		Topic:          "architecture",
		Summary:        "Decidido manter SQLite.",
		Kind:           "decision",
		Importance:     6,
		Source:         "conversation",
	}
	second := Note{
		ConversationID: "conv-1",
		Topic:          "architecture",
		Summary:        "Vamos usar facts e notes.",
		Kind:           "decision",
		Importance:     8,
		Source:         "conversation",
	}

	if err := mm.AddNote(ctx, first); err != nil {
		t.Fatalf("AddNote(first) error = %v", err)
	}
	if err := mm.AddNote(ctx, second); err != nil {
		t.Fatalf("AddNote(second) error = %v", err)
	}

	notes, err := mm.ListRecentNotes(ctx, "conv-1", 10)
	if err != nil {
		t.Fatalf("ListRecentNotes() error = %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected one consolidated note, got %d", len(notes))
	}
	if !strings.Contains(notes[0].Summary, "Decidido manter SQLite.") || !strings.Contains(notes[0].Summary, "Vamos usar facts e notes.") {
		t.Fatalf("expected merged summary, got %q", notes[0].Summary)
	}
	if notes[0].Importance != 8 {
		t.Fatalf("expected importance 8, got %d", notes[0].Importance)
	}
}

func TestAddAndListArchiveEntries(t *testing.T) {
	mm := setupTestDB(t)
	ctx := context.Background()

	err := mm.AddArchiveEntry(ctx, ArchiveEntry{
		ConversationID: "conv-1",
		SessionID:      "session-a",
		Role:           "user",
		Content:        "Quero manter isso minimalista.",
		MessageType:    "chat",
	})
	if err != nil {
		t.Fatalf("AddArchiveEntry() error = %v", err)
	}

	err = mm.AddArchiveEntry(ctx, ArchiveEntry{
		ConversationID: "conv-1",
		SessionID:      "session-a",
		Role:           "assistant",
		Content:        "Entendido.",
		MessageType:    "chat",
	})
	if err != nil {
		t.Fatalf("AddArchiveEntry() second error = %v", err)
	}

	entries, err := mm.ListArchiveEntries(ctx, "conv-1", 10)
	if err != nil {
		t.Fatalf("ListArchiveEntries() error = %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 archive entries, got %d", len(entries))
	}
	if entries[0].Content != "Quero manter isso minimalista." {
		t.Fatalf("unexpected first archive content %q", entries[0].Content)
	}
	if entries[1].Role != "assistant" {
		t.Fatalf("unexpected second archive role %q", entries[1].Role)
	}
}
