package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/pkg/llm"
	"golang.org/x/term"
)

const (
	colorBlue  = "\x1b[94m"
	colorReset = "\x1b[0m"
)

type onboardStep int

const (
	stepLLMProvider onboardStep = iota
	stepOpenAIAuthMode
	stepOpenAICodexLogin
	stepLLMKey
	stepLLMModel
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
	cfg             config.EditableConfig
	step            onboardStep
	menuIndex       int
	input           string
	message         string
	modelSource     string
	allModelOptions []llm.ModelOption
	modelOptions    []llm.ModelOption
	modelFilter     string
	modelCapability modelCapabilityFilter
	reviewOptions   []string
	pendingAction   string
}

type modelCapabilityFilter int

const (
	modelCapabilityAll modelCapabilityFilter = iota
	modelCapabilityVision
	modelCapabilityTools
	modelCapabilityFree
)

var llmModelCatalog = llm.ListModels

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

	current.LLMProvider, _ = promptChoice(reader, stdout, "LLM provider", current.LLMProvider, llmProviderChoices())
	if current.LLMProvider == "openai" {
		current.OpenAIAuthMode, _ = promptChoice(reader, stdout, "OpenAI auth mode", current.OpenAIAuthMode, []string{"api_key", "codex"})
	}
	if err := writef(stdout, "STT provider [%s]: %s\n\n", current.STTProvider, "Groq"); err != nil {
		return err
	}

	if usesOpenAICodex(*current) {
		if err := writef(stdout, "OpenAI auth mode: codex CLI (experimental). Starting device auth now.\n\n"); err != nil {
			return err
		}
		if err := runOpenAIDeviceAuthCommand(stdin, stdout); err != nil {
			return err
		}
	} else {
		currentKey := currentLLMKey(*current)
		currentKey, _ = promptString(reader, stdout, llmKeyLabel(current.LLMProvider), currentKey, true)
		setCurrentLLMKey(current, currentKey)
	}
	current.LLMModel, _ = promptLLMModel(reader, stdout, current)
	current.GroqAPIKey, _ = promptString(reader, stdout, "Groq API key", current.GroqAPIKey, true)
	current.TelegramBotToken, _ = promptString(reader, stdout, "Telegram bot token", current.TelegramBotToken, true)
	current.TelegramAllowedUserIDs, _ = promptInt64List(reader, stdout, "Telegram allowed user IDs (comma-separated)", current.TelegramAllowedUserIDs)
	current.MaxIterations, _ = promptInt(reader, stdout, "Max iterations", current.MaxIterations)
	current.MemoryWindowSize, _ = promptInt(reader, stdout, "Memory window size", current.MemoryWindowSize)

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
		if _, err := io.WriteString(stdout, rawTerminalFrame(ui.View(resolver))); err != nil {
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
		if action := ui.consumePendingAction(); action != "" {
			switch action {
			case "openai_codex_login":
				if err := term.Restore(int(stdin.Fd()), oldState); err != nil {
					return fmt.Errorf("restore terminal mode: %w", err)
				}
				clearScreen(stdout)
				if err := runOpenAIDeviceAuthCommand(stdin, stdout); err != nil {
					ui.message = err.Error()
				}
				oldState, err = term.MakeRaw(int(stdin.Fd()))
				if err != nil {
					return fmt.Errorf("re-enable raw terminal mode: %w", err)
				}
				reader = bufio.NewReader(stdin)
			}
		}
	}
}

func newOnboardingUI(cfg config.EditableConfig) *onboardingUI {
	if cfg.LLMProvider == "" {
		cfg.LLMProvider = "kimi"
	}
	if cfg.LLMModel == "" {
		cfg.LLMModel = config.DefaultEditableConfig().LLMModel
	}
	if cfg.OpenAIAuthMode == "" {
		cfg.OpenAIAuthMode = "api_key"
	}
	if cfg.STTProvider == "" {
		cfg.STTProvider = "groq"
	}
	modelOptions, modelSource := resolveModelOptions(cfg)
	return &onboardingUI{
		cfg:             cfg,
		allModelOptions: append([]llm.ModelOption(nil), modelOptions...),
		modelOptions:    append([]llm.ModelOption(nil), modelOptions...),
		modelSource:     modelSource,
		step:            stepLLMProvider,
		reviewOptions:   []string{"Save config", "Back", "Cancel"},
	}
}

