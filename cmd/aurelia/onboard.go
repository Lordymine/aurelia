package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

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
	stepAnthropicAuthMode
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
	allModelOptions []ModelOption
	modelOptions    []ModelOption
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

var llmModelCatalog = listModels

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
	if config.NormalizeProvider(current.LLMProvider) == "anthropic" {
		current.AnthropicAuthMode, _ = promptChoice(reader, stdout, "Anthropic auth mode", current.AnthropicAuthMode, []string{"api_key", "subscription"})
	}
	if err := writef(stdout, "STT provider [%s]: %s\n\n", current.STTProvider, "Groq"); err != nil {
		return err
	}

	if usesAnthropicSubscription(*current) {
		if err := writef(stdout, "Anthropic auth mode: subscription (no API key needed).\n\n"); err != nil {
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
			_ = action
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
	if cfg.AnthropicAuthMode == "" {
		cfg.AnthropicAuthMode = "api_key"
	}
	if cfg.STTProvider == "" {
		cfg.STTProvider = "groq"
	}
	modelOptions, modelSource := resolveModelOptions(cfg)
	return &onboardingUI{
		cfg:             cfg,
		allModelOptions: append([]ModelOption(nil), modelOptions...),
		modelOptions:    append([]ModelOption(nil), modelOptions...),
		modelSource:     modelSource,
		step:            stepLLMProvider,
		reviewOptions:   []string{"Save config", "Back", "Cancel"},
	}
}
