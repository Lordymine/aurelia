package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/pkg/llm"
)

func nextOnboardStep(cfg config.EditableConfig, step onboardStep) onboardStep {
	switch step {
	case stepLLMProvider:
		if cfg.LLMProvider == "openai" {
			return stepOpenAIAuthMode
		}
		return stepLLMKey
	case stepOpenAIAuthMode:
		if usesOpenAICodex(cfg) {
			return stepOpenAICodexLogin
		}
		return stepLLMKey
	case stepOpenAICodexLogin:
		return stepLLMModel
	case stepLLMKey:
		return stepLLMModel
	case stepLLMModel:
		return stepSTTProvider
	case stepSTTProvider:
		return stepSTTKey
	case stepSTTKey:
		return stepTelegramToken
	case stepTelegramToken:
		return stepTelegramUsers
	case stepTelegramUsers:
		return stepRuntimeMaxIterations
	case stepRuntimeMaxIterations:
		return stepRuntimeMemoryWindow
	case stepRuntimeMemoryWindow:
		return stepReview
	default:
		return stepReview
	}
}

func previousOnboardStep(cfg config.EditableConfig, step onboardStep) onboardStep {
	switch step {
	case stepOpenAIAuthMode:
		return stepLLMProvider
	case stepOpenAICodexLogin:
		return stepOpenAIAuthMode
	case stepLLMKey:
		if cfg.LLMProvider == "openai" {
			return stepOpenAIAuthMode
		}
		return stepLLMProvider
	case stepLLMModel:
		if cfg.LLMProvider == "openai" && usesOpenAICodex(cfg) {
			return stepOpenAICodexLogin
		}
		return stepLLMKey
	case stepSTTProvider:
		return stepLLMModel
	case stepSTTKey:
		return stepSTTProvider
	case stepTelegramToken:
		return stepSTTKey
	case stepTelegramUsers:
		return stepTelegramToken
	case stepRuntimeMaxIterations:
		return stepTelegramUsers
	case stepRuntimeMemoryWindow:
		return stepRuntimeMaxIterations
	case stepReview:
		return stepRuntimeMemoryWindow
	default:
		return stepLLMProvider
	}
}

func wrapIndex(index, size int) int {
	if size <= 0 {
		return 0
	}
	if index < 0 {
		return size - 1
	}
	if index >= size {
		return 0
	}
	return index
}

func selectedProviderIndex(provider string) int {
	for i, option := range llmProviderChoices() {
		if option == llm.NormalizeProvider(provider) {
			return i
		}
	}
	return 0
}

func llmProviderChoices() []string {
	return llm.ProviderChoices()
}

func llmProviderLabels() []string {
	return llm.ProviderLabels()
}

func llmKeyLabel(provider string) string {
	spec, ok := llm.Provider(provider)
	if !ok {
		spec, _ = llm.Provider("kimi")
	}
	return spec.APIKeyLabel
}

func usesOpenAICodex(cfg config.EditableConfig) bool {
	return llm.NormalizeProvider(cfg.LLMProvider) == "openai" && cfg.OpenAIAuthMode == "codex"
}

func llmKeyHelp(provider string) string {
	spec, ok := llm.Provider(provider)
	if !ok {
		spec, _ = llm.Provider("kimi")
	}
	return spec.APIKeyHelp
}

func currentLLMKey(cfg config.EditableConfig) string {
	return cfg.LLMAPIKey(cfg.LLMProvider)
}

func setCurrentLLMKey(cfg *config.EditableConfig, value string) {
	cfg.SetLLMAPIKey(cfg.LLMProvider, value)
}

func runOpenAIDeviceAuthCommand(stdin io.Reader, stdout io.Writer) error {
	return runCodexLoginCommand(stdin, stdout, "--device-auth")
}