func (u *onboardingUI) View(resolver *runtime.PathResolver) string {
	var b strings.Builder
	b.WriteString("\x1b[2J\x1b[H")
	b.WriteString(renderOnboardingHeader())
	_, _ = fmt.Fprintf(&b, "Config file: %s\n", resolver.AppConfig())
	_, _ = fmt.Fprintf(&b, "Step %d/12\n\n", int(u.step)+1)
	if u.message != "" {
		b.WriteString(colorize("! "+u.message, colorBlue))
		b.WriteString("\n\n")
	}

	switch u.step {
	case stepLLMProvider:
		b.WriteString("LLM Provider\n")
		b.WriteString("Select the main chat model provider.\n\n")
		b.WriteString(renderMenu(llmProviderLabels(), u.menuIndex))
		b.WriteString("\nUse ↑/↓ and Enter.\n")
	case stepOpenAIAuthMode:
		b.WriteString("OpenAI Auth Mode\n")
		b.WriteString("Choose whether OpenAI should use an API key or the local Codex CLI.\n\n")
		b.WriteString(renderMenu([]string{"API key", "Codex CLI (experimental)"}, u.menuIndex))
		b.WriteString("\nUse arrows and Enter. Use left to go back.\n")
	case stepOpenAICodexLogin:
		b.WriteString("OpenAI Codex Login\n")
		b.WriteString("Launch the Codex device-auth flow now to get the link and verification code.\n\n")
		b.WriteString(renderMenu([]string{"Launch login now", "Skip for now", "Back"}, u.menuIndex))
		b.WriteString("\nUse arrows and Enter.\n")
	case stepLLMKey:
		b.WriteString(u.renderInputStep(llmKeyLabel(u.cfg.LLMProvider), llmKeyHelp(u.cfg.LLMProvider), true))
	case stepLLMModel:
		b.WriteString("LLM Model\n")
		b.WriteString("Select the model for the chosen provider.\n\n")
		if usesProviderModelSearch(u.cfg) {
			_, _ = fmt.Fprintf(&b, "Search: %s\n", u.modelFilter)
		}
		_, _ = fmt.Fprintf(&b, "Capability filter: %s\n", u.modelCapabilityLabel())
		_, _ = fmt.Fprintf(&b, "Showing %d of %d models\n\n", len(u.modelOptions), len(u.allModelOptions))
		b.WriteString(renderModelMenu(u.modelOptions, u.menuIndex))
		_, _ = fmt.Fprintf(&b, "\nCatalog source: %s\n", u.modelSource)
		if usesProviderModelSearch(u.cfg) {
			b.WriteString("\nType to filter by model or provider. Use right to cycle capability filters. Use arrows and Enter. Backspace removes filter. Use left to go back.\n")
		} else {
			b.WriteString("\nUse right to cycle capability filters. Use arrows and Enter. Use left to go back.\n")
		}
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
		if u.cfg.LLMProvider == "openai" {
			_, _ = fmt.Fprintf(&b, "OpenAI auth mode: %s\n", u.cfg.OpenAIAuthMode)
		}
		_, _ = fmt.Fprintf(&b, "LLM model: %s\n", u.cfg.LLMModel)
		if usesOpenAICodex(u.cfg) {
			_, _ = fmt.Fprintf(&b, "OpenAI Codex login: run `aurelia auth openai`\n")
		} else {
			_, _ = fmt.Fprintf(&b, "%s: %s\n", llmKeyLabel(u.cfg.LLMProvider), maskSecret(currentLLMKey(u.cfg)))
		}
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
		return u.handleMenuKey(ev, llmProviderChoices(), nextOnboardStep(u.cfg, stepLLMProvider), stepLLMProvider)
	case stepOpenAIAuthMode:
		return u.handleOpenAIAuthModeMenuKey(ev)
	case stepOpenAICodexLogin:
		return u.handleOpenAICodexLoginKey(ev)
	case stepLLMModel:
		return u.handleModelMenuKey(ev)
	case stepSTTProvider:
		return u.handleMenuKey(ev, []string{"groq"}, stepSTTKey, stepLLMModel)
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
		targetStep := next
		switch u.step {
		case stepLLMProvider:
			u.cfg.LLMProvider = values[u.menuIndex]
			targetStep = nextOnboardStep(u.cfg, stepLLMProvider)
		case stepSTTProvider:
			u.cfg.STTProvider = values[u.menuIndex]
		}
		u.setStep(targetStep)
	case keyLeft:
		if u.step != prev {
			u.setStep(prev)
		}
	case keyQuit:
		return false, true, nil
	}
	return false, false, nil
}

