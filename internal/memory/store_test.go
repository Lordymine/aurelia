package memory

import (
	"context"
	"strings"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()

	store, err := NewStore(":memory:", NewMockEmbedder(128))
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	t.Cleanup(func() {
		_ = store.Close()
	})

	return store
}

func TestStore_SaveAndSearch(t *testing.T) {
	t.Parallel()

	store := newTestStore(t)
	ctx := context.Background()

	if err := store.Save(ctx, "Go is a compiled language", "fact", ""); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if err := store.Save(ctx, "Python is interpreted", "fact", ""); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	results, err := store.Search(ctx, "compiled programming", 1)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.Contains(results[0].Content, "Go") {
		t.Fatalf("expected result to contain 'Go', got %q", results[0].Content)
	}
}

func TestStore_SearchReturnsMultiple(t *testing.T) {
	t.Parallel()

	store := newTestStore(t)
	ctx := context.Background()

	if err := store.Save(ctx, "Go is a compiled language", "fact", ""); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if err := store.Save(ctx, "Rust is also compiled", "fact", ""); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if err := store.Save(ctx, "Python is interpreted", "fact", ""); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	results, err := store.Search(ctx, "compiled programming language", 3)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestStore_Inject(t *testing.T) {
	t.Parallel()

	store := newTestStore(t)
	ctx := context.Background()

	if err := store.Save(ctx, "user prefers dark mode", "preference", ""); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	block, err := store.Inject(ctx, "what theme does the user like", 5)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !strings.Contains(block, "dark mode") {
		t.Fatalf("expected inject to contain 'dark mode', got %q", block)
	}
	if !strings.Contains(block, "## Relevant Memories") {
		t.Fatalf("expected inject to contain header, got %q", block)
	}
	if !strings.Contains(block, "[preference]") {
		t.Fatalf("expected inject to contain category, got %q", block)
	}
}

func TestStore_InjectEmpty(t *testing.T) {
	t.Parallel()

	store := newTestStore(t)
	ctx := context.Background()

	block, err := store.Inject(ctx, "anything", 5)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if block != "" {
		t.Fatalf("expected empty inject for empty store, got %q", block)
	}
}

func TestStore_EmptySearch(t *testing.T) {
	t.Parallel()

	store := newTestStore(t)
	ctx := context.Background()

	results, err := store.Search(ctx, "anything", 10)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results for empty store, got %d", len(results))
	}
}

func TestStore_CategoryPreserved(t *testing.T) {
	t.Parallel()

	store := newTestStore(t)
	ctx := context.Background()

	categories := []string{"fact", "conversation", "decision", "preference"}
	for _, cat := range categories {
		if err := store.Save(ctx, "memory for "+cat, cat, ""); err != nil {
			t.Fatalf("Save(%s) error = %v", cat, err)
		}
	}

	results, err := store.Search(ctx, "memory for", 10)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	foundCategories := make(map[string]bool)
	for _, r := range results {
		foundCategories[r.Category] = true
	}
	for _, cat := range categories {
		if !foundCategories[cat] {
			t.Errorf("category %q not found in results", cat)
		}
	}
}

func TestStore_AgentField(t *testing.T) {
	t.Parallel()

	store := newTestStore(t)
	ctx := context.Background()

	if err := store.Save(ctx, "agent-specific memory", "fact", "prospector"); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	results, err := store.Search(ctx, "agent-specific memory", 1)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Agent != "prospector" {
		t.Fatalf("expected agent 'prospector', got %q", results[0].Agent)
	}
}

func TestStore_LimitRespectsCount(t *testing.T) {
	t.Parallel()

	store := newTestStore(t)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		if err := store.Save(ctx, "memory number", "fact", ""); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	results, err := store.Search(ctx, "memory", 3)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestSerializeDeserializeEmbedding(t *testing.T) {
	t.Parallel()

	original := []float32{0.1, -0.5, 3.14, 0, -1.0}
	blob := serializeEmbedding(original)
	restored := deserializeEmbedding(blob)

	if len(restored) != len(original) {
		t.Fatalf("expected %d elements, got %d", len(original), len(restored))
	}
	for i := range original {
		if original[i] != restored[i] {
			t.Fatalf("mismatch at index %d: %f != %f", i, original[i], restored[i])
		}
	}
}
