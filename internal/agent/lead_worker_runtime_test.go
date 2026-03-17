package agent

import (
	"context"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type fakeTaskExecutor struct {
	lastTaskID  string
	lastWorkdir string
	result      string
	err         error
}

func (f *fakeTaskExecutor) ExecuteTask(ctx context.Context, task TeamTask) (string, error) {
	f.lastTaskID = task.ID
	f.lastWorkdir, _ = WorkdirFromContext(ctx)
	return f.result, f.err
}

func newRuntimeTestManager(t *testing.T) TeamManager {
	t.Helper()

	store, err := NewSQLiteTaskStore(filepath.Join(t.TempDir(), "runtime.db"))
	if err != nil {
		t.Fatalf("NewSQLiteTaskStore() error = %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	manager, err := NewTeamManager(store)
	if err != nil {
		t.Fatalf("NewTeamManager() error = %v", err)
	}

	return manager
}

func TestWorkerRuntime_RunOnce_CompletesTaskAndReportsToLead(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-rt-1", "user-rt-1", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Implements feature work"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-rt-1",
		TeamID: teamID,
		Title:  "Implement feature",
		Prompt: "Write the first implementation draft",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	executor := &fakeTaskExecutor{result: "feature implemented"}
	worker := NewWorkerRuntime("worker-a", manager, executor)

	processed, err := worker.RunOnce(ctx, teamID)
	if err != nil {
		t.Fatalf("RunOnce() error = %v", err)
	}
	if !processed {
		t.Fatalf("expected worker to process one task")
	}
	if executor.lastTaskID != "task-rt-1" {
		t.Fatalf("expected worker executor to run task-rt-1, got %q", executor.lastTaskID)
	}

	inbox, err := manager.PullMessages(ctx, teamID, "master", 10)
	if err != nil {
		t.Fatalf("PullMessages() error = %v", err)
	}
	if len(inbox) != 1 {
		t.Fatalf("expected lead inbox to receive one result message, got %d", len(inbox))
	}
	if inbox[0].Kind != "result" {
		t.Fatalf("expected worker completion to post a result message, got %q", inbox[0].Kind)
	}
	if inbox[0].TaskID == nil || *inbox[0].TaskID != "task-rt-1" {
		t.Fatalf("expected result message to point to task-rt-1, got %#v", inbox[0].TaskID)
	}
}

func TestWorkerRuntime_RunOnce_InjectsTaskWorkdirIntoContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-rt-workdir", "user-rt-workdir", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Implements feature work"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:      "task-rt-workdir",
		TeamID:  teamID,
		Title:   "Implement feature",
		Prompt:  "Write the first implementation draft",
		Workdir: `C:\projetos\alvo`,
		Status:  TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	executor := &fakeTaskExecutor{result: "feature implemented"}
	worker := NewWorkerRuntime("worker-a", manager, executor)

	processed, err := worker.RunOnce(ctx, teamID)
	if err != nil {
		t.Fatalf("RunOnce() error = %v", err)
	}
	if !processed {
		t.Fatalf("expected worker to process one task")
	}
	if executor.lastWorkdir != `C:\projetos\alvo` {
		t.Fatalf("expected injected task workdir, got %q", executor.lastWorkdir)
	}
}

func TestLeadRuntime_CollectInbox_ReturnsStructuredWorkerUpdates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-rt-2", "user-rt-2", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	taskID := "task-mail-1"
	if err := manager.PostMessage(ctx, MailMessage{
		ID:        "mail-status-1",
		TeamID:    teamID,
		FromAgent: "worker-a",
		ToAgent:   "master",
		TaskID:    &taskID,
		Kind:      "status_update",
		Body:      "halfway done",
	}); err != nil {
		t.Fatalf("PostMessage(status) error = %v", err)
	}

	if err := manager.PostMessage(ctx, MailMessage{
		ID:        "mail-result-1",
		TeamID:    teamID,
		FromAgent: "worker-a",
		ToAgent:   "master",
		TaskID:    &taskID,
		Kind:      "result",
		Body:      "finished successfully",
	}); err != nil {
		t.Fatalf("PostMessage(result) error = %v", err)
	}

	lead := NewLeadRuntime(manager, "master")

	updates, err := lead.CollectInbox(ctx, teamID, 10)
	if err != nil {
		t.Fatalf("CollectInbox() error = %v", err)
	}
	if len(updates) != 2 {
		t.Fatalf("expected two structured updates for lead, got %d", len(updates))
	}
	if updates[0].Kind != "status_update" {
		t.Fatalf("expected first update to preserve ordering, got %q", updates[0].Kind)
	}
	if updates[1].Kind != "result" {
		t.Fatalf("expected second update to be result, got %q", updates[1].Kind)
	}
}

func TestLeadRuntime_CollectEvents_ReturnsOrderedTaskEvents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-rt-events", "user-rt-events", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes eventful work"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-evt-1",
		TeamID: teamID,
		Title:  "Event task",
		Prompt: "emit events",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	task, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask() error = %v", err)
	}
	if task == nil {
		t.Fatalf("expected claimed task")
	}

	if err := manager.FailTask(ctx, teamID, task.ID, "worker-a", "boom"); err != nil {
		t.Fatalf("FailTask() error = %v", err)
	}

	lead := NewLeadRuntime(manager, "master")
	events, err := lead.CollectEvents(ctx, teamID, 10)
	if err != nil {
		t.Fatalf("CollectEvents() error = %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	if events[0].EventType != "task_created" || events[1].EventType != "task_claimed" || events[2].EventType != "task_failed" {
		t.Fatalf("unexpected event types: %#v", events)
	}
}

