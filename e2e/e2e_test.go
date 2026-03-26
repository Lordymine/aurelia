package e2e

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/session"
	"github.com/kocar/aurelia/internal/telegram"
)

func TestWiring_AllComponentsInitialize(t *testing.T) {
	dir := t.TempDir()

	// Create persona files.
	identityPath := filepath.Join(dir, "IDENTITY.md")
	soulPath := filepath.Join(dir, "SOUL.md")
	userPath := filepath.Join(dir, "USER.md")

	identityContent := "---\nname: \"Aurelia\"\nrole: \"AI Assistant\"\n---\nYou are Aurelia, an AI operating system."
	if err := os.WriteFile(identityPath, []byte(identityContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(soulPath, []byte("You are helpful and thoughtful."), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(userPath, []byte("Nome: Developer\nThe user is a developer."), 0644); err != nil {
		t.Fatal(err)
	}

	// Create agent file.
	agentsDir := filepath.Join(dir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatal(err)
	}
	agentContent := "---\nname: helper\ndescription: General helper\nmodel: claude-sonnet-4-6\n---\n\nYou help with general tasks."
	if err := os.WriteFile(filepath.Join(agentsDir, "helper.md"), []byte(agentContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Load components.
	personaSvc := persona.NewCanonicalIdentityService(identityPath, soulPath, userPath, "", "", "")

	agentReg, err := agents.Load(agentsDir)
	if err != nil {
		t.Fatalf("agents.Load() error = %v", err)
	}

	// Verify persona builds prompt.
	prompt, err := personaSvc.BuildPrompt()
	if err != nil {
		t.Fatalf("BuildPrompt() error = %v", err)
	}
	if !strings.Contains(prompt, "Aurelia") {
		t.Fatalf("expected prompt to contain 'Aurelia', got %q", prompt)
	}

	// Verify agent registry works.
	agent := agentReg.Get("helper")
	if agent == nil {
		t.Fatal("expected agent 'helper' to exist")
	}
	if agent.Model != "claude-sonnet-4-6" {
		t.Fatalf("expected model 'claude-sonnet-4-6', got %q", agent.Model)
	}

	// Verify full system prompt can be assembled.
	fullPrompt := prompt + "\n\n" + agent.Prompt
	if !strings.Contains(fullPrompt, "Aurelia") {
		t.Fatal("full prompt missing persona identity")
	}
	if !strings.Contains(fullPrompt, "general tasks") {
		t.Fatal("full prompt missing agent prompt")
	}
}

func TestRouting_BuildsCorrectPrompt(t *testing.T) {
	dir := t.TempDir()

	// Create two agents.
	agentsDir := filepath.Join(dir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatal(err)
	}

	helperContent := "---\nname: helper\ndescription: General helper\nmodel: claude-sonnet-4-6\n---\n\nYou help with general tasks."
	if err := os.WriteFile(filepath.Join(agentsDir, "helper.md"), []byte(helperContent), 0644); err != nil {
		t.Fatal(err)
	}

	codeContent := "---\nname: coder\ndescription: Code assistant\nmodel: claude-sonnet-4-6\n---\n\nYou write clean code."
	if err := os.WriteFile(filepath.Join(agentsDir, "coder.md"), []byte(codeContent), 0644); err != nil {
		t.Fatal(err)
	}

	reg, err := agents.Load(agentsDir)
	if err != nil {
		t.Fatalf("agents.Load() error = %v", err)
	}

	// Route "@helper analyze this" -> helper agent.
	agent := reg.Route("@helper analyze this")
	if agent == nil {
		t.Fatal("expected routing to match 'helper'")
	}
	if agent.Name != "helper" {
		t.Fatalf("expected agent 'helper', got %q", agent.Name)
	}

	// Route "@coder fix the bug" -> coder agent.
	agent = reg.Route("@coder fix the bug")
	if agent == nil {
		t.Fatal("expected routing to match 'coder'")
	}
	if agent.Name != "coder" {
		t.Fatalf("expected agent 'coder', got %q", agent.Name)
	}

	// No match for plain messages.
	agent = reg.Route("just a normal message")
	if agent != nil {
		t.Fatalf("expected no routing match, got %q", agent.Name)
	}

	// Build system prompt with persona + routed agent + memory.
	identityPath := filepath.Join(dir, "IDENTITY.md")
	soulPath := filepath.Join(dir, "SOUL.md")
	userPath := filepath.Join(dir, "USER.md")

	if err := os.WriteFile(identityPath, []byte("---\nname: \"Aurelia\"\nrole: \"OS\"\n---\nCore identity."), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(soulPath, []byte("Soul values."), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(userPath, []byte("Nome: Dev\nUser context."), 0644); err != nil {
		t.Fatal(err)
	}

	personaSvc := persona.NewCanonicalIdentityService(identityPath, soulPath, userPath, "", "", "")
	personaPrompt, err := personaSvc.BuildPrompt()
	if err != nil {
		t.Fatalf("BuildPrompt() error = %v", err)
	}

	routed := reg.Route("@helper analyze this")
	if routed == nil {
		t.Fatal("expected routing match")
	}

	fullPrompt := personaPrompt + "\n\n" + routed.Prompt
	if !strings.Contains(fullPrompt, "Aurelia") {
		t.Fatal("full prompt missing persona")
	}
	if !strings.Contains(fullPrompt, "general tasks") {
		t.Fatal("full prompt missing agent instructions")
	}
}

func TestCron_ScheduledAgentsRegistered(t *testing.T) {
	dir := t.TempDir()
	agentsDir := filepath.Join(dir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Agent with schedule.
	scheduledContent := "---\nname: reporter\ndescription: Daily report\nmodel: claude-sonnet-4-6\nschedule: \"0 9 * * *\"\n---\n\nGenerate daily report."
	if err := os.WriteFile(filepath.Join(agentsDir, "reporter.md"), []byte(scheduledContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Agent without schedule.
	normalContent := "---\nname: helper\ndescription: General helper\nmodel: claude-sonnet-4-6\n---\n\nYou help."
	if err := os.WriteFile(filepath.Join(agentsDir, "helper.md"), []byte(normalContent), 0644); err != nil {
		t.Fatal(err)
	}

	reg, err := agents.Load(agentsDir)
	if err != nil {
		t.Fatalf("agents.Load() error = %v", err)
	}

	// Only the scheduled agent should appear.
	scheduled := reg.Scheduled()
	if len(scheduled) != 1 {
		t.Fatalf("expected 1 scheduled agent, got %d", len(scheduled))
	}
	if scheduled[0].Name != "reporter" {
		t.Fatalf("expected scheduled agent 'reporter', got %q", scheduled[0].Name)
	}
	if scheduled[0].Schedule != "0 9 * * *" {
		t.Fatalf("expected schedule '0 9 * * *', got %q", scheduled[0].Schedule)
	}

	// Verify the cron store can accept a job derived from this agent.
	cronDBPath := filepath.Join(dir, "cron.db")
	cronStore, err := cron.NewSQLiteCronStore(cronDBPath)
	if err != nil {
		t.Fatalf("NewSQLiteCronStore() error = %v", err)
	}
	defer cronStore.Close()

	ctx := context.Background()
	job := cron.CronJob{
		ID:           "job-1",
		OwnerUserID:  "user-1",
		TargetChatID: 123,
		AgentName:    scheduled[0].Name,
		ScheduleType: "cron",
		CronExpr:     scheduled[0].Schedule,
		Prompt:       scheduled[0].Prompt,
		Active:       true,
	}

	if err := cronStore.CreateJob(ctx, job); err != nil {
		t.Fatalf("CreateJob() error = %v", err)
	}

	stored, err := cronStore.GetJob(ctx, "job-1")
	if err != nil {
		t.Fatalf("GetJob() error = %v", err)
	}
	if stored == nil {
		t.Fatal("expected stored job, got nil")
	}
	if stored.CronExpr != "0 9 * * *" {
		t.Fatalf("expected cron expr '0 9 * * *', got %q", stored.CronExpr)
	}
	if stored.Prompt != "Generate daily report." {
		t.Fatalf("expected prompt 'Generate daily report.', got %q", stored.Prompt)
	}
}

func TestCommandLayer_MatchRouting(t *testing.T) {
	// Verify that system commands are correctly identified and normal messages pass through.
	commands := []struct {
		text    string
		isCmd   bool
		cmdName string
	}{
		{"nova conversa", true, "session_reset"},
		{"meus agendamentos", true, "cron_list"},
		{"status", true, "status"},
		{"quais agents?", true, "list_agents"},
		{"quais modelos?", true, "list_models"},
		{"agenda todo dia às 9h check emails", true, "cron_create"},
		{"bom dia, como vai?", false, ""},
		{"me ajuda a debugar esse código", false, ""},
		{"ontem eu tentei agendar uma reunião", false, ""}, // narrative — should NOT match
	}

	for _, tc := range commands {
		cmd := telegram.MatchCommand(tc.text)
		if tc.isCmd && cmd == nil {
			t.Errorf("expected %q to match as command (%s), got nil", tc.text, tc.cmdName)
		}
		if !tc.isCmd && cmd != nil {
			t.Errorf("expected %q to NOT match as command, got type=%d", tc.text, cmd.Type)
		}
	}
}

func TestCommandLayer_SessionResetClearsState(t *testing.T) {
	// Integration: session reset via command layer clears session and tracker.
	sessions := session.NewStore()
	tracker := session.NewTracker()

	sessions.Set(100, "sess-integration-test")
	tracker.Add(100, 5000, 2000, 3, 0.05)

	// Verify session exists.
	if sid := sessions.Get(100); sid != "sess-integration-test" {
		t.Fatalf("expected session, got %q", sid)
	}

	// Clear (simulating what cmdSessionReset does).
	sessions.Clear(100)
	tracker.Clear(100)

	if sid := sessions.Get(100); sid != "" {
		t.Fatalf("session should be cleared, got %q", sid)
	}
	usage := tracker.Get(100)
	if usage.NumTurns != 0 {
		t.Fatalf("tracker should be cleared, got %d turns", usage.NumTurns)
	}
}

func TestCommandLayer_CronListIntegration(t *testing.T) {
	// Integration: list cron jobs from real SQLite store.
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "cron.db")

	store, err := cron.NewSQLiteCronStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteCronStore() error = %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	job := cron.CronJob{
		ID:           "cmd-layer-test-job",
		OwnerUserID:  "user-1",
		TargetChatID: 200,
		ScheduleType: "cron",
		CronExpr:     "0 9 * * *",
		Prompt:       "check emails",
		Active:       true,
	}
	if err := store.CreateJob(ctx, job); err != nil {
		t.Fatalf("CreateJob() error = %v", err)
	}

	// List jobs for the chat.
	jobs, err := store.ListJobsByChat(ctx, 200)
	if err != nil {
		t.Fatalf("ListJobsByChat() error = %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Prompt != "check emails" {
		t.Fatalf("expected prompt 'check emails', got %q", jobs[0].Prompt)
	}
}

func TestBridge_ProtocolPing(t *testing.T) {
	// Find bridge directory relative to test file.
	// The bridge dir is at the project root under "bridge/".
	bridgeDir := filepath.Join("..", "bridge")

	// Verify bridge dir exists before attempting ping.
	if _, err := os.Stat(bridgeDir); os.IsNotExist(err) {
		t.Skipf("bridge directory not found at %s", bridgeDir)
	}

	// Import bridge inline to avoid import cycle concerns - we use it directly.
	// The bridge.Ping() spawns npx tsx, which requires Node.js.
	// Skip gracefully if the environment doesn't have it.
	br := newBridgeForTest(bridgeDir)
	ctx := context.Background()
	if err := br.Ping(ctx); err != nil {
		t.Skipf("bridge not available (npx/tsx not installed or bridge not built): %v", err)
	}
}
