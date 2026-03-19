package memory

import (
	"context"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/pkg/llm"
)

const (
	ContextSummaryTopic = "conversation_context"
	ContextSummaryKind  = "context_summary"
	contextSoftRatio    = 0.70
	contextHardRatio    = 0.85
	contextReplyBuffer  = 1024
)

type ContextAction string

const (
	ContextActionNone      ContextAction = "none"
	ContextActionSummarize ContextAction = "summarize"
	ContextActionRotate    ContextAction = "rotate"
)

type ContextPolicy struct {
	memory *MemoryManager
}

func NewContextPolicy(memory *MemoryManager) *ContextPolicy {
	return &ContextPolicy{memory: memory}
}

func (p *ContextPolicy) ManageConversation(ctx context.Context, provider, model, conversationID string, memoryWindowSize int) error {
	if p == nil || p.memory == nil {
		return nil
	}

	capacity := llm.ContextWindow(provider, model)
	if capacity <= 0 {
		return nil
	}

	scanLimit := maxInt(memoryWindowSize*3, 24)
	messages, err := p.memory.ListMessages(ctx, conversationID, scanLimit)
	if err != nil {
		return err
	}
	if len(messages) <= 1 {
		return nil
	}

	action := DecideContextAction(capacity, EstimateStoredMessagesTokens(messages)+contextReplyBuffer)
	if action == ContextActionNone {
		return nil
	}

	keepRecent := maxInt(8, memoryWindowSize/2)
	if action == ContextActionRotate {
		keepRecent = 2
	}
	if len(messages) <= keepRecent {
		return nil
	}

	olderMessages := messages[:len(messages)-keepRecent]
	summary := SummarizeMessages(olderMessages)
	if strings.TrimSpace(summary) == "" {
		return nil
	}

	if err := p.memory.AddNote(ctx, Note{
		ConversationID: conversationID,
		Topic:          ContextSummaryTopic,
		Kind:           ContextSummaryKind,
		Summary:        summary,
		Importance:     9,
		Source:         "context_policy",
	}); err != nil {
		return err
	}
	if err := p.memory.AddArchiveEntry(ctx, ArchiveEntry{
		ConversationID: conversationID,
		SessionID:      conversationID,
		Role:           "system",
		Content:        "Context summary: " + summary,
		MessageType:    "summary",
	}); err != nil {
		return err
	}

	return p.memory.TrimMessages(ctx, conversationID, keepRecent)
}

func (p *ContextPolicy) SummaryPrefix(ctx context.Context, conversationID string) (*agent.Message, error) {
	if p == nil || p.memory == nil {
		return nil, nil
	}

	note, ok, err := p.memory.GetLatestNote(ctx, conversationID, ContextSummaryTopic, ContextSummaryKind)
	if err != nil {
		return nil, fmt.Errorf("get latest context summary: %w", err)
	}
	if !ok || strings.TrimSpace(note.Summary) == "" {
		return nil, nil
	}

	return &agent.Message{
		Role:    "assistant",
		Content: "Resumo de contexto anterior para continuidade:\n" + note.Summary,
	}, nil
}

func DecideContextAction(capacity, estimatedTokens int) ContextAction {
	if capacity <= 0 || estimatedTokens <= 0 {
		return ContextActionNone
	}
	ratio := float64(estimatedTokens) / float64(capacity)
	switch {
	case ratio >= contextHardRatio:
		return ContextActionRotate
	case ratio >= contextSoftRatio:
		return ContextActionSummarize
	default:
		return ContextActionNone
	}
}

func EstimateStoredMessagesTokens(messages []Message) int {
	total := 0
	for _, message := range messages {
		total += estimateTextTokens(message.Role)
		total += estimateTextTokens(message.Content)
		total += 12
	}
	return total
}

func SummarizeMessages(messages []Message) string {
	if len(messages) == 0 {
		return ""
	}

	const maxLines = 12
	lines := make([]string, 0, maxLines)
	for _, message := range messages {
		if len(lines) >= maxLines {
			break
		}
		content := compactMessageContent(message.Content, 220)
		if content == "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("- %s: %s", normalizeSummaryRole(message.Role), content))
	}
	if len(lines) == 0 {
		return ""
	}
	if len(messages) > len(lines) {
		lines = append(lines, fmt.Sprintf("- resumo parcial de %d mensagens anteriores preservado para continuidade", len(messages)))
	}
	return strings.Join(lines, "\n")
}

func estimateTextTokens(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}
	return len([]rune(text))/4 + 1
}

func compactMessageContent(content string, maxRunes int) string {
	content = strings.TrimSpace(strings.ReplaceAll(content, "\n", " "))
	content = strings.Join(strings.Fields(content), " ")
	if content == "" {
		return ""
	}

	runes := []rune(content)
	if len(runes) <= maxRunes {
		return content
	}
	return strings.TrimSpace(string(runes[:maxRunes])) + "..."
}

func normalizeSummaryRole(role string) string {
	switch strings.TrimSpace(strings.ToLower(role)) {
	case "assistant":
		return "assistente"
	case "tool":
		return "tool"
	default:
		return "usuario"
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