func TestWorkerRuntime_RunUntilIdle_ProcessesUnlockedChain(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-rt-3", "user-rt-3", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes chained work"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "chain-1",
		TeamID: teamID,
		Title:  "Step 1",
		Prompt: "First step",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask(chain-1) error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "chain-2",
		TeamID: teamID,
		Title:  "Step 2",
		Prompt: "Second step",
		Status: TaskPending,
	}, []string{"chain-1"}); err != nil {
		t.Fatalf("CreateTask(chain-2) error = %v", err)
	}

	executor := &sequencedTaskExecutor{
		results: map[string]string{
			"chain-1": "step one complete",
			"chain-2": "step two complete",
		},
	}
	worker := NewWorkerRuntime("worker-a", manager, executor)

	processedCount, err := worker.RunUntilIdle(ctx, teamID)
	if err != nil {
		t.Fatalf("RunUntilIdle() error = %v", err)
	}
	if processedCount != 2 {
		t.Fatalf("expected worker to process two chained tasks, got %d", processedCount)
	}
	if len(executor.seen) != 2 || executor.seen[0] != "chain-1" || executor.seen[1] != "chain-2" {
		t.Fatalf("expected sequential execution order [chain-1 chain-2], got %#v", executor.seen)
	}
}

func TestMasterTeamService_Spawn_NotifiesMasterLeadWhenWorkerFinishes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{
		response: &ModelResponse{Content: "worker finished successfully"},
	}
	registry := NewToolRegistry()

	notifications := make(chan string, 1)
	service, err := NewMasterTeamService(manager, llm, registry, 3, func(teamKey string, message string) {
		notifications <- teamKey + "|" + message
	})
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	_, err = service.Spawn(ctx, "conversation-1", "user-5", "researcher", "Researches the task", "Summarize the current architecture")
	if err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	select {
	case notification := <-notifications:
		if notification == "" {
			t.Fatalf("expected non-empty master notification")
		}
		if !strings.Contains(notification, "pending=") {
			t.Fatalf("expected master notification to include team snapshot, got %q", notification)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("expected master lead to receive a worker update")
	}
}

func TestMasterTeamService_Spawn_FromTaskContextCreatesDependentSubtask(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{
		response: &ModelResponse{Content: "child worker completed"},
	}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	teamID, err := service.ensureTeam(ctx, "conversation-dep", "user-dep")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "parent-worker", "Parent task owner"); err != nil {
		t.Fatalf("RegisterTeammate(parent-worker) error = %v", err)
	}

	parentTask := TeamTask{
		ID:     "parent-task",
		TeamID: teamID,
		Title:  "Parent",
		Prompt: "parent prompt",
		Status: TaskPending,
	}
	if err := manager.CreateTask(ctx, parentTask, nil); err != nil {
		t.Fatalf("CreateTask(parent) error = %v", err)
	}

	claimedParent, err := manager.ClaimNextTask(ctx, teamID, "parent-worker")
	if err != nil {
		t.Fatalf("ClaimNextTask(parent) error = %v", err)
	}
	if claimedParent == nil || claimedParent.ID != "parent-task" {
		t.Fatalf("expected parent-task to be running, got %#v", claimedParent)
	}

	spawnCtx := WithTaskContext(WithTeamContext(ctx, "conversation-dep", "user-dep"), teamID, "parent-task")
	childTaskID, err := service.Spawn(spawnCtx, "conversation-dep", "user-dep", "child-worker", "Child specialist", "child prompt")
	if err != nil {
		t.Fatalf("Spawn(child) error = %v", err)
	}
	if childTaskID == "" {
		t.Fatalf("expected child task id")
	}

	noChildYet, err := manager.ClaimNextTask(ctx, teamID, "child-worker")
	if err != nil {
		t.Fatalf("ClaimNextTask(child before parent done) error = %v", err)
	}
	if noChildYet != nil {
		t.Fatalf("expected child task to remain blocked until parent completes, got %#v", noChildYet)
	}

	if err := manager.CompleteTask(ctx, teamID, "parent-task", "parent-worker", "parent done"); err != nil {
		t.Fatalf("CompleteTask(parent) error = %v", err)
	}

	childTask, err := manager.ClaimNextTask(ctx, teamID, "child-worker")
	if err != nil {
		t.Fatalf("ClaimNextTask(child after parent done) error = %v", err)
	}
	if childTask == nil || childTask.ID != childTaskID {
		t.Fatalf("expected dependent child task to unlock after parent completion, got %#v", childTask)
	}
}

func TestMasterTeamService_ReusesPersistentWorkerLoopPerAgent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{
		response: &ModelResponse{Content: "done"},
	}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	if _, err := service.Spawn(ctx, "conversation-persist", "user-persist", "builder", "Builds things", "task one"); err != nil {
		t.Fatalf("Spawn(task one) error = %v", err)
	}
	if _, err := service.Spawn(ctx, "conversation-persist", "user-persist", "builder", "Builds things", "task two"); err != nil {
		t.Fatalf("Spawn(task two) error = %v", err)
	}

	service.mu.Lock()
	defer service.mu.Unlock()
	if len(service.workerLoops) != 1 {
		t.Fatalf("expected one persistent worker loop for same team/agent, got %d", len(service.workerLoops))
	}
}

