# Session Continuation + Hybrid Agents + Scheduler Improvements

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add session continuation to reduce token usage, pass agents to the SDK for native delegation, and harden the scheduler with full cron parsing and run locking.

**Architecture:** Three independent changes that touch distinct layers. Session continuation changes the Telegram session store + bridge protocol. Hybrid agents adds an `Agents` field to bridge requests and builds SDK agent definitions from the registry. Scheduler improvements swap the manual cron parser for `robfig/cron/v3` and add run locking + cost tracking.

**Tech Stack:** Go 1.25 + `robfig/cron/v3` + `@anthropic-ai/claude-agent-sdk` (existing, verified: `continue?: boolean` at sdk.d.ts:827, `agents?: Record<string, AgentDefinition>` at sdk.d.ts:811)

**Pre-requisite:** The cloud MCP fix (bundle.ts rename + loadCloudMCPs) must be committed first. This plan assumes that work is already merged.

---

## File Map

### Session Continuation
| Action | File | Responsibility |
|--------|------|----------------|
| Modify | `internal/bridge/protocol.go` | Add `Continue` field to `RequestOptions` |
| Modify | `bridge/index.ts` | Map `continue` option to SDK |
| Modify | `internal/bridge/bundle.ts` | Mirror of index.ts |
| Modify | `internal/telegram/sessions.go` | Add `IsActive` flag to track warm sessions |
| Modify | `internal/telegram/input_pipeline.go` | Use `continue` for warm sessions, `resume` for cold |
| Modify | `internal/bridge/bridge_test.go` | Test continue field passthrough |

### Hybrid Agents
| Action | File | Responsibility |
|--------|------|----------------|
| Modify | `internal/bridge/protocol.go` | Add `Agents` field to `RequestOptions` |
| Modify | `bridge/index.ts` | Map `agents` option to SDK |
| Modify | `internal/bridge/bundle.ts` | Mirror of index.ts |
| Create | `internal/agents/sdk.go` | Convert Registry agents to SDK format |
| Create | `internal/agents/sdk_test.go` | Test SDK conversion |
| Modify | `internal/telegram/input_pipeline.go` | Populate `Agents` in request |

### Scheduler Improvements
| Action | File | Responsibility |
|--------|------|----------------|
| Modify | `go.mod` | Add `robfig/cron/v3` |
| Modify | `internal/cron/scheduler.go` | Replace `computeNextRun` with robfig, add run locking |
| Create | `internal/cron/scheduler_test.go` | Test new cron parsing + locking |
| Modify | `internal/cron/store_schema.go` | Add `cost_usd`, `tokens_used`, `session_id` to executions |
| Modify | `internal/cron/store_executions.go` | Persist new fields |
| Modify | `internal/cron/types.go` | Add fields to `CronExecution` struct |
| Modify | `internal/cron/runtime.go` | Extract cost/session from bridge result event |

---

## Task 1: Session Continuation — Protocol

**Files:**
- Modify: `internal/bridge/protocol.go`
- Modify: `bridge/index.ts`
- Modify: `internal/bridge/bundle.ts`

- [ ] **Step 1: Add Continue field to Go protocol**

In `internal/bridge/protocol.go`, add to `RequestOptions`:

```go
Continue bool `json:"continue,omitempty"`
```

- [ ] **Step 2: Map continue in bridge TypeScript**

In `bridge/index.ts`, add `continue?: boolean` to the `RequestOptions` interface.

In `buildSDKOptions`, add after the resume mapping:

```typescript
if (opts.continue) sdkOpts.continue = opts.continue;
```

- [ ] **Step 3: Sync bundle.ts**

Copy `bridge/index.ts` to `internal/bridge/bundle.ts`.

- [ ] **Step 4: Verify build**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 5: Commit**

```
feat(bridge): add continue option to request protocol
```

---

## Task 2: Session Continuation — Session Store

**Files:**
- Modify: `internal/telegram/sessions.go`

- [ ] **Step 1: Rewrite sessions.go with active tracking**

Replace the entire `sessions.go` file. The key changes: `sessions` map value becomes `*sessionEntry` (not `string`), `newSessionStore` returns the new type, `GetWithState` is added, `Clear` is updated.

