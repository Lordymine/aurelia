package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/skill"
	"github.com/kocar/aurelia/pkg/llm"
	"gopkg.in/telebot.v3"
)

const defaultSystemPrompt = "Voce e o agente pessoal Aurelia. Siga as ordens do usuario com precisao. Retorne markdown estruturado. Se editar codigo, valide com run_command antes de concluir."

type inputSession struct {
	senderID string
	convID   string
	ctx      context.Context
	text     string
	message  agent.Message
}

func (bc *BotController) processInput(c telebot.Context, text string, parts []agent.ContentPart, requiresAudio bool) error {
	text = strings.ReplaceAll(text, "\x00", "")

	if state, ok := bc.popPendingBootstrap(c.Sender().ID); ok {
		return bc.completeBootstrapProfile(c, state, text)
	}

	session := newInputSession(c, text, bc.attachRecentMediaIfRelevant(c, text, parts))
	if handled, err := bc.handleMemoryCommand(c, session); handled {
		return err
	}
	if handled, err := bc.handleOpsCommand(c, session); handled {
		return err
	}
	bc.storeRecentMedia(session)

	if err := bc.persistIncomingContext(session, c.Sender().ID); err != nil {
		log.Printf("Input persistence warning: %v\n", err)
	}
	if err := bc.manageConversationContext(session.ctx, session.convID); err != nil {
		log.Printf("Context management warning: %v\n", err)
	}

	activeSkill, history, systemPrompt, allowedTools, err := bc.prepareExecution(session)
	if err != nil {
		_ = SendError(bc.bot, c.Chat(), err.Error())
		return nil
	}

	finalAnswer, err := bc.executeConversation(c, session, activeSkill, history, systemPrompt, allowedTools)
	if err != nil {
		var visionErr llm.VisionUnsupportedError
		if errors.As(err, &visionErr) {
			return SendContextText(c, visionErr.Error())
		}
		_ = SendError(bc.bot, c.Chat(), err.Error())
		return nil
	}

	if err := bc.persistAssistantAnswer(session, finalAnswer); err != nil {
		log.Printf("Assistant persistence warning: %v\n", err)
	}
	return bc.deliverFinalAnswer(c, finalAnswer, requiresAudio)
}

func newInputSession(c telebot.Context, text string, parts []agent.ContentPart) inputSession {
	senderID := fmt.Sprintf("%d", c.Sender().ID)
	convID := senderID
	ctx := agent.WithRunContext(agent.WithTeamContext(context.Background(), convID, senderID), uuid.NewString())
	message := agent.Message{Role: "user", Content: text, Parts: append([]agent.ContentPart(nil), parts...)}
	if len(message.Parts) == 0 {
		message.Parts = []agent.ContentPart{{Type: agent.ContentPartText, Text: text}}
	}
	return inputSession{senderID: senderID, convID: convID, ctx: ctx, text: text, message: message}
}

func (bc *BotController) attachRecentMediaIfRelevant(c telebot.Context, text string, parts []agent.ContentPart) []agent.ContentPart {
	if len(parts) != 0 {
		return parts
	}
	if !shouldReuseRecentMedia(text) {
		return nil
	}

	senderID := fmt.Sprintf("%d", c.Sender().ID)
	media, ok := bc.loadRecentMedia(senderID)
	if !ok {
		return nil
	}

	attached := make([]agent.ContentPart, 0, len(media.parts)+1)
	attached = append(attached, agent.ContentPart{Type: agent.ContentPartText, Text: text})
	attached = append(attached, media.parts...)
	return attached
}

func shouldReuseRecentMedia(text string) bool {
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return false
	}
	tokens := []string{"imagem", "foto", "print", "screenshot", "anexo", "arquivo", "pdf", "planilha", "excel", "word", "doc"}
	for _, token := range tokens {
		if strings.Contains(text, token) {
			return true
		}
	}
	return false
}

func (bc *BotController) storeRecentMedia(session inputSession) {
	mediaParts := extractMediaParts(session.message.Parts)
	if len(mediaParts) == 0 {
		return
	}

	bc.mediaMu.Lock()
	defer bc.mediaMu.Unlock()
	bc.recentMedia[session.convID] = recentMedia{
		parts:     mediaParts,
		updatedAt: time.Now(),
	}
}

func (bc *BotController) loadRecentMedia(conversationID string) (recentMedia, bool) {
	bc.mediaMu.Lock()
	defer bc.mediaMu.Unlock()

	media, ok := bc.recentMedia[conversationID]
	if !ok {
		return recentMedia{}, false
	}
	if time.Since(media.updatedAt) > 3*time.Minute {
		delete(bc.recentMedia, conversationID)
		return recentMedia{}, false
	}
	media.parts = append([]agent.ContentPart(nil), media.parts...)
	return media, true
}