func TestMasterTeamService_NotifiesMasterWhenWorkerFails(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{
		err: context.DeadlineExceeded,
	}
	registry := NewToolRegistry()

	notifications := make(chan string, 1)
	service, err := NewMasterTeamService(manager, llm, registry, 3, func(teamKey string, message string) {
		notifications <- message
	})
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	if _, err := service.Spawn(ctx, "conversation-fail", "user-fail", "reviewer", "Reviews work", "task failing"); err != nil {
		t.Fatalf("Spawn(failing task) error = %v", err)
	}

	select {
	case msg := <-notifications:
		if !strings.Contains(msg, "falhou") {
			t.Fatalf("expected failure notification to mention falha, got %q", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("expected master failure notification")
	}
}

type sequencedTaskExecutor struct {
	results map[string]string
	seen    []string
}

func (s *sequencedTaskExecutor) ExecuteTask(ctx context.Context, task TeamTask) (string, error) {
	s.seen = append(s.seen, task.ID)
	return s.results[task.ID], nil
}

type queuedLLMProvider struct {
	mu        sync.Mutex
	responses []*ModelResponse
	errors    []error
}

func (q *queuedLLMProvider) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
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
		return &ModelResponse{Content: "default"}, nil
	}

	resp := q.responses[0]
	q.responses = q.responses[1:]
	return resp, nil
}

type capturingToolProvider struct {
	mu        sync.Mutex
	toolNames []string
	response  *ModelResponse
}

func (c *capturingToolProvider) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.toolNames = c.toolNames[:0]
	for _, tool := range tools {
		c.toolNames = append(c.toolNames, tool.Name)
	}

	if c.response != nil {
		return c.response, nil
	}
	return &ModelResponse{Content: "done"}, nil
}

func (c *capturingToolProvider) CapturedTools() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]string(nil), c.toolNames...)
}

func TestMasterTeamService_Spawn_PersistsAndAppliesAllowedToolsForWorker(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &capturingToolProvider{response: &ModelResponse{Content: "filtered execution complete"}}
	registry := NewToolRegistry()
	registry.Register(Tool{Name: "web_search", Description: "searches web"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "read_file", Description: "reads files"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "write_file", Description: "writes files"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})

	notifications := make(chan string, 1)
	service, err := NewMasterTeamService(manager, llm, registry, 3, func(teamKey string, message string) {
		notifications <- message
	})
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	taskID, err := service.Spawn(ctx, "conversation-filtered-tools", "user-filtered-tools", "researcher", "Pesquisa externa", "Busque referencias atuais", "web_search", "read_file")
	if err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	select {
	case <-notifications:
	case <-time.After(2 * time.Second):
		t.Fatalf("expected worker completion notification")
	}

	teamID, err := service.ensureTeam(ctx, "conversation-filtered-tools", "user-filtered-tools")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}

	task, err := manager.GetTask(ctx, teamID, taskID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if task == nil {
		t.Fatal("expected task to exist")
	}
	if len(task.AllowedTools) != 2 || task.AllowedTools[0] != "web_search" || task.AllowedTools[1] != "read_file" {
		t.Fatalf("unexpected persisted allowed tools: %#v", task.AllowedTools)
	}

	gotTools := llm.CapturedTools()
	if len(gotTools) != 2 || gotTools[0] != "web_search" || gotTools[1] != "read_file" {
		t.Fatalf("unexpected worker tool list: %#v", gotTools)
	}
}

func TestMasterTeamService_Spawn_PersistsWorkdirForWorkerAndRecovery(t *testing.T) {
	t.Parallel()

	ctx := WithWorkdirContext(context.Background(), `C:\projetos\api-alvo`)
	manager := newRuntimeTestManager(t)
	llm := &queuedLLMProvider{
		errors: []error{context.DeadlineExceeded, nil},
		responses: []*ModelResponse{
			{Content: "recovery completed"},
		},
	}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	taskID, err := service.Spawn(ctx, "conversation-workdir-recovery", "user-workdir-recovery", "builder", "Implementa setup", "executar setup do projeto")
	if err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	time.Sleep(400 * time.Millisecond)

	teamID, err := service.ensureTeam(ctx, "conversation-workdir-recovery", "user-workdir-recovery")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}

	task, err := manager.GetTask(ctx, teamID, taskID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if task == nil {
		t.Fatal("expected original task to exist")
	}
	if task.Workdir != `C:\projetos\api-alvo` {
		t.Fatalf("expected original task workdir to persist, got %q", task.Workdir)
	}

	tasks, err := manager.ListTasks(ctx, teamID)
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	var recovery *TeamTask
	for i := range tasks {
		if strings.HasPrefix(tasks[i].Title, recoveryTaskPrefix) {
			recovery = &tasks[i]
			break
		}
	}
	if recovery == nil {
		t.Fatal("expected recovery task to be created")
	}
	if recovery.Workdir != `C:\projetos\api-alvo` {
		t.Fatalf("expected recovery task to inherit workdir, got %q", recovery.Workdir)
	}
}

func TestMasterTeamService_OnFailure_CreatesOneRecoveryCycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &queuedLLMProvider{
		errors: []error{context.DeadlineExceeded, nil},
		responses: []*ModelResponse{
			{Content: "recovery completed"},
		},
	}
	registry := NewToolRegistry()

	notifications := make(chan string, 4)
	service, err := NewMasterTeamService(manager, llm, registry, 3, func(teamKey string, message string) {
		notifications <- message
	})
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	if _, err := service.Spawn(ctx, "conversation-recovery", "user-recovery", "reviewer", "Reviews work", "task failing first"); err != nil {
		t.Fatalf("Spawn(recovery flow) error = %v", err)
	}

	deadline := time.After(3 * time.Second)
	var sawFailure bool
	var sawRecovery bool
	for !sawFailure || !sawRecovery {
		select {
		case msg := <-notifications:
			if strings.Contains(msg, "falhou") {
				sawFailure = true
			}
			if strings.Contains(msg, "recovery completed") {
				sawRecovery = true
			}
		case <-deadline:
			t.Fatalf("expected both failure and recovery notifications, got failure=%v recovery=%v", sawFailure, sawRecovery)
		}
	}
}

