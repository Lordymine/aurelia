package skill

import (
	"context"

	"github.com/kocar/aurelia/internal/agent"
)

// Executor manages the runtime context injection limits
type Executor struct {
	loop *agent.Loop
}

// NewExecutor constructs an executor
func NewExecutor(loop *agent.Loop) *Executor {
	return &Executor{loop: loop}
}

// Execute wraps the ReAct loop injecting the targeted skill file cleanly without polluting history
func (e *Executor) Execute(ctx context.Context, baseSystemPrompt string, skill *Skill, history []agent.Message, allowedTools []string) ([]agent.Message, string, error) {
	// Dynamically inject the huge skill.md markdown
	finalSystemPrompt := baseSystemPrompt
	if skill != nil {
		finalSystemPrompt += "\n\n# ACTIVE SKILL CONTEXT:\n" + skill.Content
	}

	return e.loop.Run(ctx, finalSystemPrompt, history, allowedTools)
}


