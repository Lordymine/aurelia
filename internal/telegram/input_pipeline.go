package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/bridge"
)

func (bc *BotController) processInput(c telebot.Context, text string) error {
	if state, ok := bc.popPendingBootstrap(c.Sender().ID); ok {
		switch state.Step {
		case bootstrapStepAssistant:
			return bc.completeBootstrapAssistant(c, state, text)
		default:
			return bc.completeBootstrapProfile(c, state, text)
		}
	}

	// 0. Command layer — intercept system commands before LLM
	if cmd := MatchCommand(text); cmd != nil {
		return bc.handleCommand(c, cmd)
	}

	// 1. Route to agent (sync — fast)
	agent := bc.routeAgent(text)

	// Strip @agent prefix from user text if agent was routed
	userText := text
	if agent != nil {
		if idx := strings.IndexByte(text[1:], ' '); idx != -1 {
			userText = strings.TrimSpace(text[idx+2:])
		} else {
			userText = ""
		}
	}

	if userText == "" {
		userText = text
	}

	// 2. Build system prompt (sync — fast)
	chatID := c.Chat().ID
	messageID := c.Message().ID
	systemPrompt, err := bc.buildSystemPrompt(userText, agent, chatID, messageID)
	if err != nil {
		log.Printf("Failed to build system prompt: %v", err)
		return SendError(bc.bot, c.Chat(), "Falha ao montar o prompt de sistema.")
	}

	// 3. Build bridge request (sync)
	req := bc.buildBridgeRequest(userText, systemPrompt, agent, chatID)

	// 4. Launch async execution — don't block the handler
	go bc.executeAsync(chatID, messageID, req, userText)

	return nil
}

const (
	classifyTimeout         = 15 * time.Second
	typingIndicatorInterval = 4 * time.Second
	bridgeExecutionTimeout  = 10 * time.Minute
)

// routeAgent resolves which agent should handle the message, first by @name
// prefix, then by LLM classification if agents are configured.
func (bc *BotController) routeAgent(text string) *agents.Agent {
	agent := bc.agents.Route(text)
	if agent != nil {
		return agent
	}
	if bc.agents == nil {
		return nil
	}
	classifyCtx, classifyCancel := context.WithTimeout(context.Background(), classifyTimeout)
	defer classifyCancel()
	return bc.agents.Classify(classifyCtx, text, bc.classifyFunc())
}

func (bc *BotController) classifyFunc() agents.ClassifyFunc {
	return func(ctx context.Context, system, prompt string) (string, error) {
		result, err := bc.bridge.ExecuteSync(ctx, bridge.Request{
			Command: "query",
			Prompt:  prompt,
			Options: bridge.RequestOptions{
				Model:          bc.config.DefaultModel,
				SystemPrompt:   system,
				MaxTurns:       1,
				PermissionMode: "bypassPermissions",
			},
		})
		if err != nil {
			return "", err
		}
		return result.Content, nil
	}
}

// buildBridgeRequest assembles the bridge.Request with agent overrides, session
// resume, and working directory.
func (bc *BotController) buildBridgeRequest(userText, systemPrompt string, agent *agents.Agent, chatID int64) bridge.Request {
	req := bridge.Request{
		Command: "query",
		Prompt:  userText,
		Options: bridge.RequestOptions{
			Model:          bc.config.DefaultModel,
			SystemPrompt:   systemPrompt,
			MaxTurns:       bc.config.MaxIterations,
			PermissionMode: "bypassPermissions",
			DisabledTools: bridge.TelegramPluginTools,
		},
	}

	if agent != nil {
		if agent.Model != "" {
			req.Options.Model = agent.Model
		}
		if agent.Cwd != "" {
			req.Options.Cwd = agent.Cwd
		}
		if len(agent.MCPServers) > 0 {
			req.Options.MCPServers = agent.MCPServers
		}
		if len(agent.AllowedTools) > 0 {
			req.Options.AllowedTools = agent.AllowedTools
		}
	}

	// Pass all agents to SDK for native delegation
	if sdkAgents := agents.BuildSDKAgents(bc.agents); sdkAgents != nil {
		req.Options.Agents = sdkAgents
	}

	// Continue warm sessions (same process), resume cold ones (restored from disk)
	if sessionID, active := bc.sessions.GetWithState(chatID); sessionID != "" {
		if active {
			req.Options.Continue = true
			log.Printf("session: chat=%d mode=continue", chatID)
		} else {
			req.Options.Resume = sessionID
			log.Printf("session: chat=%d mode=resume sid=%s", chatID, sessionID[:8])
		}
	} else {
		log.Printf("session: chat=%d mode=new", chatID)
	}

	// Apply chat-level cwd if no agent overrides it
	if req.Options.Cwd == "" {
		if chatCwd := bc.sessions.GetCwd(chatID); chatCwd != "" {
			req.Options.Cwd = chatCwd
		}
	}

	return req
}

