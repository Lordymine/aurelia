package main

import (
	"bufio"
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
)

func TestRunOnboard_SavesInteractiveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	input := strings.Join([]string{
		"",
		"kimi-key",
		"kimi-k2-thinking",
		"groq-key",
		"telegram-token",
		"101,202",
		"700",
		"33",
		"",
	}, "\n")

	var out bytes.Buffer
	if err := runOnboard(strings.NewReader(input), &out); err != nil {
		t.Fatalf("runOnboard() error = %v", err)
	}

	resolver, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() error = %v", err)
	}
	cfg, err := config.Load(resolver)
	if err != nil {
		t.Fatalf("config.Load() error = %v", err)
	}

	if cfg.LLMProvider != "kimi" {
		t.Fatalf("LLMProvider = %q", cfg.LLMProvider)
	}
	if cfg.STTProvider != "groq" {
		t.Fatalf("STTProvider = %q", cfg.STTProvider)
	}
	if cfg.LLMModel != "kimi-k2-thinking" {
		t.Fatalf("LLMModel = %q", cfg.LLMModel)
	}
	if cfg.TelegramBotToken != "telegram-token" {
		t.Fatalf("TelegramBotToken = %q", cfg.TelegramBotToken)
	}
	if got, want := cfg.TelegramAllowedUserIDs, []int64{101, 202}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("TelegramAllowedUserIDs = %v", got)
	}
	if cfg.KimiAPIKey != "kimi-key" {
		t.Fatalf("KimiAPIKey = %q", cfg.KimiAPIKey)
	}
	if cfg.GroqAPIKey != "groq-key" {
		t.Fatalf("GroqAPIKey = %q", cfg.GroqAPIKey)
	}
	if cfg.MaxIterations != 700 {
		t.Fatalf("MaxIterations = %d", cfg.MaxIterations)
	}
	if cfg.MemoryWindowSize != 33 {
		t.Fatalf("MemoryWindowSize = %d", cfg.MemoryWindowSize)
	}
	if cfg.DBPath != filepath.Join(tmpDir, "data", "aurelia.db") {
		t.Fatalf("DBPath = %q", cfg.DBPath)
	}
}

func TestRunOnboard_PreservesExistingValuesOnBlankInput(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	resolver, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() error = %v", err)
	}
	if err := runtime.Bootstrap(resolver); err != nil {
		t.Fatalf("runtime.Bootstrap() error = %v", err)
	}
	if err := config.SaveEditable(resolver, config.EditableConfig{
		LLMProvider:            "kimi",
		LLMModel:               "moonshot-v1-32k",
		STTProvider:            "groq",
		TelegramBotToken:       "old-telegram",
		TelegramAllowedUserIDs: []int64{42},
		AnthropicAPIKey:        "old-anthropic",
		GoogleAPIKey:           "old-google",
		KiloAPIKey:             "old-kilo",
		KimiAPIKey:             "old-kimi",
		OpenRouterAPIKey:       "old-openrouter",
		ZAIAPIKey:              "old-zai",
		AlibabaAPIKey:          "old-alibaba",
		OpenAIAPIKey:           "old-openai",
		GroqAPIKey:             "old-groq",
		MaxIterations:          600,
		MemoryWindowSize:       21,
	}); err != nil {
		t.Fatalf("config.SaveEditable() error = %v", err)
	}

	var out bytes.Buffer
	if err := runOnboard(strings.NewReader("\n\n\n\n\n\n\n\n"), &out); err != nil {
		t.Fatalf("runOnboard() error = %v", err)
	}

	cfg, err := config.Load(resolver)
	if err != nil {
		t.Fatalf("config.Load() error = %v", err)
	}
	if cfg.TelegramBotToken != "old-telegram" || cfg.KimiAPIKey != "old-kimi" || cfg.AnthropicAPIKey != "old-anthropic" || cfg.GoogleAPIKey != "old-google" || cfg.KiloAPIKey != "old-kilo" || cfg.OpenRouterAPIKey != "old-openrouter" || cfg.ZAIAPIKey != "old-zai" || cfg.AlibabaAPIKey != "old-alibaba" || cfg.OpenAIAPIKey != "old-openai" || cfg.GroqAPIKey != "old-groq" {
		t.Fatalf("expected secrets to be preserved, got %+v", cfg)
	}
	if cfg.LLMProvider != "kimi" || cfg.LLMModel != "moonshot-v1-32k" || cfg.STTProvider != "groq" {
		t.Fatalf("expected providers to be preserved, got llm=%q model=%q stt=%q", cfg.LLMProvider, cfg.LLMModel, cfg.STTProvider)
	}
	if len(cfg.TelegramAllowedUserIDs) != 1 || cfg.TelegramAllowedUserIDs[0] != 42 {
		t.Fatalf("expected allowed user IDs to be preserved, got %v", cfg.TelegramAllowedUserIDs)
	}
	if cfg.MaxIterations != 600 || cfg.MemoryWindowSize != 21 {
		t.Fatalf("expected numeric fields to be preserved, got max=%d memory=%d", cfg.MaxIterations, cfg.MemoryWindowSize)
	}
}