func (u *onboardingUI) handleOpenAIAuthModeMenuKey(ev keyEvent) (bool, bool, error) {
	options := []string{"api_key", "codex"}
	switch ev.code {
	case keyUp:
		u.menuIndex = wrapIndex(u.menuIndex-1, len(options))
	case keyDown:
		u.menuIndex = wrapIndex(u.menuIndex+1, len(options))
	case keyEnter:
		u.cfg.OpenAIAuthMode = options[u.menuIndex]
		u.setStep(nextOnboardStep(u.cfg, stepOpenAIAuthMode))
	case keyLeft:
		u.setStep(stepLLMProvider)
		u.menuIndex = selectedProviderIndex(u.cfg.LLMProvider)
	case keyQuit:
		return false, true, nil
	}
	return false, false, nil
}

func (u *onboardingUI) handleOpenAICodexLoginKey(ev keyEvent) (bool, bool, error) {
	options := []string{"launch", "skip", "back"}
	switch ev.code {
	case keyUp:
		u.menuIndex = wrapIndex(u.menuIndex-1, len(options))
	case keyDown:
		u.menuIndex = wrapIndex(u.menuIndex+1, len(options))
	case keyEnter:
		switch options[u.menuIndex] {
		case "launch":
			u.pendingAction = "openai_codex_login"
			u.setStep(stepLLMModel)
		case "skip":
			u.setStep(stepLLMModel)
		case "back":
			u.setStep(stepOpenAIAuthMode)
			u.menuIndex = 1
			return false, false, nil
		}
	case keyLeft:
		u.setStep(stepOpenAIAuthMode)
		u.menuIndex = 1
	case keyQuit:
		return false, true, nil
	}
	return false, false, nil
}

