package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"
)

const (
	openRouterChatCompletionsURL = "https://openrouter.ai/api/v1/chat/completions"
	openRouterModelsURL          = "https://openrouter.ai/api/v1/models"
	openRouterReferer            = "https://github.com/Lordymine/aurelia"
	openRouterTitle              = "Aurelia"
)

var openRouterVisionCache = struct {
	mu     sync.RWMutex
	loaded bool
	models map[string]bool
}{
	models: map[string]bool{},
}

func NewOpenRouterProvider(apiKey string, model string) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		Provider:  "openrouter",
		APIKey:    apiKey,
		BaseURL:   openRouterChatCompletionsURL,
		Model:     model,
		UserAgent: "Aurelia/1.0",
		Headers: map[string]string{
			"HTTP-Referer": openRouterReferer,
			"X-Title":      openRouterTitle,
		},
	})
}

func listOpenRouterModels(ctx context.Context, apiKey string, baseURL string, client *http.Client) ([]ModelOption, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("HTTP-Referer", openRouterReferer)
	req.Header.Set("X-Title", openRouterTitle)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openrouter models list returned %s", resp.Status)
	}

	var payload struct {
		Data []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Pricing struct {
				Prompt     string `json:"prompt"`
				Completion string `json:"completion"`
			} `json:"pricing"`
			SupportedParameters []string `json:"supported_parameters"`
			Architecture        struct {
				InputModalities []string `json:"input_modalities"`
			} `json:"architecture"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	models := []ModelOption{
		{ID: "openrouter/auto", Name: "OpenRouter Auto"},
		{ID: "openrouter/free", Name: "OpenRouter Free Router", IsFree: true},
	}
	for _, model := range payload.Data {
		models = append(models, ModelOption{
			ID:                 model.ID,
			Name:               model.Name,
			SupportsImageInput: slices.Contains(model.Architecture.InputModalities, "image"),
			SupportsTools:      supportsToolCalling(model.SupportedParameters),
			IsFree:             isOpenRouterFreeModel(model.ID, model.Pricing.Prompt, model.Pricing.Completion),
		})
	}
	cacheOpenRouterVisionModels(models)
	return models, nil
}

func supportsToolCalling(parameters []string) bool {
	for _, parameter := range parameters {
		if strings.EqualFold(parameter, "tools") {
			return true
		}
	}
	return false
}

func isOpenRouterFreeModel(modelID, promptPrice, completionPrice string) bool {
	if strings.Contains(strings.ToLower(modelID), ":free") {
		return true
	}
	return normalizePrice(promptPrice) == "0" && normalizePrice(completionPrice) == "0"
}

func normalizePrice(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "$")
	value = strings.TrimLeft(value, "0")
	value = strings.TrimPrefix(value, ".")
	if value == "" {
		return "0"
	}
	if strings.Trim(value, "0") == "" {
		return "0"
	}
	return value
}

func cacheOpenRouterVisionModels(models []ModelOption) {
	openRouterVisionCache.mu.Lock()
	defer openRouterVisionCache.mu.Unlock()

	openRouterVisionCache.models = make(map[string]bool, len(models))
	for _, model := range models {
		openRouterVisionCache.models[model.ID] = model.SupportsImageInput
	}
	openRouterVisionCache.loaded = true
}

func openRouterModelSupportsVision(modelID string) (bool, bool) {
	openRouterVisionCache.mu.RLock()
	if openRouterVisionCache.loaded {
		value, ok := openRouterVisionCache.models[modelID]
		openRouterVisionCache.mu.RUnlock()
		return value, ok
	}
	openRouterVisionCache.mu.RUnlock()
	return false, false
}