func TestParseInt64List_RejectsInvalidInput(t *testing.T) {
	if _, err := parseInt64List("123,abc"); err == nil {
		t.Fatal("expected parseInt64List() to fail on invalid input")
	}
}

func TestRenderOnboardingHeader_IncludesBannerAndColor(t *testing.T) {
	header := renderOnboardingHeader()

	if !strings.Contains(header, "$$$$$$\\") {
		t.Fatal("expected ASCII banner in onboarding header")
	}
	if !strings.Contains(header, colorBlue) || !strings.Contains(header, colorReset) {
		t.Fatal("expected ANSI blue color codes in onboarding header")
	}
	if !strings.Contains(header, "Local onboarding for runtime config") {
		t.Fatal("expected onboarding subtitle in header")
	}
}

func TestRawTerminalFrame_UsesCRLFLineEndings(t *testing.T) {
	frame := rawTerminalFrame("line1\nline2\nline3")

	if strings.Contains(frame, "line1\nline2") {
		t.Fatal("expected line feeds to be normalized for raw terminal output")
	}
	if want := "line1\r\nline2\r\nline3"; frame != want {
		t.Fatalf("frame = %q, want %q", frame, want)
	}
}

func TestRawTerminalFrame_DoesNotDuplicateExistingCRLF(t *testing.T) {
	frame := rawTerminalFrame("line1\r\nline2\nline3\r\n")

	if strings.Contains(frame, "\r\r\n") {
		t.Fatal("expected CRLF normalization to avoid duplicated carriage returns")
	}
	if want := "line1\r\nline2\r\nline3\r\n"; frame != want {
		t.Fatalf("frame = %q, want %q", frame, want)
	}
}

func TestOnboardingUI_MenuFlowAndBack(t *testing.T) {
	ui := newOnboardingUI(config.DefaultEditableConfig())

	_, _, err := ui.HandleKey(keyEvent{code: keyEnter})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.step != stepLLMKey {
		t.Fatalf("step = %v, want %v", ui.step, stepLLMKey)
	}
	ui.input = "kimi-key"
	_, _, err = ui.HandleKey(keyEvent{code: keyEnter})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.step != stepLLMModel {
		t.Fatalf("step = %v, want %v", ui.step, stepLLMModel)
	}
	_, _, err = ui.HandleKey(keyEvent{code: keyLeft})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.step != stepLLMKey {
		t.Fatalf("step = %v, want %v after back", ui.step, stepLLMKey)
	}
}

func TestOnboardingUI_ModelSelectionPersistsChoice(t *testing.T) {
	ui := newOnboardingUI(config.DefaultEditableConfig())
	ui.step = stepLLMModel
	ui.modelOptions = []ModelOption{
		{ID: "kimi-k2-thinking", Name: "Kimi K2 Thinking"},
		{ID: "moonshot-v1-32k", Name: "Moonshot v1 32K"},
	}
	ui.menuIndex = 1

	_, _, err := ui.HandleKey(keyEvent{code: keyEnter})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.cfg.LLMModel != "moonshot-v1-32k" {
		t.Fatalf("LLMModel = %q", ui.cfg.LLMModel)
	}
	if ui.step != stepSTTProvider {
		t.Fatalf("step = %v, want %v", ui.step, stepSTTProvider)
	}
}

func TestOnboardingUI_AnthropicKeyInputTargetsAnthropicSecret(t *testing.T) {
	ui := newOnboardingUI(config.EditableConfig{
		LLMProvider: "anthropic",
		LLMModel:    "claude-sonnet-4-6",
	})
	ui.step = stepLLMKey
	ui.input = "anthropic-key"

	_, _, err := ui.HandleKey(keyEvent{code: keyEnter})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.cfg.AnthropicAPIKey != "anthropic-key" {
		t.Fatalf("AnthropicAPIKey = %q", ui.cfg.AnthropicAPIKey)
	}
	if ui.cfg.KimiAPIKey != "" {
		t.Fatalf("KimiAPIKey = %q", ui.cfg.KimiAPIKey)
	}
}

func TestOnboardingUI_KiloKeyInputTargetsKiloSecret(t *testing.T) {
	ui := newOnboardingUI(config.EditableConfig{
		LLMProvider: "kilo",
		LLMModel:    "gpt-5.4",
	})
	ui.step = stepLLMKey
	ui.input = "kilo-key"

	_, _, err := ui.HandleKey(keyEvent{code: keyEnter})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.cfg.KiloAPIKey != "kilo-key" {
		t.Fatalf("KiloAPIKey = %q", ui.cfg.KiloAPIKey)
	}
	if ui.cfg.KimiAPIKey != "" {
		t.Fatalf("KimiAPIKey = %q", ui.cfg.KimiAPIKey)
	}
}

