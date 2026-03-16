package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/agent"
)

const (
	defaultCommandTimeout = 120 * time.Second
	maxCommandTimeout     = 600 * time.Second
	maxCommandOutputChars = 4000
)

var blockedCommandPatterns = []string{
	"rm -rf /",
	"del /f /s /q",
	"shutdown",
	"reboot",
	"mkfs",
	"diskpart",
	"git reset --hard",
	"git clean -fdx",
}

type commandResult struct {
	Command    string `json:"command"`
	Workdir    string `json:"workdir"`
	ExitCode   int    `json:"exit_code"`
	DurationMS int64  `json:"duration_ms"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	TimedOut   bool   `json:"timed_out,omitempty"`
}

func RunCommandHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	command := optionalStringArg(args, "command")
	if command == "" {
		return marshalCommandResult(commandResult{
			ExitCode: -2,
			Stderr:   "blocked: command is required",
		}), nil
	}

	workdir := optionalStringArg(args, "workdir")
	workspaceRoot := optionalStringArg(args, "workspace_root")
	if workdir == "" {
		workdir, _ = agent.WorkdirFromContext(ctx)
	}
	if workdir == "" {
		if _, _, ok := agent.TaskContextFromContext(ctx); ok {
			return marshalCommandResult(commandResult{
				Command:  command,
				ExitCode: -2,
				Stderr:   "blocked: task requires explicit workdir or task default workdir",
			}), nil
		}
		cwd, err := os.Getwd()
		if err == nil {
			workdir = cwd
		} else {
			workdir = "."
		}
	}
	workdir = filepath.Clean(workdir)

	timeout := defaultCommandTimeout
	if rawTimeout, ok := args["timeout_seconds"].(float64); ok && rawTimeout > 0 {
		timeout = time.Duration(rawTimeout * float64(time.Second))
		if timeout > maxCommandTimeout {
			timeout = maxCommandTimeout
		}
	}

	result := commandResult{
		Command: command,
		Workdir: workdir,
	}

	if isBlockedCommand(command) {
		result.ExitCode = -2
		result.Stderr = "blocked: command matches a forbidden pattern"
		return marshalCommandResult(result), nil
	}
	if workspaceRoot != "" && !isWithinWorkspace(workdir, workspaceRoot) {
		result.ExitCode = -2
		result.Stderr = "blocked: workdir is outside workspace root"
		return marshalCommandResult(result), nil
	}

	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	cmd := exec.CommandContext(runCtx, "powershell", "-NoProfile", "-Command", command)
	cmd.Dir = workdir

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	result.DurationMS = time.Since(start).Milliseconds()
	result.Stdout = truncateOutput(stdoutBuf.String())
	result.Stderr = truncateOutput(stderrBuf.String())

	if runCtx.Err() == context.DeadlineExceeded {
		result.ExitCode = -1
		result.TimedOut = true
		if strings.TrimSpace(result.Stderr) == "" {
			result.Stderr = "process killed by timeout"
		}
		return marshalCommandResult(result), nil
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = 1
			if strings.TrimSpace(result.Stderr) == "" {
				result.Stderr = err.Error()
			}
		}
		return marshalCommandResult(result), nil
	}

	result.ExitCode = 0
	return marshalCommandResult(result), nil
}

func isBlockedCommand(command string) bool {
	normalized := strings.ToLower(strings.TrimSpace(command))
	if normalized == "format" || strings.HasPrefix(normalized, "format ") {
		return true
	}
	for _, pattern := range blockedCommandPatterns {
		if strings.Contains(normalized, pattern) {
			return true
		}
	}
	return false
}

func isWithinWorkspace(workdir, workspaceRoot string) bool {
	workdir = filepath.Clean(workdir)
	workspaceRoot = filepath.Clean(workspaceRoot)

	rel, err := filepath.Rel(workspaceRoot, workdir)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, "..") && rel != "..")
}

func truncateOutput(text string) string {
	text = strings.TrimSpace(text)
	if len(text) <= maxCommandOutputChars {
		return text
	}
	return text[:maxCommandOutputChars] + "\n[truncated]"
}

func marshalCommandResult(result commandResult) string {
	payload, err := json.Marshal(result)
	if err != nil {
		fallback := commandResult{
			Command:  result.Command,
			Workdir:  result.Workdir,
			ExitCode: 1,
			Stderr:   "failed to marshal command result",
		}
		payload, _ = json.Marshal(fallback)
	}
	return string(payload)
}