func promptString(reader *bufio.Reader, stdout io.Writer, label, current string, secret bool) (string, error) {
	if err := writef(stdout, "%s", label); err != nil {
		return "", err
	}
	if current != "" {
		display := current
		if secret {
			display = maskSecret(current)
		}
		if err := writef(stdout, " [%s]", display); err != nil {
			return "", err
		}
	}
	if err := writeString(stdout, ": "); err != nil {
		return "", err
	}

	line, err := readLine(reader)
	if err != nil {
		return "", err
	}
	if line == "" {
		return current, nil
	}
	return line, nil
}

func promptChoice(reader *bufio.Reader, stdout io.Writer, label, current string, options []string) (string, error) {
	if err := writef(stdout, "%s [%s] (%s): ", label, current, strings.Join(options, "/")); err != nil {
		return "", err
	}

	line, err := readLine(reader)
	if err != nil {
		return "", err
	}
	if line == "" {
		return current, nil
	}

	line = strings.ToLower(strings.TrimSpace(line))
	for _, option := range options {
		if line == option {
			return line, nil
		}
	}
	return current, fmt.Errorf("%s must be one of: %s", label, strings.Join(options, ", "))
}

func promptLLMModel(reader *bufio.Reader, stdout io.Writer, current *config.EditableConfig) (string, error) {
	options, source := resolveModelOptions(*current)
	if err := writef(stdout, "LLM model catalog: %s\n", source); err != nil {
		return "", err
	}
	for _, option := range options {
		if err := writef(stdout, "- %s\n", option.Label()); err != nil {
			return "", err
		}
	}
	return promptString(reader, stdout, "LLM model", current.LLMModel, false)
}

func promptInt(reader *bufio.Reader, stdout io.Writer, label string, current int) (int, error) {
	if err := writef(stdout, "%s [%d]: ", label, current); err != nil {
		return 0, err
	}

	line, err := readLine(reader)
	if err != nil {
		return 0, err
	}
	if line == "" {
		return current, nil
	}

	return parsePositiveInt(line, label)
}

func promptInt64List(reader *bufio.Reader, stdout io.Writer, label string, current []int64) ([]int64, error) {
	if err := writef(stdout, "%s", label); err != nil {
		return nil, err
	}
	if len(current) != 0 {
		if err := writef(stdout, " [%s]", formatInt64List(current)); err != nil {
			return nil, err
		}
	}
	if err := writeString(stdout, ": "); err != nil {
		return nil, err
	}

	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	if line == "" {
		return append([]int64(nil), current...), nil
	}
	return parseInt64List(line)
}

func readKey(reader *bufio.Reader) (keyEvent, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return keyEvent{}, err
	}

	switch b {
	case '\r', '\n':
		return keyEvent{code: keyEnter}, nil
	case 8, 127:
		return keyEvent{code: keyBackspace}, nil
	case 27:
		seq := make([]byte, 2)
		if _, err := io.ReadFull(reader, seq); err != nil {
			return keyEvent{code: keyUnknown}, nil
		}
		if seq[0] == '[' {
			switch seq[1] {
			case 'A':
				return keyEvent{code: keyUp}, nil
			case 'B':
				return keyEvent{code: keyDown}, nil
			case 'C':
				return keyEvent{code: keyRight}, nil
			case 'D':
				return keyEvent{code: keyLeft}, nil
			}
		}
	case 3:
		return keyEvent{code: keyQuit}, nil
	default:
		if b >= 32 && b <= 126 {
			return keyEvent{code: keyRune, r: rune(b)}, nil
		}
	}
	return keyEvent{code: keyUnknown}, nil
}

func parseInt64List(raw string) ([]int64, error) {
	parts := strings.Split(raw, ",")
	values := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		value, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid telegram user id %q", part)
		}
		values = append(values, value)
	}
	return values, nil
}

func formatInt64List(values []int64) string {
	if len(values) == 0 {
		return "(empty)"
	}
	return formatInt64CSV(values)
}

func formatInt64CSV(values []int64) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, strconv.FormatInt(value, 10))
	}
	return strings.Join(parts, ",")
}

func parsePositiveInt(raw, label string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", label)
	}
	return value, nil
}

