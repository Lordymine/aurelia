package llm

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/agent"
)

// KimiProvider implements LLMProvider using direct HTTP requests to Kimi API.
// This avoids go-openai missing fields like reasoning_content for tool calls.
type KimiProvider struct {
	client *http.Client
	model  string
	apiKey string
}

func NewKimiProvider(apiKey string, modelName string) *KimiProvider {
	return &KimiProvider{
		client: &http.Client{
			Timeout: 8 * time.Minute,
		},
		model:  modelName,
		apiKey: apiKey,
	}
}

func (p *KimiProvider) Close() {}

func (p *KimiProvider) GenerateContent(
	ctx context.Context,
	systemPrompt string,
	history []agent.Message,
	tools []agent.Tool,
) (*agent.ModelResponse, error) {
	reqBody, err := p.buildChatCompletionRequest(systemPrompt, history, tools)
	if err != nil {
		return nil, err
	}

	respBody, err := p.doChatCompletionRequest(ctx, reqBody)
	if err != nil {
		return nil, err
	}

	result, err := parseChatCompletionResponse(respBody)
	if err != nil {
		return nil, err
	}

	if err := p.applyFallbackToolCalls(result); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *KimiProvider) applyFallbackToolCalls(result *agent.ModelResponse) error {
	if len(result.ToolCalls) == 0 && strings.HasPrefix(strings.TrimSpace(result.Content), "Calling tools:") {
		fallbackCalls, cleanedContent := extractToolCallsFromContent(result.Content)
		if len(fallbackCalls) > 0 {
			result.ToolCalls = fallbackCalls
			result.Content = cleanedContent
			return nil
		}
		return malformedToolCallContentError(result.Content)
	}
	return nil
}
