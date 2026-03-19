package mcp

import "testing"

func TestBuildToolDescription_CompactsVerboseDescription(t *testing.T) {
	t.Parallel()

	got := buildToolDescription("supabase", "execute_sql", "Execute SQL queries against the active database. This extra explanation should not survive.")
	want := "Execute SQL queries against the active database."
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBuildToolDescription_UsesShortFallback(t *testing.T) {
	t.Parallel()

	got := buildToolDescription("supabase", "execute_sql", "")
	want := "MCP tool execute_sql."
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
