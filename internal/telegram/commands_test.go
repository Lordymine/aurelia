package telegram

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/session"
)

func TestMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want *CommandType // nil means no match
	}{
		// --- cron_create ---
		{name: "agenda with details", text: "agenda uma reunião amanhã às 10h", want: cmdPtr(CmdCronCreate)},
		{name: "agendar keyword", text: "agendar lembrete pra sexta", want: cmdPtr(CmdCronCreate)},
		{name: "cria um lembrete", text: "cria um lembrete pra amanhã", want: cmdPtr(CmdCronCreate)},
		{name: "me lembra", text: "me lembra de revisar o PR às 15h", want: cmdPtr(CmdCronCreate)},

		// --- cron_list ---
		{name: "meus agendamentos", text: "meus agendamentos", want: cmdPtr(CmdCronList)},
		{name: "o que ta agendado", text: "o que tá agendado?", want: cmdPtr(CmdCronList)},
		{name: "lista agendamentos", text: "lista agendamentos", want: cmdPtr(CmdCronList)},

		// --- cron_cancel ---
		{name: "cancela agendamento", text: "cancela o agendamento abc123", want: cmdPtr(CmdCronCancel)},
		{name: "cancele agendamento", text: "cancele o agendamento das 7h", want: cmdPtr(CmdCronCancel)},
		{name: "remove lembrete", text: "remove o lembrete de reunião", want: cmdPtr(CmdCronCancel)},
		{name: "desativa agendamento", text: "desativa agendamento abc", want: cmdPtr(CmdCronCancel)},
		{name: "exclui agendamento", text: "exclui agendamento abc123", want: cmdPtr(CmdCronCancel)},
		{name: "apaga agendamento", text: "apaga agendamento abc123", want: cmdPtr(CmdCronCancel)},

		// --- session_reset ---
		{name: "nova conversa", text: "nova conversa", want: cmdPtr(CmdSessionReset)},
		{name: "limpa o contexto", text: "limpa o contexto", want: cmdPtr(CmdSessionReset)},
		{name: "reset", text: "reset", want: cmdPtr(CmdSessionReset)},
		{name: "comeca de novo", text: "começa de novo", want: cmdPtr(CmdSessionReset)},

		// --- status ---
		{name: "status", text: "status", want: cmdPtr(CmdStatus)},

		// --- list_agents ---
		{name: "quais agents", text: "quais agents?", want: cmdPtr(CmdListAgents)},
		{name: "lista agents", text: "lista agents", want: cmdPtr(CmdListAgents)},
		{name: "meus agents", text: "meus agents", want: cmdPtr(CmdListAgents)},

		// --- list_models ---
		{name: "quais modelos", text: "quais modelos?", want: cmdPtr(CmdListModels)},
		{name: "lista modelos", text: "lista modelos", want: cmdPtr(CmdListModels)},
		{name: "lista provedores", text: "lista provedores", want: cmdPtr(CmdListModels)},

		// --- NO match: normal conversation ---
		{name: "greeting", text: "bom dia", want: nil},
		{name: "question", text: "como funciona o bridge?", want: nil},
		{name: "code request", text: "escreve um teste pro handler", want: nil},

		// --- NO match: narrative context (anti-false-positive) ---
		{name: "narrative agendar", text: "ontem eu tentei agendar uma reunião", want: nil},
		{name: "narrative lembrete", text: "ele me lembrou de fazer o deploy", want: nil},
		{name: "narrative status", text: "o status do PR tá verde", want: nil},
		{name: "narrative reset", text: "depois do reset o servidor voltou", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := MatchCommand(tt.text)

			if tt.want == nil {
				if got != nil {
					t.Fatalf("MatchCommand(%q) = %v, want nil", tt.text, got.Type)
				}
				return
			}

			if got == nil {
				t.Fatalf("MatchCommand(%q) = nil, want %v", tt.text, *tt.want)
			}
			if got.Type != *tt.want {
				t.Fatalf("MatchCommand(%q).Type = %v, want %v", tt.text, got.Type, *tt.want)
			}
		})
	}
}

func cmdPtr(c CommandType) *CommandType { return &c }

// --- T4: session_reset tests ---

func TestCmdSessionReset(t *testing.T) {
	t.Parallel()

	sessions := session.NewStore()
	tracker := session.NewTracker()
	sessions.Set(42, "sess-abc")
	tracker.Add(42, 1000, 500, 1, 0.01)

	bc := &BotController{
		config:   &config.AppConfig{Providers: map[string]config.ProviderConfig{}},
		sessions: sessions,
		tracker:  tracker,
	}

	reply, err := bc.cmdSessionReset(42)
	if err != nil {
		t.Fatalf("cmdSessionReset() error = %v", err)
	}

	// Session should be cleared
	if sid := sessions.Get(42); sid != "" {
		t.Fatalf("session should be cleared, got %q", sid)
	}

	// Tracker should be cleared
	usage := tracker.Get(42)
	if usage.NumTurns != 0 {
		t.Fatalf("tracker should be cleared, got %d turns", usage.NumTurns)
	}

	if reply == "" {
		t.Fatal("expected non-empty reply")
	}
}

