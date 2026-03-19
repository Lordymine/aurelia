package llm

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
)

type VisionUnsupportedError struct {
	provider string
	model    string
}

func (e VisionUnsupportedError) Error() string {
	if e.model == "" {
		return fmt.Sprintf("o provider %s nao suporta imagens neste runtime atual", e.provider)
	}
	return fmt.Sprintf("o modelo %s do provider %s nao suporta entrada de imagem nesta configuracao atual", e.model, e.provider)
}

func ensureVisionSupport(provider, model string, history []agent.Message) error {
	if !historyHasImages(history) {
		return nil
	}
	if modelSupportsVision(provider, model) {
		return nil
	}
	return VisionUnsupportedError{provider: provider, model: model}
}

func historyHasImages(history []agent.Message) bool {
	for _, msg := range history {
		for _, part := range msg.Parts {
			if part.Type == agent.ContentPartImage && len(part.Data) != 0 {
				return true
			}
		}
	}
	return false
}

func modelSupportsVision(provider, model string) bool {
	model = strings.ToLower(model)
	if model == "" {
		return false
	}
	blockedTokens := []string{"embedding", "moderation", "tts", "transcription", "whisper", "rerank"}
	for _, token := range blockedTokens {
		if strings.Contains(model, token) {
			return false
		}
	}

	switch provider {
	case "kilo":
		if supported, ok := kiloModelSupportsVision(model); ok {
			return supported
		}
		return false
	case "openrouter":
		if supported, ok := openRouterModelSupportsVision(model); ok {
			return supported
		}
		return heuristicVisionSupport(model)
	case "openai":
		return looksLikeOpenAIVisionModel(model)
	case "zai":
		return strings.Contains(model, "4.6v") || strings.Contains(model, "vision")
	case "alibaba":
		return strings.Contains(model, "vl")
	case "kimi":
		return strings.Contains(model, "vision")
	case "anthropic", "google":
		return true
	default:
		return false
	}
}

func openAIImageURL(part agent.ContentPart) string {
	return fmt.Sprintf("data:%s;base64,%s", part.MIMEType, base64.StdEncoding.EncodeToString(part.Data))
}

func heuristicVisionSupport(model string) bool {
	positiveTokens := []string{"vision", "vl", "4.6v", "gpt-", "gemini", "claude", "image"}
	for _, token := range positiveTokens {
		if strings.Contains(model, token) {
			return true
		}
	}
	return false
}
