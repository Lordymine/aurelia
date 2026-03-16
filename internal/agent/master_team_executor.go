package agent

import (
	"context"
	"fmt"
)

type loopTaskExecutor struct {
	llm           LLMProvider
	registry      *ToolRegistry
	maxIterations int
	agentName     string
	roleDesc      string
}

func (e *loopTaskExecutor) ExecuteTask(ctx context.Context, task TeamTask) (string, error) {
	loop := NewLoop(e.llm, e.registry, e.maxIterations)
	systemPrompt := fmt.Sprintf(`Voce e o worker "%s".

Identidade e metodo:
%s

Tarefa:
%s

Diretorio de trabalho canonico da task:
%s

Regras:
1. O master e sempre o lider. Voce nunca responde ao usuario final.
2. Execute a tarefa com as ferramentas explicitamente liberadas para esta task.
3. Considere a secao RUNTIME CAPABILITIES como a fonte canonica das suas capacidades reais nesta execucao.
4. Se precisar delegar subtarefas, use spawn_agent apenas se essa tool estiver liberada.
5. Se 'read_team_inbox' estiver disponivel e a tarefa depender de contexto de outros agentes, leia sua inbox antes de decidir que falta informacao.
6. Se 'send_team_message' estiver disponivel, use a mailbox para pedir dados, alinhar dependencias ou avisar bloqueios a outros agentes da equipe.
7. Nao use a mailbox para conversa inutil; envie mensagens objetivas, acionaveis e curtas.
8. Se um diretorio de trabalho canonico estiver definido, reuse esse mesmo workdir nas tools locais e nunca caia por padrao na pasta do Aurelia.
9. Quando concluir, retorne um resumo objetivo do que foi feito e do que ainda depende de outros agentes, se houver.`, e.agentName, e.roleDesc, task.Prompt, displayTaskWorkdir(task.Workdir))

	allowedTools := append([]string(nil), task.AllowedTools...)
	if len(allowedTools) == 0 {
		allowedTools = make([]string, 0, len(e.registry.GetDefinitions()))
		for _, tool := range e.registry.GetDefinitions() {
			allowedTools = append(allowedTools, tool.Name)
		}
	}

	_, finalAnswer, err := loop.Run(ctx, systemPrompt, nil, allowedTools)
	return finalAnswer, err
}

func displayTaskWorkdir(workdir string) string {
	if workdir == "" {
		return "nao definido"
	}
	return workdir
}


