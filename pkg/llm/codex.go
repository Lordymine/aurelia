package llm

import "os/exec"

var codexLookPath = exec.LookPath

// CodexLookPathForTest overrides the codex lookup function for testing.
func CodexLookPathForTest(fn func(string) (string, error)) func() {
	previous := codexLookPath
	codexLookPath = fn
	return func() {
		codexLookPath = previous
	}
}

// UseNoopCodexCallerForTest is a no-op stub (codex runtime was removed).
func UseNoopCodexCallerForTest() func() {
	return func() {}
}

// EnsureCodexCLIAvailable checks that the codex CLI binary is reachable.
func EnsureCodexCLIAvailable() error {
	_, err := codexLookPath("codex")
	return err
}

// CodexCLIProvider is a stub — the real provider was removed.
type CodexCLIProvider struct{}