// --- T5: cron_list tests ---

func TestCmdCronList_WithJobs(t *testing.T) {
	t.Parallel()

	service := &fakeCronCommandService{
		jobs: []cron.CronJob{
			{ID: "abc12345-full-uuid", ScheduleType: "cron", CronExpr: "0 9 * * *", Prompt: "bom dia", Active: true, LastStatus: "idle"},
			{ID: "def67890-full-uuid", ScheduleType: "once", Prompt: "lembrete", Active: true, LastStatus: "pending"},
		},
	}

	bc := &BotController{
		config:      &config.AppConfig{Providers: map[string]config.ProviderConfig{}},
		cronHandler: NewCronCommandHandler(service),
	}

	reply, err := bc.cmdCronList(42)
	if err != nil {
		t.Fatalf("cmdCronList() error = %v", err)
	}

	if reply == "" {
		t.Fatal("expected non-empty reply")
	}
	// Should contain job info
	if !contains(reply, "abc12345") || !contains(reply, "bom dia") {
		t.Fatalf("reply should contain job info, got %q", reply)
	}
}

func TestCmdCronList_Empty(t *testing.T) {
	t.Parallel()

	service := &fakeCronCommandService{jobs: nil}
	bc := &BotController{
		config:      &config.AppConfig{Providers: map[string]config.ProviderConfig{}},
		cronHandler: NewCronCommandHandler(service),
	}

	reply, err := bc.cmdCronList(42)
	if err != nil {
		t.Fatalf("cmdCronList() error = %v", err)
	}
	if reply == "" {
		t.Fatal("expected non-empty reply even when no jobs")
	}
}

// --- T6: cron_cancel tests ---

func TestCmdCronCancel_WithID(t *testing.T) {
	t.Parallel()

	service := &fakeCronCommandService{}
	bc := &BotController{
		config:      &config.AppConfig{Providers: map[string]config.ProviderConfig{}},
		cronHandler: NewCronCommandHandler(service),
	}

	reply, err := bc.cmdCronCancel(42, "cancela agendamento abc123")
	if err != nil {
		t.Fatalf("cmdCronCancel() error = %v", err)
	}

	if len(service.deleteCalls) != 1 {
		t.Fatalf("expected 1 delete call, got %d", len(service.deleteCalls))
	}
	if service.deleteCalls[0] != "abc123" {
		t.Fatalf("expected delete of 'abc123', got %q", service.deleteCalls[0])
	}
	if reply == "" {
		t.Fatal("expected non-empty reply")
	}
}

func TestCmdCronCancel_NoID(t *testing.T) {
	t.Parallel()

	service := &fakeCronCommandService{}
	bc := &BotController{
		config:      &config.AppConfig{Providers: map[string]config.ProviderConfig{}},
		cronHandler: NewCronCommandHandler(service),
	}

	reply, err := bc.cmdCronCancel(42, "cancela agendamento")
	if err != nil {
		t.Fatalf("cmdCronCancel() error = %v", err)
	}

	// Should not attempt delete without an ID
	if len(service.deleteCalls) != 0 {
		t.Fatalf("expected no delete calls, got %d", len(service.deleteCalls))
	}
	if reply == "" {
		t.Fatal("expected guidance reply")
	}
}

// --- T8: cron_create tests ---

func TestParseCronCreateResponse_RecurringJSON(t *testing.T) {
	t.Parallel()

	raw := `{"type":"cron","cron_expr":"0 9 * * *","prompt":"revisar emails"}`
	parsed, err := parseCronCreateResponse(raw)
	if err != nil {
		t.Fatalf("parseCronCreateResponse() error = %v", err)
	}
	if parsed.Type != "cron" || parsed.CronExpr != "0 9 * * *" || parsed.Prompt != "revisar emails" {
		t.Fatalf("unexpected parsed: %+v", parsed)
	}
}

func TestParseCronCreateResponse_OnceJSON(t *testing.T) {
	t.Parallel()

	raw := `{"type":"once","run_at":"2026-03-27T15:00:00-03:00","prompt":"fazer deploy"}`
	parsed, err := parseCronCreateResponse(raw)
	if err != nil {
		t.Fatalf("parseCronCreateResponse() error = %v", err)
	}
	if parsed.Type != "once" || parsed.RunAt != "2026-03-27T15:00:00-03:00" || parsed.Prompt != "fazer deploy" {
		t.Fatalf("unexpected parsed: %+v", parsed)
	}
}

func TestParseCronCreateResponse_MarkdownFences(t *testing.T) {
	t.Parallel()

	raw := "```json\n{\"type\":\"cron\",\"cron_expr\":\"0 9 * * 1\",\"prompt\":\"standup\"}\n```"
	parsed, err := parseCronCreateResponse(raw)
	if err != nil {
		t.Fatalf("parseCronCreateResponse() error = %v", err)
	}
	if parsed.Type != "cron" || parsed.CronExpr != "0 9 * * 1" {
		t.Fatalf("unexpected parsed: %+v", parsed)
	}
}

