package e2e_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/tools"
)

type scriptedLLM struct {
	mu       sync.Mutex
	calls    int
	handler  func(systemPrompt string, history []agent.Message, tools []agent.Tool) (*agent.ModelResponse, error)
}

func (s *scriptedLLM) GenerateContent(ctx context.Context, systemPrompt string, history []agent.Message, defs []agent.Tool) (*agent.ModelResponse, error) {
	s.mu.Lock()
	s.calls++
	s.mu.Unlock()
	return s.handler(systemPrompt, history, defs)
}

type queuedLLM struct {
	mu        sync.Mutex
	responses []*agent.ModelResponse
	errors    []error
}

func (q *queuedLLM) GenerateContent(ctx context.Context, systemPrompt string, history []agent.Message, tools []agent.Tool) (*agent.ModelResponse, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.errors) > 0 {
		err := q.errors[0]
		q.errors = q.errors[1:]
		if err != nil {
			return nil, err
		}
	}

	if len(q.responses) == 0 {
		return &agent.ModelResponse{Content: "ok"}, nil
	}

	resp := q.responses[0]
	q.responses = q.responses[1:]
	return resp, nil
}

type loopExecutorAdapter struct {
	loop *agent.Loop
}

func (a *loopExecutorAdapter) Execute(ctx context.Context, systemPrompt string, history []agent.Message, allowedTools []string) ([]agent.Message, string, error) {
	return a.loop.Run(ctx, systemPrompt, history, allowedTools)
}

func TestE2E_PersonaLoopWithRealTools(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mem, canonical := newCanonicalServiceForE2E(t)

	if err := canonical.ApplyFacts(ctx, "42", map[string]memory.Fact{
		"user.name":                         {Scope: "user", EntityID: "42", Key: "user.name", Value: "Rafael", Source: "bootstrap"},
		"user.preference.response_style":   {Scope: "user", EntityID: "42", Key: "user.preference.response_style", Value: "direto", Source: "conversation"},
		"project.memory.strategy":          {Scope: "project", EntityID: "default", Key: "project.memory.strategy", Value: "sqlite + facts + notes", Source: "conversation"},
	}) ; err != nil {
		t.Fatalf("ApplyFacts() error = %v", err)
	}
	if err := mem.AddNote(ctx, memory.Note{
		ConversationID: "42",
		Topic:          "architecture",
		Kind:           "decision",
		Summary:        "Manter o monolito modular com tools locais.",
		Importance:     9,
		Source:         "conversation",
	}); err != nil {
		t.Fatalf("AddNote() error = %v", err)
	}

	prompt, allowedTools, err := canonical.BuildPromptForQuery(ctx, "42", "42", "verifique o workspace local")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}

	workspace := t.TempDir()
	targetFile := filepath.Join(workspace, "hello.txt")
	if err := os.WriteFile(targetFile, []byte("hello e2e"), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", targetFile, err)
	}

	llm := &scriptedLLM{
		handler: func(systemPrompt string, history []agent.Message, defs []agent.Tool) (*agent.ModelResponse, error) {
			if len(history) == 1 {
				if !strings.Contains(systemPrompt, "Nome canonico do usuario: Rafael") {
					t.Fatalf("expected canonical user name in prompt, got %q", systemPrompt)
				}
				if !strings.Contains(systemPrompt, "# CANONICAL IDENTITY") {
					t.Fatalf("expected canonical identity block in prompt, got %q", systemPrompt)
					}
				return &agent.ModelResponse{
					Content: "Vou verificar o arquivo local.",
					ToolCalls: []agent.ToolCall{
						{
							ID:   "call-1",
							Name: "run_command",
							Arguments: map[string]interface{}{
								"command": "Get-Content hello.txt",
								"workdir": workspace,
							},
						},
					},
				}, nil
			}

			toolOutput := history[len(history)-1].Content
			if !strings.Contains(toolOutput, "hello e2e") {
				t.Fatalf("expected run_command output in history, got %q", toolOutput)
			}

			return &agent.ModelResponse{
				Content: "Verificacao concluida: o arquivo hello.txt contem `hello e2e`.",
			}, nil
		},
	}

	registry := agent.NewToolRegistry()
	registry.Register(agent.Tool{Name: "read_file"}, tools.ReadFileHandler)
	registry.Register(agent.Tool{Name: "list_dir"}, tools.ListDirHandler)
	registry.Register(agent.Tool{Name: "run_command"}, tools.RunCommandHandler)

	loop := agent.NewLoop(llm, registry, 4)
	history, answer, err := loop.Run(ctx, prompt, []agent.Message{{Role: "user", Content: "verifique o workspace local"}}, allowedTools)
	if err != nil {
		t.Fatalf("Loop.Run() error = %v", err)
	}
	if !strings.Contains(answer, "hello e2e") {
		t.Fatalf("expected final answer to contain tool observation, got %q", answer)
	}
	if len(history) < 3 || history[len(history)-2].Role != "tool" {
		t.Fatalf("expected tool observation in loop history, got %#v", history)
	}
}