func (u *onboardingUI) handleModelMenuKey(ev keyEvent) (bool, bool, error) {
	if len(u.modelOptions) == 0 {
		u.refreshModelOptions()
	}

	switch ev.code {
	case keyUp:
		u.menuIndex = wrapIndex(u.menuIndex-1, len(u.modelOptions))
	case keyDown:
		u.menuIndex = wrapIndex(u.menuIndex+1, len(u.modelOptions))
	case keyRight:
		u.modelCapability = nextModelCapabilityFilter(u.modelCapability)
		u.applyModelFilter()
	case keyRune:
		if usesProviderModelSearch(u.cfg) {
			u.modelFilter += string(ev.r)
			u.applyModelFilter()
		}
	case keyBackspace:
		if usesProviderModelSearch(u.cfg) && len(u.modelFilter) > 0 {
			u.modelFilter = u.modelFilter[:len(u.modelFilter)-1]
			u.applyModelFilter()
		}
	case keyEnter:
		if len(u.modelOptions) == 0 {
			u.message = "no models available for the selected provider"
			return false, false, nil
		}
		u.cfg.LLMModel = u.modelOptions[u.menuIndex].ID
		u.setStep(stepSTTProvider)
	case keyLeft:
		u.setStep(previousOnboardStep(u.cfg, stepLLMModel))
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
		u.setStep(previousOnboardStep(u.cfg, u.step))
	case keyEnter:
		if err := u.commitInput(); err != nil {
			u.message = err.Error()
			return false, false, nil
		}
		u.setStep(nextOnboardStep(u.cfg, u.step))
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
		u.setStep(stepRuntimeMemoryWindow)
	case keyEnter:
		switch u.menuIndex {
		case 0:
			return true, false, nil
		case 1:
			u.setStep(stepRuntimeMemoryWindow)
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
		setCurrentLLMKey(&u.cfg, strings.TrimSpace(u.input))
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
		return currentLLMKey(u.cfg)
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

func selectedProviderIndex(provider string) int {
	for i, option := range llmProviderChoices() {
		if option == provider {
			return i
		}
	}
	return 0
}

func llmProviderChoices() []string {
	return []string{"kimi", "anthropic", "google", "kilo", "openrouter", "zai", "alibaba", "openai"}
}

func llmProviderLabels() []string {
	return []string{"Kimi", "Anthropic", "Google", "Kilo Code", "OpenRouter", "Z.ai", "Alibaba", "OpenAI"}
}

func llmKeyLabel(provider string) string {
	switch provider {
	case "anthropic":
		return "Anthropic API key"
	case "google":
		return "Google API key"
	case "kilo":
		return "Kilo API key"
	case "openrouter":
		return "OpenRouter API key"
	case "zai":
		return "Z.ai Coding Plan API key"
	case "alibaba":
		return "Alibaba Coding Plan API key"
	case "openai":
		return "OpenAI API key"
	default:
		return "Kimi API key"
	}
}

func usesOpenAICodex(cfg config.EditableConfig) bool {
	return cfg.LLMProvider == "openai" && cfg.OpenAIAuthMode == "codex"
}

func llmKeyHelp(provider string) string {
	switch provider {
	case "anthropic":
		return "Used for the Anthropic LLM runtime."
	case "google":
		return "Used for the Google Gemini LLM runtime."
	case "kilo":
		return "Used for the Kilo Gateway LLM runtime."
	case "openrouter":
		return "Used for the OpenRouter LLM runtime."
	case "zai":
		return "Used for the Z.ai GLM Coding Plan runtime."
	case "alibaba":
		return "Used for the Alibaba Coding Plan runtime."
	case "openai":
		return "Used for the OpenAI LLM runtime."
	default:
		return "Used for the main LLM runtime."
	}
}

func currentLLMKey(cfg config.EditableConfig) string {
	switch cfg.LLMProvider {
	case "anthropic":
		return cfg.AnthropicAPIKey
	case "google":
		return cfg.GoogleAPIKey
	case "kilo":
		return cfg.KiloAPIKey
	case "openrouter":
		return cfg.OpenRouterAPIKey
	case "zai":
		return cfg.ZAIAPIKey
	case "alibaba":
		return cfg.AlibabaAPIKey
	case "openai":
		return cfg.OpenAIAPIKey
	default:
		return cfg.KimiAPIKey
	}
}

func setCurrentLLMKey(cfg *config.EditableConfig, value string) {
	switch cfg.LLMProvider {
	case "anthropic":
		cfg.AnthropicAPIKey = value
	case "google":
		cfg.GoogleAPIKey = value
	case "kilo":
		cfg.KiloAPIKey = value
	case "openrouter":
		cfg.OpenRouterAPIKey = value
	case "zai":
		cfg.ZAIAPIKey = value
	case "alibaba":
		cfg.AlibabaAPIKey = value
	case "openai":
		cfg.OpenAIAPIKey = value
	default:
		cfg.KimiAPIKey = value
	}
}

func renderModelMenu(options []llm.ModelOption, selected int) string {
	if len(options) == 0 {
		return "  No models available.\n"
	}

	var labels []string
	for _, option := range options {
		labels = append(labels, option.Label())
	}
	return renderMenu(labels, selected)
}

func loadModelOptions(cfg config.EditableConfig) []llm.ModelOption {
	options, _ := llmModelCatalog(context.Background(), cfg.LLMProvider, modelCatalogCredentials(cfg))
	if len(options) != 0 {
		return options
	}
	return llm.FallbackModels(cfg.LLMProvider)
}

func resolveModelOptions(cfg config.EditableConfig) ([]llm.ModelOption, string) {
	options, err := llmModelCatalog(context.Background(), cfg.LLMProvider, modelCatalogCredentials(cfg))
	if err == nil && len(options) != 0 {
		return options, "provider catalog"
	}

	fallback := llm.FallbackModels(cfg.LLMProvider)
	if len(fallback) != 0 {
		return fallback, "curated fallback"
	}
	return nil, "no catalog available"
}

func filterModelOptions(cfg config.EditableConfig, options []llm.ModelOption, filter string, capability modelCapabilityFilter) []llm.ModelOption {
	filter = strings.ToLower(strings.TrimSpace(filter))

	filtered := make([]llm.ModelOption, 0, len(options))
	for _, option := range options {
		if !matchesCapabilityFilter(option, capability) {
			continue
		}
		if filter != "" && usesProviderModelSearch(cfg) && !matchesModelFilter(option, filter) {
			continue
		}
		filtered = append(filtered, option)
	}
	return filtered
}

func matchesModelFilter(option llm.ModelOption, filter string) bool {
	candidates := []string{
		strings.ToLower(option.ID),
		strings.ToLower(option.Name),
		strings.ToLower(option.Label()),
		strings.ToLower(openRouterProviderName(option.ID)),
	}
	for _, candidate := range candidates {
		if strings.Contains(candidate, filter) {
			return true
		}
	}
	return false
}

func openRouterProviderName(modelID string) string {
	prefix, _, ok := strings.Cut(modelID, "/")
	if !ok {
		return ""
	}
	return prefix
}

func selectedModelIndex(options []llm.ModelOption, current string) int {
	for i, option := range options {
		if option.ID == current {
			return i
		}
	}
	return 0
}

func modelCatalogCredentials(cfg config.EditableConfig) llm.ModelCatalogCredentials {
	return llm.ModelCatalogCredentials{
		AnthropicAPIKey:  cfg.AnthropicAPIKey,
		GoogleAPIKey:     cfg.GoogleAPIKey,
		KiloAPIKey:       cfg.KiloAPIKey,
		KimiAPIKey:       cfg.KimiAPIKey,
		OpenRouterAPIKey: cfg.OpenRouterAPIKey,
		ZAIAPIKey:        cfg.ZAIAPIKey,
		AlibabaAPIKey:    cfg.AlibabaAPIKey,
		OpenAIAPIKey:     cfg.OpenAIAPIKey,
		OpenAIAuthMode:   cfg.OpenAIAuthMode,
	}
}

func usesProviderModelSearch(cfg config.EditableConfig) bool {
	switch cfg.LLMProvider {
	case "openrouter", "kilo":
		return true
	default:
		return false
	}
}

func (u *onboardingUI) consumePendingAction() string {
	action := u.pendingAction
	u.pendingAction = ""
	return action
}

func (u *onboardingUI) refreshModelOptions() {
	options, source := resolveModelOptions(u.cfg)
	u.allModelOptions = append([]llm.ModelOption(nil), options...)
	u.modelSource = source
	u.applyModelFilter()
}

func (u *onboardingUI) applyModelFilter() {
	u.modelOptions = filterModelOptions(u.cfg, u.allModelOptions, u.modelFilter, u.modelCapability)
	if len(u.modelOptions) == 0 {
		u.menuIndex = 0
		return
	}
	if u.menuIndex >= len(u.modelOptions) {
		u.menuIndex = len(u.modelOptions) - 1
	}
	if u.menuIndex < 0 {
		u.menuIndex = 0
	}
}

func (u *onboardingUI) setStep(step onboardStep) {
	u.step = step
	u.input = u.currentInputValue()
	if step == stepLLMModel {
		u.modelFilter = ""
		u.modelCapability = modelCapabilityAll
		u.refreshModelOptions()
		u.menuIndex = selectedModelIndex(u.modelOptions, u.cfg.LLMModel)
		return
	}
	u.menuIndex = 0
}

func (u *onboardingUI) modelCapabilityLabel() string {
	switch u.modelCapability {
	case modelCapabilityVision:
		return "vision"
	case modelCapabilityTools:
		return "tools"
	case modelCapabilityFree:
		return "free"
	default:
		return "all"
	}
}

func nextModelCapabilityFilter(current modelCapabilityFilter) modelCapabilityFilter {
	switch current {
	case modelCapabilityAll:
		return modelCapabilityVision
	case modelCapabilityVision:
		return modelCapabilityTools
	case modelCapabilityTools:
		return modelCapabilityFree
	default:
		return modelCapabilityAll
	}
}

func matchesCapabilityFilter(option llm.ModelOption, capability modelCapabilityFilter) bool {
	switch capability {
	case modelCapabilityVision:
		return option.SupportsImageInput
	case modelCapabilityTools:
		return option.SupportsTools
	case modelCapabilityFree:
		return option.IsFree
	default:
		return true
	}
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
