package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"gopkg.in/telebot.v3"

	"github.com/kocar/aurelia/internal/bridge"
)

// CommandType identifies a system command that can be handled locally without LLM.
type CommandType int

const (
	CmdCronCreate   CommandType = iota
	CmdCronList
	CmdCronCancel
	CmdSessionReset
	CmdStatus
	CmdListAgents
	CmdListModels
)

// MatchedCommand represents a message that was identified as a system command.
type MatchedCommand struct {
	Type CommandType
	Text string // original message text
}

// commandRule defines a pattern for matching a command.
type commandRule struct {
	cmdType  CommandType
	phrases  []string // phrase matches (checked in order)
	exact    bool     // if true, the entire message must equal one of the phrases
}

// rules are checked in order; first match wins.
var commandRules = []commandRule{
	// cron_list (check before cron_create to avoid "lista agendamentos" matching create)
	{CmdCronList, []string{
		"meus agendamentos",
		"o que tá agendado", "o que ta agendado",
		"lista agendamentos", "listar agendamentos",
	}, false},
	// cron_cancel (check before cron_create to avoid "cancela" partial matching)
	{CmdCronCancel, []string{
		"cancela o agendamento", "cancela agendamento", "cancelar agendamento",
		"cancele o agendamento", "cancele agendamento",
		"remove o agendamento", "remove agendamento",
		"remove o lembrete", "remove lembrete", "remover lembrete",
		"desativa agendamento", "desativar agendamento",
		"exclui agendamento", "excluir agendamento",
		"deleta agendamento", "deletar agendamento",
		"apaga agendamento", "apagar agendamento",
	}, false},
	// cron_create
	{CmdCronCreate, []string{
		"agenda ", "agendar ", "agende ",
		"cria um lembrete", "cria um agendamento", "criar lembrete", "criar agendamento",
		"crie um lembrete", "crie um agendamento",
		"me lembra ", "me lembre ",
	}, false},
	// session_reset ("reset" is exact-only but other phrases use substring)
	{CmdSessionReset, []string{
		"nova conversa",
		"limpa o contexto", "limpar contexto", "limpa contexto",
		"começa de novo", "comeca de novo",
	}, false},
	{CmdSessionReset, []string{
		"reset",
	}, true},
	// status (exact match only — "status" is too common as a word)
	{CmdStatus, []string{
		"status",
		"tá funcionando", "ta funcionando",
	}, true},
	// list_agents
	{CmdListAgents, []string{
		"quais agents", "lista agents", "listar agents",
		"meus agents",
	}, false},
	// list_models
	{CmdListModels, []string{
		"quais modelos", "lista modelos", "listar modelos",
		"quais provedores", "lista provedores", "listar provedores",
	}, false},
}

// MatchCommand checks if a message is a system command. Returns nil if no match.
// Uses keyword matching with a narrative-context heuristic to avoid false positives.
func MatchCommand(text string) *MatchedCommand {
	lower := strings.ToLower(strings.TrimSpace(text))
	if lower == "" {
		return nil
	}

	for _, rule := range commandRules {
		for _, phrase := range rule.phrases {
			if rule.exact {
				// Strip trailing punctuation for exact match comparison
				cleaned := strings.TrimRightFunc(lower, func(r rune) bool {
					return unicode.IsPunct(r) || unicode.IsSpace(r)
				})
				if cleaned == phrase {
					return &MatchedCommand{Type: rule.cmdType, Text: text}
				}
				continue
			}

			idx := strings.Index(lower, phrase)
			if idx < 0 {
				continue
			}
			// Anti-false-positive: if the phrase is not at the start,
			// check if the preceding text looks like narrative context.
			if idx > 0 && looksNarrative(lower[:idx]) {
				continue
			}
			return &MatchedCommand{Type: rule.cmdType, Text: text}
		}
	}
	return nil
}