// bridgeOutcome indicates how processBridgeEventsAsync terminated.
type bridgeOutcome int

const (
	outcomeSuccess      bridgeOutcome = iota // terminal "result" event
	outcomeLLMError                          // terminal "error" event
	outcomeProcessDeath                      // channel closed without terminal event
)

// bridgeFailureTracker tracks consecutive bridge failures to implement cooldown.
type bridgeFailureTracker struct {
	mu       sync.Mutex
	failures []time.Time // timestamps of recent failures
}

const (
	failureWindowMax   = 3                // max failures before cooldown
	failureWindowDur   = 1 * time.Minute  // window to count failures
	cooldownDuration   = 30 * time.Second // cooldown period after max failures
)

// record adds a failure timestamp and returns true if in cooldown.
func (t *bridgeFailureTracker) record() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	t.failures = append(t.failures, now)

	// Trim failures outside the window
	cutoff := now.Add(-failureWindowDur)
	start := 0
	for start < len(t.failures) && t.failures[start].Before(cutoff) {
		start++
	}
	t.failures = t.failures[start:]

	return len(t.failures) >= failureWindowMax
}

// inCooldown returns true if we're in cooldown (recent failures >= max).
func (t *bridgeFailureTracker) inCooldown() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.failures) < failureWindowMax {
		return false
	}

	// In cooldown if last failure was within cooldown duration
	last := t.failures[len(t.failures)-1]
	return time.Since(last) < cooldownDuration
}

// reset clears the failure history after a successful execution.
func (t *bridgeFailureTracker) reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.failures = t.failures[:0]
}

// executeAsync runs the bridge execution in a goroutine with its own typing
// indicator and progress reporter. Errors are sent directly to the chat since
// the original handler has already returned.
func (bc *BotController) executeAsync(chatID int64, messageID int, req bridge.Request, userText string) {
	chat := &telebot.Chat{ID: chatID}

	// Start typing indicator
	stopTyping := startChatActionLoop(bc.bot, chat, telebot.Typing, typingIndicatorInterval)
	defer stopTyping()

	// Progress reporter
	progress := newProgressReporter(bc.bot, chat)
	defer progress.Delete()

	// Execute via bridge
	ctx, cancel := context.WithTimeout(context.Background(), bridgeExecutionTimeout)
	defer cancel()

	ch, err := bc.bridge.Execute(ctx, req)
	if err != nil {
		log.Printf("Bridge execute error: %v", err)
		if err := SendError(bc.bot, chat, "Falha ao conectar com o processador."); err != nil {
			log.Printf("Failed to send error to chat %d: %v", chat.ID, err)
		}
		return
	}

	// Process events — first attempt
	outcome := bc.processBridgeEventsAsync(chat, ch, progress, userText, messageID)

	if outcome == outcomeSuccess {
		bc.bridgeFailures.reset()
		return
	}
	if outcome != outcomeProcessDeath {
		return
	}

	// --- RETRY PATH: bridge died mid-request ---
	bc.bridgeFailures.record()
	log.Printf("bridge: process died mid-request, retrying for chat=%d", chatID)

	// P3: Check cooldown before retrying
	if bc.bridgeFailures.inCooldown() {
		log.Printf("bridge: in cooldown, skipping retry for chat=%d", chatID)
		_ = SendError(bc.bot, chat, "Processador temporariamente indisponível. Tente novamente em alguns segundos.")
		return
	}

	// P2: Send reconnection feedback
	var reconnectMsg *telebot.Message
	reconnectMsg, _ = bc.bot.Send(chat, "⚡ Reconectando...")

	retryReq := req
	retryReq.Options.Continue = false
	retryReq.RequestID = ""
	if sid := bc.sessions.Get(chatID); sid != "" {
		retryReq.Options.Resume = sid
		log.Printf("bridge: retry with resume sid=%s", sid[:8])
	}

	ch, err = bc.bridge.Execute(ctx, retryReq)
	bc.deleteMessage(reconnectMsg) // Bridge restarted (or failed) — remove feedback immediately
	if err != nil {
		log.Printf("bridge: retry failed for chat=%d: %v", chatID, err)
		_ = SendError(bc.bot, chat, "Processador reiniciado mas não conseguiu completar. Tente novamente.")
		return
	}

	// Second attempt — no more retries
	outcome = bc.processBridgeEventsAsync(chat, ch, progress, userText, messageID)

	if outcome == outcomeSuccess {
		bc.bridgeFailures.reset()
	} else if outcome == outcomeProcessDeath {
		bc.bridgeFailures.record()
		_ = SendError(bc.bot, chat, "Processador reiniciado mas não conseguiu completar. Tente novamente.")
	}
}

