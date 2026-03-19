package observability

import (
	"context"
	"testing"
)

func TestSQLiteStore_RecordAndListOperations(t *testing.T) {
	t.Parallel()

	store := setupStore(t)
	ctx := context.Background()

	if err := store.RecordOperation(ctx, Operation{
		RunID:      "run-1",
		Component:  "agent.loop",
		Operation:  "llm_generate",
		Status:     "ok",
		DurationMS: 123,
		Summary:    "tool_calls=0",
	}); err != nil {
		t.Fatalf("RecordOperation(ok) error = %v", err)
	}
	if err := store.RecordOperation(ctx, Operation{
		RunID:      "run-2",
		Component:  "agent.tool",
		Operation:  "run_command",
		Status:     "error",
		DurationMS: 77,
		Summary:    "command failed",
	}); err != nil {
		t.Fatalf("RecordOperation(error) error = %v", err)
	}

	recent, err := store.ListRecentOperations(ctx, 10)
	if err != nil {
		t.Fatalf("ListRecentOperations() error = %v", err)
	}
	if len(recent) != 2 {
		t.Fatalf("expected 2 recent operations, got %d", len(recent))
	}

	failures, err := store.ListFailedOperations(ctx, 10)
	if err != nil {
		t.Fatalf("ListFailedOperations() error = %v", err)
	}
	if len(failures) != 1 {
		t.Fatalf("expected 1 failed operation, got %d", len(failures))
	}
	if failures[0].Operation != "run_command" {
		t.Fatalf("unexpected failed operation %+v", failures[0])
	}
}

func TestSQLiteStore_RecordArtifact(t *testing.T) {
	t.Parallel()

	store := setupStore(t)
	ctx := context.Background()

	id, err := store.RecordArtifact(ctx, Artifact{
		RunID:     "run-3",
		Component: "agent.tool",
		Operation: "read_file",
		Content:   "full raw output",
	})
	if err != nil {
		t.Fatalf("RecordArtifact() error = %v", err)
	}
	if id == 0 {
		t.Fatalf("expected non-zero artifact id")
	}
}

func setupStore(t *testing.T) *SQLiteStore {
	t.Helper()

	store, err := NewSQLiteStore(t.TempDir() + "\\ops.db")
	if err != nil {
		t.Fatalf("NewSQLiteStore() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}