func TestMasterTeamService_OnRepeatedFailure_EscalatesAfterRetryLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &queuedLLMProvider{
		errors: []error{context.DeadlineExceeded, context.DeadlineExceeded},
	}
	registry := NewToolRegistry()

	notifications := make(chan string, 8)
	service, err := NewMasterTeamService(manager, llm, registry, 3, func(teamKey string, message string) {
		notifications <- message
	})
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	service.maxRecoveryAttempts = 1

	if _, err := service.Spawn(ctx, "conversation-escalate", "user-escalate", "reviewer", "Reviews work", "task keeps failing"); err != nil {
		t.Fatalf("Spawn(escalation flow) error = %v", err)
	}

	deadline := time.After(4 * time.Second)
	var sawEscalation bool
	for !sawEscalation {
		select {
		case msg := <-notifications:
			if strings.Contains(strings.ToLower(msg), "escalated") {
				sawEscalation = true
			}
		case <-deadline:
			t.Fatalf("expected escalation notification after retry limit")
		}
	}
}

func TestMasterTeamService_OnFailure_NotifiesMasterAboutReplanning(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{
		err: context.DeadlineExceeded,
	}
	registry := NewToolRegistry()

	notifications := make(chan string, 4)
	service, err := NewMasterTeamService(manager, llm, registry, 3, func(teamKey string, message string) {
		notifications <- message
	})
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	if _, err := service.Spawn(ctx, "conversation-replan-msg", "user-replan-msg", "reviewer", "Reviews work", "task failing"); err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	deadline := time.After(3 * time.Second)
	for {
		select {
		case msg := <-notifications:
			if strings.Contains(strings.ToLower(msg), "replanejamento") {
				return
			}
		case <-deadline:
			t.Fatalf("expected replanning notification")
		}
	}
}

func TestMasterTeamService_OnFailure_ReassignsRecoveryToAnotherWorkerWhenAvailable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &queuedLLMProvider{
		errors: []error{context.DeadlineExceeded, nil},
		responses: []*ModelResponse{
			{Content: "reassigned recovery completed"},
		},
	}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	teamID, err := service.ensureTeam(ctx, "conversation-reassign", "user-reassign")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}
	if err := manager.RegisterTeammate(ctx, teamID, "reviewer", "Reviews"); err != nil {
		t.Fatalf("RegisterTeammate(reviewer) error = %v", err)
	}
	if err := manager.RegisterTeammate(ctx, teamID, "fixer", "Fixes"); err != nil {
		t.Fatalf("RegisterTeammate(fixer) error = %v", err)
	}
	service.mu.Lock()
	service.teamByKey["conversation-reassign"] = teamID
	service.userByKey["conversation-reassign"] = "user-reassign"
	service.memberSeen[teamID] = map[string]bool{MasterAgentName: true, "reviewer": true, "fixer": true}
	service.mu.Unlock()

	if _, err := service.Spawn(ctx, "conversation-reassign", "user-reassign", "reviewer", "Reviews work", "task fails first"); err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	time.Sleep(300 * time.Millisecond)

	tasks, err := manager.ListTasks(ctx, teamID)
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	var recovery *TeamTask
	for i := range tasks {
		if strings.HasPrefix(tasks[i].Title, recoveryTaskPrefix) {
			recovery = &tasks[i]
			break
		}
	}
	if recovery == nil {
		t.Fatal("expected recovery task to be created")
	}
	if recovery.AssignedAgent == nil || *recovery.AssignedAgent != "fixer" {
		t.Fatalf("expected recovery to be reassigned to fixer, got %#v", recovery.AssignedAgent)
	}
}

func TestMasterTeamService_OnFailure_MarksTaskTerminalAfterRetryLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &queuedLLMProvider{
		errors: []error{context.DeadlineExceeded, context.DeadlineExceeded},
	}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}
	service.maxRecoveryAttempts = 1

	teamID, err := service.ensureTeam(ctx, "conversation-terminal", "user-terminal")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}

	if _, err := service.Spawn(ctx, "conversation-terminal", "user-terminal", "reviewer", "Reviews work", "task keeps failing"); err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	time.Sleep(400 * time.Millisecond)

	tasks, err := manager.ListTasks(ctx, teamID)
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	var failedCount int
	for _, task := range tasks {
		if task.Status == TaskFailed {
			failedCount++
		}
	}
	if failedCount < 2 {
		t.Fatalf("expected original task and recovery task to end failed, got %d failed tasks", failedCount)
	}

	events, err := manager.ListEvents(ctx, teamID, 50)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}
	var recoveryCreated int
	for _, event := range events {
		if event.EventType == "task_created" && strings.Contains(event.Payload, recoveryTaskPrefix) {
			recoveryCreated++
		}
	}
	if recoveryCreated != 1 {
		t.Fatalf("expected exactly one recovery task before terminal stop, got %d", recoveryCreated)
	}
}

func TestMasterTeamService_OnFailure_BuildsRecoveryTaskFromEventHistory(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{
		err: context.DeadlineExceeded,
	}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	_, err = service.Spawn(ctx, "conversation-replan-events", "user-replan-events", "reviewer", "Reviews work", "task failing for history")
	if err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	teamID, err := service.ensureTeam(ctx, "conversation-replan-events", "user-replan-events")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}

	var recovery *TeamTask
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		tasks, err := manager.ListTasks(ctx, teamID)
		if err != nil {
			t.Fatalf("ListTasks() error = %v", err)
		}
		for i := range tasks {
			if strings.HasPrefix(tasks[i].Title, recoveryTaskPrefix) {
				candidate := tasks[i]
				recovery = &candidate
				break
			}
		}
		if recovery != nil {
			break
		}
		time.Sleep(25 * time.Millisecond)
	}
	if recovery == nil {
		t.Fatalf("expected recovery task to be created")
	}

	if !strings.Contains(recovery.Prompt, "Historico recente") {
		t.Fatalf("expected recovery prompt to include event history, got %q", recovery.Prompt)
	}
	if !strings.Contains(recovery.Prompt, "task_failed") {
		t.Fatalf("expected recovery prompt to mention failure event type, got %q", recovery.Prompt)
	}
	if !strings.Contains(recovery.Prompt, "context deadline exceeded") {
		t.Fatalf("expected recovery prompt to include failure reason, got %q", recovery.Prompt)
	}
}