// deleteMessage removes a Telegram message if it exists. Used to clean up
// temporary feedback messages like "Reconectando...".
func (bc *BotController) deleteMessage(msg *telebot.Message) {
	if msg != nil && bc.bot != nil {
		if err := bc.bot.Delete(msg); err != nil {
			log.Printf("Failed to delete reconnect message: %v", err)
		}
	}
}

// processBridgeEventsAsync reads bridge events and sends responses to the
// Telegram chat. Returns the outcome so the caller can decide whether to retry.
func (bc *BotController) processBridgeEventsAsync(chat *telebot.Chat, ch <-chan bridge.Event, progress *progressReporter, userText string, messageID int) bridgeOutcome {
	var assistantText strings.Builder

	for ev := range ch {
		switch ev.Type {
		case "system":
			if ev.SessionID != "" {
				bc.sessions.Set(chat.ID, ev.SessionID)
			}

		case "tool_use":
			toolName := ev.Name
			if toolName == "" {
				toolName = "tool"
			}
			progress.ReportTool(toolName)

		case "assistant":
			content := ev.Text
			if content == "" {
				content = ev.Content
			}
			assistantText.WriteString(content)

		case "result":
			content := ev.Text
			if content == "" {
				content = ev.Content
			}
			if content != "" {
				assistantText.Reset()
				assistantText.WriteString(content)
			}

			if ev.CostUSD > 0 || ev.NumTurns > 0 {
				if bc.tracker.RecordUsage(chat.ID, ev.NumTurns, ev.CostUSD, bc.config.MaxSessionTokens) {
					log.Printf("session auto-reset: chat=%d threshold=%d", chat.ID, bc.config.MaxSessionTokens)
					bc.sessions.Clear(chat.ID)
					bc.tracker.Clear(chat.ID)
				} else {
					usage := bc.tracker.Get(chat.ID)
					log.Printf("session usage: chat=%d %s", chat.ID, usage)
				}
			}

			finalText := strings.TrimSpace(assistantText.String())
			if finalText == "" {
				finalText = "(sem resposta)"
			}

			if err := SendTextReply(bc.bot, chat, finalText, messageID); err != nil {
				log.Printf("Failed to send reply to chat %d: %v", chat.ID, err)
			}
			return outcomeSuccess

		case "error":
			errMsg := ev.Message
			if errMsg == "" {
				errMsg = ev.Content
			}
			if errMsg == "" {
				errMsg = "Erro desconhecido no processador."
			}
			log.Printf("Bridge error: %s", errMsg)
			if err := SendError(bc.bot, chat, errMsg); err != nil {
				log.Printf("Failed to send error to chat %d: %v", chat.ID, err)
			}
			return outcomeLLMError

		default:
			log.Printf("Bridge event (ignored): %s", ev.Type)
		}
	}

	// Channel closed without terminal event — process died
	return outcomeProcessDeath
}