func maskSecret(value string) string {
	if value == "" {
		return "(empty)"
	}
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

func maskForInput(value string) string {
	if value == "" {
		return ""
	}
	return strings.Repeat("*", len(value))
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func clearScreen(w io.Writer) {
	_, _ = io.WriteString(w, "\x1b[2J\x1b[H")
}

func rawTerminalFrame(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return strings.ReplaceAll(text, "\n", "\r\n")
}

func renderSavedSummary(stdout io.Writer, resolver *runtime.PathResolver, current *config.EditableConfig) error {
	if err := writeString(stdout, renderOnboardingHeader()); err != nil {
		return err
	}
	if err := writef(stdout, "Saved config to %s\n", resolver.AppConfig()); err != nil {
		return err
	}
	if err := writef(stdout, "LLM provider: %s\n", strings.ToUpper(current.LLMProvider)); err != nil {
		return err
	}
	if current.LLMProvider == "openai" {
		if err := writef(stdout, "OpenAI auth mode: %s\n", current.OpenAIAuthMode); err != nil {
			return err
		}
	}
	if err := writef(stdout, "LLM model: %s\n", current.LLMModel); err != nil {
		return err
	}
	if usesOpenAICodex(*current) {
		if err := writef(stdout, "OpenAI Codex login: run `aurelia auth openai`\n"); err != nil {
			return err
		}
	} else {
		if err := writef(stdout, "%s: %s\n", llmKeyLabel(current.LLMProvider), maskSecret(currentLLMKey(*current))); err != nil {
			return err
		}
	}
	if err := writef(stdout, "STT provider: %s\n", strings.ToUpper(current.STTProvider)); err != nil {
		return err
	}
	if err := writef(stdout, "Groq API key: %s\n", maskSecret(current.GroqAPIKey)); err != nil {
		return err
	}
	if err := writef(stdout, "Telegram bot token: %s\n", maskSecret(current.TelegramBotToken)); err != nil {
		return err
	}
	if err := writef(stdout, "Telegram allowed user IDs: %s\n", formatInt64List(current.TelegramAllowedUserIDs)); err != nil {
		return err
	}
	if err := writef(stdout, "Max iterations: %d\n", current.MaxIterations); err != nil {
		return err
	}
	if err := writef(stdout, "Memory window size: %d\n", current.MemoryWindowSize); err != nil {
		return err
	}
	return nil
}

func renderOnboardingHeader() string {
	jellyfish := colorize(`
            .-.
         .-(   )-.
        (___.__)__)
         / /   \ \
        /_/     \_\
         \ \   / /
          \_\ /_/
`, colorBlue)

	banner := colorize(`
 $$$$$$\  $$\   $$\ $$$$$$$\  $$$$$$$$\ $$\       $$$$$$\  $$$$$$\  
$$  __$$\ $$ |  $$ |$$  __$$\ $$  _____|$$ |      \_$$  _|$$  __$$\ 
$$ /  $$ |$$ |  $$ |$$ |  $$ |$$ |      $$ |        $$ |  $$ /  $$ |
$$$$$$$$ |$$ |  $$ |$$$$$$$  |$$$$$\    $$ |        $$ |  $$$$$$$$ |
$$  __$$ |$$ |  $$ |$$  __$$< $$  __|   $$ |        $$ |  $$  __$$ |
$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |      $$ |        $$ |  $$ |  $$ |
$$ |  $$ |\$$$$$$  |$$ |  $$ |$$$$$$$$\ $$$$$$$$\ $$$$$$\ $$ |  $$ |
\__|  \__| \______/ \__|  \__|\________|\________|\______|\__|  \__|
`, colorBlue)

	return jellyfish + banner + "Local onboarding for runtime config\n\n"
}

func colorize(text, color string) string {
	return color + text + colorReset
}

func writeString(w io.Writer, text string) error {
	_, err := io.WriteString(w, text)
	return err
}

func writef(w io.Writer, format string, args ...any) error {
	_, err := fmt.Fprintf(w, format, args...)
	return err
}

func writeln(w io.Writer, text string) error {
	_, err := fmt.Fprintln(w, text)
	return err
}