func extractMediaParts(parts []agent.ContentPart) []agent.ContentPart {
	media := make([]agent.ContentPart, 0, len(parts))
	for _, part := range parts {
		if part.Type != agent.ContentPartImage || len(part.Data) == 0 {
			continue
		}
		media = append(media, agent.ContentPart{
			Type:     part.Type,
			MIMEType: part.MIMEType,
			Data:     append([]byte(nil), part.Data...),
		})
	}
	return media
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

func (bc *BotController) handleOpsCommand(c telebot.Context, session inputSession) (bool, error) {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(session.text)), "/ops") {
		return false, nil
	}

	reply, err := NewOpsCommandHandler(bc.ops).HandleText(context.Background(), session.text)
	if err != nil {
		_ = SendError(bc.bot, c.Chat(), err.Error())
		return true, nil
	}
	return true, SendContextText(c, reply)
}

func (bc *BotController) persistIncomingContext(session inputSession, senderUserID int64) error {
	if err := bc.memory.EnsureConversation(session.ctx, session.convID, senderUserID, bc.config.LLMProvider); err != nil {
		return err
	}
	var persistErr error
	if bc.canonical != nil {
		persistErr = errors.Join(persistErr, bc.canonical.ApplyFactsAndSync(session.ctx, session.senderID, extractFactsFromConversation(session.text, session.senderID)))
	}
	persistErr = errors.Join(persistErr, persistConversationNote(session.ctx, bc.memory, session.convID, session.text))
	persistErr = errors.Join(persistErr, bc.memory.AddMessage(session.ctx, session.convID, "user", session.persistedContent()))
	persistErr = errors.Join(persistErr, bc.memory.AddArchiveEntry(session.ctx, memory.ArchiveEntry{
		ConversationID: session.convID,
		SessionID:      session.convID,
		Role:           "user",
		Content:        session.persistedContent(),
		MessageType:    "chat",
	}))
	return persistErr
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

	agentHistory := make([]agent.Message, 0, len(history)+1)
	summary, err := bc.summaryPrefix(session.ctx, session.convID)
	if err != nil {
		return nil, err
	}
	if summary != nil {
		agentHistory = append(agentHistory, *summary)
	}
	for _, m := range history {
		agentHistory = append(agentHistory, agent.Message{Role: m.Role, Content: m.Content})
	}
	if len(agentHistory) != 0 && agentHistory[len(agentHistory)-1].Role == "user" {
		agentHistory[len(agentHistory)-1] = session.message
	} else {
		agentHistory = append(agentHistory, session.message)
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

func (bc *BotController) manageConversationContext(ctx context.Context, conversationID string) error {
	if bc.contextPolicy == nil {
		return nil
	}
	return bc.contextPolicy.ManageConversation(ctx, bc.config.LLMProvider, bc.config.LLMModel, conversationID, bc.config.MemoryWindowSize)
}

func (bc *BotController) summaryPrefix(ctx context.Context, conversationID string) (*agent.Message, error) {
	if bc.contextPolicy == nil {
		return nil, nil
	}
	return bc.contextPolicy.SummaryPrefix(ctx, conversationID)
}

func (s inputSession) persistedContent() string {
	if !s.message.HasMedia() {
		return s.text
	}
	imageCount := 0
	for _, part := range s.message.Parts {
		if part.Type == agent.ContentPartImage {
			imageCount++
		}
	}
	if imageCount == 0 {
		return s.text
	}
	text := s.text
	if strings.TrimSpace(text) == "" {
		text = "Imagem recebida."
	}
	return fmt.Sprintf("%s [imagem anexada: %d]", text, imageCount)
}

func (bc *BotController) resolveExecutionPrompt(session inputSession) (string, []string) {
	if bc.canonical == nil {
		return defaultSystemPrompt, agent.ResolveAllowedToolsForQuery(session.text, nil)
	}

	prompt, tools, err := bc.canonical.BuildPromptForQuery(session.ctx, session.senderID, session.convID, session.text)
	if err != nil {
		log.Printf("Persona files not found or invalid. Using default prompt. Error: %v\n", err)
		return defaultSystemPrompt, agent.ResolveAllowedToolsForQuery(session.text, nil)
	}
	return prompt, agent.ResolveAllowedToolsForQuery(session.text, tools)
}

func (bc *BotController) executeConversation(c telebot.Context, session inputSession, activeSkill *skill.Skill, history []agent.Message, systemPrompt string, allowedTools []string) (string, error) {
	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, 4*time.Second)
	defer stopTyping()

	_, finalAnswer, err := bc.executor.Execute(session.ctx, systemPrompt, activeSkill, history, allowedTools)
	return finalAnswer, err
}

func (bc *BotController) persistAssistantAnswer(session inputSession, finalAnswer string) error {
	return errors.Join(
		bc.memory.AddMessage(session.ctx, session.convID, "assistant", finalAnswer),
		bc.memory.AddArchiveEntry(session.ctx, memory.ArchiveEntry{
			ConversationID: session.convID,
			SessionID:      session.convID,
			Role:           "assistant",
			Content:        finalAnswer,
			MessageType:    "chat",
		}),
	)
}

func (bc *BotController) deliverFinalAnswer(c telebot.Context, finalAnswer string, requiresAudio bool) error {
	if requiresAudio {
		return SendAudio(bc.bot, c.Chat(), finalAnswer)
	}
	return SendText(bc.bot, c.Chat(), finalAnswer)
}