func TestE2E_MasterTeamRecoveryFlow(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store, err := agent.NewSQLiteTaskStore(filepath.Join(t.TempDir(), "teams.db"))
	if err != nil {
		t.Fatalf("NewSQLiteTaskStore() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	manager, err := agent.NewTeamManager(store)
	if err != nil {
		t.Fatalf("NewTeamManager() error = %v", err)
	}

	llm := &queuedLLM{
		errors: []error{context.DeadlineExceeded, nil},
		responses: []*agent.ModelResponse{
			{Content: "recovery completed"},
		},
	}

	notifications := make(chan string, 8)
	service, err := agent.NewMasterTeamService(manager, llm, agent.NewToolRegistry(), 3, func(teamKey string, message string) {
		notifications <- message
	})
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	if _, err := service.Spawn(ctx, "conversation-e2e", "user-e2e", "reviewer", "Reviews work", "analyze and recover if needed"); err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	var sawFailure bool
	var sawFinal bool
	deadline := time.After(4 * time.Second)
	for !sawFailure || !sawFinal {
		select {
		case msg := <-notifications:
			if strings.Contains(strings.ToLower(msg), "falhou") {
				sawFailure = true
			}
			if strings.Contains(strings.ToLower(msg), "consolidei o que saiu deste run") && strings.Contains(msg, "recovery completed") {
				sawFinal = true
			}
		case <-deadline:
			t.Fatalf("expected both failure and final recovery notification, got failure=%v final=%v", sawFailure, sawFinal)
		}
	}

	snapshot, err := service.BuildExecutionStatusSnapshot(ctx, "conversation-e2e", "")
	if err != nil {
		t.Fatalf("BuildExecutionStatusSnapshot() error = %v", err)
	}
	if snapshot.Completed == 0 {
		t.Fatalf("expected completed execution snapshot after recovery, got %#v", snapshot)
	}
}

func TestE2E_CronScheduleLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store, err := cron.NewSQLiteCronStore(filepath.Join(t.TempDir(), "cron.db"))
	if err != nil {
		t.Fatalf("NewSQLiteCronStore() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	service := cron.NewService(store, nil)
	createTool := tools.NewCreateScheduleTool(service)

	workspace := t.TempDir()
	llm := &scriptedLLM{
		handler: func(systemPrompt string, history []agent.Message, defs []agent.Tool) (*agent.ModelResponse, error) {
			if len(history) == 1 {
				return &agent.ModelResponse{
					Content: "Executando rotina agendada.",
					ToolCalls: []agent.ToolCall{
						{
							ID:   "cron-call-1",
							Name: "run_command",
							Arguments: map[string]interface{}{
								"command": "Write-Output 'rotina executada'",
								"workdir": workspace,
							},
						},
					},
				}, nil
			}

			if got := history[len(history)-1].Content; !strings.Contains(got, "rotina executada") {
				t.Fatalf("expected cron tool output in history, got %q", got)
			}
			return &agent.ModelResponse{Content: "Resumo da rotina: rotina executada."}, nil
		},
	}

	registry := agent.NewToolRegistry()
	registry.Register(agent.Tool{Name: "run_command"}, tools.RunCommandHandler)
	loop := agent.NewLoop(llm, registry, 4)
	runtime := cron.NewAgentCronRuntime(&loopExecutorAdapter{loop: loop}, "cron system prompt", []string{"run_command"})
	scheduler, err := cron.NewScheduler(store, runtime, nil, cron.SchedulerConfig{PollInterval: time.Millisecond})
	if err != nil {
		t.Fatalf("NewScheduler() error = %v", err)
	}

	runAt := time.Now().UTC().Add(-1 * time.Minute).Format(time.RFC3339)
	createCtx := agent.WithTeamContext(ctx, "12345", "user-cron")
	if _, err := createTool.Execute(createCtx, map[string]interface{}{
		"schedule_type": "once",
		"run_at":        runAt,
		"prompt":        "Execute a rotina e traga um resumo",
	}); err != nil {
		t.Fatalf("createTool.Execute() error = %v", err)
	}

	jobs, err := service.ListJobs(ctx, 12345)
	if err != nil {
		t.Fatalf("ListJobs() error = %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected one created cron job, got %d", len(jobs))
	}

	processed, err := scheduler.RunDueJobs(ctx)
	if err != nil {
		t.Fatalf("RunDueJobs() error = %v", err)
	}
	if processed != 1 {
		t.Fatalf("expected one due cron job to be processed, got %d", processed)
	}

	executions, err := store.ListExecutionsByJob(ctx, jobs[0].ID)
	if err != nil {
		t.Fatalf("ListExecutionsByJob() error = %v", err)
	}
	if len(executions) != 1 {
		t.Fatalf("expected one cron execution record, got %d", len(executions))
	}
	if executions[0].Status != "success" || !strings.Contains(executions[0].OutputSummary, "rotina executada") {
		t.Fatalf("unexpected execution record: %#v", executions[0])
	}

	job, err := store.GetJob(ctx, jobs[0].ID)
	if err != nil {
		t.Fatalf("GetJob() error = %v", err)
	}
	if job == nil || job.Active {
		t.Fatalf("expected once job to be deactivated after execution, got %#v", job)
	}
}

func newCanonicalServiceForE2E(t *testing.T) (*memory.MemoryManager, *persona.CanonicalIdentityService) {
	t.Helper()

	mem, err := memory.NewMemoryManager(filepath.Join(t.TempDir(), "e2e-memory.db"), 8)
	if err != nil {
		t.Fatalf("NewMemoryManager() error = %v", err)
	}
	t.Cleanup(func() { _ = mem.Close() })

	dir := t.TempDir()
	identityPath := filepath.Join(dir, "IDENTITY.md")
	soulPath := filepath.Join(dir, "SOUL.md")
	userPath := filepath.Join(dir, "USER.md")

	identityContent := `---
name: "Lex"
role: "Team Lead"
memory_window_size: 10
tools:
  - read_file
  - list_dir
---
IDENTITY_BODY`
	if err := os.WriteFile(identityPath, []byte(identityContent), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", identityPath, err)
	}
	if err := os.WriteFile(soulPath, []byte("# Soul\nBase.\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", soulPath, err)
	}
	if err := os.WriteFile(userPath, []byte("# User\nNome: Nao definido\nFuso horario: Relativo a sua localidade.\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", userPath, err)
	}

	return mem, persona.NewCanonicalIdentityService(mem, identityPath, soulPath, userPath, "", "", "")
}


