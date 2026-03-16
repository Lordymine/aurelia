package tools

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
)

func TestReadFileHandler_ResolvesRelativePathFromWorkdir(t *testing.T) {
	t.Parallel()

	workdir := t.TempDir()
	targetPath := filepath.Join(workdir, "src", "file.txt")
	if _, err := WriteFileHandler(context.Background(), map[string]interface{}{
		"path":    "src/file.txt",
		"content": "hello",
		"workdir": workdir,
	}); err != nil {
		t.Fatalf("WriteFileHandler() error = %v", err)
	}

	got, err := ReadFileHandler(context.Background(), map[string]interface{}{
		"path":    "src/file.txt",
		"workdir": workdir,
	})
	if err != nil {
		t.Fatalf("ReadFileHandler() error = %v", err)
	}
	if got != "hello" {
		t.Fatalf("unexpected file content from %q: %q", targetPath, got)
	}
}

func TestWriteFileHandler_WritesRelativePathFromWorkdir(t *testing.T) {
	t.Parallel()

	workdir := t.TempDir()
	_, err := WriteFileHandler(context.Background(), map[string]interface{}{
		"path":    "nested/output.txt",
		"content": "from-workdir",
		"workdir": workdir,
	})
	if err != nil {
		t.Fatalf("WriteFileHandler() error = %v", err)
	}

	got, err := ReadFileHandler(context.Background(), map[string]interface{}{
		"path": filepath.Join(workdir, "nested", "output.txt"),
	})
	if err != nil {
		t.Fatalf("ReadFileHandler(abs) error = %v", err)
	}
	if got != "from-workdir" {
		t.Fatalf("unexpected file content: %q", got)
	}
}

func TestListDirHandler_ListsRelativePathFromWorkdir(t *testing.T) {
	t.Parallel()

	workdir := t.TempDir()
	if _, err := WriteFileHandler(context.Background(), map[string]interface{}{
		"path":    "docs/a.txt",
		"content": "a",
		"workdir": workdir,
	}); err != nil {
		t.Fatalf("WriteFileHandler(a) error = %v", err)
	}
	if _, err := WriteFileHandler(context.Background(), map[string]interface{}{
		"path":    "docs/b.txt",
		"content": "b",
		"workdir": workdir,
	}); err != nil {
		t.Fatalf("WriteFileHandler(b) error = %v", err)
	}

	got, err := ListDirHandler(context.Background(), map[string]interface{}{
		"path":    "docs",
		"workdir": workdir,
	})
	if err != nil {
		t.Fatalf("ListDirHandler() error = %v", err)
	}
	if !strings.Contains(got, "a.txt") || !strings.Contains(got, "b.txt") {
		t.Fatalf("expected listed files in output, got %q", got)
	}
}

func TestReadFileHandler_ReturnsFriendlyErrorForMissingFile(t *testing.T) {
	t.Parallel()

	workdir := t.TempDir()
	got, err := ReadFileHandler(context.Background(), map[string]interface{}{
		"path":    "missing.txt",
		"workdir": workdir,
	})
	if err != nil {
		t.Fatalf("ReadFileHandler() error = %v", err)
	}
	if !strings.Contains(strings.ToLower(got), "error reading file") {
		t.Fatalf("expected friendly missing file error, got %q", got)
	}
}

func TestWriteFileHandler_CreatesParentDirectoriesRecursively(t *testing.T) {
	t.Parallel()

	workdir := t.TempDir()
	_, err := WriteFileHandler(context.Background(), map[string]interface{}{
		"path":    "deep/nested/tree/output.txt",
		"content": "recursive",
		"workdir": workdir,
	})
	if err != nil {
		t.Fatalf("WriteFileHandler() error = %v", err)
	}

	got, err := ReadFileHandler(context.Background(), map[string]interface{}{
		"path": filepath.Join(workdir, "deep", "nested", "tree", "output.txt"),
	})
	if err != nil {
		t.Fatalf("ReadFileHandler() error = %v", err)
	}
	if got != "recursive" {
		t.Fatalf("expected recursive file content, got %q", got)
	}
}

func TestReadFileHandler_BlocksRelativePathForTaskWithoutWorkdir(t *testing.T) {
	t.Parallel()

	ctx := agent.WithTaskContext(context.Background(), "team-1", "task-1")
	_, err := ReadFileHandler(ctx, map[string]interface{}{
		"path": "relative.txt",
	})
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "workdir") {
		t.Fatalf("expected workdir error for task-relative read, got %v", err)
	}
}


