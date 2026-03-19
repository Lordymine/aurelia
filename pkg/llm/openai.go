package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	openAIChatCompletionsURL = "https://api.openai.com/v1/chat/completions"
	openAIModelsURL          = "https://api.openai.com/v1/models"
)

func NewOpenAIProvider(apiKey string, model string) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		Provider:  "openai",
		APIKey:    apiKey,
		BaseURL:   openAIChatCompletionsURL,
		Model:     model,
		UserAgent: "Aurelia/1.0",
	})
}

func listOpenAIModels(ctx context.Context, apiKey string, baseURL string, client *http.Client) ([]ModelOption, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai models list returned %s", resp.Status)
	}

	var payload struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	var models []ModelOption
	for _, model := range payload.Data {
		if !looksLikeOpenAIChatModel(model.ID) {
			continue
		}
		models = append(models, ModelOption{
			ID:                 model.ID,
			Name:               model.ID,
			SupportsImageInput: looksLikeOpenAIVisionModel(model.ID),
		})
	}
	return models, nil
}

func looksLikeOpenAIChatModel(modelID string) bool {
	prefixes := []string{"gpt-", "o", "chatgpt-", "codex-", "gpt-oss-"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(modelID, prefix) {
			return true
		}
	}
	return false
}

func looksLikeOpenAIVisionModel(modelID string) bool {
	modelID = strings.ToLower(modelID)
	if strings.Contains(modelID, "audio") || strings.Contains(modelID, "codex") {
		return false
	}
	prefixes := []string{"gpt-", "o", "chatgpt-"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(modelID, prefix) {
			return true
		}
	}
	return false
}
