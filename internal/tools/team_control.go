package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
)

type TeamController interface {
	Pause(ctx context.Context, teamKey string) error
	Resume(ctx context.Context, teamKey string) error
	Cancel(ctx context.Context, teamKey, reason string) error
	BuildStatusSnapshot(ctx context.Context, teamKey string) (agent.TeamStatusSnapshot, error)
}

type PauseTeamTool struct {
	Controller TeamController
}

func NewPauseTeamTool(controller TeamController) *PauseTeamTool {
	return &PauseTeamTool{Controller: controller}
}

func (t *PauseTeamTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "pause_team",
		Description: "Pausa a distribuicao de novas tasks para a equipe atual.",
		JSONSchema:  objectSchema(map[string]any{}),
	}
}

func (t *PauseTeamTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	if t.Controller == nil {
		return "", fmt.Errorf("team controller is not configured")
	}
	teamKey, _, ok := agent.TeamContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing team context for pause_team")
	}
	if err := t.Controller.Pause(ctx, teamKey); err != nil {
		return "", err
	}
	return "Pausei a equipe atual. Nenhuma nova task sera iniciada ate voce mandar retomar.", nil
}

type ResumeTeamTool struct {
	Controller TeamController
}

func NewResumeTeamTool(controller TeamController) *ResumeTeamTool {
	return &ResumeTeamTool{Controller: controller}
}

func (t *ResumeTeamTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "resume_team",
		Description: "Retoma a distribuicao de tasks da equipe atual.",
		JSONSchema:  objectSchema(map[string]any{}),
	}
}

func (t *ResumeTeamTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	if t.Controller == nil {
		return "", fmt.Errorf("team controller is not configured")
	}
	teamKey, _, ok := agent.TeamContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing team context for resume_team")
	}
	if err := t.Controller.Resume(ctx, teamKey); err != nil {
		return "", err
	}
	return "Retomei a equipe atual. Vou continuar distribuindo as tasks pendentes.", nil
}

type CancelTeamTool struct {
	Controller TeamController
}

func NewCancelTeamTool(controller TeamController) *CancelTeamTool {
	return &CancelTeamTool{Controller: controller}
}

func (t *CancelTeamTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "cancel_team",
		Description: "Cancela a operacao da equipe atual, interrompendo o que estiver em andamento e marcando tasks abertas como canceladas.",
		JSONSchema: objectSchema(map[string]any{
			"reason": stringProperty("Motivo opcional do cancelamento."),
		}),
	}
}

func (t *CancelTeamTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	if t.Controller == nil {
		return "", fmt.Errorf("team controller is not configured")
	}
	teamKey, _, ok := agent.TeamContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing team context for cancel_team")
	}
	reason := strings.TrimSpace(optionalStringArg(args, "reason"))
	if err := t.Controller.Cancel(ctx, teamKey, reason); err != nil {
		return "", err
	}
	if reason == "" {
		return "Cancelei a operacao da equipe atual.", nil
	}
	return fmt.Sprintf("Cancelei a operacao da equipe atual. Motivo registrado: %s.", reason), nil
}

type TeamStatusTool struct {
	Controller TeamController
}

func NewTeamStatusTool(controller TeamController) *TeamStatusTool {
	return &TeamStatusTool{Controller: controller}
}

func (t *TeamStatusTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "team_status",
		Description: "Mostra um resumo do estado atual da equipe e das tasks em aberto.",
		JSONSchema:  objectSchema(map[string]any{}),
	}
}

func (t *TeamStatusTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	if t.Controller == nil {
		return "", fmt.Errorf("team controller is not configured")
	}
	teamKey, _, ok := agent.TeamContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing team context for team_status")
	}
	snapshot, err := t.Controller.BuildStatusSnapshot(ctx, teamKey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"Equipe atual: status=%s | pendentes=%d | rodando=%d | bloqueadas=%d | concluidas=%d | falhas=%d | canceladas=%d | total=%d",
		snapshot.TeamStatus,
		snapshot.Pending,
		snapshot.Running,
		snapshot.Blocked,
		snapshot.Completed,
		snapshot.Failed,
		snapshot.Cancelled,
		snapshot.TotalTasks,
	), nil
}