```go
type sessionEntry struct {
	sessionID string
	active    bool // true if session was used since last process start
}

type sessionStore struct {
	mu       sync.RWMutex
	sessions map[int64]*sessionEntry
	cwds     map[int64]string
}

func newSessionStore() *sessionStore {
	return &sessionStore{
		sessions: make(map[int64]*sessionEntry),
		cwds:     make(map[int64]string),
	}
}

func (s *sessionStore) Get(chatID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e := s.sessions[chatID]
	if e == nil {
		return ""
	}
	return e.sessionID
}

func (s *sessionStore) GetWithState(chatID int64) (sessionID string, active bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e := s.sessions[chatID]
	if e == nil {
		return "", false
	}
	return e.sessionID, e.active
}

func (s *sessionStore) Set(chatID int64, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[chatID] = &sessionEntry{sessionID: sessionID, active: true}
}

func (s *sessionStore) Clear(chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, chatID)
	delete(s.cwds, chatID)
}

// GetCwd and SetCwd remain unchanged.
```

- [ ] **Step 2: Verify build**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 3: Commit**

```
feat(telegram): track active session state per chat
```

---

## Task 3: Session Continuation — Input Pipeline

**Files:**
- Modify: `internal/telegram/input_pipeline.go`

- [ ] **Step 1: Update buildBridgeRequest to use continue**

Find the section that sets `Resume` (around line 122-125). Replace with:

```go
sessionID, active := bc.sessions.GetWithState(chatID)
if sessionID != "" {
	if active {
		req.Options.Continue = true
	} else {
		req.Options.Resume = sessionID
	}
}
```

- [ ] **Step 2: Verify session Set is already called**

In `processBridgeEventsAsync` at line 175-177 of `input_pipeline.go`, the code already calls `bc.sessions.Set(chat.ID, ev.SessionID)` on `"system"` events. With the updated `Set` from Task 2, this now automatically marks the session as `active: true`. No code change needed here — just verify this path exists.

- [ ] **Step 3: Verify build**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 4: Manual smoke test**

Start Aurelia, send two messages to same chat. First message should use `resume` (or no session). Second should log `continue: true`. Check bridge stderr for `query start` logs.

- [ ] **Step 5: Commit**

```
feat(telegram): use continue for warm sessions, resume for cold
```

---

## Task 4: Hybrid Agents — Protocol + Bridge

**Files:**
- Modify: `internal/bridge/protocol.go`
- Modify: `bridge/index.ts`
- Modify: `internal/bridge/bundle.ts`

- [ ] **Step 1: Add Agents field to Go protocol**

In `internal/bridge/protocol.go`, add to `RequestOptions`:

```go
Agents map[string]any `json:"agents,omitempty"`
```

- [ ] **Step 2: Map agents in bridge TypeScript**

In `bridge/index.ts`, add to `RequestOptions` interface:

```typescript
agents?: Record<string, unknown>;
```

In `buildSDKOptions`, add:

```typescript
if (opts.agents && Object.keys(opts.agents).length > 0) {
  sdkOpts.agents = opts.agents;
}
```

- [ ] **Step 3: Sync bundle.ts**

Copy `bridge/index.ts` to `internal/bridge/bundle.ts`.

- [ ] **Step 4: Verify build**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 5: Commit**

```
feat(bridge): add agents option to request protocol
```

---

## Task 5: Hybrid Agents — SDK Conversion

**Files:**
- Create: `internal/agents/sdk.go`
- Create: `internal/agents/sdk_test.go`

- [ ] **Step 1: Write test for SDK conversion**

Create `internal/agents/sdk_test.go`:

```go
package agents

import "testing"

func TestBuildSDKAgents_Empty(t *testing.T) {
	r := &Registry{agents: map[string]*Agent{}}
	got := BuildSDKAgents(r)
	if len(got) != 0 {
		t.Errorf("expected empty map, got %d entries", len(got))
	}
}

func TestBuildSDKAgents_SingleAgent(t *testing.T) {
	r := &Registry{agents: map[string]*Agent{
		"coder": {
			Name:        "coder",
			Description: "writes code",
			Model:       "claude-sonnet-4-6",
			Prompt:      "You are a coder.",
			AllowedTools: []string{"Read", "Edit"},
		},
	}}
	got := BuildSDKAgents(r)
	if len(got) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(got))
	}
	a, ok := got["coder"]
	if !ok {
		t.Fatal("missing 'coder' key")
	}
	m := a.(map[string]any)
	if m["description"] != "writes code" {
		t.Errorf("description = %q", m["description"])
	}
	if m["prompt"] != "You are a coder." {
		t.Errorf("prompt = %q", m["prompt"])
	}
	if m["model"] != "claude-sonnet-4-6" {
		t.Errorf("model = %q", m["model"])
	}
	tools, ok := m["tools"]
	if !ok {
		t.Fatal("missing 'tools' key")
	}
	toolSlice := tools.([]string)
	if len(toolSlice) != 2 || toolSlice[0] != "Read" {
		t.Errorf("tools = %v", toolSlice)
	}
}

func TestBuildSDKAgents_OmitsEmptyFields(t *testing.T) {
	r := &Registry{agents: map[string]*Agent{
		"simple": {
			Name:        "simple",
			Description: "a simple agent",
			Prompt:      "Be helpful.",
		},
	}}
	got := BuildSDKAgents(r)
	m := got["simple"].(map[string]any)
	if _, ok := m["model"]; ok {
		t.Error("model should be omitted when empty")
	}
	if _, ok := m["tools"]; ok {
		t.Error("tools should be omitted when empty")
	}
}
```

