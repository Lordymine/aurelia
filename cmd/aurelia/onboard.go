package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
	"golang.org/x/term"
)

const (
	colorBlue  = "\x1b[94m"
	colorReset = "\x1b[0m"
)

type onboardStep int

const (
	stepLLMProvider onboardStep = iota
	stepLLMKey
	stepSTTProvider
	stepSTTKey
	stepTelegramToken
	stepTelegramUsers
	stepRuntimeMaxIterations
	stepRuntimeMemoryWindow
	stepReview
)

type keyCode int

const (
	keyUnknown keyCode = iota
	keyUp
	keyDown
	keyLeft
	keyRight
	keyEnter
	keyBackspace
	keyRune
	keyQuit
)

type keyEvent struct {
	code keyCode
	r    rune
}

type onboardingUI struct {
	cfg           config.EditableConfig
	step          onboardStep
	menuIndex     int
	input         string
	message       string
	reviewOptions []string
}

func runOnboard(stdin io.Reader, stdout io.Writer) error {
	resolver, err := runtime.New()
	if err != nil {
		return fmt.Errorf("resolve instance root: %w", err)
	}
	if err := runtime.Bootstrap(resolver); err != nil {
		return fmt.Errorf("bootstrap instance directory: %w", err)
	}

	current, err := config.LoadEditable(resolver)
	if err != nil {
		return fmt.Errorf("load editable config: %w", err)
	}

	inFile, inOK := stdin.(*os.File)
	outFile, outOK := stdout.(*os.File)
	if inOK && outOK && term.IsTerminal(int(inFile.Fd())) && term.IsTerminal(int(outFile.Fd())) {
		if err := runOnboardTUI(inFile, outFile, resolver, current); err != nil {
			return err
		}
		return nil
	}

	return runOnboardPrompt(stdin, stdout, resolver, current)
}

func runOnboardPrompt(stdin io.Reader, stdout io.Writer, resolver *runtime.PathResolver, current *config.EditableConfig) error {
	reader := bufio.NewReader(stdin)

	if err := writeString(stdout, renderOnboardingHeader()); err != nil {
		return err
	}
	if err := writef(stdout, "Config file: %s\n", resolver.AppConfig()); err != nil {
		return err
	}
	if err := writeln(stdout, "Press Enter to keep the current value."); err != nil {
		return err
	}
	if err := writeln(stdout, ""); err != nil {
		return err
	}

	if err := writef(stdout, "LLM provider [%s]: %s\n", current.LLMProvider, "Kimi"); err != nil {
		return err
	}
	if err := writef(stdout, "STT provider [%s]: %s\n\n", current.STTProvider, "Groq"); err != nil {
		return err
	}

	current.KimiAPIKey, _ = promptString(reader, stdout, "Kimi API key", current.KimiAPIKey, true)
	current.GroqAPIKey, _ = promptString(reader, stdout, "Groq API key", current.GroqAPIKey, true)
	current.TelegramBotToken, _ = promptString(reader, stdout, "Telegram bot token", current.TelegramBotToken, true)
	current.TelegramAllowedUserIDs, _ = promptInt64List(reader, stdout, "Telegram allowed user IDs (comma-separated)", current.TelegramAllowedUserIDs)
	current.MaxIterations, _ = promptInt(reader, stdout, "Max iterations", current.MaxIterations)
	current.MemoryWindowSize, _ = promptInt(reader, stdout, "Memory window size", current.MemoryWindowSize)

	current.LLMProvider = "kimi"
	current.STTProvider = "groq"

	if err := config.SaveEditable(resolver, *current); err != nil {
		return fmt.Errorf("save app config: %w", err)
	}

	return renderSavedSummary(stdout, resolver, current)
}

func runOnboardTUI(stdin *os.File, stdout *os.File, resolver *runtime.PathResolver, current *config.EditableConfig) error {
	oldState, err := term.MakeRaw(int(stdin.Fd()))
	if err != nil {
		return fmt.Errorf("enable raw terminal mode: %w", err)
	}
	defer func() { _ = term.Restore(int(stdin.Fd()), oldState) }()

	ui := newOnboardingUI(*current)
	reader := bufio.NewReader(stdin)

	for {
		if _, err := io.WriteString(stdout, ui.View(resolver)); err != nil {
			return err
		}

		ev, err := readKey(reader)
		if err != nil {
			return err
		}

		saved, cancelled, err := ui.HandleKey(ev)
		if err != nil {
			return err
		}
		if cancelled {
			clearScreen(stdout)
			if err := writeln(stdout, "Onboarding canceled."); err != nil {
				return err
			}
			return nil
		}
		if saved {
			if err := config.SaveEditable(resolver, ui.cfg); err != nil {
				return fmt.Errorf("save app config: %w", err)
			}
			clearScreen(stdout)
			return renderSavedSummary(stdout, resolver, &ui.cfg)
		}
	}
}

