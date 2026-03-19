package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/observability"
)

// Loop executes the ReAct logic
type Loop struct {
	llm           LLMProvider
	registry      *ToolRegistry
	maxIterations int
	observer      observability.Recorder
}

// NewLoop constructs an agent loop. Pass -1 for unlimited iterations.
func NewLoop(llm LLMProvider, registry *ToolRegistry, maxIterations int) *Loop {
	if maxIterations == 0 {
		maxIterations = 5
	}
	return &Loop{
		llm:           llm,
		registry:      registry,
		maxIterations: maxIterations,
	}
}

func NewLoopWithObserver(llm LLMProvider, registry *ToolRegistry, maxIterations int, observer observability.Recorder) *Loop {
	loop := NewLoop(llm, registry, maxIterations)
	loop.observer = observer
	return loop
}

// Run executes the agent resolving loop on a given state of messages
func (l *Loop) Run(ctx context.Context, systemPrompt string, history []Message, allowedTools []string) ([]Message, string, error) {
	if _, ok := RunContextFromContext(ctx); !ok {
		ctx = WithRunContext(ctx, uuid.NewString())
	}

	currentHistory := make([]Message, len(history))
	copy(currentHistory, history)

	tools := l.registry.FilterDefinitions(allowedTools)
	tools = CompactToolsForPrompt(tools)
	systemPrompt = augmentSystemPromptWithToolGuidance(systemPrompt, tools)
	systemPrompt = augmentSystemPromptWithRuntimeCapabilities(systemPrompt, tools)
	toolMetrics := MeasureToolPayload(tools)

	var toolNames []string
	for _, t := range tools {
		toolNames = append(toolNames, t.Name)
	}
	observability.Log("info", "agent.loop", "starting loop run", observability.MergeFields(ContextFields(ctx), map[string]string{
		"tools":              strings.Join(toolNames, ","),
		"tool_count":         fmt.Sprintf("%d", toolMetrics.Count),
		"tool_payload_bytes": fmt.Sprintf("%d", toolMetrics.SerializedBytes),
	}))
	observability.Observe(ctx, l.observer, observability.Operation{
		RunID:      ContextFields(ctx)["run_id"],
		TeamID:     ContextFields(ctx)["team_id"],
		TaskID:     ContextFields(ctx)["task_id"],
		AgentName:  ContextFields(ctx)["agent"],
		Component:  "agent.loop",
		Operation:  "tool_context",
		Status:     "ok",
		DurationMS: 0,
		Summary:    fmt.Sprintf("tool_count=%d tool_payload_bytes=%d", toolMetrics.Count, toolMetrics.SerializedBytes),
	})

	iterations := 0
	for l.maxIterations < 0 || iterations < l.maxIterations {
		iterations++
		iterationFields := observability.MergeFields(ContextFields(ctx), map[string]string{
			"iteration": fmt.Sprintf("%d", iterations),
		})
		if l.maxIterations >= 0 {
			iterationFields["max_iterations"] = fmt.Sprintf("%d", l.maxIterations)
		}
		observability.Log("info", "agent.loop", "loop iteration", iterationFields)

		if ctx.Err() != nil {
			observability.Log("warn", "agent.loop", "loop cancelled before provider call", ContextFields(ctx))
			return currentHistory, "", fmt.Errorf("context cancelled by timer: %w", ctx.Err())
		}

		startedAt := time.Now()
		resp, err := l.llm.GenerateContent(ctx, systemPrompt, currentHistory, tools)
		duration := time.Since(startedAt)
		if err != nil {
			observability.Log("error", "agent.loop", "provider call failed", observability.MergeFields(ContextFields(ctx), map[string]string{
				"duration_ms": fmt.Sprintf("%d", duration.Milliseconds()),
				"error":       err.Error(),
			}))
			observability.Observe(ctx, l.observer, observability.Operation{
				RunID:      ContextFields(ctx)["run_id"],
				TeamID:     ContextFields(ctx)["team_id"],
				TaskID:     ContextFields(ctx)["task_id"],
				AgentName:  ContextFields(ctx)["agent"],
				Component:  "agent.loop",
				Operation:  "llm_generate",
				Status:     "error",
				DurationMS: duration.Milliseconds(),
				Summary:    err.Error(),
			})
			return currentHistory, "", fmt.Errorf("provider error: %w", err)
		}
		observability.Observe(ctx, l.observer, observability.Operation{
			RunID:      ContextFields(ctx)["run_id"],
			TeamID:     ContextFields(ctx)["team_id"],
			TaskID:     ContextFields(ctx)["task_id"],
			AgentName:  ContextFields(ctx)["agent"],
			Component:  "agent.loop",
			Operation:  "llm_generate",
			Status:     "ok",
			DurationMS: duration.Milliseconds(),
			Summary:    fmt.Sprintf("tool_calls=%d content_chars=%d", len(resp.ToolCalls), len(resp.Content)),
		})

		if len(resp.ToolCalls) == 0 {
			if resp.Content != "" || resp.ReasoningContent != "" {
				currentHistory = append(currentHistory, Message{
					Role:             "assistant",
					Content:          resp.Content,
					ReasoningContent: resp.ReasoningContent,
				})
			}
			return currentHistory, resp.Content, nil
		}

		currentHistory = append(currentHistory, Message{
			Role:             "assistant",
			Content:          resp.Content,
			ReasoningContent: resp.ReasoningContent,
			ToolCalls:        resp.ToolCalls,
		})

		for _, call := range resp.ToolCalls {
			observability.Log("info", "agent.tool", "executing tool", observability.MergeFields(ContextFields(ctx), map[string]string{
				"tool": call.Name,
			}))

			toolStartedAt := time.Now()
			resultStr, toolErr := l.registry.Execute(ctx, call.Name, call.Arguments)
			toolDuration := time.Since(toolStartedAt)

			if toolErr != nil {
				observability.Log("error", "agent.tool", "tool execution failed", observability.MergeFields(ContextFields(ctx), map[string]string{
					"tool":        call.Name,
					"duration_ms": fmt.Sprintf("%d", toolDuration.Milliseconds()),
					"error":       toolErr.Error(),
				}))
				observability.Observe(ctx, l.observer, observability.Operation{
					RunID:      ContextFields(ctx)["run_id"],
					TeamID:     ContextFields(ctx)["team_id"],
					TaskID:     ContextFields(ctx)["task_id"],
					AgentName:  ContextFields(ctx)["agent"],
					Component:  "agent.tool",
					Operation:  call.Name,
					Status:     "error",
					DurationMS: toolDuration.Milliseconds(),
					Summary:    toolErr.Error(),
				})
				errorPayload, _ := json.Marshal(map[string]string{
					"error": toolErr.Error(),
				})
				resultStr = string(errorPayload)
			} else {
				compacted := CompactToolOutputForHistory(ctx, l.observer, call.Name, resultStr)
				summary := fmt.Sprintf("raw_chars=%d compacted_chars=%d", compacted.RawChars, compacted.CompactedChars)
				if compacted.Oversized {
					summary += fmt.Sprintf(" oversized=true threshold_chars=%d", OversizedToolOutputThresholdChars)
				}
				if compacted.ArtifactID > 0 {
					summary += fmt.Sprintf(" artifact_id=%d", compacted.ArtifactID)
				}
				observability.Observe(ctx, l.observer, observability.Operation{
					RunID:      ContextFields(ctx)["run_id"],
					TeamID:     ContextFields(ctx)["team_id"],
					TaskID:     ContextFields(ctx)["task_id"],
					AgentName:  ContextFields(ctx)["agent"],
					Component:  "agent.tool",
					Operation:  call.Name,
					Status:     "ok",
					DurationMS: toolDuration.Milliseconds(),
					Summary:    summary,
				})
				resultStr = compacted.Content
			}

			currentHistory = append(currentHistory, Message{
				Role:       "tool",
				Content:    resultStr,
				ToolCallID: call.ID,
			})
		}
	}

	observability.Log("warn", "agent.loop", "loop hit max iterations", ContextFields(ctx))
	return currentHistory, "Desculpe, desisti ou deu timeout no processamento pois falhei nas chamadas em MAX iteracoes.", fmt.Errorf("max iterations reached")
}

