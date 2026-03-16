package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/persona"
)

type MemoryCommandService interface {
	DebugLongTermMemory(ctx context.Context, userID, conversationID, query string) (*persona.LongTermMemoryDebugReport, error)
}

type MemoryCommandHandler struct {
	service MemoryCommandService
}

func NewMemoryCommandHandler(service MemoryCommandService) *MemoryCommandHandler {
	return &MemoryCommandHandler{service: service}
}

func (h *MemoryCommandHandler) HandleText(ctx context.Context, userID, conversationID, text string) (string, error) {
	query := strings.TrimSpace(text)
	if strings.HasPrefix(query, "/memory") {
		query = strings.TrimSpace(strings.TrimPrefix(query, "/memory"))
	}
	if strings.HasPrefix(strings.ToLower(query), "debug") {
		query = strings.TrimSpace(query[len("debug"):])
	}

	report, err := h.service.DebugLongTermMemory(ctx, userID, conversationID, query)
	if err != nil {
		return "", err
	}
	return formatMemoryDebugReport(report), nil
}

func formatMemoryDebugReport(report *persona.LongTermMemoryDebugReport) string {
	if report == nil {
		return "Memoria indisponivel."
	}

	var lines []string
	if strings.TrimSpace(report.Query) == "" {
		lines = append(lines, "# Memory Debug", "Consulta: contexto recente")
	} else {
		lines = append(lines, "# Memory Debug", fmt.Sprintf("Consulta: %s", report.Query))
	}

	if len(report.Tokens) > 0 {
		lines = append(lines, fmt.Sprintf("Tokens: %s", strings.Join(report.Tokens, ", ")))
	}

	if len(report.SelectedFacts) == 0 && len(report.SelectedNotes) == 0 {
		lines = append(lines, "Nenhuma memoria relevante foi recuperada.")
		return strings.Join(lines, "\n")
	}

	if len(report.SelectedFacts) > 0 {
		lines = append(lines, "Facts:")
		for _, fact := range report.SelectedFacts {
			lines = append(lines, fmt.Sprintf("- %s=%s (score=%d)", fact.Fact.Key, fact.Fact.Value, fact.Score))
		}
	}

	if len(report.SelectedNotes) > 0 {
		lines = append(lines, "Notes:")
		for _, note := range report.SelectedNotes {
			lines = append(lines, fmt.Sprintf("- %s/%s: %s (score=%d)", note.Note.Topic, note.Note.Kind, note.Note.Summary, note.Score))
		}
	}

	return strings.Join(lines, "\n")
}


