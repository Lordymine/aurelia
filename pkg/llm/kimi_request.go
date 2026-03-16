package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kocar/aurelia/internal/agent"
)

func (p *KimiProvider) buildChatCompletionRequest(
	systemPrompt string,
	history []agent.Message,
	tools []agent.Tool,
) (map[string]any, error) {
	messages := append([]chatMessage{{
		Role:    "system",
		Content: systemPrompt,
	}}, buildChatHistory(history)...)

	reqBody := map[string]any{
		"model":    p.model,
		"messages": messages,
	}

	if len(tools) > 0 {
		chatTools, err := buildChatTools(tools)
		if err != nil {
			return nil, err
		}
		reqBody["tools"] = chatTools
	}

	return reqBody, nil
}

func buildChatHistory(history []agent.Message) []chatMessage {
	messages := make([]chatMessage, 0, len(history))
	for _, msg := range history {
		messages = append(messages, mapAgentMessage(msg))
	}
	return messages
}

func mapAgentMessage(msg agent.Message) chatMessage {
	cMsg := chatMessage{
		Role:    mapChatRole(msg.Role),
		Content: msg.Content,
	}

	if msg.Role == "tool" {
		cMsg.ToolCallID = msg.ToolCallID
		cMsg.Name = msg.ToolCallID
		return cMsg
	}

	if msg.Role == "assistant" && len(msg.ToolCalls) > 0 {
		cMsg.ReasoningContent = msg.ReasoningContent
		cMsg.ToolCalls = mapToolCalls(msg.ToolCalls)
	}

	return cMsg
}

func mapChatRole(role string) string {
	switch role {
	case "assistant":
		return "assistant"
	case "tool":
		return "tool"
	default:
		return "user"
	}
}

func mapToolCalls(calls []agent.ToolCall) []chatToolCall {
	toolCalls := make([]chatToolCall, 0, len(calls))
	for _, call := range calls {
		argsRaw, _ := json.Marshal(call.Arguments)
		toolCalls = append(toolCalls, chatToolCall{
			ID:   call.ID,
			Type: "function",
			Function: chatFunctionCall{
				Name:      call.Name,
				Arguments: string(argsRaw),
			},
		})
	}
	return toolCalls
}

func buildChatTools(tools []agent.Tool) ([]chatTool, error) {
	chatTools := make([]chatTool, 0, len(tools))
	for _, t := range tools {
		schemaRaw, err := json.Marshal(t.JSONSchema)
		if err != nil {
			return nil, fmt.Errorf("marshal tool schema %s: %w", t.Name, err)
		}
		if string(schemaRaw) == "null" {
			schemaRaw = []byte(`{"type":"object","properties":{}}`)
		}

		chatTools = append(chatTools, chatTool{
			Type: "function",
			Function: chatToolDef{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  json.RawMessage(schemaRaw),
			},
		})
	}
	return chatTools, nil
}

func (p *KimiProvider) doChatCompletionRequest(ctx context.Context, reqBody map[string]any) ([]byte, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.kimi.com/coding/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("User-Agent", "RooCode/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kimi API error: %d, response: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}