func augmentSystemPromptWithToolGuidance(systemPrompt string, tools []Tool) string {
	toolNames := make(map[string]bool, len(tools))
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	var sections []string

	if toolNames["run_command"] {
		sections = append(sections, "Se o usuario pedir para rodar, testar, iniciar, buildar, validar, verificar healthcheck ou inspecionar um projeto local, voce deve tentar usar `run_command` antes de responder com passos manuais. So ofereca execucao manual se `run_command` falhar, for bloqueado ou nao existir.")
		sections = append(sections, "Quando operar em Windows com `run_command`, prefira sintaxe de PowerShell em vez de comandos Unix como `find`, `grep` ou `ls`.")
		sections = append(sections, "Nao diga que o ambiente esta bloqueado, que nao consegue executar processos ou que a execucao deve ser manual sem antes receber esse resultado explicitamente de uma tool. Se `run_command` nao retornou bloqueio ou erro, continue usando ferramentas.")
		sections = append(sections, "Se a tarefa exigir varias etapas locais, execute em sequencia: por exemplo subir o servico com `run_command`, depois testar endpoint com outro `run_command`, depois sintetizar o resultado observado.")
	}

	hasFilesystem := toolNames["read_file"] || toolNames["write_file"] || toolNames["list_dir"]
	if hasFilesystem {
		sections = append(sections, "As tools `read_file`, `write_file` e `list_dir` aceitam `workdir`. Sempre que estiver trabalhando em outro projeto ou pasta fora da raiz atual, informe `workdir` e use caminhos relativos a esse diretorio.")
	}

	if toolNames["run_command"] && hasFilesystem {
		sections = append(sections, "Se voce descobrir um diretorio de projeto via `run_command`, reutilize o mesmo `workdir` nas tools de filesystem para nao ler ou escrever no repositorio errado.")
	}
	if toolNames["spawn_agent"] {
		sections = append(sections, "Ao delegar trabalho com `spawn_agent` para outro projeto, passe o `workdir` canonico do projeto alvo. Nunca deixe subagente assumir por padrao a pasta do Aurelia como diretorio de trabalho.")
		sections = append(sections, "Se o usuario quiser interromper, pausar, retomar ou inspecionar a operacao do time, prefira usar `cancel_team`, `pause_team`, `resume_team` e `team_status` em vez de responder apenas em texto.")
	}

	if toolNames["create_schedule"] {
		sections = append(sections, "Se o usuario pedir lembretes, rotinas, tarefas recorrentes, avisos futuros ou qualquer acao para acontecer depois, voce deve considerar usar `create_schedule` em vez de apenas responder com texto.")
	}
	if toolNames["list_schedules"] || toolNames["pause_schedule"] || toolNames["resume_schedule"] || toolNames["delete_schedule"] {
		sections = append(sections, "Se o usuario perguntar quais agendamentos existem ou pedir para pausar, retomar ou remover uma rotina, use `list_schedules`, `pause_schedule`, `resume_schedule` e `delete_schedule` conforme a intencao.")
	}
	if toolNames["create_schedule"] || toolNames["list_schedules"] {
		sections = append(sections, "Nao exija comandos como `/cron`. A interface correta e linguagem natural; as tools de scheduling existem para voce transformar a intencao do usuario em operacoes reais.")
	}

	if len(sections) == 0 {
		return systemPrompt
	}

	return strings.TrimSpace(systemPrompt) + "\n\n# TOOL USAGE GUIDE\n" + strings.Join(sections, "\n")
}

func augmentSystemPromptWithRuntimeCapabilities(systemPrompt string, tools []Tool) string {
	var lines []string
	lines = append(lines, "# RUNTIME CAPABILITIES")
	if len(tools) == 0 {
		lines = append(lines, "Nenhuma tool esta disponivel neste runtime para esta execucao.")
		lines = append(lines, "Se houver duvida sobre capacidades, considere esta secao como fonte canonica em vez de assumir ferramentas inexistentes.")
	} else {
		lines = append(lines, "Tools disponiveis nesta execucao:")
		names := make([]string, 0, len(tools))
		for _, tool := range tools {
			names = append(names, tool.Name)
		}
		sort.Strings(names)
		for _, name := range names {
			lines = append(lines, "- "+name)
		}
		lines = append(lines, "Considere esta lista como a fonte canonica das capacidades reais deste runtime.")
	}

	base := strings.TrimSpace(systemPrompt)
	if base == "" {
		return strings.Join(lines, "\n")
	}
	return base + "\n\n" + strings.Join(lines, "\n")
}
