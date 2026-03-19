package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/observability"
)

type OpsCommandService interface {
	ListRecentOperations(ctx context.Context, limit int) ([]observability.Operation, error)
	ListFailedOperations(ctx context.Context, limit int) ([]observability.Operation, error)
}

type OpsCommandHandler struct {
	service OpsCommandService
}

func NewOpsCommandHandler(service OpsCommandService) *OpsCommandHandler {
	return &OpsCommandHandler{service: service}
}

func (h *OpsCommandHandler) HandleText(ctx context.Context, text string) (string, error) {
	if h.service == nil {
		return "Observabilidade operacional indisponivel.", nil
	}

	failures, err := h.service.ListFailedOperations(ctx, 5)
	if err != nil {
		return "", err
	}
	recent, err := h.service.ListRecentOperations(ctx, 5)
	if err != nil {
		return "", err
	}

	var lines []string
	lines = append(lines, "# Ops Debug")
	if len(failures) == 0 {
		lines = append(lines, "Falhas recentes: nenhuma")
	} else {
		lines = append(lines, "Falhas recentes:")
		for _, operation := range failures {
			lines = append(lines, formatOperation(operation))
		}
	}

	if len(recent) == 0 {
		lines = append(lines, "Eventos recentes: nenhum")
	} else {
		lines = append(lines, "Eventos recentes:")
		for _, operation := range recent {
			lines = append(lines, formatOperation(operation))
		}
	}
	return strings.Join(lines, "\n"), nil
}

func formatOperation(operation observability.Operation) string {
	line := fmt.Sprintf("- %s/%s status=%s duration_ms=%d", operation.Component, operation.Operation, operation.Status, operation.DurationMS)
	if strings.TrimSpace(operation.RunID) != "" {
		line += " run_id=" + operation.RunID
	}
	if strings.TrimSpace(operation.Summary) != "" {
		line += " summary=" + strings.TrimSpace(operation.Summary)
	}
	return line
}
