package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/agent"
)

type TeamMessenger interface {
	PostMessage(ctx context.Context, msg agent.MailMessage) error
	PullMessages(ctx context.Context, teamID, agentName string, limit int) ([]agent.MailMessage, error)
}

type SendTeamMessageTool struct {
	Manager TeamMessenger
}

func NewSendTeamMessageTool(manager TeamMessenger) *SendTeamMessageTool {
	return &SendTeamMessageTool{Manager: manager}
}

func (t *SendTeamMessageTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "send_team_message",
		Description: "Envia uma mensagem para outro agente da mesma equipe usando a mailbox interna.",
		JSONSchema: objectSchema(
			map[string]any{
				"to_agent": stringProperty("Nome do agente destinatario."),
				"body":     stringProperty("Conteudo objetivo da mensagem."),
				"kind":     stringProperty("Tipo opcional da mensagem, por exemplo status_update, question ou result."),
			},
			"to_agent",
			"body",
		),
	}
}

func (t *SendTeamMessageTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	teamID, taskID, ok := agent.TaskContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing task context for send_team_message")
	}
	fromAgent, ok := agent.AgentContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing agent context for send_team_message")
	}
	toAgent, err := requireStringArg(args, "to_agent")
	if err != nil {
		return "", err
	}
	body, err := requireStringArg(args, "body")
	if err != nil {
		return "", err
	}
	kind := optionalStringArg(args, "kind")
	if kind == "" {
		kind = "status_update"
	}
	if t.Manager == nil {
		return "", fmt.Errorf("team mailbox manager is not configured")
	}

	if err := t.Manager.PostMessage(ctx, agent.MailMessage{
		ID:        uuid.NewString(),
		TeamID:    teamID,
		FromAgent: fromAgent,
		ToAgent:   toAgent,
		TaskID:    &taskID,
		Kind:      kind,
		Body:      body,
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("Mensagem enviada para `%s` com kind `%s`.", toAgent, kind), nil
}

type ReadTeamInboxTool struct {
	Manager TeamMessenger
}

func NewReadTeamInboxTool(manager TeamMessenger) *ReadTeamInboxTool {
	return &ReadTeamInboxTool{Manager: manager}
}

func (t *ReadTeamInboxTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "read_team_inbox",
		Description: "Le e consome mensagens pendentes da inbox do agente atual na equipe.",
		JSONSchema: objectSchema(
			map[string]any{
				"limit": numberProperty("Quantidade maxima de mensagens para ler."),
			},
		),
	}
}

func (t *ReadTeamInboxTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	teamID, _, ok := agent.TaskContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing task context for read_team_inbox")
	}
	agentName, ok := agent.AgentContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing agent context for read_team_inbox")
	}
	if t.Manager == nil {
		return "", fmt.Errorf("team mailbox manager is not configured")
	}

	limit := 10
	if raw, ok := args["limit"].(float64); ok && raw > 0 {
		limit = int(raw)
	}

	messages, err := t.Manager.PullMessages(ctx, teamID, agentName, limit)
	if err != nil {
		return "", err
	}
	if len(messages) == 0 {
		return "[]", nil
	}

	payload, err := json.Marshal(messages)
	if err != nil {
		return "", fmt.Errorf("marshal inbox messages: %w", err)
	}
	return string(payload), nil
}