func newOnboardingUI(cfg config.EditableConfig) *onboardingUI {
	if cfg.LLMProvider == "" {
		cfg.LLMProvider = "kimi"
	}
	if cfg.STTProvider == "" {
		cfg.STTProvider = "groq"
	}
	return &onboardingUI{
		cfg:           cfg,
		step:          stepLLMProvider,
		reviewOptions: []string{"Save config", "Back", "Cancel"},
	}
}

func (u *onboardingUI) View(resolver *runtime.PathResolver) string {
	var b strings.Builder
	b.WriteString("\x1b[2J\x1b[H")
	b.WriteString(renderOnboardingHeader())
	_, _ = fmt.Fprintf(&b, "Config file: %s\n", resolver.AppConfig())
	_, _ = fmt.Fprintf(&b, "Step %d/9\n\n", int(u.step)+1)
	if u.message != "" {
		b.WriteString(colorize("! "+u.message, colorBlue))
		b.WriteString("\n\n")
	}

	switch u.step {
	case stepLLMProvider:
		b.WriteString("LLM Provider\n")
		b.WriteString("Select the main chat model provider.\n\n")
		b.WriteString(renderMenu([]string{"Kimi"}, u.menuIndex))
		b.WriteString("\nUse ↑/↓ and Enter.\n")
	case stepLLMKey:
		b.WriteString(u.renderInputStep("Kimi API key", "Used for the main LLM runtime.", true))
	case stepSTTProvider:
		b.WriteString("STT Provider\n")
		b.WriteString("Select the speech-to-text provider.\n\n")
		b.WriteString(renderMenu([]string{"Groq"}, u.menuIndex))
		b.WriteString("\nUse ↑/↓ and Enter. Use ← to go back.\n")
	case stepSTTKey:
		b.WriteString(u.renderInputStep("Groq API key", "Used for speech transcription.", true))
	case stepTelegramToken:
		b.WriteString(u.renderInputStep("Telegram bot token", "Used by the Telegram bot interface.", true))
	case stepTelegramUsers:
		b.WriteString(u.renderInputStep("Telegram allowed user IDs", "Comma-separated list, e.g. 123,456.", false))
	case stepRuntimeMaxIterations:
		b.WriteString(u.renderInputStep("Max iterations", "Maximum loop iterations per run.", false))
	case stepRuntimeMemoryWindow:
		b.WriteString(u.renderInputStep("Memory window size", "How many recent messages stay in the working window.", false))
	case stepReview:
		b.WriteString("Review & Save\n")
		b.WriteString("Check the config before saving.\n\n")
		_, _ = fmt.Fprintf(&b, "LLM provider: %s\n", strings.ToUpper(u.cfg.LLMProvider))
		_, _ = fmt.Fprintf(&b, "Kimi API key: %s\n", maskSecret(u.cfg.KimiAPIKey))
		_, _ = fmt.Fprintf(&b, "STT provider: %s\n", strings.ToUpper(u.cfg.STTProvider))
		_, _ = fmt.Fprintf(&b, "Groq API key: %s\n", maskSecret(u.cfg.GroqAPIKey))
		_, _ = fmt.Fprintf(&b, "Telegram bot token: %s\n", maskSecret(u.cfg.TelegramBotToken))
		_, _ = fmt.Fprintf(&b, "Telegram allowed user IDs: %s\n", formatInt64List(u.cfg.TelegramAllowedUserIDs))
		_, _ = fmt.Fprintf(&b, "Max iterations: %d\n", u.cfg.MaxIterations)
		_, _ = fmt.Fprintf(&b, "Memory window size: %d\n\n", u.cfg.MemoryWindowSize)
		b.WriteString(renderMenu(u.reviewOptions, u.menuIndex))
		b.WriteString("\nUse ↑/↓ and Enter. Use ← to go back. Press Ctrl+C to cancel.\n")
	}

	return b.String()
}

func (u *onboardingUI) renderInputStep(label, help string, secret bool) string {
	var b strings.Builder
	b.WriteString(label)
	b.WriteString("\n")
	b.WriteString(help)
	b.WriteString("\n\n")
	display := u.input
	if secret {
		display = maskForInput(display)
	}
	b.WriteString("> ")
	b.WriteString(display)
	b.WriteString("\n\nType and press Enter. Use ← to go back. Press Ctrl+C to cancel.\n")
	return b.String()
}

func (u *onboardingUI) HandleKey(ev keyEvent) (saved bool, cancelled bool, err error) {
	u.message = ""

	switch u.step {
	case stepLLMProvider:
		return u.handleMenuKey(ev, []string{"kimi"}, stepLLMKey, stepLLMProvider)
	case stepSTTProvider:
		return u.handleMenuKey(ev, []string{"groq"}, stepSTTKey, stepLLMKey)
	case stepReview:
		return u.handleReviewKey(ev)
	default:
		return u.handleInputKey(ev)
	}
}

