package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/skill"
	"gopkg.in/telebot.v3"
)

const defaultSystemPrompt = "Voce e o agente pessoal Aurelia. Siga as ordens do usuario com precisao. Retorne markdown estruturado. Se editar codigo, valide com run_command antes de concluir."

type inputSession struct {
	senderID string
	convID   string
	ctx      context.Context
	text     string
}

func (bc *BotController) processInput(c telebot.Context, text string, requiresAudio bool) error {
	text = strings.ReplaceAll(text, "\x00", "")

	if state, ok := bc.popPendingBootstrap(c.Sender().ID); ok {
		return bc.completeBootstrapProfile(c, state, text)
	}

	session := newInputSession(c, text)
	if handled, err := bc.handleMemoryCommand(c, session); handled {
		return err
	}

	if err := bc.persistIncomingContext(session, c.Sender().ID); err != nil {
		log.Printf("Input persistence warning: %v\n", err)
	}

	activeSkill, history, systemPrompt, allowedTools, err := bc.prepareExecution(session)
	if err != nil {
		_ = SendError(bc.bot, c.Chat(), err.Error())
		return nil
	}

	finalAnswer, err := bc.executeConversation(c, session, activeSkill, history, systemPrompt, allowedTools)
	if err != nil {
		_ = SendError(bc.bot, c.Chat(), err.Error())
		return nil
	}

	bc.persistAssistantAnswer(session, finalAnswer)
	return bc.deliverFinalAnswer(c, finalAnswer, requiresAudio)
}

func newInputSession(c telebot.Context, text string) inputSession {
	senderID := fmt.Sprintf("%d", c.Sender().ID)
	convID := senderID
	ctx := agent.WithRunContext(agent.WithTeamContext(context.Background(), convID, senderID), uuid.NewString())
	return inputSession{senderID: senderID, convID: convID, ctx: ctx, text: text}
}

func (bc *BotController) handleMemoryCommand(c telebot.Context, session inputSession) (bool, error) {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(session.text)), "/memory") {
		return false, nil
	}
	if bc.canonical == nil {
		return true, SendContextText(c, "Memoria longa indisponivel neste runtime.")
	}

	reply, err := NewMemoryCommandHandler(bc.canonical).HandleText(context.Background(), session.senderID, session.convID, session.text)
	if err != nil {
		_ = SendError(bc.bot, c.Chat(), err.Error())
		return true, nil
	}
	return true, SendContextText(c, reply)
}

func (bc *BotController) persistIncomingContext(session inputSession, senderUserID int64) error {
	if err := bc.memory.EnsureConversation(session.ctx, session.convID, senderUserID, "gemini"); err != nil {
		return err
	}
	if bc.canonical != nil {
		_ = bc.canonical.ApplyFactsAndSync(session.ctx, session.senderID, extractFactsFromConversation(session.text, session.senderID))
	}
	_ = persistConversationNote(session.ctx, bc.memory, session.convID, session.text)
	_ = bc.memory.AddMessage(session.ctx, session.convID, "user", session.text)
	_ = bc.memory.AddArchiveEntry(session.ctx, memory.ArchiveEntry{
		ConversationID: session.convID,
		SessionID:      session.convID,
		Role:           "user",
		Content:        session.text,
		MessageType:    "chat",
	})
	return nil
}

func (bc *BotController) prepareExecution(session inputSession) (*skill.Skill, []agent.Message, string, []string, error) {
	skills, _ := bc.loader.LoadAll()
	targetSkill, err := bc.router.Route(session.ctx, session.text, skills)
	if err != nil {
		log.Println("Router non-fatal error:", err)
	}

	history, err := bc.buildAgentHistory(session)
	if err != nil {
		return nil, nil, "", nil, err
	}
	activeSkill := resolveActiveSkill(skills, targetSkill)
	systemPrompt, allowedTools := bc.resolveExecutionPrompt(session)
	return activeSkill, history, systemPrompt, allowedTools, nil
}

func (bc *BotController) buildAgentHistory(session inputSession) ([]agent.Message, error) {
	history, err := bc.memory.GetRecentMessages(session.ctx, session.convID)
	if err != nil {
		return nil, err
	}

	agentHistory := make([]agent.Message, 0, len(history))
	for _, m := range history {
		agentHistory = append(agentHistory, agent.Message{Role: m.Role, Content: m.Content})
	}
	return agentHistory, nil
}

func resolveActiveSkill(skills map[string]skill.Skill, targetSkill string) *skill.Skill {
	if targetSkill == "" {
		return nil
	}
	log.Printf("Router decided skill: %s\n", targetSkill)
	s, ok := skills[targetSkill]
	if !ok {
		return nil
	}
	return &s
}

func (bc *BotController) resolveExecutionPrompt(session inputSession) (string, []string) {
	if bc.canonical == nil {
		return defaultSystemPrompt, nil
	}

	prompt, tools, err := bc.canonical.BuildPromptForQuery(session.ctx, session.senderID, session.convID, session.text)
	if err != nil {
		log.Printf("Persona files not found or invalid. Using default prompt. Error: %v\n", err)
		return defaultSystemPrompt, nil
	}
	return prompt, tools
}

func (bc *BotController) executeConversation(c telebot.Context, session inputSession, activeSkill *skill.Skill, history []agent.Message, systemPrompt string, allowedTools []string) (string, error) {
	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, 4*time.Second)
	defer stopTyping()

	_, finalAnswer, err := bc.executor.Execute(session.ctx, systemPrompt, activeSkill, history, allowedTools)
	return finalAnswer, err
}

func (bc *BotController) persistAssistantAnswer(session inputSession, finalAnswer string) {
	_ = bc.memory.AddMessage(session.ctx, session.convID, "assistant", finalAnswer)
	_ = bc.memory.AddArchiveEntry(session.ctx, memory.ArchiveEntry{
		ConversationID: session.convID,
		SessionID:      session.convID,
		Role:           "assistant",
		Content:        finalAnswer,
		MessageType:    "chat",
	})
}

func (bc *BotController) deliverFinalAnswer(c telebot.Context, finalAnswer string, requiresAudio bool) error {
	if requiresAudio {
		return SendAudio(bc.bot, c.Chat(), finalAnswer)
	}
	return SendText(bc.bot, c.Chat(), finalAnswer)
}


