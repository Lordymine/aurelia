package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
)

func decodeCommandResult(t *testing.T, raw string) commandResult {
	t.Helper()

	var result commandResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("failed to decode command result json: %v\nraw=%s", err, raw)
	}
	return result
}

func TestRunCommandHandler_Success(t *testing.T) {
	t.Parallel()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": "Write-Output 'ok'",
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", result.ExitCode)
	}
	if !strings.Contains(result.Stdout, "ok") {
		t.Fatalf("expected stdout to contain ok, got %q", result.Stdout)
	}
}

func TestRunCommandHandler_FailureCapturesStderr(t *testing.T) {
	t.Parallel()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": "Write-Error 'boom'; exit 3",
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != 3 {
		t.Fatalf("expected exit code 3, got %d", result.ExitCode)
	}
	if !strings.Contains(strings.ToLower(result.Stderr), "boom") {
		t.Fatalf("expected stderr to contain boom, got %q", result.Stderr)
	}
}

func TestRunCommandHandler_Timeout(t *testing.T) {
	t.Parallel()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command":         "Start-Sleep -Seconds 2",
		"timeout_seconds": float64(1),
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if !result.TimedOut {
		t.Fatalf("expected timeout flag to be true")
	}
	if result.ExitCode != -1 {
		t.Fatalf("expected exit code -1 on timeout, got %d", result.ExitCode)
	}
}

func TestRunCommandHandler_UsesWorkdir(t *testing.T) {
	t.Parallel()

	workdir := t.TempDir()
	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": "Get-Location | Select-Object -ExpandProperty Path",
		"workdir": workdir,
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if filepath.Clean(strings.TrimSpace(result.Stdout)) != filepath.Clean(workdir) {
		t.Fatalf("expected stdout workdir %q, got %q", workdir, result.Stdout)
	}
}

func TestRunCommandHandler_UsesTaskWorkdirFromContext(t *testing.T) {
	t.Parallel()

	workdir := t.TempDir()
	ctx := agent.WithWorkdirContext(agent.WithTaskContext(context.Background(), "team-ctx", "task-ctx"), workdir)
	raw, err := RunCommandHandler(ctx, map[string]interface{}{
		"command": "Get-Location | Select-Object -ExpandProperty Path",
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if filepath.Clean(strings.TrimSpace(result.Stdout)) != filepath.Clean(workdir) {
		t.Fatalf("expected task workdir %q, got %q", workdir, result.Stdout)
	}
}

func TestRunCommandHandler_BlocksTaskCommandWithoutWorkdir(t *testing.T) {
	t.Parallel()

	ctx := agent.WithTaskContext(context.Background(), "team-ctx", "task-ctx")
	raw, err := RunCommandHandler(ctx, map[string]interface{}{
		"command": "Write-Output 'ok'",
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != -2 {
		t.Fatalf("expected blocked exit code -2, got %d", result.ExitCode)
	}
	if !strings.Contains(strings.ToLower(result.Stderr), "workdir") {
		t.Fatalf("expected workdir block message, got %q", result.Stderr)
	}
}

func TestRunCommandHandler_TruncatesLongOutput(t *testing.T) {
	t.Parallel()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": "[string]::new('a', 12000)",
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if len(result.Stdout) > maxCommandOutputChars+64 {
		t.Fatalf("expected stdout to be truncated, got len=%d", len(result.Stdout))
	}
	if !strings.Contains(result.Stdout, "[truncated]") {
		t.Fatalf("expected stdout truncation marker, got %q", result.Stdout)
	}
}

func TestRunCommandHandler_BlocksDangerousCommands(t *testing.T) {
	t.Parallel()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": "git reset --hard HEAD",
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != -2 {
		t.Fatalf("expected blocked command exit code -2, got %d", result.ExitCode)
	}
	if !strings.Contains(strings.ToLower(result.Stderr), "blocked") {
		t.Fatalf("expected blocked command stderr, got %q", result.Stderr)
	}
}

func TestRunCommandHandler_AllowsCommandOutsidePreviousAllowlist(t *testing.T) {
	t.Parallel()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": "Get-Date -Format o",
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != 0 {
		t.Fatalf("expected command to be allowed, got exit=%d stderr=%q", result.ExitCode, result.Stderr)
	}
}

func TestRunCommandHandler_BlocksWorkdirOutsideWorkspace(t *testing.T) {
	t.Parallel()

	workspace := t.TempDir()
	outside := filepath.Dir(workspace)
	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command":        "Write-Output 'ok'",
		"workdir":        outside,
		"workspace_root": workspace,
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != -2 {
		t.Fatalf("expected workspace block exit code -2, got %d", result.ExitCode)
	}
	if !strings.Contains(strings.ToLower(result.Stderr), "workspace") {
		t.Fatalf("expected workspace block message, got %q", result.Stderr)
	}
}

func TestRunCommandHandler_AllowsInvokeRestMethodForLocalhost(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": "Invoke-RestMethod -Uri '" + server.URL + "' -Method GET | ConvertTo-Json -Compress",
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q", result.ExitCode, result.Stderr)
	}
	if !strings.Contains(result.Stdout, "ok") {
		t.Fatalf("expected stdout to contain response payload, got %q", result.Stdout)
	}
}

func TestRunCommandHandler_AllowsInvokeRestMethodForNonLocalhostSyntax(t *testing.T) {
	t.Parallel()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": "Write-Output \"Invoke-RestMethod -Uri 'https://example.com' -Method GET\"",
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != 0 {
		t.Fatalf("expected command to be allowed, got exit=%d stderr=%q", result.ExitCode, result.Stderr)
	}
}

func TestRunCommandHandler_AllowsSafeStartProcessForNodeOrNpm(t *testing.T) {
	t.Parallel()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": `Start-Process -FilePath "node" -ArgumentList "--version" -PassThru | Select-Object Id`,
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != 0 {
		t.Fatalf("expected Start-Process command to be allowed, got exit=%d stderr=%q", result.ExitCode, result.Stderr)
	}
}

func TestRunCommandHandler_AllowsStartProcessSyntaxInYoloMode(t *testing.T) {
	t.Parallel()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": `Write-Output 'Start-Process -FilePath "powershell" -ArgumentList "-NoProfile" -PassThru'`,
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != 0 {
		t.Fatalf("expected command to be allowed, got exit=%d stderr=%q", result.ExitCode, result.Stderr)
	}
}

func TestRunCommandHandler_AllowsNetstatForLocalDiagnostics(t *testing.T) {
	t.Parallel()

	raw, err := RunCommandHandler(context.Background(), map[string]interface{}{
		"command": "netstat -an",
	})
	if err != nil {
		t.Fatalf("RunCommandHandler() error = %v", err)
	}

	result := decodeCommandResult(t, raw)
	if result.ExitCode != 0 {
		t.Fatalf("expected netstat command to be allowed, got exit=%d stderr=%q", result.ExitCode, result.Stderr)
	}
}