// looksNarrative returns true if the prefix text suggests the keyword appears
// inside a narrative sentence rather than as a command. We check if there are
// word characters before the keyword (excluding small connectors like "me", "um").
func looksNarrative(prefix string) bool {
	trimmed := strings.TrimRightFunc(prefix, unicode.IsSpace)
	if trimmed == "" {
		return false
	}

	// Count significant words (3+ chars) in the prefix
	words := strings.Fields(trimmed)
	significant := 0
	for _, w := range words {
		if len(w) >= 3 {
			significant++
		}
	}

	// If there are 2+ significant words before the keyword, it's likely narrative.
	// "ontem eu tentei agendar" → 3 significant words → narrative
	// "me lembra" → 0 significant words (only "me") → not narrative
	return significant >= 2
}

// handleCommand dispatches a matched command to the appropriate handler.
// Returns nil if the command was handled (response sent to Telegram).
func (bc *BotController) handleCommand(c telebot.Context, cmd *MatchedCommand) error {
	chatID := c.Chat().ID
	messageID := c.Message().ID
	log.Printf("command: type=%d chat=%d text=%q", cmd.Type, chatID, cmd.Text)

	var reply string
	var err error

	switch cmd.Type {
	case CmdSessionReset:
		reply, err = bc.cmdSessionReset(chatID)
	case CmdCronList:
		reply, err = bc.cmdCronList(chatID)
	case CmdCronCancel:
		reply, err = bc.cmdCronCancel(chatID, cmd.Text)
	case CmdCronCreate:
		reply, err = bc.cmdCronCreate(c, cmd.Text)
	case CmdStatus:
		reply, err = bc.cmdStatus(chatID)
	case CmdListAgents:
		reply, err = bc.cmdListAgents()
	case CmdListModels:
		reply, err = bc.cmdListModels()
	default:
		return fmt.Errorf("unknown command type: %d", cmd.Type)
	}

	if err != nil {
		log.Printf("command error: type=%d err=%v", cmd.Type, err)
		return SendError(bc.bot, c.Chat(), fmt.Sprintf("Erro ao executar comando: %v", err))
	}

	return SendTextReply(bc.bot, c.Chat(), reply, messageID)
}

// --- P1 handlers ---

func (bc *BotController) cmdSessionReset(chatID int64) (string, error) {
	bc.sessions.Clear(chatID)
	bc.tracker.Clear(chatID)
	log.Printf("command: session reset for chat=%d", chatID)
	return "Contexto limpo. Próxima mensagem começa uma conversa nova.", nil
}

func (bc *BotController) cmdCronList(chatID int64) (string, error) {
	if bc.cronHandler == nil {
		return "Sistema de agendamentos não disponível.", nil
	}
	ctx := context.Background()
	jobs, err := bc.cronHandler.service.ListJobs(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("falha ao listar agendamentos: %w", err)
	}
	if len(jobs) == 0 {
		return "Nenhum agendamento encontrado.", nil
	}
	return formatCronJobs(jobs), nil
}

func (bc *BotController) cmdCronCancel(chatID int64, text string) (string, error) {
	if bc.cronHandler == nil {
		return "Sistema de agendamentos não disponível.", nil
	}

	// Extract job ID: last word of the message (should be a UUID prefix)
	jobID := extractLastWord(text)
	if jobID == "" || !looksLikeJobID(jobID) {
		return "Não encontrei o ID do agendamento. Use 'meus agendamentos' pra ver a lista com os IDs.", nil
	}

	ctx := context.Background()
	if err := bc.cronHandler.service.DeleteJob(ctx, jobID); err != nil {
		return "", fmt.Errorf("falha ao cancelar agendamento: %w", err)
	}

	return fmt.Sprintf("Agendamento `%s` removido.", jobID), nil
}

// extractLastWord returns the last whitespace-separated token from text.
func extractLastWord(text string) string {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return ""
	}
	return fields[len(fields)-1]
}

// looksLikeJobID checks if a string looks like a UUID prefix (alphanumeric, 4+ chars).
func looksLikeJobID(s string) bool {
	if len(s) < 4 {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'f') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return true
}

const cronParseTimeout = 10 * time.Second

