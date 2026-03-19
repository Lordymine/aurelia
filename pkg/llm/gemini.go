package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"

	"github.com/kocar/aurelia/internal/agent"
)

// GeminiProvider implements agent.LLMProvider using Google's generative AI SDK
type GeminiProvider struct {
	client *genai.Client
	model  *genai.GenerativeModel
	name   string
}

// NewGeminiProvider creates a new provider
func NewGeminiProvider(ctx context.Context, apiKey string, modelName string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	model := client.GenerativeModel(modelName)
	return &GeminiProvider{client: client, model: model, name: modelName}, nil
}

func NewGeminiProviderWithTokenSource(ctx context.Context, tokenSource oauth2.TokenSource, modelName string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	model := client.GenerativeModel(modelName)
	return &GeminiProvider{client: client, model: model, name: modelName}, nil
}

// Close cleans up
func (p *GeminiProvider) Close() {
	_ = p.client.Close()
}

// GenerateContent maps our internal representation to Gemini's
func (p *GeminiProvider) GenerateContent(ctx context.Context, systemPrompt string, history []agent.Message, tools []agent.Tool) (*agent.ModelResponse, error) {
	if err := ensureVisionSupport("google", p.name, history); err != nil {
		return nil, err
	}
	// Set system instruction
	p.model.SystemInstruction = genai.NewUserContent(genai.Text(systemPrompt))

	// Register tools if any
	if len(tools) > 0 {
		var geminiTools []*genai.Tool
		geminiFuncs := make([]*genai.FunctionDeclaration, 0, len(tools))
		for _, t := range tools {
			// In a real app we'd convert the JSONSchema to genai.Schema
			// For simplicity we create a blank schema or mock it here
			// Full struct parsing requires mapping `type`, `properties`, etc.
			geminiFuncs = append(geminiFuncs, &genai.FunctionDeclaration{
				Name:        t.Name,
				Description: t.Description,
				// Schema should be mapped from t.JSONSchema
			})
		}
		geminiTools = append(geminiTools, &genai.Tool{FunctionDeclarations: geminiFuncs})
		p.model.Tools = geminiTools
	} else {
		p.model.Tools = nil // clear tools if none
	}

	var contents []*genai.Content
	var lastPart genai.Part
	chat := p.model.StartChat()

	for i, msg := range history {
		// Map our Role to theirs
		role := "user"
		if msg.Role == "assistant" || msg.Role == "tool" {
			role = "model"
		}

		var parts []genai.Part
		if msg.Role == "tool" {
			parts = []genai.Part{genai.FunctionResponse{
				Name:     msg.ToolCallID,
				Response: map[string]any{"result": msg.Content},
			}}
		} else {
			parts = buildGeminiParts(msg)
		}

		if i == len(history)-1 && role == "user" {
			if len(parts) > 0 {
				lastPart = parts[0]
				if len(parts) > 1 {
					resp, err := chat.SendMessage(ctx, parts...)
					if err != nil {
						return nil, fmt.Errorf("generate content error: %w", err)
					}
					return parseGeminiResponse(resp)
				}
			}
			break // Don't append to history, send it!
		} else if i == len(history)-1 && role == "model" {
			// This is strange but we can just append and send empty, or send it.
			// ReAct usually ends on user or tool. If tool (mapped to model? actually tool is user role in Gemini usually, wait. Gemini tool response is "function" or "user"? The docs say FunctionResponse is "user" role or "function" role, but in genai it's part of Content.
			// Let's just treat the last one as the send payload
			if len(parts) > 0 {
				lastPart = parts[0]
			}
			break
		}

		contents = append(contents, &genai.Content{
			Role:  role,
			Parts: parts,
		})
	}

	chat.History = contents

	if lastPart == nil {
		lastPart = genai.Text("")
	}
	resp, err := chat.SendMessage(ctx, lastPart)
	if err != nil {
		return nil, fmt.Errorf("generate content error: %w", err)
	}

	return parseGeminiResponse(resp)
}

func buildGeminiParts(msg agent.Message) []genai.Part {
	if len(msg.Parts) == 0 {
		return []genai.Part{genai.Text(msg.Content)}
	}
	parts := make([]genai.Part, 0, len(msg.Parts))
	for _, part := range msg.Parts {
		switch part.Type {
		case agent.ContentPartImage:
			format := strings.TrimPrefix(part.MIMEType, "image/")
			parts = append(parts, genai.ImageData(format, part.Data))
		default:
			parts = append(parts, genai.Text(part.Text))
		}
	}
	if len(parts) == 0 {
		parts = append(parts, genai.Text(msg.Content))
	}
	return parts
}

func parseGeminiResponse(resp *genai.GenerateContentResponse) (*agent.ModelResponse, error) {

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates returned")
	}

	candidate := resp.Candidates[0]
	var result agent.ModelResponse

	// Extract parts
	if candidate.Content != nil {
		for _, part := range candidate.Content.Parts {
			switch p := part.(type) {
			case genai.Text:
				result.Content += string(p)
			case genai.FunctionCall:
				args := make(map[string]interface{})
				for k, v := range p.Args {
					args[k] = v
				}
				result.ToolCalls = append(result.ToolCalls, agent.ToolCall{
					ID:        p.Name, // Using name as ID for simplicity
					Name:      p.Name,
					Arguments: args,
				})
			}
		}
	}

	return &result, nil
}