- [ ] **Step 2: Run test, verify failure**

Run: `go test ./internal/agents/ -run TestBuildSDKAgents -v`
Expected: FAIL (BuildSDKAgents not defined)

- [ ] **Step 3: Implement BuildSDKAgents**

Create `internal/agents/sdk.go`:

```go
package agents

// BuildSDKAgents converts all agents in the registry to the format expected
// by the Claude Agent SDK's "agents" query option. Each agent becomes a map
// with keys: description, prompt, model (optional), tools (optional).
func BuildSDKAgents(r *Registry) map[string]any {
	all := r.Agents()
	if len(all) == 0 {
		return nil
	}
	result := make(map[string]any, len(all))
	for _, a := range all {
		def := map[string]any{
			"description": a.Description,
			"prompt":      a.Prompt,
		}
		if a.Model != "" {
			def["model"] = a.Model
		}
		if len(a.AllowedTools) > 0 {
			def["tools"] = a.AllowedTools
		}
		result[a.Name] = def
	}
	return result
}
```

- [ ] **Step 4: Run tests, verify pass**

Run: `go test ./internal/agents/ -run TestBuildSDKAgents -v`
Expected: PASS

- [ ] **Step 5: Commit**

```
feat(agents): add BuildSDKAgents to convert registry to SDK format
```

---

## Task 6: Hybrid Agents — Wire into Input Pipeline

**Files:**
- Modify: `internal/telegram/input_pipeline.go`

- [ ] **Step 1: Populate Agents in buildBridgeRequest**

In `buildBridgeRequest`, after setting agent-specific overrides, add:

```go
// Pass all agents to SDK for native delegation
if sdkAgents := agents.BuildSDKAgents(bc.agents); sdkAgents != nil {
	req.Options.Agents = sdkAgents
}
```

This must be imported: the `agents` package is already imported in input_pipeline.go.

- [ ] **Step 2: Verify build**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 3: Commit**

```
feat(telegram): pass agent definitions to SDK for native delegation
```

---

## Task 7: Scheduler — Add robfig/cron dependency

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Add dependency**

Run: `go get github.com/robfig/cron/v3`

- [ ] **Step 2: Tidy**

Run: `go mod tidy`

- [ ] **Step 3: Commit**

```
chore: add robfig/cron/v3 for full cron expression parsing
```

---

## Task 8: Scheduler — Replace computeNextRun

**Files:**
- Modify: `internal/cron/scheduler.go`
- Create: `internal/cron/scheduler_test.go`

- [ ] **Step 1: Write tests for new cron parsing**

Create `internal/cron/scheduler_test.go`:

```go
package cron

import (
	"testing"
	"time"
)

func TestComputeNextRun_EveryFiveMinutes(t *testing.T) {
	base := time.Date(2026, 3, 22, 10, 0, 0, 0, time.UTC)
	next, err := computeNextRun("*/5 * * * *", base)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 3, 22, 10, 5, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("got %v, want %v", next, want)
	}
}

func TestComputeNextRun_DailyAt9AM(t *testing.T) {
	base := time.Date(2026, 3, 22, 10, 0, 0, 0, time.UTC)
	next, err := computeNextRun("0 9 * * *", base)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 3, 23, 9, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("got %v, want %v", next, want)
	}
}

func TestComputeNextRun_WeekdaysOnly(t *testing.T) {
	// Friday 2026-03-27
	base := time.Date(2026, 3, 27, 18, 0, 0, 0, time.UTC)
	next, err := computeNextRun("0 9 * * 1-5", base)
	if err != nil {
		t.Fatal(err)
	}
	// Should skip Sat+Sun → Monday 2026-03-30
	if next.Weekday() != time.Monday {
		t.Errorf("expected Monday, got %v (%v)", next.Weekday(), next)
	}
}

func TestComputeNextRun_InvalidExpr(t *testing.T) {
	_, err := computeNextRun("not a cron", time.Now())
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}
```

