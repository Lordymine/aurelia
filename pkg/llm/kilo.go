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
	kiloChatCompletionsURL = "https://api.kilo.ai/api/gateway/chat/completions"
	kiloModelsURL          = "https://api.kilo.ai/api/gateway/models"
)

var kiloModelsCatalogURL = kiloModelsURL

var kiloVisionCache = struct {
	mu     sync.RWMutex
	loaded bool
	models map[string]bool
}{
	models: map[string]bool{},
}

func NewKiloProvider(apiKey string, model string) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		Provider:  "kilo",
		APIKey:    apiKey,
		BaseURL:   kiloChatCompletionsURL,
		Model:     model,
		UserAgent: "Aurelia/1.0",
	})
}

func listKiloModels(ctx context.Context, apiKey string, baseURL string, client *http.Client) ([]ModelOption, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kilo models list returned %s", resp.Status)
	}

	var payload struct {
		Data []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			OwnedBy      string `json:"owned_by"`
			Architecture struct {
				InputModalities []string `json:"input_modalities"`
			} `json:"architecture"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	models := make([]ModelOption, 0, len(payload.Data))
	for _, model := range payload.Data {
		if model.ID == "" {
			continue
		}
		name := model.Name
		if name == "" {
			name = model.ID
		}
		if model.OwnedBy != "" && !containsFold(name, model.OwnedBy) {
			name = fmt.Sprintf("%s · %s", name, model.OwnedBy)
		}
		models = append(models, ModelOption{
			ID:                 model.ID,
			Name:               name,
			SupportsImageInput: slices.Contains(model.Architecture.InputModalities, "image"),
			IsFree:             strings.Contains(strings.ToLower(model.ID), ":free"),
		})
	}
	cacheKiloVisionModels(models)
	return models, nil
}

func kiloModelsURLForTest(url string) func() {
	original := kiloModelsCatalogURL
	kiloModelsCatalogURL = url
	return func() {
		kiloModelsCatalogURL = original
	}
}

func containsFold(value string, needle string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}

func cacheKiloVisionModels(models []ModelOption) {
	kiloVisionCache.mu.Lock()
	defer kiloVisionCache.mu.Unlock()

	kiloVisionCache.models = make(map[string]bool, len(models))
	for _, model := range models {
		kiloVisionCache.models[model.ID] = model.SupportsImageInput
	}
	kiloVisionCache.loaded = true
}

func kiloModelSupportsVision(modelID string) (bool, bool) {
	kiloVisionCache.mu.RLock()
	if kiloVisionCache.loaded {
		value, ok := kiloVisionCache.models[modelID]
		kiloVisionCache.mu.RUnlock()
		return value, ok
	}
	kiloVisionCache.mu.RUnlock()

	models, err := listKiloModels(context.Background(), "", kiloModelsCatalogURL, http.DefaultClient)
	if err != nil {
		return false, false
	}
	for _, model := range models {
		if model.ID == modelID {
			return model.SupportsImageInput, true
		}
	}
	return false, false
}
