package main

import (
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/pkg/llm"
)

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
		b.WriteString("\nUse arrows and Enter.\n")
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
		b.WriteString("\nUse arrows and Enter. Use left to go back.\n")
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
		b.WriteString("\nUse arrows and Enter. Use left to go back. Press Ctrl+C to cancel.\n")
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
	b.WriteString("\n\nType and press Enter. Use left to go back. Press Ctrl+C to cancel.\n")
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
		value, err := parsePositiveInt(strings.TrimSpace(u.input), "max iterations")
		if err != nil {
			return err
		}
		u.cfg.MaxIterations = value
	case stepRuntimeMemoryWindow:
		value, err := parsePositiveInt(strings.TrimSpace(u.input), "memory window size")
		if err != nil {
			return err
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
		return fmt.Sprintf("%d", u.cfg.MaxIterations)
	case stepRuntimeMemoryWindow:
		return fmt.Sprintf("%d", u.cfg.MemoryWindowSize)
	default:
		return ""
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

func renderModelMenu(options []llm.ModelOption, selected int) string {
	if len(options) == 0 {
		return "  No models available.\n"
	}

	labels := make([]string, 0, len(options))
	for _, option := range options {
		labels = append(labels, option.Label())
	}
	return renderMenu(labels, selected)
}

func usesProviderModelSearch(cfg config.EditableConfig) bool {
	switch cfg.LLMProvider {
	case "openrouter", "kilo":
		return true
	default:
		return false
	}
}
