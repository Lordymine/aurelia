package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
)

func runOpenAIAuthLogin(stdin io.Reader, stdout io.Writer) error {
	resolver, err := runtime.New()
	if err != nil {
		return fmt.Errorf("resolve instance root: %w", err)
	}
	if err := runtime.Bootstrap(resolver); err != nil {
		return fmt.Errorf("bootstrap instance directory: %w", err)
	}

	if err := ensureCodexCLIAvailable(); err != nil {
		return fmt.Errorf("%w; install Codex CLI first", err)
	}

	editable, err := config.LoadEditable(resolver)
	if err != nil {
		return fmt.Errorf("load editable config: %w", err)
	}
	editable.LLMProvider = "openai"
	editable.OpenAIAuthMode = "codex"
	if err := config.SaveEditable(resolver, *editable); err != nil {
		return fmt.Errorf("save config before auth: %w", err)
	}

	if err := runCodexLoginCommand(stdin, stdout, "--device-auth"); err != nil {
		return fmt.Errorf("run codex login: %w", err)
	}

	_, _ = fmt.Fprintln(stdout, "OpenAI Codex login complete. Aurelia is configured to use Codex CLI mode.")
	return nil
}

func runCodexLoginCommand(stdin io.Reader, stdout io.Writer, args ...string) error {
	cmdArgs := append([]string{"login"}, args...)
	cmd := exec.Command("codex", cmdArgs...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