func TestParseCronCreateResponse_MissingPrompt(t *testing.T) {
	t.Parallel()

	raw := `{"type":"cron","cron_expr":"0 9 * * *"}`
	_, err := parseCronCreateResponse(raw)
	if err == nil {
		t.Fatal("expected error for missing prompt")
	}
}

func TestParseCronCreateResponse_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := parseCronCreateResponse("not json at all")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// --- T9: status tests ---

func TestCmdStatus(t *testing.T) {
	t.Parallel()

	service := &fakeCronCommandService{
		jobs: []cron.CronJob{
			{ID: "j1", Active: true},
			{ID: "j2", Active: false},
		},
	}
	sessions := session.NewStore()
	sessions.Set(42, "sess-abc-12345678")

	bc := &BotController{
		config: &config.AppConfig{
			DefaultModel: "kimi-k2-thinking",
			Providers:    map[string]config.ProviderConfig{},
		},
		cronHandler: NewCronCommandHandler(service),
		sessions:    sessions,
		tracker:     session.NewTracker(),
	}

	reply, err := bc.cmdStatus(42)
	if err != nil {
		t.Fatalf("cmdStatus() error = %v", err)
	}
	if !strings.Contains(reply, "kimi-k2-thinking") {
		t.Fatalf("expected model in status, got %q", reply)
	}
	if !strings.Contains(reply, "1") { // 1 active job
		t.Fatalf("expected active job count in status, got %q", reply)
	}
	if !strings.Contains(reply, "sess-abc") {
		t.Fatalf("expected session ID in status, got %q", reply)
	}
}

// --- T10: list_agents tests ---

func TestCmdListAgents_WithAgents(t *testing.T) {
	t.Parallel()

	// Create a temp dir with agent files
	reg := buildTestRegistry(t, map[string]string{
		"coder":      "Writes and debugs code",
		"prospector": "Busca leads e prospecta clientes",
	})

	bc := &BotController{
		config: &config.AppConfig{Providers: map[string]config.ProviderConfig{}},
		agents: reg,
	}

	reply, err := bc.cmdListAgents()
	if err != nil {
		t.Fatalf("cmdListAgents() error = %v", err)
	}
	if !strings.Contains(reply, "coder") || !strings.Contains(reply, "prospector") {
		t.Fatalf("expected agent names, got %q", reply)
	}
}

func TestCmdListAgents_Empty(t *testing.T) {
	t.Parallel()

	bc := &BotController{
		config: &config.AppConfig{Providers: map[string]config.ProviderConfig{}},
	}

	reply, err := bc.cmdListAgents()
	if err != nil {
		t.Fatalf("cmdListAgents() error = %v", err)
	}
	if !strings.Contains(reply, "Nenhum") {
		t.Fatalf("expected 'nenhum' message, got %q", reply)
	}
}

// buildTestRegistry creates a Registry with agents for testing.
func buildTestRegistry(t *testing.T, agentMap map[string]string) *agents.Registry {
	t.Helper()
	dir := t.TempDir()
	for name, desc := range agentMap {
		content := fmt.Sprintf("---\nname: %s\ndescription: %s\n---\nYou are %s.", name, desc, name)
		path := dir + "/" + name + ".md"
		if err := writeTestFile(path, content); err != nil {
			t.Fatalf("failed to write agent file: %v", err)
		}
	}
	reg, err := agents.Load(dir)
	if err != nil {
		t.Fatalf("agents.Load() error = %v", err)
	}
	return reg
}

func writeTestFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// --- T11: list_models tests ---

func TestCmdListModels_WithProviders(t *testing.T) {
	t.Parallel()

	bc := &BotController{
		config: &config.AppConfig{
			DefaultModel: "kimi-k2-thinking",
			Providers: map[string]config.ProviderConfig{
				"anthropic": {APIKey: "sk-test"},
				"kimi":      {APIKey: "kimi-key"},
				"google":    {APIKey: ""},
			},
		},
	}

	reply, err := bc.cmdListModels()
	if err != nil {
		t.Fatalf("cmdListModels() error = %v", err)
	}
	if !strings.Contains(reply, "anthropic") || !strings.Contains(reply, "kimi") {
		t.Fatalf("expected provider names, got %q", reply)
	}
	if !strings.Contains(reply, "sem API key") {
		t.Fatalf("expected 'sem API key' for google, got %q", reply)
	}
}

func TestCmdListModels_NoProviders(t *testing.T) {
	t.Parallel()

	bc := &BotController{
		config: &config.AppConfig{
			Providers: map[string]config.ProviderConfig{},
		},
	}

	reply, err := bc.cmdListModels()
	if err != nil {
		t.Fatalf("cmdListModels() error = %v", err)
	}
	if !strings.Contains(reply, "Nenhum") {
		t.Fatalf("expected 'nenhum' message, got %q", reply)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && stringContains(s, substr))
}

func stringContains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
