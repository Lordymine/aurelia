package agent

import (
	"sort"
	"strings"
	"unicode"
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

	return sortedKeys(selected)
}

func ResolveAllowedToolsForQueryWithDefinitions(query string, explicit []string, defs []Tool, recentMCPServers []string) []string {
	allowed := ResolveAllowedToolsForQuery(query, explicit)
	selected := make(map[string]bool, len(allowed))
	for _, tool := range allowed {
		selected[tool] = true
	}

	for _, tool := range resolveRelevantMCPTools(query, defs, recentMCPServers) {
		selected[tool] = true
	}

	return sortedKeys(selected)
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
	if matchesAny(query, "arquivo", "ler", "editar", "escrever", "pasta", "diretorio", "codigo", "source", "repo", "repositorio") ||
		queryHasAnyToken(query, "file", "read", "edit", "write", "dir") {
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

func resolveRelevantMCPTools(query string, defs []Tool, recentMCPServers []string) []string {
	query = normalizeIntentQuery(query)
	if query == "" || len(defs) == 0 {
		return nil
	}

	servers := groupMCPToolsByServer(defs)
	selectedServers := make(map[string]bool)
	for serverName, tools := range servers {
		if queryMentionsMCPServer(query, serverName, tools) {
			selectedServers[serverName] = true
		}
	}

	if queryLooksLikeContinuation(query) {
		for _, serverName := range recentMCPServers {
			serverName = normalizeMCPToken(serverName)
			if serverName == "" {
				continue
			}
			if _, ok := servers[serverName]; ok {
				selectedServers[serverName] = true
			}
		}
	}

	selected := make(map[string]bool)
	for serverName := range selectedServers {
		for _, tool := range servers[serverName] {
			selected[tool.Name] = true
		}
	}
	return sortedKeys(selected)
}

func queryMentionsMCPServer(query, serverName string, tools []Tool) bool {
	serverName = normalizeMCPToken(serverName)
	if serverName == "" {
		return false
	}

	if queryContainsServerName(query, serverName) {
		return true
	}

	for _, tool := range tools {
		toolText := normalizeIntentQuery(strings.TrimSpace(tool.Name + " " + tool.Description))
		if toolText == "" {
			continue
		}
		if strings.Contains(query, toolText) || strings.Contains(toolText, query) {
			return true
		}
	}

	return false
}

func SplitMCPToolName(toolName string) (serverName string, remoteName string) {
	if !strings.HasPrefix(toolName, "mcp_") {
		return "", ""
	}

	parts := strings.Split(strings.TrimPrefix(toolName, "mcp_"), "_")
	if len(parts) < 2 {
		return "", ""
	}

	serverName = parts[0]
	remoteName = strings.Join(parts[1:], "_")
	return serverName, remoteName
}

func groupMCPToolsByServer(defs []Tool) map[string][]Tool {
	grouped := make(map[string][]Tool)
	for _, def := range defs {
		serverName, _ := SplitMCPToolName(def.Name)
		if serverName == "" {
			continue
		}
		grouped[serverName] = append(grouped[serverName], def)
	}
	return grouped
}

func queryContainsServerName(query, serverName string) bool {
	serverName = normalizeMCPToken(serverName)
	if serverName == "" {
		return false
	}

	collapsedQuery := collapseAlphaNumeric(query)
	if strings.Contains(collapsedQuery, serverName) {
		return true
	}

	for _, token := range tokenizeAlphaNumeric(query) {
		if token == serverName {
			return true
		}
		distance := levenshteinDistance(token, serverName)
		if len(serverName) >= 8 && distance <= 3 {
			return true
		}
		if len(serverName) >= 5 && distance <= 1 {
			return true
		}
	}

	return false
}

func queryLooksLikeContinuation(query string) bool {
	return matchesAny(query,
		"agora",
		"nessa",
		"nesta",
		"nesse",
		"nessa janela",
		"nessa aba",
		"nessa tela",
		"continue",
		"continua",
		"prossiga",
		"vai",
		"entao",
		"depois",
		"em seguida",
		"segue",
	)
}

func queryHasAnyToken(query string, needles ...string) bool {
	tokens := tokenizeAlphaNumeric(query)
	if len(tokens) == 0 {
		return false
	}
	seen := make(map[string]bool, len(tokens))
	for _, token := range tokens {
		seen[token] = true
	}
	for _, needle := range needles {
		needle = normalizeMCPToken(needle)
		if needle == "" {
			continue
		}
		if seen[needle] {
			return true
		}
	}
	return false
}

func normalizeMCPToken(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, " ", "")
	return value
}

func collapseAlphaNumeric(value string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func tokenizeAlphaNumeric(value string) []string {
	fields := strings.FieldsFunc(strings.ToLower(value), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	out := make([]string, 0, len(fields))
	for _, field := range fields {
		field = normalizeMCPToken(field)
		if field == "" {
			continue
		}
		out = append(out, field)
	}
	return out
}

func levenshteinDistance(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return len([]rune(b))
	}
	if b == "" {
		return len([]rune(a))
	}

	ar := []rune(a)
	br := []rune(b)
	prev := make([]int, len(br)+1)
	for j := range prev {
		prev[j] = j
	}

	for i, ra := range ar {
		curr := make([]int, len(br)+1)
		curr[0] = i + 1
		for j, rb := range br {
			cost := 0
			if ra != rb {
				cost = 1
			}
			curr[j+1] = minInt(
				prev[j+1]+1,
				curr[j]+1,
				prev[j]+cost,
			)
		}
		prev = curr
	}
	return prev[len(br)]
}

func minInt(values ...int) int {
	best := 0
	for i, value := range values {
		if i == 0 || value < best {
			best = value
		}
	}
	return best
}