// buildSystemPrompt assembles the system prompt from persona, agent, cron/telegram instructions, and memory.
func (bc *BotController) buildSystemPrompt(userText string, agent *agents.Agent, chatID int64, messageID int) (string, error) {
	var sections []string
	var personaLen, agentLen, cronLen, telegramLen int

	// Persona prompt
	if bc.persona != nil {
		personaPrompt, err := bc.persona.BuildPrompt()
		if err != nil {
			log.Printf("Persona prompt error (non-fatal): %v", err)
		} else if personaPrompt != "" {
			personaLen = len(personaPrompt)
			sections = append(sections, personaPrompt)
		}
	}

	// Agent-specific prompt
	if agent != nil && agent.Prompt != "" {
		agentSection := "# Agent Instructions\n\n" + agent.Prompt
		agentLen = len(agentSection)
		sections = append(sections, agentSection)
	}

	// Cron scheduling instructions
	cronSection := bc.buildCronInstructions(chatID)
	cronLen = len(cronSection)
	sections = append(sections, cronSection)

	// Telegram interaction instructions
	telegramSection := bc.buildTelegramInstructions(chatID, messageID)
	telegramLen = len(telegramSection)
	sections = append(sections, telegramSection)

	result := strings.Join(sections, "\n\n")
	log.Printf("system prompt breakdown: persona=%d agent=%d cron=%d telegram=%d total=%d chars",
		personaLen, agentLen, cronLen, telegramLen, len(result))

	return result, nil
}

// buildTelegramInstructions returns instructions for interacting with the Telegram chat.
func (bc *BotController) buildTelegramInstructions(chatID int64, messageID int) string {
	bin := "aurelia"
	if bc.exePath != "" {
		bin = bc.exePath
	}

	return fmt.Sprintf(`## Telegram Context

You ARE the Telegram bot. The user is talking to you via Telegram chat %d.
The current message ID is %d.

You can interact with the chat using the Aurelia CLI via Bash:

React to a message with emoji:
`+"`%s telegram react %d %d <emoji>`"+`

Available emojis: 👍 👎 ❤️ 🔥 🎉 🤩 😱 😁 😢 💩 🤮 🥰 🤯 🤔 🤬 👏 🙏 👌 😍 💯 ⚡️ 🏆

Use reactions naturally and contextually — react when it adds to the conversation, not on every message.
DO NOT use the Telegram MCP plugin for reactions or replies — use the Aurelia CLI above.`,
		chatID, messageID, bin, chatID, messageID)
}

// buildCronInstructions returns the system prompt section that teaches the agent
// how to create and manage cron jobs via the aurelia CLI.
func (bc *BotController) buildCronInstructions(chatID int64) string {
	bin := "aurelia"
	if bc.exePath != "" {
		bin = bc.exePath
	}
	chatFlag := fmt.Sprintf("--chat-id %d", chatID)

	return fmt.Sprintf(`## Scheduling Tasks — MANDATORY

CRITICAL: You MUST use the Aurelia cron CLI for ALL scheduling. NEVER use your internal scheduling tools — they die with the session. The Aurelia cron is persistent and delivers results to Telegram automatically.

Recurring schedule:
`+"```bash\n%s cron add \"<cron-expression>\" \"<prompt>\" %s\n```"+`

One-time schedule:
`+"```bash\n%s cron once \"<ISO-timestamp>\" \"<prompt>\" %s\n```"+`

List schedules:
`+"```bash\n%s cron list %s\n```"+`

Delete: `+"`%s cron del <job-id>`"+`
Pause: `+"`%s cron pause <job-id>`"+`
Resume: `+"`%s cron resume <job-id>`"+`

Cron expressions: "30 8 * * *" = daily 8:30 | "0 9 * * 1" = Monday 9:00 | "0 */2 * * *" = every 2h

The --chat-id flag is REQUIRED — it ensures results are delivered to this Telegram chat.

CRITICAL RULES FOR CRON PROMPTS:
1. The prompt is an INSTRUCTION, not content. It tells the agent WHAT TO DO when the job fires.
2. It executes in an isolated session with NO conversation history.
3. The agent will execute the prompt and its text output is delivered directly to Telegram.
4. NEVER paste content/data into the prompt. Write an ACTION instruction.
5. Bad: "Envie esta newsletter: [conteúdo colado aqui]"
6. Good: "Pesquise as principais notícias de tech e IA da última semana usando WebSearch. Para cada notícia inclua título, resumo de 1 linha e link. Formate como newsletter com emojis e seções: IA, Hardware, Dev, Segurança, Regulação. Encerre com um insight da semana e hashtags."`,
		bin, chatFlag,
		bin, chatFlag,
		bin, chatFlag,
		bin,
		bin, bin,
	)
}

