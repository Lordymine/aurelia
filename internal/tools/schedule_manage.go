package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	cronpkg "github.com/kocar/aurelia/internal/cron"
)

type ScheduleLister interface {
	ListJobs(ctx context.Context, chatID int64) ([]cronpkg.CronJob, error)
}

type ScheduleController interface {
	PauseJob(ctx context.Context, jobID string) error
	ResumeJob(ctx context.Context, jobID string) error
	DeleteJob(ctx context.Context, jobID string) error
}

type ListSchedulesTool struct {
	service ScheduleLister
}

func NewListSchedulesTool(service ScheduleLister) *ListSchedulesTool {
	return &ListSchedulesTool{service: service}
}

func (t *ListSchedulesTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "list_schedules",
		Description: "Lista os agendamentos do chat atual.",
		JSONSchema:  objectSchema(map[string]any{}),
	}
}

func (t *ListSchedulesTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	if t.service == nil {
		return "", fmt.Errorf("schedule lister service is not configured")
	}
	teamKey, _, ok := agent.TeamContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing chat context for list_schedules")
	}
	chatID, err := strconv.ParseInt(teamKey, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid chat context for list_schedules: %w", err)
	}

	jobs, err := t.service.ListJobs(ctx, chatID)
	if err != nil {
		return "", err
	}
	if len(jobs) == 0 {
		return "Nenhum agendamento encontrado para este chat.", nil
	}

	var lines []string
	for _, job := range jobs {
		schedule := job.CronExpr
		if job.ScheduleType == "once" && job.RunAt != nil {
			schedule = job.RunAt.Format(time.RFC3339)
		}
		lines = append(lines, fmt.Sprintf("- `%s` [%s] active=%t status=%s schedule=%s prompt=%s", job.ID, job.ScheduleType, job.Active, job.LastStatus, schedule, job.Prompt))
	}
	return strings.Join(lines, "\n"), nil
}

type PauseScheduleTool struct {
	service ScheduleController
}

func NewPauseScheduleTool(service ScheduleController) *PauseScheduleTool {
	return &PauseScheduleTool{service: service}
}

func (t *PauseScheduleTool) Definition() agent.Tool {
	return scheduleControlDefinition("pause_schedule", "Pausa um agendamento existente.")
}

func (t *PauseScheduleTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	return executeScheduleControl(ctx, t.service, "pause", args)
}

type ResumeScheduleTool struct {
	service ScheduleController
}

func NewResumeScheduleTool(service ScheduleController) *ResumeScheduleTool {
	return &ResumeScheduleTool{service: service}
}

func (t *ResumeScheduleTool) Definition() agent.Tool {
	return scheduleControlDefinition("resume_schedule", "Retoma um agendamento pausado.")
}

func (t *ResumeScheduleTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	return executeScheduleControl(ctx, t.service, "resume", args)
}

type DeleteScheduleTool struct {
	service ScheduleController
}

func NewDeleteScheduleTool(service ScheduleController) *DeleteScheduleTool {
	return &DeleteScheduleTool{service: service}
}

func (t *DeleteScheduleTool) Definition() agent.Tool {
	return scheduleControlDefinition("delete_schedule", "Remove um agendamento existente.")
}

func (t *DeleteScheduleTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	return executeScheduleControl(ctx, t.service, "delete", args)
}

func scheduleControlDefinition(name, description string) agent.Tool {
	return agent.Tool{
		Name:        name,
		Description: description,
		JSONSchema: objectSchema(
			map[string]any{
				"job_id": stringProperty(""),
			},
			"job_id",
		),
	}
}

func executeScheduleControl(ctx context.Context, service ScheduleController, action string, args map[string]interface{}) (string, error) {
	if service == nil {
		return "", fmt.Errorf("schedule controller service is not configured")
	}
	jobID, _ := args["job_id"].(string)
	jobID = strings.TrimSpace(jobID)
	if jobID == "" {
		return "", fmt.Errorf("job_id is required")
	}

	var err error
	switch action {
	case "pause":
		err = service.PauseJob(ctx, jobID)
	case "resume":
		err = service.ResumeJob(ctx, jobID)
	case "delete":
		err = service.DeleteJob(ctx, jobID)
	default:
		return "", fmt.Errorf("unsupported schedule action %q", action)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Agendamento `%s` executou a acao `%s` com sucesso.", jobID, action), nil
}


