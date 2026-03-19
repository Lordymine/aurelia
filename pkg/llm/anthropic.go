package llm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/kocar/aurelia/internal/agent"
)

const anthropicDefaultMaxTokens = 4096

// AnthropicProvider implements agent.LLMProvider using Anthropic's official Go SDK.
type AnthropicProvider struct {
	client anthropic.Client
	model  string
}

func NewAnthropicProvider(apiKey string, modelName string, opts ...option.RequestOption) *AnthropicProvider {
	requestOptions := []option.RequestOption{option.WithAPIKey(apiKey)}
	requestOptions = append(requestOptions, opts...)

	return &AnthropicProvider{
		client: anthropic.NewClient(requestOptions...),
		model:  modelName,
	}
}

func (p *AnthropicProvider) Close() {}

func (p *AnthropicProvider) GenerateContent(
	ctx context.Context,
	systemPrompt string,
	history []agent.Message,
	tools []agent.Tool,
) (*agent.ModelResponse, error) {
	if err := ensureVisionSupport("anthropic", p.model, history); err != nil {
		return nil, err
	}
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		MaxTokens: anthropicDefaultMaxTokens,
		Messages:  buildAnthropicMessages(history),
	}

	if systemPrompt != "" {
		params.System = []anthropic.TextBlockParam{{Text: systemPrompt}}
	}
	if len(tools) != 0 {
		params.Tools = buildAnthropicTools(tools)
		params.ToolChoice = anthropic.ToolChoiceUnionParam{
			OfAuto: &anthropic.ToolChoiceAutoParam{},
		}
	}

	message, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("generate content error: %w", err)
	}

	return parseAnthropicMessage(message)
}

func buildAnthropicMessages(history []agent.Message) []anthropic.MessageParam {
	params := make([]anthropic.MessageParam, 0, len(history))
	for _, msg := range history {
		converted := convertAnthropicMessage(msg)
		if len(converted.Content) == 0 {
			continue
		}
		params = append(params, converted)
	}
	return params
}

func convertAnthropicMessage(msg agent.Message) anthropic.MessageParam {
	switch msg.Role {
	case "assistant":
		var blocks []anthropic.ContentBlockParamUnion
		if msg.Content != "" {
			blocks = append(blocks, anthropic.NewTextBlock(msg.Content))
		}
		for _, call := range msg.ToolCalls {
			payload, _ := json.Marshal(call.Arguments)
			blocks = append(blocks, anthropic.ContentBlockParamUnion{
				OfToolUse: &anthropic.ToolUseBlockParam{
					ID:    call.ID,
					Name:  call.Name,
					Input: payload,
				},
			})
		}
		return anthropic.NewAssistantMessage(blocks...)
	case "tool":
		return anthropic.NewUserMessage(anthropic.ContentBlockParamUnion{
			OfToolResult: &anthropic.ToolResultBlockParam{
				ToolUseID: msg.ToolCallID,
				Content: []anthropic.ToolResultBlockParamContentUnion{
					{OfText: &anthropic.TextBlockParam{Text: msg.Content}},
				},
			},
		})
	default:
		return anthropic.NewUserMessage(buildAnthropicUserBlocks(msg)...)
	}
}

func buildAnthropicUserBlocks(msg agent.Message) []anthropic.ContentBlockParamUnion {
	if len(msg.Parts) == 0 {
		return []anthropic.ContentBlockParamUnion{anthropic.NewTextBlock(msg.Content)}
	}
	blocks := make([]anthropic.ContentBlockParamUnion, 0, len(msg.Parts))
	for _, part := range msg.Parts {
		switch part.Type {
		case agent.ContentPartImage:
			blocks = append(blocks, anthropic.NewImageBlockBase64(part.MIMEType, base64.StdEncoding.EncodeToString(part.Data)))
		default:
			if part.Text != "" {
				blocks = append(blocks, anthropic.NewTextBlock(part.Text))
			}
		}
	}
	if len(blocks) == 0 {
		blocks = append(blocks, anthropic.NewTextBlock(msg.Content))
	}
	return blocks
}

func buildAnthropicTools(tools []agent.Tool) []anthropic.ToolUnionParam {
	params := make([]anthropic.ToolUnionParam, 0, len(tools))
	for _, tool := range tools {
		params = append(params, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type:       "object",
					Properties: tool.JSONSchema["properties"],
					Required:   stringSlice(tool.JSONSchema["required"]),
				},
				Strict: anthropic.Bool(true),
				Type:   anthropic.ToolTypeCustom,
			},
		})
	}
	return params
}

func parseAnthropicMessage(message *anthropic.Message) (*agent.ModelResponse, error) {
	if message == nil {
		return nil, fmt.Errorf("anthropic returned no message")
	}

	result := &agent.ModelResponse{}
	for _, block := range message.Content {
		switch variant := block.AsAny().(type) {
		case anthropic.TextBlock:
			result.Content += variant.Text
		case anthropic.ThinkingBlock:
			result.ReasoningContent += variant.Thinking
		case anthropic.ToolUseBlock:
			args := make(map[string]any)
			if len(variant.Input) != 0 {
				if err := json.Unmarshal(variant.Input, &args); err != nil {
					return nil, fmt.Errorf("decode anthropic tool input for %s: %w", variant.Name, err)
				}
			}
			result.ToolCalls = append(result.ToolCalls, agent.ToolCall{
				ID:        variant.ID,
				Name:      variant.Name,
				Arguments: args,
			})
		}
	}

	return result, nil
}

func stringSlice(value any) []string {
	raw, ok := value.([]interface{})
	if !ok {
		return nil
	}
	items := make([]string, 0, len(raw))
	for _, item := range raw {
		if str, ok := item.(string); ok {
			items = append(items, str)
		}
	}
	return items
}
