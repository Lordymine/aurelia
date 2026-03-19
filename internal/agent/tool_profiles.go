package agent

import (
	"sort"
	"strings"
)

const (
	ToolProfileChatOnly    = "chat_only"
	ToolProfileLocalFiles  = "local_files"
	ToolProfileLocalExec   = "local_exec"
	ToolProfileWebResearch = "web_research"
	ToolProfileScheduler   = "scheduler"
	ToolProfileTeamOps     = "team_ops"
	ToolProfileSkillAdmin  = "skill_admin"
)

var toolProfiles = map[string][]string{
	ToolProfileChatOnly:    {},
	ToolProfileLocalFiles:  {"read_file", "write_file", "list_dir"},
	ToolProfileLocalExec:   {"read_file", "write_file", "list_dir", "run_command"},
	ToolProfileWebResearch: {"web_search"},
	ToolProfileScheduler: {
		"create_schedule",
		"list_schedules",
		"pause_schedule",
		"resume_schedule",
		"delete_schedule",
	},
	ToolProfileTeamOps: {
		"spawn_agent",
		"pause_team",
		"resume_team",
		"cancel_team",
		"team_status",
		"send_team_message",
		"read_team_inbox",
	},
	ToolProfileSkillAdmin: {"install_skill"},
}

func ToolProfileDefinitions() map[string][]string {
	out := make(map[string][]string, len(toolProfiles))
	for key, values := range toolProfiles {
		out[key] = append([]string(nil), values...)
	}
	return out
}

func ResolveAllowedToolsForQuery(query string, explicit []string) []string {
	selected := make(map[string]bool)
	for _, tool := range explicit {
		tool = strings.TrimSpace(tool)
		if tool == "" {
			continue
		}
		selected[tool] = true
	}

	profiles := selectToolProfiles(query)
	for _, profile := range profiles {
		for _, tool := range toolProfiles[profile] {
			selected[tool] = true
		}
	}

	allowed := make([]string, 0, len(selected))
	for tool := range selected {
		allowed = append(allowed, tool)
	}
	sort.Strings(allowed)
	return allowed
}

func ResolveAllowedToolsForWorker(agentName, roleDescription, taskPrompt string, explicit []string) []string {
	if len(explicit) != 0 {
		selected := make(map[string]bool, len(explicit)+2)
		for _, tool := range explicit {
			tool = strings.TrimSpace(tool)
			if tool == "" {
				continue
			}
			selected[tool] = true
		}
		selected["send_team_message"] = true
		selected["read_team_inbox"] = true
		return sortedKeys(selected)
	}

	text := normalizeIntentQuery(strings.Join([]string{agentName, roleDescription, taskPrompt}, "\n"))
	selected := map[string]bool{
		"send_team_message": true,
		"read_team_inbox":   true,
	}

	switch {
	case matchesAny(text, "research", "pesquisa", "researcher", "buscar", "busca", "docs atuais", "documentacao", "internet", "web"):
		selected["web_search"] = true
		selected["read_file"] = true
	case matchesAny(text, "implement", "executor", "builder", "codar", "codigo", "feature", "fix", "corrigir", "refactor", "refatorar"):
		selected["read_file"] = true
		selected["write_file"] = true
		selected["run_command"] = true
	case matchesAny(text, "review", "revisor", "reviewer", "validar", "verification", "verificar", "auditar", "auditor", "checker", "qa", "teste", "test"):
		selected["read_file"] = true
		selected["run_command"] = true
	case matchesAny(text, "plan", "planner", "roadmap", "arquitet", "architecture", "requirements", "requisito", "design"):
		selected["read_file"] = true
		selected["write_file"] = true
	default:
		selected["read_file"] = true
		selected["write_file"] = true
		selected["run_command"] = true
	}

	return sortedKeys(selected)
}

func selectToolProfiles(query string) []string {
	query = normalizeIntentQuery(query)
	if query == "" {
		return []string{ToolProfileChatOnly}
	}

	selected := []string{}
	if matchesAny(query, "/ops", "/memory", "equipe", "time", "subagente", "subagent", "deleg", "worker", "task", "tasks", "mailbox") {
		selected = append(selected, ToolProfileTeamOps)
	}
	if matchesAny(query, "/cron", "agend", "lembrete", "lembrar", "rotina", "periodic", "schedule", "scheduled", "recorr") {
		selected = append(selected, ToolProfileScheduler)
	}
	if matchesAny(query, "instal", "skill", "habilidade") {
		selected = append(selected, ToolProfileSkillAdmin)
	}
	if matchesAny(query, "pesquis", "buscar na internet", "web", "google", "duckduckgo", "online", "documentacao oficial", "docs", "latest", "noticia") {
		selected = append(selected, ToolProfileWebResearch)
	}
	if matchesAny(query, "rodar", "execut", "run ", "run_command", "test", "teste", "build", "compil", "lint", "healthcheck", "endpoint", "servidor", "server", "npm", "go test", "pytest", "cargo", "uv ", "docker") {
		selected = append(selected, ToolProfileLocalExec)
	}
	if matchesAny(query, "arquivo", "file", "ler", "read", "editar", "edit", "escrever", "write", "pasta", "dir", "diretorio", "codigo", "código", "source", "repo", "repositorio") {
		selected = append(selected, ToolProfileLocalFiles)
	}

	if len(selected) == 0 {
		return []string{ToolProfileChatOnly}
	}

	return uniqueStrings(selected)
}

func normalizeIntentQuery(query string) string {
	query = strings.ToLower(strings.TrimSpace(query))
	return strings.Join(strings.Fields(query), " ")
}

func matchesAny(query string, needles ...string) bool {
	for _, needle := range needles {
		needle = strings.TrimSpace(strings.ToLower(needle))
		if needle == "" {
			continue
		}
		if strings.Contains(query, needle) {
			return true
		}
	}
	return false
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func sortedKeys(values map[string]bool) []string {
	out := make([]string, 0, len(values))
	for value := range values {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
