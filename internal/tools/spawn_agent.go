package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
)

type TeamSpawner interface {
	Spawn(ctx context.Context, teamKey, userID, agentName, roleDescription, taskPrompt string, allowedTools ...string) (string, error)
}

type SpawnAgentTool struct {
	Spawner TeamSpawner
}

func NewSpawnAgentTool(spawner TeamSpawner) *SpawnAgentTool {
	return &SpawnAgentTool{Spawner: spawner}
}

func (t *SpawnAgentTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "spawn_agent",
		Description: "Cria uma task para um worker especialista dentro da equipe liderada pelo master. O master continua sendo o lider e acompanha os resultados.",
		JSONSchema: objectSchema(
			map[string]any{
				"agent_name":       stringProperty("Nome do agente especialista."),
				"role_description": stringProperty("Papel e metodologia exata do especialista."),
				"task_prompt":      stringProperty("Tarefa pratica a ser executada pelo especialista."),
				"workdir":          stringProperty("Diretorio de trabalho canonico do projeto alvo para esse worker."),
				"allowed_tools": map[string]any{
					"type":        "array",
					"description": "Whitelist opcional de tools para esse worker. Se omitido, o worker recebe um perfil automatico compativel com a task.",
					"items": map[string]any{
						"type": "string",
					},
				},
			},
			"agent_name",
			"role_description",
			"task_prompt",
		),
	}
}

func (t *SpawnAgentTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	agentName, _ := args["agent_name"].(string)
	roleDesc, _ := args["role_description"].(string)
	taskPrompt, _ := args["task_prompt"].(string)
	workdir := optionalStringArg(args, "workdir")
	allowedTools := readStringArrayArg(args["allowed_tools"])
	if len(allowedTools) == 0 {
		allowedTools = inferAllowedToolsProfile(agentName, roleDesc, taskPrompt)
	}

	if t.Spawner == nil {
		return "", fmt.Errorf("master team service is not configured")
	}

	teamKey, userID, ok := agent.TeamContextFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing master team context for spawn_agent")
	}
	if workdir != "" {
		ctx = agent.WithWorkdirContext(ctx, workdir)
	}

	taskID, err := t.Spawner.Spawn(ctx, teamKey, userID, agentName, roleDesc, taskPrompt, allowedTools...)
	if err != nil {
		return "", fmt.Errorf("falha ao criar task do sub-agente: %v", err)
	}

	return fmt.Sprintf("Acionei o especialista `%s` para cuidar de: %s.\nTask aberta: `%s`.\nVou acompanhar o progresso do time e te atualizar quando houver avancos relevantes.", agentName, taskPrompt, taskID), nil
}

func inferAllowedToolsProfile(agentName, roleDescription, taskPrompt string) []string {
	text := strings.ToLower(strings.Join([]string{agentName, roleDescription, taskPrompt}, "\n"))

	mailbox := []string{"send_team_message", "read_team_inbox"}

	if matchesAny(text, "research", "pesquisa", "researcher", "buscar", "busca", "docs atuais", "documentacao", "internet", "web") {
		return append([]string{"web_search", "read_file"}, mailbox...)
	}

	if matchesAny(text, "implement", "executor", "builder", "codar", "codigo", "feature", "fix", "corrigir", "refactor", "refatorar") {
		return append([]string{"read_file", "write_file", "run_command"}, mailbox...)
	}

	if matchesAny(text, "review", "revisor", "reviewer", "validar", "verification", "verificar", "auditar", "auditor", "checker", "qa", "teste", "test") {
		return append([]string{"read_file", "run_command"}, mailbox...)
	}

	if matchesAny(text, "plan", "planner", "roadmap", "arquitet", "architecture", "requirements", "requisito", "design") {
		return append([]string{"read_file", "write_file"}, mailbox...)
	}

	return append([]string{"read_file", "write_file", "run_command"}, mailbox...)
}

func matchesAny(text string, patterns ...string) bool {
	for _, pattern := range patterns {
		if strings.Contains(text, pattern) {
			return true
		}
	}
	return false
}


