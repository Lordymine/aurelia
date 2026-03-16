package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/agent"
	cronpkg "github.com/kocar/aurelia/internal/cron"
)

type ScheduleCreator interface {
	CreateJob(ctx context.Context, job cronpkg.CronJob) (string, error)
}

type CreateScheduleTool struct {
	service ScheduleCreator
}

func NewCreateScheduleTool(service ScheduleCreator) *CreateScheduleTool {
	return &CreateScheduleTool{service: service}
}

func (t *CreateScheduleTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "create_schedule",
		Description: "Cria um agendamento recorrente ou pontual para o chat atual ou para um target_chat_id especifico.",
		JSONSchema: objectSchema(
			map[string]any{
				"schedule_type":  stringProperty("cron ou once"),
				"cron_expr":      stringProperty("Expressao cron quando schedule_type=cron"),
				"run_at":         stringProperty("Timestamp RFC3339 quando schedule_type=once"),
				"prompt":         stringProperty("Prompt que sera executado futuramente"),
				"target_chat_id": numberProperty("Chat alvo opcional; por default usa o chat atual"),
			},
			"schedule_type",
			"prompt",
		),
	}
}

func (t *CreateScheduleTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	if t.service == nil {
		return "", fmt.Errorf("schedule creator service is not configured")
	}

	teamKey, userID, ok := agent.TeamContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing chat context for create_schedule")
	}

	targetChatID, err := parseTargetChatID(teamKey, args["target_chat_id"])
	if err != nil {
		return "", err
	}

	scheduleType, _ := args["schedule_type"].(string)
	scheduleType = strings.TrimSpace(strings.ToLower(scheduleType))
	prompt, _ := args["prompt"].(string)
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return "", fmt.Errorf("prompt is required")
	}

	job := cronpkg.CronJob{
		ID:           uuid.NewString(),
		OwnerUserID:  userID,
		TargetChatID: targetChatID,
		ScheduleType: scheduleType,
		Prompt:       prompt,
		Active:       true,
		LastStatus:   "idle",
	}

	switch scheduleType {
	case "cron":
		cronExpr, _ := args["cron_expr"].(string)
		cronExpr = strings.TrimSpace(cronExpr)
		if cronExpr == "" {
			return "", fmt.Errorf("cron_expr is required for recurring schedules")
		}
		job.CronExpr = cronExpr
	case "once":
		runAtRaw, _ := args["run_at"].(string)
		runAtRaw = strings.TrimSpace(runAtRaw)
		if runAtRaw == "" {
			return "", fmt.Errorf("run_at is required for once schedules")
		}
		runAt, err := time.Parse(time.RFC3339, runAtRaw)
		if err != nil {
			return "", fmt.Errorf("invalid run_at timestamp: %w", err)
		}
		job.RunAt = &runAt
		job.NextRunAt = &runAt
	default:
		return "", fmt.Errorf("unsupported schedule_type %q", scheduleType)
	}

	jobID, err := t.service.CreateJob(ctx, job)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Agendamento criado com sucesso com id `%s`.", jobID), nil
}

func parseTargetChatID(teamKey string, raw any) (int64, error) {
	if raw == nil {
		return strconv.ParseInt(teamKey, 10, 64)
	}
	if value, ok := raw.(float64); ok {
		return int64(value), nil
	}
	if value, ok := raw.(string); ok && strings.TrimSpace(value) != "" {
		return strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	}
	return 0, fmt.Errorf("invalid target_chat_id")
}