- [ ] **Step 2: Run tests, verify failure**

Run: `go test ./internal/cron/ -run TestComputeNextRun -v`
Expected: FAIL (weekdays test fails with old parser)

- [ ] **Step 3: Replace computeNextRun implementation**

In `internal/cron/scheduler.go`, replace the `computeNextRun` function and remove `isExactNumber` and `parseIntCron` helpers:

```go
import (
	"fmt"
	"time"

	robfigcron "github.com/robfig/cron/v3"
)

// computeNextRun calculates the next run time for a standard cron expression.
func computeNextRun(expr string, after time.Time) (time.Time, error) {
	parser := robfigcron.NewParser(robfigcron.Minute | robfigcron.Hour | robfigcron.Dom | robfigcron.Month | robfigcron.Dow)
	sched, err := parser.Parse(expr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron expression %q: %w", expr, err)
	}
	return sched.Next(after), nil
}
```

- [ ] **Step 4: Run tests, verify pass**

Run: `go test ./internal/cron/ -run TestComputeNextRun -v`
Expected: PASS

- [ ] **Step 5: Commit**

```
feat(cron): replace manual parser with robfig/cron for full expression support
```

---

## Task 9: Scheduler — Run Locking

**Files:**
- Modify: `internal/cron/scheduler.go`

- [ ] **Step 1: Add running set to Scheduler**

Add field to `Scheduler` struct:

```go
running sync.Map // jobID → struct{} to prevent concurrent execution
```

- [ ] **Step 2: Guard runSingleJob with lock**

At the start of `runSingleJob`:

```go
if _, loaded := s.running.LoadOrStore(job.ID, struct{}{}); loaded {
	return // already running
}
defer s.running.Delete(job.ID)
```

- [ ] **Step 3: Verify build**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 4: Commit**

```
fix(cron): prevent concurrent execution of same job
```

---

## Task 10: Scheduler — Cost Tracking + Session ID

**Files:**
- Modify: `internal/cron/types.go`
- Modify: `internal/cron/store_schema.go`
- Modify: `internal/cron/store_executions.go`
- Modify: `internal/cron/runtime.go`

- [ ] **Step 1: Add fields to CronExecution**

In `internal/cron/types.go`, add to `CronExecution`:

```go
SessionID  string  // bridge session ID for audit trail
CostUSD    float64 // API cost for this execution
TokensUsed int     // total tokens consumed
```

- [ ] **Step 2: Update schema**