func TestLLMProviderChoicesMatchLabels(t *testing.T) {
	choices := llmProviderChoices()
	labels := llmProviderLabels()

	if len(choices) != len(labels) {
		t.Fatalf("choices=%d labels=%d", len(choices), len(labels))
	}
	if len(choices) == 0 {
		t.Fatal("expected provider choices")
	}
}

func TestOnboardingUI_OpenAICodexSkipsAPIKeyStep(t *testing.T) {
	ui := newOnboardingUI(config.EditableConfig{
		LLMProvider:      "openai",
		LLMModel:         "gpt-5.4",
		OpenAIAuthMode:   "api_key",
		MemoryWindowSize: 20,
		MaxIterations:    500,
		STTProvider:      "groq",
	})

	ui.step = stepOpenAIAuthMode
	ui.menuIndex = 1

	_, _, err := ui.HandleKey(keyEvent{code: keyEnter})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.cfg.OpenAIAuthMode != "codex" {
		t.Fatalf("OpenAIAuthMode = %q", ui.cfg.OpenAIAuthMode)
	}
	if ui.step != stepOpenAICodexLogin {
		t.Fatalf("step = %v, want %v", ui.step, stepOpenAICodexLogin)
	}
}

func TestFilterModelOptions_OpenRouterMatchesProviderAndModel(t *testing.T) {
	options := []ModelOption{
		{ID: "openrouter/auto", Name: "OpenRouter Auto"},
		{ID: "anthropic/claude-sonnet-4", Name: "Claude Sonnet 4"},
		{ID: "google/gemini-2.5-flash", Name: "Gemini 2.5 Flash"},
	}

	cfg := config.EditableConfig{LLMProvider: "openrouter"}

	filteredByProvider := filterModelOptions(cfg, options, "anthropic", modelCapabilityAll)
	if len(filteredByProvider) != 1 || filteredByProvider[0].ID != "anthropic/claude-sonnet-4" {
		t.Fatalf("filteredByProvider = %+v", filteredByProvider)
	}

	filteredByModel := filterModelOptions(cfg, options, "gemini", modelCapabilityAll)
	if len(filteredByModel) != 1 || filteredByModel[0].ID != "google/gemini-2.5-flash" {
		t.Fatalf("filteredByModel = %+v", filteredByModel)
	}
}

func TestOnboardingUI_OpenRouterModelSearchFiltersResults(t *testing.T) {
	ui := newOnboardingUI(config.EditableConfig{
		LLMProvider: "openrouter",
		LLMModel:    "openrouter/auto",
	})
	ui.step = stepLLMModel
	ui.allModelOptions = []ModelOption{
		{ID: "openrouter/auto", Name: "OpenRouter Auto"},
		{ID: "anthropic/claude-sonnet-4", Name: "Claude Sonnet 4"},
		{ID: "google/gemini-2.5-flash", Name: "Gemini 2.5 Flash"},
	}
	ui.modelOptions = append([]ModelOption(nil), ui.allModelOptions...)

	_, _, err := ui.HandleKey(keyEvent{code: keyRune, r: 'a'})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	_, _, err = ui.HandleKey(keyEvent{code: keyRune, r: 'n'})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.modelFilter != "an" {
		t.Fatalf("modelFilter = %q", ui.modelFilter)
	}
	if len(ui.modelOptions) != 1 || ui.modelOptions[0].ID != "anthropic/claude-sonnet-4" {
		t.Fatalf("modelOptions = %+v", ui.modelOptions)
	}
}