const cronParseSystemPrompt = `You are a scheduling assistant. Extract scheduling parameters from the user's message.

Respond with ONLY a JSON object (no markdown, no explanation):

For recurring schedules:
{"type":"cron","cron_expr":"<cron expression>","prompt":"<what to do>"}

For one-time schedules:
{"type":"once","run_at":"<ISO 8601 timestamp>","prompt":"<what to do>"}

Rules:
- cron_expr uses standard 5-field cron: minute hour day month weekday
- run_at must be ISO 8601 with timezone (use -03:00 for BRT unless specified)
- prompt is the ACTION to perform, not the scheduling part
- If the user says "amanhã" or relative dates, calculate from current time
- If no time specified, default to 09:00

Examples:
"agenda todo dia às 9h revisar emails" → {"type":"cron","cron_expr":"0 9 * * *","prompt":"revisar emails"}
"me lembra amanhã às 15h de fazer deploy" → {"type":"once","run_at":"2026-03-27T15:00:00-03:00","prompt":"fazer deploy"}
"agendar toda segunda e quarta às 10h standup" → {"type":"cron","cron_expr":"0 10 * * 1,3","prompt":"standup"}`

// cronCreateParsed holds the extracted scheduling parameters from LLM response.
type cronCreateParsed struct {
	Type     string `json:"type"`
	CronExpr string `json:"cron_expr,omitempty"`
	RunAt    string `json:"run_at,omitempty"`
	Prompt   string `json:"prompt"`
}

// parseCronCreateResponse extracts JSON from the LLM response, tolerating markdown fences.
func parseCronCreateResponse(raw string) (*cronCreateParsed, error) {
	// Strip markdown code fences if present
	cleaned := strings.TrimSpace(raw)
	if strings.HasPrefix(cleaned, "```") {
		// Remove first line (```json) and last line (```)
		lines := strings.Split(cleaned, "\n")
		if len(lines) >= 3 {
			cleaned = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}
	cleaned = strings.TrimSpace(cleaned)

	var parsed cronCreateParsed
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	if parsed.Prompt == "" {
		return nil, fmt.Errorf("missing prompt in response")
	}

	switch parsed.Type {
	case "cron":
		if parsed.CronExpr == "" {
			return nil, fmt.Errorf("missing cron_expr for cron type")
		}
	case "once":
		if parsed.RunAt == "" {
			return nil, fmt.Errorf("missing run_at for once type")
		}
	default:
		return nil, fmt.Errorf("unknown type %q", parsed.Type)
	}

	return &parsed, nil
}

func (bc *BotController) cmdCronCreate(c telebot.Context, text string) (string, error) {
	if bc.cronHandler == nil {
		return "Sistema de agendamentos não disponível.", nil
	}
	if bc.bridge == nil {
		return "Processador não disponível para interpretar o agendamento.", nil
	}

	// Use a focused LLM call to extract scheduling parameters from natural language
	ctx, cancel := context.WithTimeout(context.Background(), cronParseTimeout)
	defer cancel()

	result, err := bc.bridge.ExecuteSync(ctx, bridge.Request{
		Command: "query",
		Prompt:  text,
		Options: bridge.RequestOptions{
			Model:          bc.config.DefaultModel,
			SystemPrompt:   cronParseSystemPrompt,
			MaxTurns:       1,
			PermissionMode: "bypassPermissions",
		},
	})
	if err != nil {
		return "Não consegui interpretar o agendamento. Tente de novo com mais detalhes.", nil
	}

	parsed, parseErr := parseCronCreateResponse(result.Content)
	if parseErr != nil {
		log.Printf("command: cron_create parse error: %v (raw: %q)", parseErr, result.Content)
		return "Não entendi o agendamento. Tente algo como: \"agenda todo dia às 9h revisar emails\"", nil
	}

	chatID := c.Chat().ID
	userID := fmt.Sprintf("%d", c.Sender().ID)

	var jobID string
	switch parsed.Type {
	case "cron":
		jobID, err = bc.cronHandler.service.AddRecurringJob(ctx, userID, chatID, parsed.CronExpr, parsed.Prompt)
	case "once":
		jobID, err = bc.cronHandler.service.AddOnceJob(ctx, userID, chatID, parsed.RunAt, parsed.Prompt)
	default:
		return "Não consegui determinar o tipo de agendamento.", nil
	}

	if err != nil {
		return "", fmt.Errorf("falha ao criar agendamento: %w", err)
	}

	switch parsed.Type {
	case "cron":
		return fmt.Sprintf("Agendamento recorrente criado (`%s`)\nSchedule: `%s`\nPrompt: %s", jobID, parsed.CronExpr, parsed.Prompt), nil
	default:
		return fmt.Sprintf("Agendamento pontual criado (`%s`)\nQuando: %s\nPrompt: %s", jobID, parsed.RunAt, parsed.Prompt), nil
	}
}

func (bc *BotController) cmdStatus(chatID int64) (string, error) {
	var lines []string
	lines = append(lines, "**Status do Sistema**\n")

	// Bridge status
	bridgeStatus := "desligado"
	if bc.bridge != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := bc.bridge.Ping(ctx); err == nil {
			bridgeStatus = "online"
		} else {
			bridgeStatus = "offline"
		}
	}
	lines = append(lines, fmt.Sprintf("Processador: **%s**", bridgeStatus))

	// Agents
	agentCount := 0
	if bc.agents != nil {
		agentCount = len(bc.agents.Agents())
	}
	lines = append(lines, fmt.Sprintf("Agents: **%d**", agentCount))

	// Cron jobs
	if bc.cronHandler != nil {
		ctx := context.Background()
		jobs, err := bc.cronHandler.service.ListJobs(ctx, chatID)
		if err == nil {
			active := 0
			for _, j := range jobs {
				if j.Active {
					active++
				}
			}
			lines = append(lines, fmt.Sprintf("Agendamentos ativos: **%d**", active))
		}
	}

	// Model
	if bc.config != nil && bc.config.DefaultModel != "" {
		lines = append(lines, fmt.Sprintf("Modelo padrão: **%s**", bc.config.DefaultModel))
	}

	// Session
	if sid, active := bc.sessions.GetWithState(chatID); sid != "" {
		state := "cold"
		if active {
			state = "warm"
		}
		lines = append(lines, fmt.Sprintf("Sessão: `%s` (%s)", sid[:8], state))
	} else {
		lines = append(lines, "Sessão: nenhuma")
	}

	return strings.Join(lines, "\n"), nil
}