In `internal/cron/store_schema.go`, add columns to the `cron_executions` CREATE TABLE (use `IF NOT EXISTS` for the table — existing columns won't break). Also add an `ALTER TABLE` migration block after table creation:

```go
// Migration: add columns if missing (safe to re-run)
for _, col := range []string{
	"ALTER TABLE cron_executions ADD COLUMN session_id TEXT DEFAULT ''",
	"ALTER TABLE cron_executions ADD COLUMN cost_usd REAL DEFAULT 0",
	"ALTER TABLE cron_executions ADD COLUMN tokens_used INTEGER DEFAULT 0",
} {
	_, err := db.Exec(col)
	if err != nil && !strings.Contains(err.Error(), "duplicate column") {
		return fmt.Errorf("migration: %w", err)
	}
}
```

- [ ] **Step 3: Update RecordExecution**

In `internal/cron/store_executions.go`, add the new columns to the INSERT statement:

```sql
INSERT INTO cron_executions (id, job_id, started_at, finished_at, status, output_summary, error_message, session_id, cost_usd, tokens_used)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
```

And pass `exec.SessionID`, `exec.CostUSD`, `exec.TokensUsed`.

- [ ] **Step 4: Update runtime to capture cost**

In `internal/cron/runtime.go`, the `ExecuteJob` method calls `b.bridge.Execute()` which returns events. The `result` event has `CostUSD`, `SessionID`, and `NumTurns`. Capture these and return them.

The current `ExecuteJob` returns `(string, error)`. It uses `ExecuteSync` which returns an `*Event`. Extract cost data from the result event and populate the `CronExecution` in `scheduler.go`'s `runSingleJob`.

Change `Runtime` interface to return richer data. In `types.go`:

```go
type ExecutionResult struct {
	Output    string
	SessionID string
	CostUSD   float64
	NumTurns  int
}

type Runtime interface {
	ExecuteJob(ctx context.Context, job CronJob) (*ExecutionResult, error)
}
```

In `runtime.go`, update `BridgeCronRuntime.ExecuteJob`:

```go
func (r *BridgeCronRuntime) ExecuteJob(ctx context.Context, job CronJob) (*ExecutionResult, error) {
	// ... existing system prompt building code stays the same ...

	ev, err := r.bridge.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("bridge execute: %w", err)
	}
	if ev.Type == "error" {
		return nil, fmt.Errorf("bridge error: %s", ev.Message)
	}
	// ... existing memory save stays the same ...

	return &ExecutionResult{
		Output:    ev.Content,
		SessionID: ev.SessionID,
		CostUSD:   ev.CostUSD,
		NumTurns:  ev.NumTurns,
	}, nil
}
```

Update `DeliveryFunc` signature:

```go
type DeliveryFunc func(ctx context.Context, job CronJob, result *ExecutionResult, execErr error) error
```

Update `NotifyingRuntime.ExecuteJob`:

```go
func (r *NotifyingRuntime) ExecuteJob(ctx context.Context, job CronJob) (*ExecutionResult, error) {
	if r.inner == nil {
		return nil, fmt.Errorf("inner runtime is required")
	}
	result, err := r.inner.ExecuteJob(ctx, job)
	if r.deliver != nil {
		if deliverErr := r.deliver(ctx, job, result, err); deliverErr != nil {
			return result, deliverErr
		}
	}
	return result, err
}
```

Update `scheduler.go`'s `runSingleJob` to populate CronExecution from result:

```go
result, execErr := s.runtime.ExecuteJob(ctx, job)

exec := CronExecution{
	ID:        uuid.NewString(),
	JobID:     job.ID,
	StartedAt: startedAt,
	FinishedAt: s.clock.Now(),
	Status:    "completed",
}
if result != nil {
	exec.OutputSummary = result.Output
	exec.SessionID = result.SessionID
	exec.CostUSD = result.CostUSD
	exec.TokensUsed = result.NumTurns // best available proxy
}
if execErr != nil {
	exec.Status = "failed"
	exec.ErrorMessage = execErr.Error()
}
```

Update caller in `cmd/aurelia/app.go` that creates the `DeliveryFunc` — change signature to accept `*cron.ExecutionResult` instead of `string`:

```go
deliver := func(ctx context.Context, job cron.CronJob, result *cron.ExecutionResult, execErr error) error {
	output := ""
	if result != nil {
		output = result.Output
	}
	// ... existing delivery logic using output ...
}
```

Update `ListExecutionsByJob` in `store_executions.go` to scan new columns:

```go
func (s *SQLiteCronStore) ListExecutionsByJob(ctx context.Context, jobID string) ([]CronExecution, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, job_id, started_at, finished_at, status, output_summary, error_message,
		        COALESCE(session_id, ''), COALESCE(cost_usd, 0), COALESCE(tokens_used, 0)
		 FROM cron_executions WHERE job_id = ? ORDER BY started_at DESC`, jobID)
	// ... scan into all fields including SessionID, CostUSD, TokensUsed ...
}
```

- [ ] **Step 5: Verify build**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 6: Run existing tests**

Run: `go test ./internal/cron/ -v`
Expected: PASS (or no tests yet — that's ok)

- [ ] **Step 7: Commit**

```
feat(cron): track cost, tokens, and session ID per execution
```

---

## Task 11: Final Sync + Integration Test

- [ ] **Step 1: Copy bridge to bundle**

```bash
cp bridge/index.ts internal/bridge/bundle.ts
```

- [ ] **Step 2: Full build**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 3: Run all tests**

Run: `go test ./... -count=1`
Expected: all PASS

- [ ] **Step 4: Manual smoke test**

1. Start Aurelia
2. Send message via Telegram → verify response (session continuation)
3. Send second message → check bridge logs for `continue` instead of `resume`
4. Create agent in `~/.aurelia/agents/coder.md` with basic frontmatter
5. Send message without `@coder` → verify SDK may use the agent
6. Create cron job with complex expression (e.g., `0 9 * * 1-5`) → verify next_run is correct

- [ ] **Step 5: Final commit**

```
chore: sync bundle.ts and verify integration
```
