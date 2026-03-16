package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/cron"
)

type CronCommandService interface {
	AddRecurringJob(ctx context.Context, userID string, chatID int64, expr, prompt string) (string, error)
	AddOnceJob(ctx context.Context, userID string, chatID int64, timestamp, prompt string) (string, error)
	ListJobs(ctx context.Context, chatID int64) ([]cron.CronJob, error)
	PauseJob(ctx context.Context, jobID string) error
	ResumeJob(ctx context.Context, jobID string) error
	DeleteJob(ctx context.Context, jobID string) error
}

type CronCommandHandler struct {
	service CronCommandService
}

func NewCronCommandHandler(service CronCommandService) *CronCommandHandler {
	return &CronCommandHandler{service: service}
}

func (h *CronCommandHandler) HandleText(ctx context.Context, userID string, chatID int64, text string) (string, error) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "/cron") {
		return cronUsage(), nil
	}

	rest := strings.TrimSpace(strings.TrimPrefix(text, "/cron"))
	if rest == "" {
		return cronUsage(), nil
	}

	switch {
	case rest == "list":
		jobs, err := h.service.ListJobs(ctx, chatID)
		if err != nil {
			return "", err
		}
		return formatCronJobs(jobs), nil
	case strings.HasPrefix(rest, "add "):
		args := parseQuotedArgs(strings.TrimSpace(strings.TrimPrefix(rest, "add")))
		if len(args) != 2 {
			return cronUsage(), nil
		}
		jobID, err := h.service.AddRecurringJob(ctx, userID, chatID, args[0], args[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Job cron criado com sucesso: `%s`", jobID), nil
	case strings.HasPrefix(rest, "once "):
		args := parseQuotedArgs(strings.TrimSpace(strings.TrimPrefix(rest, "once")))
		if len(args) != 2 {
			return cronUsage(), nil
		}
		jobID, err := h.service.AddOnceJob(ctx, userID, chatID, args[0], args[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Job pontual criado com sucesso: `%s`", jobID), nil
	case strings.HasPrefix(rest, "pause "):
		jobID := strings.TrimSpace(strings.TrimPrefix(rest, "pause"))
		if jobID == "" {
			return cronUsage(), nil
		}
		if err := h.service.PauseJob(ctx, jobID); err != nil {
			return "", err
		}
		return fmt.Sprintf("Job `%s` pausado.", jobID), nil
	case strings.HasPrefix(rest, "resume "):
		jobID := strings.TrimSpace(strings.TrimPrefix(rest, "resume"))
		if jobID == "" {
			return cronUsage(), nil
		}
		if err := h.service.ResumeJob(ctx, jobID); err != nil {
			return "", err
		}
		return fmt.Sprintf("Job `%s` retomado.", jobID), nil
	case strings.HasPrefix(rest, "del "):
		jobID := strings.TrimSpace(strings.TrimPrefix(rest, "del"))
		if jobID == "" {
			return cronUsage(), nil
		}
		if err := h.service.DeleteJob(ctx, jobID); err != nil {
			return "", err
		}
		return fmt.Sprintf("Job `%s` removido.", jobID), nil
	default:
		return cronUsage(), nil
	}
}

func cronUsage() string {
	return "Uso: /cron list | /cron add \"<expr>\" \"<prompt>\" | /cron once \"<timestamp>\" \"<prompt>\" | /cron pause <id> | /cron resume <id> | /cron del <id>"
}

func parseQuotedArgs(input string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false

	for _, r := range input {
		switch r {
		case '"':
			if inQuotes {
				args = append(args, current.String())
				current.Reset()
			}
			inQuotes = !inQuotes
		default:
			if inQuotes {
				current.WriteRune(r)
			}
		}
	}

	return args
}

func formatCronJobs(jobs []cron.CronJob) string {
	if len(jobs) == 0 {
		return "Nenhum job cron cadastrado neste chat."
	}

	var lines []string
	for _, job := range jobs {
		schedule := job.CronExpr
		if job.ScheduleType == "once" && job.RunAt != nil {
			schedule = job.RunAt.Format(time.RFC3339)
		}
		lines = append(lines, fmt.Sprintf("- `%s` [%s] active=%t status=%s schedule=%s prompt=%s", job.ID, job.ScheduleType, job.Active, job.LastStatus, schedule, job.Prompt))
	}
	return strings.Join(lines, "\n")
}