func TestMasterTeamService_BuildStatusSnapshot_ReturnsAggregatedCounts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{response: &ModelResponse{Content: "done"}}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	teamID, err := service.ensureTeam(ctx, "conversation-status", "user-status")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{ID: "task-pending", TeamID: teamID, Title: "Pending", Prompt: "pending", Status: TaskPending}, nil); err != nil {
		t.Fatalf("CreateTask(pending) error = %v", err)
	}
	if err := manager.CreateTask(ctx, TeamTask{ID: "task-running", TeamID: teamID, Title: "Running", Prompt: "running", Status: TaskPending, AssignedAgent: strPtr("worker-a")}, nil); err != nil {
		t.Fatalf("CreateTask(running) error = %v", err)
	}
	claimed, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil || claimed == nil {
		t.Fatalf("ClaimNextTask() error = %v claimed=%#v", err, claimed)
	}
	if err := manager.CreateTask(ctx, TeamTask{ID: "task-done", TeamID: teamID, Title: "Done", Prompt: "done", Status: TaskPending}, nil); err != nil {
		t.Fatalf("CreateTask(done) error = %v", err)
	}
	claimedDone, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil || claimedDone == nil {
		t.Fatalf("ClaimNextTask(done) error = %v claimed=%#v", err, claimedDone)
	}
	if err := manager.CompleteTask(ctx, teamID, claimedDone.ID, "worker-a", "done"); err != nil {
		t.Fatalf("CompleteTask(done) error = %v", err)
	}

	snapshot, err := service.BuildStatusSnapshot(ctx, "conversation-status")
	if err != nil {
		t.Fatalf("BuildStatusSnapshot() error = %v", err)
	}
	if snapshot.Pending != 1 || snapshot.Running != 1 || snapshot.Completed != 1 {
		t.Fatalf("unexpected snapshot counts: %#v", snapshot)
	}
}

func strPtr(v string) *string {
	return &v
}

type slowTaskExecutor struct {
	started chan struct{}
	release chan struct{}
	result  string
}

func (s *slowTaskExecutor) ExecuteTask(ctx context.Context, task TeamTask) (string, error) {
	close(s.started)
	<-s.release
	return s.result, nil
}

func TestMasterTeamService_FinalNotification_WhenTeamConverges(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{
		response: &ModelResponse{Content: "implementation complete"},
	}
	registry := NewToolRegistry()

	notifications := make(chan string, 2)
	service, err := NewMasterTeamService(manager, llm, registry, 3, func(teamKey string, message string) {
		notifications <- message
	})
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	if _, err := service.Spawn(ctx, "conversation-final", "user-final", "builder", "Builds the final output", "Finish the feature"); err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	select {
	case msg := <-notifications:
		if !strings.Contains(strings.ToLower(msg), "fechei este ciclo do time") {
			t.Fatalf("expected final master response, got %q", msg)
		}
		if !strings.Contains(msg, "completed=1") {
			t.Fatalf("expected final snapshot in message, got %q", msg)
		}
		if !strings.Contains(strings.ToLower(msg), "encerrando a operacao deste ciclo") {
			t.Fatalf("expected final message to mention cycle cleanup, got %q", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("expected final master notification")
	}
}

func TestMasterTeamService_BuildExecutionStatusSnapshot_IgnoresHistoricalRuns(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{response: &ModelResponse{Content: "done"}}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	teamID, err := service.ensureTeam(ctx, "conversation-run-snapshot", "user-run-snapshot")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:            "old-task",
		TeamID:        teamID,
		RunID:         "run-old",
		Title:         "Old task",
		Prompt:        "old",
		Status:        TaskPending,
		AssignedAgent: strPtr("worker-a"),
	}, nil); err != nil {
		t.Fatalf("CreateTask(old-task) error = %v", err)
	}

	oldClaimed, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil || oldClaimed == nil {
		t.Fatalf("ClaimNextTask(old) error = %v claimed=%#v", err, oldClaimed)
	}
	if err := manager.FailTask(ctx, teamID, oldClaimed.ID, "worker-a", "old failure"); err != nil {
		t.Fatalf("FailTask(old) error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:            "new-task",
		TeamID:        teamID,
		RunID:         "run-new",
		Title:         "New task",
		Prompt:        "new",
		Status:        TaskPending,
		AssignedAgent: strPtr("worker-a"),
	}, nil); err != nil {
		t.Fatalf("CreateTask(new-task) error = %v", err)
	}

	newClaimed, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil || newClaimed == nil {
		t.Fatalf("ClaimNextTask(new) error = %v claimed=%#v", err, newClaimed)
	}
	if err := manager.CompleteTask(ctx, teamID, newClaimed.ID, "worker-a", "new done"); err != nil {
		t.Fatalf("CompleteTask(new) error = %v", err)
	}

	snapshot, err := service.BuildExecutionStatusSnapshot(ctx, "conversation-run-snapshot", "run-new")
	if err != nil {
		t.Fatalf("BuildExecutionStatusSnapshot() error = %v", err)
	}
	if snapshot.Completed != 1 || snapshot.Failed != 0 || snapshot.TotalTasks != 1 {
		t.Fatalf("unexpected execution snapshot: %#v", snapshot)
	}
}

