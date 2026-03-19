package llm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
)

func parseChatCompletionResponse(respBody []byte) (*agent.ModelResponse, error) {
	var apiResp chatCompletionResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kimi response: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned by Kimi")
	}

	choice := apiResp.Choices[0].Message
	result := &agent.ModelResponse{
		Content:          choice.Content,
		ReasoningContent: choice.ReasoningContent,
	}

	for _, call := range choice.ToolCalls {
		if call.Type != "function" {
			continue
		}

		var args map[string]any
		_ = json.Unmarshal([]byte(call.Function.Arguments), &args)

		result.ToolCalls = append(result.ToolCalls, agent.ToolCall{
			ID:        call.ID,
			Name:      call.Function.Name,
			Arguments: args,
		})
	}

	return result, nil
}

func extractToolCallsFromContent(content string) ([]agent.ToolCall, string) {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "Calling tools:") {
		return nil, content
	}

	payload := strings.TrimSpace(strings.TrimPrefix(trimmed, "Calling tools:"))
	if payload == "" {
		return nil, content
	}

	braceIndex := strings.Index(payload, "{")
	if braceIndex <= 0 {
		return nil, content
	}

	header := strings.TrimSpace(payload[:braceIndex])
	jsonPart, ok := extractBalancedJSONObject(payload[braceIndex:])
	if !ok {
		return nil, content
	}

	namePart := header
	callID := ""
	if colon := strings.Index(header, ":"); colon >= 0 {
		namePart = strings.TrimSpace(header[:colon])
		callID = strings.TrimSpace(header[colon+1:])
	}
	if namePart == "" {
		return nil, content
	}
	if callID == "" {
		callID = "fallback-call"
	}

	var args map[string]any
	if err := json.Unmarshal([]byte(jsonPart), &args); err != nil {
		return nil, content
	}

	return []agent.ToolCall{{
		ID:        callID,
		Name:      namePart,
		Arguments: args,
	}}, ""
}

func extractBalancedJSONObject(text string) (string, bool) {
	depth := 0
	inString := false
	escaped := false

	for i, r := range text {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' && inString {
			escaped = true
			continue
		}
		if r == '"' {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		switch r {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return text[:i+1], true
			}
		}
	}

	return "", false
}

func malformedToolCallContentError(content string) error {
	preview := strings.TrimSpace(content)
	runes := []rune(preview)
	if len(runes) > 160 {
		preview = strings.TrimSpace(string(runes[:160])) + "..."
	}
	return fmt.Errorf("kimi returned malformed tool-call content without structured tool_calls: %s", preview)
}
