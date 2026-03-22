package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/bridge"
)

func (bc *BotController) processInput(c telebot.Context, text string, parts [][]byte, requiresAudio bool) error {
	_ = parts

	if state, ok := bc.popPendingBootstrap(c.Sender().ID); ok {
		return bc.completeBootstrapProfile(c, state, text)
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

// routeAgent resolves which agent should handle the message, first by @name
// prefix, then by LLM classification if agents are configured.
func (bc *BotController) routeAgent(text string) *agents.Agent {
	agent := bc.agents.Route(text)

	if agent == nil && bc.agents != nil && len(bc.agents.Agents()) > 0 {
		classifyPrompt := bc.agents.ClassifyPrompt(text)
		if classifyPrompt != "" {
			classifyCtx, classifyCancel := context.WithTimeout(context.Background(), 15*time.Second)
			result, err := bc.bridge.ExecuteSync(classifyCtx, bridge.Request{
				Command: "query",
				Prompt:  classifyPrompt,
				Options: bridge.RequestOptions{
					Model:          bc.config.DefaultModel,
					SystemPrompt:   "You are a message classifier. Reply with only the agent name or 'none'.",
					MaxTurns:       1,
					PermissionMode: "bypassPermissions",
				},
			})
			classifyCancel()
			if err == nil && result.Type == "result" {
				name := strings.TrimSpace(strings.ToLower(result.Content))
				if name != "none" && name != "" {
					agent = bc.agents.Get(name)
				}
			}
		}
	}

	return agent
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
			DisabledTools: []string{
				"mcp__plugin_telegram_telegram__reply",
				"mcp__plugin_telegram_telegram__react",
				"mcp__plugin_telegram_telegram__edit_message",
			},
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

	// Resume previous session for conversation continuity
	if sessionID := bc.sessions.Get(chatID); sessionID != "" {
		req.Options.Resume = sessionID
	}

	// Apply chat-level cwd if no agent overrides it
	if req.Options.Cwd == "" {
		if chatCwd := bc.sessions.GetCwd(chatID); chatCwd != "" {
			req.Options.Cwd = chatCwd
		}
	}

	return req
}

// executeAsync runs the bridge execution in a goroutine with its own typing
// indicator and progress reporter. Errors are sent directly to the chat since
// the original handler has already returned.
func (bc *BotController) executeAsync(chatID int64, messageID int, req bridge.Request, userText string) {
	chat := &telebot.Chat{ID: chatID}

	// Start typing indicator
	stopTyping := startChatActionLoop(bc.bot, chat, telebot.Typing, 4*time.Second)
	defer stopTyping()

	// Progress reporter
	progress := newProgressReporter(bc.bot, chat)
	defer progress.Delete()

	// Execute via bridge
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	ch, err := bc.bridge.Execute(ctx, req)
	if err != nil {
		log.Printf("Bridge execute error: %v", err)
		_ = SendError(bc.bot, chat, "Falha ao conectar com o processador.")
		return
	}

	// Process events
	bc.processBridgeEventsAsync(chat, ch, progress, userText, messageID)
}

// processBridgeEventsAsync reads bridge events and sends responses to the
// Telegram chat. Unlike the old processBridgeEvents, it takes a *telebot.Chat
// instead of telebot.Context, since the handler has already returned.
func (bc *BotController) processBridgeEventsAsync(chat *telebot.Chat, ch <-chan bridge.Event, progress *progressReporter, userText string, messageID int) {
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

			finalText := strings.TrimSpace(assistantText.String())
			if finalText == "" {
				finalText = "(sem resposta)"
			}

			bc.saveToMemory(userText, finalText)
			_ = SendTextReply(bc.bot, chat, finalText, messageID)
			return

		case "error":
			errMsg := ev.Message
			if errMsg == "" {
				errMsg = ev.Content
			}
			if errMsg == "" {
				errMsg = "Erro desconhecido no processador."
			}
			log.Printf("Bridge error: %s", errMsg)
			_ = SendError(bc.bot, chat, errMsg)
			return

		default:
			log.Printf("Bridge event (ignored): %s", ev.Type)
		}
	}

	// Channel closed without terminal event
	finalText := strings.TrimSpace(assistantText.String())
	if finalText != "" {
		bc.saveToMemory(userText, finalText)
		_ = SendTextReply(bc.bot, chat, finalText, messageID)
	} else {
		_ = SendError(bc.bot, chat, "O processador encerrou sem resposta.")
	}
}

// buildSystemPrompt assembles the system prompt from persona, agent, cron/telegram instructions, and memory.
func (bc *BotController) buildSystemPrompt(userText string, agent *agents.Agent, chatID int64, messageID int) (string, error) {
	var sections []string

	// Persona prompt
	if bc.persona != nil {
		personaPrompt, err := bc.persona.BuildPrompt()
		if err != nil {
			log.Printf("Persona prompt error (non-fatal): %v", err)
		} else if personaPrompt != "" {
			sections = append(sections, personaPrompt)
		}
	}

	// Agent-specific prompt
	if agent != nil && agent.Prompt != "" {
		sections = append(sections, "# Agent Instructions\n\n"+agent.Prompt)
	}

	// Cron scheduling instructions
	sections = append(sections, bc.buildCronInstructions(chatID))

	// Telegram interaction instructions
	sections = append(sections, bc.buildTelegramInstructions(chatID, messageID))

	// Memory injection
	if bc.memory != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		memoryBlock, err := bc.memory.Inject(ctx, userText, bc.config.MemoryWindowSize)
		if err != nil {
			log.Printf("Memory injection error (non-fatal): %v", err)
		} else if memoryBlock != "" {
			sections = append(sections, memoryBlock)
		}
	}

	return strings.Join(sections, "\n\n"), nil
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

Available emojis: 👍 👎 ❤️ 🔥 👀 🎉 😂 🤔 💯 🎯 ✅ ❌

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

// saveToMemory stores the conversation exchange in semantic memory.
func (bc *BotController) saveToMemory(userText, response string) {
	if bc.memory == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	content := fmt.Sprintf("User: %s\nAssistant: %s", userText, truncate(response, 500))
	if err := bc.memory.Save(ctx, content, "conversation", "telegram"); err != nil {
		log.Printf("Memory save error (non-fatal): %v", err)
	}
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