func TestFilterModelOptions_KiloMatchesProviderAndModel(t *testing.T) {
	options := []ModelOption{
		{ID: "gpt-5.4", Name: "GPT-5.4 · openai"},
		{ID: "claude-sonnet-4-6", Name: "Claude Sonnet 4.6 · anthropic"},
		{ID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro · google"},
	}

	cfg := config.EditableConfig{LLMProvider: "kilo"}

	filteredByProvider := filterModelOptions(cfg, options, "anthropic", modelCapabilityAll)
	if len(filteredByProvider) != 1 || filteredByProvider[0].ID != "claude-sonnet-4-6" {
		t.Fatalf("filteredByProvider = %+v", filteredByProvider)
	}

	filteredByModel := filterModelOptions(cfg, options, "gpt-5", modelCapabilityAll)
	if len(filteredByModel) != 1 || filteredByModel[0].ID != "gpt-5.4" {
		t.Fatalf("filteredByModel = %+v", filteredByModel)
	}
}

func TestOnboardingUI_KiloModelSearchFiltersResults(t *testing.T) {
	ui := newOnboardingUI(config.EditableConfig{
		LLMProvider: "kilo",
		LLMModel:    "gpt-5.4",
	})
	ui.step = stepLLMModel
	ui.allModelOptions = []ModelOption{
		{ID: "gpt-5.4", Name: "GPT-5.4 · openai"},
		{ID: "claude-sonnet-4-6", Name: "Claude Sonnet 4.6 · anthropic"},
		{ID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro · google"},
	}
	ui.modelOptions = append([]ModelOption(nil), ui.allModelOptions...)

	_, _, err := ui.HandleKey(keyEvent{code: keyRune, r: 'g'})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	_, _, err = ui.HandleKey(keyEvent{code: keyRune, r: 'o'})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.modelFilter != "go" {
		t.Fatalf("modelFilter = %q", ui.modelFilter)
	}
	if len(ui.modelOptions) != 1 || ui.modelOptions[0].ID != "gemini-2.5-pro" {
		t.Fatalf("modelOptions = %+v", ui.modelOptions)
	}
}

func TestFilterModelOptions_VisionOnly(t *testing.T) {
	options := []ModelOption{
		{ID: "kimi-k2-thinking", Name: "Kimi K2 Thinking"},
		{ID: "moonshot-v1-vision", Name: "Moonshot Vision", SupportsImageInput: true},
	}

	filtered := filterModelOptions(config.EditableConfig{LLMProvider: "kimi"}, options, "", modelCapabilityVision)
	if len(filtered) != 1 || filtered[0].ID != "moonshot-v1-vision" {
		t.Fatalf("filtered = %+v", filtered)
	}
}

func TestOnboardingUI_ModelVisionToggleFiltersResults(t *testing.T) {
	ui := newOnboardingUI(config.EditableConfig{
		LLMProvider: "kilo",
		LLMModel:    "openai/gpt-5.4",
	})
	ui.step = stepLLMModel
	ui.allModelOptions = []ModelOption{
		{ID: "zai/glm-5-turbo", Name: "GLM-5 Turbo"},
		{ID: "openai/gpt-5.4", Name: "GPT-5.4", SupportsImageInput: true},
	}
	ui.modelOptions = append([]ModelOption(nil), ui.allModelOptions...)

	_, _, err := ui.HandleKey(keyEvent{code: keyRight})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.modelCapability != modelCapabilityVision {
		t.Fatalf("expected vision filter, got %v", ui.modelCapability)
	}
	if len(ui.modelOptions) != 1 || ui.modelOptions[0].ID != "openai/gpt-5.4" {
		t.Fatalf("modelOptions = %+v", ui.modelOptions)
	}
}

func TestOnboardingUI_ModelCapabilityCycleFiltersToolsAndFree(t *testing.T) {
	ui := newOnboardingUI(config.EditableConfig{
		LLMProvider: "openrouter",
		LLMModel:    "openrouter/auto",
	})
	ui.step = stepLLMModel
	ui.allModelOptions = []ModelOption{
		{ID: "anthropic/claude-sonnet-4", Name: "Claude Sonnet 4", SupportsImageInput: true, SupportsTools: true},
		{ID: "meta-llama/llama-free", Name: "Llama Free", IsFree: true},
	}
	ui.modelOptions = append([]ModelOption(nil), ui.allModelOptions...)

	_, _, err := ui.HandleKey(keyEvent{code: keyRight})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	_, _, err = ui.HandleKey(keyEvent{code: keyRight})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.modelCapability != modelCapabilityTools {
		t.Fatalf("expected tools filter, got %v", ui.modelCapability)
	}
	if len(ui.modelOptions) != 1 || ui.modelOptions[0].ID != "anthropic/claude-sonnet-4" {
		t.Fatalf("tools modelOptions = %+v", ui.modelOptions)
	}

	_, _, err = ui.HandleKey(keyEvent{code: keyRight})
	if err != nil {
		t.Fatalf("HandleKey() error = %v", err)
	}
	if ui.modelCapability != modelCapabilityFree {
		t.Fatalf("expected free filter, got %v", ui.modelCapability)
	}
	if len(ui.modelOptions) != 1 || ui.modelOptions[0].ID != "meta-llama/llama-free" {
		t.Fatalf("free modelOptions = %+v", ui.modelOptions)
	}
}

func TestReadKey_TreatsQAsInputRune(t *testing.T) {
	ev, err := readKey(bufio.NewReader(strings.NewReader("q")))
	if err != nil {
		t.Fatalf("readKey() error = %v", err)
	}
	if ev.code != keyRune || ev.r != 'q' {
		t.Fatalf("expected q to be treated as input rune, got %+v", ev)
	}
}
