package main

import (
	"context"
	"log"

	"github.com/kocar/aurelia/internal/agent"
)

type loopExecutorAdapter struct {
	loop *agent.Loop
}

func (a *loopExecutorAdapter) Execute(ctx context.Context, systemPrompt string, history []agent.Message, allowedTools []string) ([]agent.Message, string, error) {
	return a.loop.Run(ctx, systemPrompt, history, allowedTools)
}

func loadCronPromptConfig(canonicalPromptLoader interface {
	BuildPrompt(ctx context.Context, userID, conversationID string) (string, []string, error)
}) (string, []string) {
	if canonicalPromptLoader != nil {
		prompt, tools, err := canonicalPromptLoader.BuildPrompt(context.Background(), "", "")
		if err == nil {
			return prompt, tools
		}
		log.Printf("Warning: failed to load canonical prompt for cron runtime, using default prompt: %v", err)
	}
	return "Voce e o agente pessoal Aurelia. Execute a tarefa agendada com precisao e retorne um resumo objetivo do resultado.", nil
}