func (u *onboardingUI) handleMenuKey(ev keyEvent, values []string, next onboardStep, prev onboardStep) (bool, bool, error) {
	switch ev.code {
	case keyUp:
		u.menuIndex = wrapIndex(u.menuIndex-1, len(values))
	case keyDown:
		u.menuIndex = wrapIndex(u.menuIndex+1, len(values))
	case keyEnter:
		switch u.step {
		case stepLLMProvider:
			u.cfg.LLMProvider = values[u.menuIndex]
		case stepSTTProvider:
			u.cfg.STTProvider = values[u.menuIndex]
		}
		u.step = next
		u.input = u.currentInputValue()
		u.menuIndex = 0
	case keyLeft:
		if u.step != prev {
			u.step = prev
			u.input = u.currentInputValue()
			u.menuIndex = 0
		}
	case keyQuit:
		return false, true, nil
	}
	return false, false, nil
}

func (u *onboardingUI) handleInputKey(ev keyEvent) (bool, bool, error) {
	switch ev.code {
	case keyRune:
		u.input += string(ev.r)
	case keyBackspace:
		if len(u.input) > 0 {
			u.input = u.input[:len(u.input)-1]
		}
	case keyLeft:
		u.step = previousStep(u.step)
		u.input = u.currentInputValue()
	case keyEnter:
		if err := u.commitInput(); err != nil {
			u.message = err.Error()
			return false, false, nil
		}
		u.step = nextStep(u.step)
		u.menuIndex = 0
		u.input = u.currentInputValue()
	case keyQuit:
		return false, true, nil
	}
	return false, false, nil
}

func (u *onboardingUI) handleReviewKey(ev keyEvent) (bool, bool, error) {
	switch ev.code {
	case keyUp:
		u.menuIndex = wrapIndex(u.menuIndex-1, len(u.reviewOptions))
	case keyDown:
		u.menuIndex = wrapIndex(u.menuIndex+1, len(u.reviewOptions))
	case keyLeft:
		u.step = stepRuntimeMemoryWindow
		u.input = u.currentInputValue()
	case keyEnter:
		switch u.menuIndex {
		case 0:
			return true, false, nil
		case 1:
			u.step = stepRuntimeMemoryWindow
			u.input = u.currentInputValue()
		case 2:
			return false, true, nil
		}
	case keyQuit:
		return false, true, nil
	}
	return false, false, nil
}

func (u *onboardingUI) commitInput() error {
	switch u.step {
	case stepLLMKey:
		u.cfg.KimiAPIKey = strings.TrimSpace(u.input)
	case stepSTTKey:
		u.cfg.GroqAPIKey = strings.TrimSpace(u.input)
	case stepTelegramToken:
		u.cfg.TelegramBotToken = strings.TrimSpace(u.input)
	case stepTelegramUsers:
		values, err := parseInt64List(u.input)
		if err != nil {
			return err
		}
		u.cfg.TelegramAllowedUserIDs = values
	case stepRuntimeMaxIterations:
		value, err := strconv.Atoi(strings.TrimSpace(u.input))
		if err != nil || value <= 0 {
			return errors.New("max iterations must be a positive integer")
		}
		u.cfg.MaxIterations = value
	case stepRuntimeMemoryWindow:
		value, err := strconv.Atoi(strings.TrimSpace(u.input))
		if err != nil || value <= 0 {
			return errors.New("memory window size must be a positive integer")
		}
		u.cfg.MemoryWindowSize = value
	}
	return nil
}

func (u *onboardingUI) currentInputValue() string {
	switch u.step {
	case stepLLMKey:
		return u.cfg.KimiAPIKey
	case stepSTTKey:
		return u.cfg.GroqAPIKey
	case stepTelegramToken:
		return u.cfg.TelegramBotToken
	case stepTelegramUsers:
		return formatInt64CSV(u.cfg.TelegramAllowedUserIDs)
	case stepRuntimeMaxIterations:
		return strconv.Itoa(u.cfg.MaxIterations)
	case stepRuntimeMemoryWindow:
		return strconv.Itoa(u.cfg.MemoryWindowSize)
	default:
		return ""
	}
}

func nextStep(step onboardStep) onboardStep {
	if step >= stepReview {
		return stepReview
	}
	return step + 1
}

func previousStep(step onboardStep) onboardStep {
	if step <= stepLLMProvider {
		return stepLLMProvider
	}
	return step - 1
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

func renderMenu(options []string, selected int) string {
	var b strings.Builder
	for i, option := range options {
		prefix := "  "
		if i == selected {
			prefix = colorize("> ", colorBlue)
		}
		b.WriteString(prefix)
		b.WriteString(option)
		b.WriteString("\n")
	}
	return b.String()
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

	value, err := strconv.Atoi(line)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", label)
	}
	return value, nil
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
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func clearScreen(w io.Writer) {
	_, _ = io.WriteString(w, "\x1b[2J\x1b[H")
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
	if err := writef(stdout, "Kimi API key: %s\n", maskSecret(current.KimiAPIKey)); err != nil {
		return err
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
