package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/pkg/llm"
)

const (
	contextSummaryTopic = "conversation_context"
	contextSummaryKind  = "context_summary"
	contextSoftRatio    = 0.70
	contextHardRatio    = 0.85
	contextReplyBuffer  = 1024
)

type contextAction string

const (
	contextActionNone      contextAction = "none"
	contextActionSummarize contextAction = "summarize"
	contextActionRotate    contextAction = "rotate"
)

func (bc *BotController) manageConversationContext(ctx context.Context, conversationID string) error {
	capacity := llm.ContextWindow(bc.config.LLMProvider, bc.config.LLMModel)
	if capacity <= 0 {
		return nil
	}

	scanLimit := maxInt(bc.config.MemoryWindowSize*3, 24)
	messages, err := bc.memory.ListMessages(ctx, conversationID, scanLimit)
	if err != nil {
		return err
	}
	if len(messages) <= 1 {
		return nil
	}

	action := decideContextAction(capacity, estimateStoredMessagesTokens(messages)+contextReplyBuffer)
	if action == contextActionNone {
		return nil
	}

	keepRecent := maxInt(8, bc.config.MemoryWindowSize/2)
	if action == contextActionRotate {
		keepRecent = 2
	}
	if len(messages) <= keepRecent {
		return nil
	}

	olderMessages := messages[:len(messages)-keepRecent]
	summary := summarizeMessages(olderMessages)
	if strings.TrimSpace(summary) == "" {
		return nil
	}

	if err := bc.memory.AddNote(ctx, memory.Note{
		ConversationID: conversationID,
		Topic:          contextSummaryTopic,
		Kind:           contextSummaryKind,
		Summary:        summary,
		Importance:     9,
		Source:         "context_policy",
	}); err != nil {
		return err
	}
	if err := bc.memory.AddArchiveEntry(ctx, memory.ArchiveEntry{
		ConversationID: conversationID,
		SessionID:      conversationID,
		Role:           "system",
		Content:        "Context summary: " + summary,
		MessageType:    "summary",
	}); err != nil {
		return err
	}

	return bc.memory.TrimMessages(ctx, conversationID, keepRecent)
}

func (bc *BotController) summaryPrefix(ctx context.Context, conversationID string) *agent.Message {
	note, ok, err := bc.memory.GetLatestNote(ctx, conversationID, contextSummaryTopic, contextSummaryKind)
	if err != nil || !ok || strings.TrimSpace(note.Summary) == "" {
		return nil
	}
	return &agent.Message{
		Role:    "assistant",
		Content: "Resumo de contexto anterior para continuidade:\n" + note.Summary,
	}
}

func decideContextAction(capacity, estimatedTokens int) contextAction {
	if capacity <= 0 || estimatedTokens <= 0 {
		return contextActionNone
	}
	ratio := float64(estimatedTokens) / float64(capacity)
	switch {
	case ratio >= contextHardRatio:
		return contextActionRotate
	case ratio >= contextSoftRatio:
		return contextActionSummarize
	default:
		return contextActionNone
	}
}

func estimateStoredMessagesTokens(messages []memory.Message) int {
	total := 0
	for _, message := range messages {
		total += estimateTextTokens(message.Role)
		total += estimateTextTokens(message.Content)
		total += 12
	}
	return total
}

func estimateTextTokens(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}
	return len([]rune(text))/4 + 1
}

func summarizeMessages(messages []memory.Message) string {
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