func TestMasterTeamService_BuildExecutionStatusSnapshot_CollapsesFailedTaskWithSuccessfulRecovery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &queuedLLMProvider{
		errors: []error{context.DeadlineExceeded, nil},
		responses: []*ModelResponse{
			{Content: "recovery completed"},
		},
	}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	runCtx := WithRunContext(ctx, "run-collapse")
	if _, err := service.Spawn(runCtx, "conversation-collapse", "user-collapse", "reviewer", "Reviews work", "task failing first"); err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	deadline := time.After(2 * time.Second)
	for {
		snapshot, err := service.BuildExecutionStatusSnapshot(ctx, "conversation-collapse", "run-collapse")
		if err != nil {
			t.Fatalf("BuildExecutionStatusSnapshot() error = %v", err)
		}
		if snapshot.Completed == 1 && snapshot.Failed == 0 && snapshot.TotalTasks == 1 {
			break
		}

		select {
		case <-deadline:
			t.Fatalf("unexpected collapsed execution snapshot: %#v", snapshot)
		case <-time.After(25 * time.Millisecond):
		}
	}
}

func TestMasterTeamService_FormatMasterNotification_ClassifiesTotalSuccess(t *testing.T) {
	t.Parallel()

	service, err := NewMasterTeamService(newRuntimeTestManager(t), &MockLLMProvider{}, NewToolRegistry(), 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	msg := service.formatMasterNotification(TeamStatusSnapshot{
		Pending:    0,
		Running:    0,
		Blocked:    0,
		Completed:  2,
		Failed:     0,
		Cancelled:  0,
		TotalTasks: 2,
	}, 2, []string{"- `builder`: all done"})

	if !strings.Contains(strings.ToLower(msg), "sucesso total") {
		t.Fatalf("expected total success classification, got %q", msg)
	}
	if !strings.Contains(strings.ToLower(msg), "limpei o estado transitorio") {
		t.Fatalf("expected final message to mention cleanup, got %q", msg)
	}
}

func TestMasterTeamService_FormatMasterNotification_ClassifiesPartialCompletion(t *testing.T) {
	t.Parallel()

	service, err := NewMasterTeamService(newRuntimeTestManager(t), &MockLLMProvider{}, NewToolRegistry(), 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	msg := service.formatMasterNotification(TeamStatusSnapshot{
		Pending:    0,
		Running:    0,
		Blocked:    0,
		Completed:  1,
		Failed:     1,
		Cancelled:  1,
		TotalTasks: 3,
	}, 3, []string{"- `builder`: partial result"})

	if !strings.Contains(strings.ToLower(msg), "conclusao parcial") {
		t.Fatalf("expected partial classification, got %q", msg)
	}
	if !strings.Contains(strings.ToLower(msg), "parei os workers") {
		t.Fatalf("expected partial final message to mention workers shutdown, got %q", msg)
	}
}

func TestMasterTeamService_FormatMasterNotification_ClassifiesTerminalBlock(t *testing.T) {
	t.Parallel()

	service, err := NewMasterTeamService(newRuntimeTestManager(t), &MockLLMProvider{}, NewToolRegistry(), 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	msg := service.formatMasterNotification(TeamStatusSnapshot{
		Pending:    0,
		Running:    0,
		Blocked:    0,
		Completed:  0,
		Failed:     1,
		Cancelled:  2,
		TotalTasks: 3,
	}, 3, []string{"- `reviewer` falhou: blocked forever"})

	if !strings.Contains(strings.ToLower(msg), "bloqueio terminal") {
		t.Fatalf("expected terminal classification, got %q", msg)
	}
	if !strings.Contains(strings.ToLower(msg), "encerrando a operacao deste ciclo") {
		t.Fatalf("expected terminal final message to mention cycle closure, got %q", msg)
	}
}

func TestMasterTeamService_BuildStatusSnapshot_CountsBlockedTasks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{response: &ModelResponse{Content: "done"}}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	teamID, err := service.ensureTeam(ctx, "conversation-blocked-snapshot", "user-blocked-snapshot")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{ID: "task-a", TeamID: teamID, Title: "A", Prompt: "A", Status: TaskPending}, nil); err != nil {
		t.Fatalf("CreateTask(task-a) error = %v", err)
	}
	if err := manager.CreateTask(ctx, TeamTask{ID: "task-b", TeamID: teamID, Title: "B", Prompt: "B", Status: TaskPending}, []string{"task-a"}); err != nil {
		t.Fatalf("CreateTask(task-b) error = %v", err)
	}

	snapshot, err := service.BuildStatusSnapshot(ctx, "conversation-blocked-snapshot")
	if err != nil {
		t.Fatalf("BuildStatusSnapshot() error = %v", err)
	}
	if snapshot.Blocked != 1 {
		t.Fatalf("expected blocked=1, got %#v", snapshot)
	}
}

func TestMasterTeamService_RehydrateWorkerLoopsFromPersistentTeams(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store, err := NewSQLiteTaskStore(filepath.Join(t.TempDir(), "rehydrate.db"))
	if err != nil {
		t.Fatalf("NewSQLiteTaskStore() error = %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	manager, err := NewTeamManager(store)
	if err != nil {
		t.Fatalf("NewTeamManager() error = %v", err)
	}

	teamID, err := manager.CreateTeam(ctx, "conversation-rehydrate", "user-rh", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}
	if err := manager.RegisterTeammate(ctx, teamID, "builder", "Builds things"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}
	if err := manager.CreateTask(ctx, TeamTask{
		ID:            "rehydrate-task",
		TeamID:        teamID,
		Title:         "Resume me",
		Prompt:        "Continue after restart",
		Status:        TaskPending,
		AssignedAgent: strPtr("builder"),
	}, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	llm := &MockLLMProvider{response: &ModelResponse{Content: "rehydrated completion"}}
	registry := NewToolRegistry()
	notifications := make(chan string, 2)

	service, err := NewMasterTeamService(manager, llm, registry, 3, func(teamKey string, message string) {
		notifications <- teamKey + "|" + message
	})
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	if err := service.Rehydrate(ctx); err != nil {
		t.Fatalf("Rehydrate() error = %v", err)
	}

	select {
	case msg := <-notifications:
		if !strings.Contains(msg, "conversation-rehydrate|") {
			t.Fatalf("expected rehydrated notification to preserve team key, got %q", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("expected worker loop to resume after rehydrate")
	}
}

func TestMasterTeamService_Rehydrate_RequeuesOrphanRunningTask(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store, err := NewSQLiteTaskStore(filepath.Join(t.TempDir(), "rehydrate-running.db"))
	if err != nil {
		t.Fatalf("NewSQLiteTaskStore() error = %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	manager, err := NewTeamManager(store)
	if err != nil {
		t.Fatalf("NewTeamManager() error = %v", err)
	}

	teamID, err := manager.CreateTeam(ctx, "conversation-running", "user-running", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}
	if err := manager.RegisterTeammate(ctx, teamID, "builder", "Builds things"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}
	if err := manager.CreateTask(ctx, TeamTask{
		ID:            "running-task",
		TeamID:        teamID,
		Title:         "Resume running task",
		Prompt:        "Continue after crash",
		Status:        TaskPending,
		AssignedAgent: strPtr("builder"),
	}, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	claimed, err := manager.ClaimNextTask(ctx, teamID, "builder")
	if err != nil || claimed == nil {
		t.Fatalf("ClaimNextTask() error = %v claimed=%#v", err, claimed)
	}

	llm := &MockLLMProvider{response: &ModelResponse{Content: "resumed completion"}}
	registry := NewToolRegistry()
	notifications := make(chan string, 2)

	service, err := NewMasterTeamService(manager, llm, registry, 3, func(teamKey string, message string) {
		notifications <- teamKey + "|" + message
	})
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	if err := service.Rehydrate(ctx); err != nil {
		t.Fatalf("Rehydrate() error = %v", err)
	}

	select {
	case <-notifications:
	case <-time.After(2 * time.Second):
		t.Fatalf("expected rehydrated running task to be resumed")
	}

	task, err := manager.GetTask(ctx, teamID, "running-task")
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if task == nil || task.Status != TaskCompleted {
		t.Fatalf("expected running task to finish after rehydrate, got %#v", task)
	}
}

func TestMasterTeamService_Rehydrate_IsIdempotentForWorkerLoops(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store, err := NewSQLiteTaskStore(filepath.Join(t.TempDir(), "rehydrate-idempotent.db"))
	if err != nil {
		t.Fatalf("NewSQLiteTaskStore() error = %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	manager, err := NewTeamManager(store)
	if err != nil {
		t.Fatalf("NewTeamManager() error = %v", err)
	}

	teamID, err := manager.CreateTeam(ctx, "conversation-idempotent", "user-idempotent", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}
	if err := manager.RegisterTeammate(ctx, teamID, "builder", "Builds things"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}
	if err := manager.CreateTask(ctx, TeamTask{
		ID:            "pending-task",
		TeamID:        teamID,
		Title:         "Pending task",
		Prompt:        "Resume after restart",
		Status:        TaskPending,
		AssignedAgent: strPtr("builder"),
	}, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	llm := &MockLLMProvider{response: &ModelResponse{Content: "done"}}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	if err := service.Rehydrate(ctx); err != nil {
		t.Fatalf("first Rehydrate() error = %v", err)
	}
	if err := service.Rehydrate(ctx); err != nil {
		t.Fatalf("second Rehydrate() error = %v", err)
	}

	time.Sleep(250 * time.Millisecond)

	service.mu.Lock()
	defer service.mu.Unlock()
	if len(service.workerLoops) > 1 {
		t.Fatalf("expected rehydrate to avoid duplicate worker loops, got %d", len(service.workerLoops))
	}
}

func TestWorkerRuntime_RunOnce_RenewsLeaseWhileTaskIsRunning(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	rawManager := newRuntimeTestManager(t)
	manager := rawManager.(*DefaultTeamManager)
	manager.store.leaseDuration = 80 * time.Millisecond

	teamID, err := manager.CreateTeam(ctx, "team-rt-heartbeat", "user-rt-heartbeat", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Runs long task"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-heartbeat",
		TeamID: teamID,
		Title:  "Long task",
		Prompt: "Take some time",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	executor := &slowTaskExecutor{
		started: make(chan struct{}),
		release: make(chan struct{}),
		result:  "finished",
	}
	worker := NewWorkerRuntime("worker-a", manager, executor)

	done := make(chan error, 1)
	go func() {
		_, err := worker.RunOnce(ctx, teamID)
		done <- err
	}()

	<-executor.started
	time.Sleep(140 * time.Millisecond)

	var leaseExpiresAt time.Time
	var status string
	err = manager.store.db.QueryRowContext(ctx,
		`SELECT status, lease_expires_at FROM team_members WHERE team_id = ? AND agent_name = ?`,
		teamID, "worker-a",
	).Scan(&status, &leaseExpiresAt)
	if err != nil {
		t.Fatalf("query team member lease error = %v", err)
	}
	if status != "active" {
		t.Fatalf("expected active worker status during task, got %q", status)
	}
	if !leaseExpiresAt.After(time.Now().UTC()) {
		t.Fatalf("expected heartbeat to renew lease into the future, got %v", leaseExpiresAt)
	}

	close(executor.release)

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("RunOnce() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("expected worker to finish after release")
	}
}

func TestMasterTeamService_PauseResumeAndCancelTeam(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	llm := &MockLLMProvider{response: &ModelResponse{Content: "done"}}
	registry := NewToolRegistry()

	service, err := NewMasterTeamService(manager, llm, registry, 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	taskID, err := service.Spawn(ctx, "conversation-control", "user-control", "builder", "Builds things", "Finish setup")
	if err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	teamID, err := service.ensureTeam(ctx, "conversation-control", "user-control")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}

	if err := service.Pause(ctx, "conversation-control"); err != nil {
		t.Fatalf("Pause() error = %v", err)
	}
	status, err := manager.GetTeamStatus(ctx, teamID)
	if err != nil {
		t.Fatalf("GetTeamStatus() error = %v", err)
	}
	if status != "paused" {
		t.Fatalf("expected paused team status, got %q", status)
	}

	if err := service.Resume(ctx, "conversation-control"); err != nil {
		t.Fatalf("Resume() error = %v", err)
	}
	status, err = manager.GetTeamStatus(ctx, teamID)
	if err != nil {
		t.Fatalf("GetTeamStatus() error = %v", err)
	}
	if status != "active" {
		t.Fatalf("expected active team status after resume, got %q", status)
	}

	if err := service.Cancel(ctx, "conversation-control", "usuario cancelou"); err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
	status, err = manager.GetTeamStatus(ctx, teamID)
	if err != nil {
		t.Fatalf("GetTeamStatus() error = %v", err)
	}
	if status != "cancelled" {
		t.Fatalf("expected cancelled team status, got %q", status)
	}

	task, err := manager.GetTask(ctx, teamID, taskID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if task == nil {
		t.Fatal("expected task to exist")
	}
	if task.Status != TaskCancelled && task.Status != TaskCompleted {
		t.Fatalf("expected task to be cancelled or already completed, got %q", task.Status)
	}
}

func TestMasterTeamService_FinalizeTeamRunIfIdle_CompletesAndCleansRuntimeState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	service, err := NewMasterTeamService(manager, &MockLLMProvider{response: &ModelResponse{Content: "done"}}, NewToolRegistry(), 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	teamID, err := service.ensureTeam(ctx, "conversation-finalize", "user-finalize")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}
	if err := service.ensureTeammate(ctx, teamID, "builder", "Builds things"); err != nil {
		t.Fatalf("ensureTeammate() error = %v", err)
	}
	if err := manager.CreateTask(ctx, TeamTask{
		ID:            "task-finalize",
		TeamID:        teamID,
		RunID:         "run-finalize",
		Title:         "Finish work",
		Prompt:        "Done",
		Status:        TaskCompleted,
		AssignedAgent: strPtr("builder"),
		ResultSummary: "ok",
	}, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	cancelled := false
	service.mu.Lock()
	service.workerLoops[teamID+"::builder"] = &workerLoopHandle{
		wake: make(chan struct{}, 1),
		cancel: func() {
			cancelled = true
		},
	}
	service.recoveryCount["task-finalize"] = 1
	service.mu.Unlock()

	snapshot := TeamStatusSnapshot{
		TeamKey:    "conversation-finalize",
		TeamID:     teamID,
		TeamStatus: TeamStatusActive,
		Completed:  1,
		TotalTasks: 1,
	}

	service.finalizeTeamRunIfIdle(ctx, "conversation-finalize", "run-finalize", &snapshot)

	status, err := manager.GetTeamStatus(ctx, teamID)
	if err != nil {
		t.Fatalf("GetTeamStatus() error = %v", err)
	}
	if status != TeamStatusCompleted {
		t.Fatalf("expected completed team status, got %q", status)
	}
	if snapshot.TeamStatus != TeamStatusCompleted {
		t.Fatalf("expected snapshot status to be updated, got %q", snapshot.TeamStatus)
	}
	if !cancelled {
		t.Fatal("expected worker loop to be cancelled")
	}

	service.mu.Lock()
	defer service.mu.Unlock()
	if len(service.workerLoops) != 0 {
		t.Fatalf("expected worker loops to be cleared, got %d", len(service.workerLoops))
	}
	if _, ok := service.recoveryCount["task-finalize"]; ok {
		t.Fatal("expected recovery count for finalized run to be cleared")
	}
}

func TestMasterTeamService_Spawn_ReactivatesCompletedTeam(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newRuntimeTestManager(t)
	service, err := NewMasterTeamService(manager, &MockLLMProvider{response: &ModelResponse{Content: "done"}}, NewToolRegistry(), 3, nil)
	if err != nil {
		t.Fatalf("NewMasterTeamService() error = %v", err)
	}

	teamID, err := service.ensureTeam(ctx, "conversation-reactivate", "user-reactivate")
	if err != nil {
		t.Fatalf("ensureTeam() error = %v", err)
	}
	if err := manager.SetTeamStatus(ctx, teamID, TeamStatusCompleted); err != nil {
		t.Fatalf("SetTeamStatus() error = %v", err)
	}

	if _, err := service.Spawn(ctx, "conversation-reactivate", "user-reactivate", "builder", "Builds things", "Start new run"); err != nil {
		t.Fatalf("Spawn() error = %v", err)
	}

	status, err := manager.GetTeamStatus(ctx, teamID)
	if err != nil {
		t.Fatalf("GetTeamStatus() error = %v", err)
	}
	if status != TeamStatusActive {
		t.Fatalf("expected team to be reactivated on spawn, got %q", status)
	}
}