func (bc *BotController) cmdListAgents() (string, error) {
	if bc.agents == nil {
		return "Nenhum agent configurado.", nil
	}

	all := bc.agents.Agents()
	if len(all) == 0 {
		return "Nenhum agent configurado.", nil
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("**Agents disponíveis** (%d)\n", len(all)))
	for _, a := range all {
		desc := a.Description
		if desc == "" {
			desc = "(sem descrição)"
		}
		line := fmt.Sprintf("- **%s**: %s", a.Name, desc)
		if a.Model != "" {
			line += fmt.Sprintf(" (modelo: %s)", a.Model)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), nil
}

func (bc *BotController) cmdListModels() (string, error) {
	if bc.config == nil || len(bc.config.Providers) == 0 {
		return "Nenhum provedor configurado.", nil
	}

	var lines []string
	lines = append(lines, "**Modelos e Provedores**\n")

	for name, prov := range bc.config.Providers {
		status := "ativo"
		if prov.APIKey == "" {
			status = "sem API key"
		}
		models := "—"
		if known, ok := providerModels[name]; ok {
			models = strings.Join(known, ", ")
		}
		lines = append(lines, fmt.Sprintf("- **%s** [%s]: %s", name, status, models))
	}

	if bc.config.DefaultModel != "" {
		lines = append(lines, fmt.Sprintf("\nModelo padrão: **%s**", bc.config.DefaultModel))
	}

	return strings.Join(lines, "\n"), nil
}

// providerModels maps known providers to their available models.
var providerModels = map[string][]string{
	"anthropic":  {"claude-sonnet-4-6", "claude-opus-4-6", "claude-haiku-4-5"},
	"google":     {"gemini-2.5-pro", "gemini-2.5-flash"},
	"kimi":       {"kimi-k2-thinking", "kimi-k2"},
	"openrouter": {"openrouter/auto"},
	"zai":        {"zai-1.5-flash"},
}
