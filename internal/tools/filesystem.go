package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
)

// ReadFileHandler returns the contents of a local file.
func ReadFileHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	pathVal, err := requireStringArg(args, "path")
	if err != nil {
		return "", err
	}

	pathVal, err = resolvePath(ctx, pathVal, args)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(pathVal)
	if err != nil {
		return fmt.Sprintf("error reading file: %v", err), nil
	}

	return string(content), nil
}

// WriteFileHandler writes full contents to a local file.
func WriteFileHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	pathVal, err := requireStringArg(args, "path")
	if err != nil {
		return "", err
	}

	contentVal, err := requireStringArg(args, "content")
	if err != nil {
		return "", err
	}

	pathVal, err = resolvePath(ctx, pathVal, args)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(pathVal)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Sprintf("error creating parent directories: %v", err), nil
	}

	if err := os.WriteFile(pathVal, []byte(contentVal), 0644); err != nil {
		return fmt.Sprintf("error writing file: %v", err), nil
	}

	return "success: file written successfully", nil
}

// ListDirHandler lists entries from a local directory.
func ListDirHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	pathVal, err := requireStringArg(args, "path")
	if err != nil {
		return "", err
	}

	pathVal, err = resolvePath(ctx, pathVal, args)
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(pathVal)
	if err != nil {
		return fmt.Sprintf("error reading directory: %v", err), nil
	}

	var lines []string
	for _, e := range entries {
		if e.IsDir() {
			lines = append(lines, fmt.Sprintf("[DIR] %s", e.Name()))
			continue
		}
		lines = append(lines, e.Name())
	}

	return strings.Join(lines, "\n"), nil
}

func resolvePath(ctx context.Context, pathVal string, args map[string]interface{}) (string, error) {
	if filepath.IsAbs(pathVal) {
		return filepath.Clean(pathVal), nil
	}

	workdir := optionalStringArg(args, "workdir")
	if workdir == "" {
		workdir, _ = agent.WorkdirFromContext(ctx)
	}
	if workdir == "" {
		if _, _, ok := agent.TaskContextFromContext(ctx); ok {
			return "", fmt.Errorf("relative path requires explicit workdir or task default workdir")
		}
		return filepath.Clean(pathVal), nil
	}

	return filepath.Clean(filepath.Join(workdir, pathVal)), nil
}


