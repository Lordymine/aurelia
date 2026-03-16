package persona

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadOptionalFile_ReturnsContentWhenExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	expected := "# Owner Playbook\nSome content here."
	if err := os.WriteFile(path, []byte(expected), 0644); err != nil {
		t.Fatal(err)
	}

	content, err := readOptionalFile(path)
	if err != nil {
		t.Fatalf("readOptionalFile() unexpected error = %v", err)
	}
	if content != expected {
		t.Fatalf("readOptionalFile() = %q, want %q", content, expected)
	}
}

func TestReadOptionalFile_ReturnsEmptyOnAbsence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.md")

	content, err := readOptionalFile(path)
	if err != nil {
		t.Fatalf("readOptionalFile() unexpected error for missing file = %v", err)
	}
	if content != "" {
		t.Fatalf("readOptionalFile() = %q, want empty string", content)
	}
}
