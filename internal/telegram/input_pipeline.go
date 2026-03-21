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

	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, 4*time.Second)
	defer stopTyping()

	// 1. Route to agent via registry
	agent := bc.agents.Route(text)

	// If no @name match and agents exist, try LLM classification
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

	// 2. Build system prompt: persona + agent prompt + memory
	systemPrompt, err := bc.buildSystemPrompt(userText, agent)
	if err != nil {
		log.Printf("Failed to build system prompt: %v", err)
		return SendError(bc.bot, c.Chat(), "Falha ao montar o prompt de sistema.")
	}

	// 3. Build bridge request
	req := bridge.Request{
		Command: "query",
		Prompt:  userText,
		Options: bridge.RequestOptions{
			Model:          bc.config.DefaultModel,
			SystemPrompt:   systemPrompt,
			MaxTurns:       bc.config.MaxIterations,
			PermissionMode: "bypassPermissions",
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

	// Resume previous session if available
	if sessionID := bc.sessions.Get(c.Chat().ID); sessionID != "" {
		req.Options.Resume = sessionID
	}

	// 4. Execute via bridge (streaming)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	ch, err := bc.bridge.Execute(ctx, req)
	if err != nil {
		log.Printf("Bridge execute error: %v", err)
		return SendError(bc.bot, c.Chat(), "Falha ao conectar com o processador.")
	}

	// 5. Process events
	return bc.processBridgeEvents(c, ch, userText)
}

// buildSystemPrompt assembles the system prompt from persona, agent, and memory.
func (bc *BotController) buildSystemPrompt(userText string, agent *agents.Agent) (string, error) {
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

// processBridgeEvents reads bridge events and sends responses to the Telegram chat.
func (bc *BotController) processBridgeEvents(c telebot.Context, ch <-chan bridge.Event, userText string) error {
	progress := newProgressReporter(bc.bot, c.Chat())
	defer progress.Delete()

	var assistantText strings.Builder

	for ev := range ch {
		switch ev.Type {
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

			// Save conversation to memory
			bc.saveToMemory(userText, finalText)

			return SendText(bc.bot, c.Chat(), finalText)

		case "error":
			errMsg := ev.Message
			if errMsg == "" {
				errMsg = ev.Content
			}
			if errMsg == "" {
				errMsg = "Erro desconhecido no processador."
			}
			log.Printf("Bridge error: %s", errMsg)
			return SendError(bc.bot, c.Chat(), errMsg)

		case "system":
			if ev.SessionID != "" {
				bc.sessions.Set(c.Chat().ID, ev.SessionID)
			}

		default:
			log.Printf("Bridge event (ignored): %s", ev.Type)
		}
	}

	// Channel closed without terminal event
	finalText := strings.TrimSpace(assistantText.String())
	if finalText != "" {
		bc.saveToMemory(userText, finalText)
		return SendText(bc.bot, c.Chat(), finalText)
	}

	return SendError(bc.bot, c.Chat(), "O processador encerrou sem resposta.")
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
