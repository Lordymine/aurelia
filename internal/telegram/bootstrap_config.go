package telegram

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/telebot.v3"
)

const (
	bootstrapStepAssistant = "assistant"
	bootstrapStepProfile   = "profile"
)

type bootstrapState struct {
	Choice string
	Step   string
}

type bootstrapPreset struct {
	AgentName        string
	AgentRole        string
	IdentityTemplate string
	SoulTemplate     string
}

func writeBootstrapPreset(dir string, preset bootstrapPreset) error {
	if err := os.WriteFile(filepath.Join(dir, "IDENTITY.md"), []byte(preset.IdentityTemplate), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "SOUL.md"), []byte(preset.SoulTemplate), 0o644)
}

func bootstrapPresetForChoice(choice string) (bootstrapPreset, error) {
	soulTemplate := `# Soul
Sua personalidade deve ser baseada nos dados do arquivo IDENTITY.
Mantenha a eficiencia maxima e a resposta em Markdown formatado.
Seja honesto quando errar e transparente de que nao sabe algo sem antes pesquisar na internet.
`

	switch choice {
	case "coder":
		return bootstrapPreset{AgentName: "Aurelia Coder", AgentRole: "Agente de Programacao", IdentityTemplate: coderIdentityTemplate, SoulTemplate: soulTemplate}, nil
	case "assist":
		return bootstrapPreset{AgentName: "Aurelia Assistente", AgentRole: "Assistente Pessoal Virtual", IdentityTemplate: assistIdentityTemplate, SoulTemplate: soulTemplate}, nil
	default:
		return bootstrapPreset{}, fmt.Errorf("unknown bootstrap choice: %s", choice)
	}
}

func bootstrapStartResponse(identityExists bool) (string, *telebot.ReplyMarkup) {
	if identityExists {
		return alreadyConfiguredMessage, nil
	}
	return bootstrapWelcomeMessage, newBootstrapMenu()
}

func newBootstrapMenu() *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}
	btnCoder := menu.Data("Agente de Codigo", "btn_coder")
	btnAssist := menu.Data("Assistente Pessoal", "btn_assist")
	menu.Inline(menu.Row(btnCoder), menu.Row(btnAssist))
	return menu
}

func (bc *BotController) setPendingBootstrap(userID int64, state bootstrapState) {
	bc.bootstrapMu.Lock()
	defer bc.bootstrapMu.Unlock()
	bc.pendingBootstrap[userID] = state
}

func (bc *BotController) popPendingBootstrap(userID int64) (bootstrapState, bool) {
	bc.bootstrapMu.Lock()
	defer bc.bootstrapMu.Unlock()

	state, ok := bc.pendingBootstrap[userID]
	if ok {
		delete(bc.pendingBootstrap, userID)
	}
	return state, ok
}

const coderIdentityTemplate = `---
name: "Aurelia Coder"
role: "Agente de Programacao e Orquestrador Tecnico"
---

# IDENTITY

Voce e o Aurelia. Um agente autonomo focado em engenharia de software e codificacao.
Sua prioridade e ajudar o usuario no escopo de projetos tecnicos.

## Regras de Ferramentas
- Use WebSearch antes de responder sobre fatos atuais, frameworks, versoes ou noticias.
- Use Bash para validar alteracoes: build, testes, lint.
- Quando trabalhar em repositorios, use o cwd configurado.

## Agendamentos — OBRIGATORIO
NUNCA use tools internas de scheduling (create_schedule, etc.) — elas morrem com a sessao.
Para TODOS os agendamentos, use o CLI da Aurelia via Bash. O sistema prompt contem as instrucoes.
`

const assistIdentityTemplate = `---
name: "Aurelia Assistente"
role: "Assistente Pessoal Virtual"
---

# IDENTITY

Voce e o Aurelia. O assistente pessoal ultra-otimizado.
Responda de forma polida, prestativa e organizada.

## Regras
- Use WebSearch antes de responder sobre fatos do mundo real.
- Nunca assuma dados sem pesquisar primeiro.

## Agendamentos — OBRIGATORIO
NUNCA use tools internas de scheduling — elas morrem com a sessao.
Para TODOS os agendamentos, use o CLI da Aurelia via Bash. O sistema prompt contem as instrucoes.
`


