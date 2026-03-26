# Codebase Cleanup — Aurelia OS

**Date:** 2026-03-22
**Status:** Approved
**Approach:** Incremental phases (B) — 4 phases, same branch, each compiles and tests before advancing

---

## Context

After a full architectural audit of the Aurelia OS codebase (12,258 LOC → 11,252 LOC post memory removal), the following categories of issues remain:

- **Dead code** still consuming config, onboard, and persona surface area
- **Bugs** — race condition in bridge, case sensitivity in agent routing
- **Boundary violations** — telegram package doing business logic (session management, token accounting, agent classification)
- **God functions** — `bootstrapApp()` with 13 responsibilities (~162 lines), `processBridgeEventsAsync()` with 83 lines and 5 concerns
- **Code smells** — magic numbers, swallowed errors, missing SQL index, non-atomic DB operations

Memory system (`internal/memory/`) and embedding config were already removed in prior session. Claude-mem MCP is the sole memory provider.

---

## Phase 1 — Dead Code + Real Bugs

### 1.1 Remove `MemoryWindowSize`

**Scope:** 41 occurrences across 10+ Go files. Documentation-only files (AGENTS.md, README.md) are out of scope — they'll be updated separately if needed.

**Files to edit:**

| File | Changes |
|------|---------|
| `internal/config/config.go` | Remove field from `AppConfig`, `fileConfig`, constant `defaultMemoryWindowSize`, references in `normalizeFileConfig`, `defaultFileConfig`, `toAppConfig` |
| `internal/config/config_editable.go` | Remove field from `EditableConfig`, `DefaultEditableConfig`, `appConfigToEditable`, `editableToFileConfig`, `sameFileConfig` |
| `internal/config/config_migration.go` | Remove field from `legacyFileConfig`, `migrateLegacy` |
| `internal/config/config_test.go` | Remove all `MemoryWindowSize` assertions and fixture values |
| `internal/persona/loader.go` | Remove `MemoryWindowSize` field from `Config` struct |
| `internal/persona/loader_test.go` | Remove `memory_window_size` from YAML fixtures |
| `internal/persona/canonical_service_test.go` | Remove `memory_window_size` from YAML fixtures |
| `internal/telegram/bootstrap_config.go` | Remove `memory_window_size` from config templates |
| `internal/telegram/input_pipeline_test.go` | Remove `MemoryWindowSize` from test BotController setup |
| `cmd/aurelia/onboard.go` | Remove memory window size prompt |
| `cmd/aurelia/onboard_ui.go` | Remove memory window size display and input handling |
| `cmd/aurelia/onboard_helpers.go` | Remove memory window size from config display |
| `cmd/aurelia/onboard_test.go` | Remove `MemoryWindowSize` assertions and fixture values |

**Verification:** `go build ./...` + `go test ./...` pass. Zero references to `MemoryWindowSize` or `memory_window_size` in codebase.

### 1.2 Fix Race Condition — `bridge.go`

**Problem:** Channel `done` is created inside `startLocked()` (line 109) but if `Stop()` is called before `Start()`, `done` is nil and `<-done` at line 193 blocks forever. Additionally, concurrent lifecycle calls (Start/Stop) can race on the channel field.

**Fix:** Make `done` a field on the Bridge struct, initialized in `New()`. Reset under lock in `startLocked()`. Guard nil-channel reads in `Stop()` with an early return if process was never started.

**Verification:** `go test -race ./internal/bridge/...` passes.

### 1.3 Fix Case Sensitivity — `agents/registry.go:98`

**Problem:** `Route()` lowercases the extracted `@name` from user input before lookup, but agent names are stored as-is from YAML (could be mixed case). If YAML defines `name: MyAgent`, looking up `"myagent"` fails.

**Fix:** Normalize agent names to lowercase in `Load()` when storing into the registry map. This ensures all lookups are case-insensitive consistently.

**Verification:** Add test case with mixed-case agent name. `go test ./internal/agents/...` passes.

### 1.4 Remove Remaining Dead Code

| Item | File | Action |
|------|------|--------|
| `seedBootstrapIdentity()` | `telegram/bootstrap_profile.go` | Delete empty stub function + remove its call site in `bootstrap.go` |
| `requiresAudio` param | `telegram/input_pipeline.go:16` | Remove from signature and all call sites |
| `parts [][]byte` param | `telegram/input_pipeline.go:16` | Remove from signature and all call sites |
| Unused `_ = state` | `telegram/bootstrap_profile.go:59` | Remove unused parameter or use it |

**Verification:** `go build ./...` + `go test ./...` pass.

---

## Phase 2 — Extract Responsibilities from Telegram

### 2.1 Create `internal/session/` Package

**Purpose:** Session management is domain logic, not Telegram I/O. Extracting it enables reuse across future channels (Slack, Discord, CLI).

**New files:**

```
internal/session/
  store.go       — session store (moved from telegram/sessions.go)
  tracker.go     — token/cost tracking (moved from telegram/session_tracker.go)
  store_test.go  — tests for session store
  tracker_test.go — tests for tracker
```

**Interfaces:**

```go
// Store manages session IDs per chat.
type Store interface {
    Get(chatID int64) string
    GetWithState(chatID int64) (sessionID string, active bool)
    Set(chatID int64, sessionID string)
    Clear(chatID int64)
    GetCwd(chatID int64) string
    SetCwd(chatID int64, cwd string)
}

// Tracker tracks token usage and cost per session.
type Tracker interface {
    Add(chatID int64, inputTokens, outputTokens, numTurns int, costUSD float64) int
    RecordUsage(chatID int64, numTurns int, costUSD float64, maxTokens int) bool
    Get(chatID int64) Usage
    Clear(chatID int64)
}
```

**BotController changes:** Replace concrete `*sessionStore` and `*sessionTracker` fields with `session.Store` and `session.Tracker` interfaces. Constructor accepts interfaces.

**Verification:** All existing session/tracker tests pass in new location. Telegram tests still pass with injected interfaces.

### 2.2 Extract Agent Classification

**Problem:** `routeAgent()` in `input_pipeline.go:62-90` does LLM classification via bridge — this is business logic, not Telegram I/O.

**Fix:** Add `Classify` method to `agents.Registry` using a callback function type to avoid importing `bridge` (keeping `agents` as a leaf package with zero internal dependencies):

```go
// ClassifyFunc sends a text prompt and returns the LLM response (agent name or "none").
// Defined in agents package — no bridge import needed.
type ClassifyFunc func(ctx context.Context, systemPrompt, userPrompt string) (string, error)

// Classify uses LLM classification to route a message to the best agent.
func (r *Registry) Classify(text string, classify ClassifyFunc) *Agent
```

The `ClassifyFunc` adapter is wired in telegram, where bridge is already imported:

```go
// In telegram package — bridges the gap between agents.ClassifyFunc and bridge.ExecuteSync
classifyFn := func(ctx context.Context, system, prompt string) (string, error) {
    result, err := bc.bridge.ExecuteSync(ctx, bridge.Request{...})
    if err != nil { return "", err }
    return result.Content, nil
}
```

**Telegram changes:** `routeAgent()` becomes a thin delegation:

```go
func (bc *BotController) routeAgent(text string) *agents.Agent {
    agent := bc.agents.Route(text)
    if agent == nil {
        agent = bc.agents.Classify(text, bc.classifyFunc())
    }
    return agent
}
```

**Key:** `agents` package stays a leaf with zero internal imports. The bridge dependency is injected via function type.

**Verification:** Add test for `Classify` in agents package using mock `ClassifyFunc`. Telegram tests still pass.

### 2.3 Extract Token Accounting to Session

**Problem:** Lines 216-230 in `processBridgeEventsAsync` mix token tracking with event processing.

**Fix:** Add method to `session.Tracker`:

```go
// RecordUsage tracks usage and returns true if session should be reset.
func (t *tracker) RecordUsage(chatID int64, turns int, cost float64, maxTokens int) bool
```

**Telegram changes:** Replace 15-line block with:

```go
if bc.tracker.RecordUsage(chat.ID, ev.NumTurns, ev.CostUSD, bc.config.MaxSessionTokens) {
    bc.sessions.Clear(chat.ID)
}
```

**Verification:** Unit test for `RecordUsage` threshold logic. Telegram tests pass.

### 2.4 Extract Cron Delivery

**Problem:** `app.go:165-184` has inline delivery callback mixing business logic with bootstrap.

**Fix:** Create `internal/cron/delivery.go`:

```go
// TelegramDelivery sends cron execution results to Telegram chats.
type TelegramDelivery struct {
    sender MessageSender
}

// MessageSender is the subset of telebot.Bot needed for delivery.
type MessageSender interface {
    Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
}

func (d *TelegramDelivery) Deliver(ctx context.Context, job CronJob, result *ExecutionResult, execErr error) error
```

**app.go changes:** Replace 20-line inline closure with `cron.NewTelegramDelivery(bot.GetBot())`.

**Verification:** Unit test for delivery formatting. `go test ./internal/cron/...` passes.

---

## Phase 3 — Break God Functions

### 3.1 Break `bootstrapApp()`

**Current:** 1 function, ~162 lines, 13 responsibilities (numbered in code comments).

**Target:** Orchestrator of ~40 lines calling focused sub-functions.

**New functions in `app.go`:**

```go
func setupBridge(cfg *config.AppConfig) (*bridge.Bridge, error)
// Handles: bridge dir discovery, auto-setup, bundle path resolution

func setupPersona(resolver *runtime.PathResolver) *persona.CanonicalIdentityService
// Handles: persona dir resolution, path assembly, service creation

func setupCronScheduler(store *cron.SQLiteCronStore, br *bridge.Bridge, agentReg *agents.Registry, personaSvc *persona.CanonicalIdentityService, bot *telegram.BotController) (*cron.Scheduler, error)
// Handles: runtime creation, delivery wiring, scheduler config, agent registration

func setupTelegramBot(cfg *config.AppConfig, br *bridge.Bridge, agentReg *agents.Registry, personaSvc *persona.CanonicalIdentityService, transcriber stt.Transcriber, cronHandler *telegram.CronCommandHandler, personasDir, exePath string, sessions session.Store, tracker session.Tracker) (*telegram.BotController, error)
// Handles: bot creation with all dependencies
```

**Cleanup pattern:** Replace scattered `_ = store.Close()` with deferred cleanup stack:

```go
var cleanups []func()
defer func() {
    if err != nil {
        for i := len(cleanups) - 1; i >= 0; i-- {
            cleanups[i]()
        }
    }
}()
```

**Verification:** `bootstrapApp()` under 50 lines. All tests pass.

### 3.2 Break `processBridgeEventsAsync()`

**Current:** 83 lines doing event dispatch + token accounting + text assembly + error handling + progress reporting.

**Target:** ~25-line loop delegating to focused handlers.

**Approach:** Token accounting already extracted in Phase 2. Remaining:
- Event type handling stays inline (switch is readable at this scale)
- Text assembly stays (builder pattern is simple)
- Net result: function drops to ~35-40 lines naturally after Phase 2 extractions

**No forced extraction** — if the function is under 50 lines after Phase 2, it's acceptable. Don't over-engineer.

**Verification:** `go test ./internal/telegram/...` passes.

### 3.3 Reduce BotController

**Current fields (13+):**
- Telegram I/O: `bot`, `exePath`
- Session: `sessions`, `tracker` (already interfaces after Phase 2)
- Albums: `pendingAlbums`, `albumMu`
- Bootstrap: `pendingBootstrap`, `bootstrapMu`
- Services: `bridge`, `agents`, `persona`, `stt`, `cronHandler`
- Config: `config`, `personasDir`

**Fix:** Extract album buffering to internal `albumBuffer` struct:

```go
type albumBuffer struct {
    mu      sync.Mutex
    pending map[string]*pendingAlbum
}
```

BotController replaces `pendingAlbums` + `albumMu` with `*albumBuffer`. Net reduction: struct is more readable, album concerns are encapsulated.

**Verification:** Album-related tests still pass.

---

## Phase 4 — Code Smells

### 4.1 Magic Numbers → Named Constants

| Current | Constant Name | File |
|---------|--------------|------|
| `3000` | `estimatedTokensPerTurn` | `session/tracker.go` (after Phase 2) |
| `4 * time.Second` (6x) | `typingIndicatorInterval` | `telegram/activity.go`, `input.go`, `input_pipeline.go` |
| `15 * time.Second` | `classifyTimeout` | `agents/registry.go` (after Phase 2) |
| `10 * time.Minute` | `bridgeExecutionTimeout` | `telegram/input_pipeline.go` |
| `16` | `eventChannelBuffer` | `bridge/bridge.go` |

**Verification:** No raw numeric literals in timeout/buffer contexts.

### 4.2 Swallowed Errors → Log

Replace all `_ = Send*()` patterns with logged errors:

```go
// Before
_ = SendError(bc.bot, chat, "message")

// After
if err := SendError(bc.bot, chat, "message"); err != nil {
    log.Printf("Failed to send error to chat %d: %v", chat.ID, err)
}
```

**Files:** `input_pipeline.go` (5 occurrences: lines 169, 235, 247, 258, 260), `app.go` cleanup paths (`_ = cronStore.Close()`).

**Verification:** Zero `_ = Send*()` and `_ = *.Close()` patterns in codebase. `go vet ./...` clean.

### 4.3 SQL Index for `ListDueJobs`

**Add to `store_schema.go`:**

```sql
CREATE INDEX IF NOT EXISTS idx_cron_jobs_due ON cron_jobs(active, next_run_at);
```

**Verification:** Index exists after store initialization. `go test ./internal/cron/...` passes.

### 4.4 Atomic Transaction in Scheduler

**Problem:** `scheduler.go:82-139` calls `RecordExecution` + `UpdateJob` without transaction. If one fails, state is inconsistent.

**Fix:** Add `WithTx` method to store interface:

```go
func (s *SQLiteCronStore) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
```

Update `runSingleJob` to wrap both operations in a transaction.

**Verification:** Test that simulates failure after RecordExecution confirms UpdateJob also rolls back. `go test ./internal/cron/...` passes.

### 4.5 Fix `shutdown(ctx)` Unused Parameter

**Problem:** `app.go` — `_ = ctx` explicitly ignoring context.

**Fix:** Use context for graceful shutdown with timeout:

```go
func (a *app) shutdown(ctx context.Context) {
    if a.cronCancel != nil {
        a.cronCancel()
    }
    if a.bot != nil {
        done := make(chan struct{})
        go func() {
            a.bot.Stop()
            close(done)
        }()
        select {
        case <-done:
        case <-ctx.Done():
            log.Println("Warning: bot shutdown timed out")
        }
    }
}
```

**Verification:** Shutdown completes within context deadline.

---

## Execution Order

```
Phase 1 (dead code + bugs)
  └→ Phase 2 (extract responsibilities)
       └→ Phase 3 (break god functions — depends on Phase 2 extractions)
            └→ Phase 4 (code smells — independent but last to minimize conflicts)
```

Each phase ends with: `go build ./...` + `go test ./...` + `go vet ./...` + commit.

## Success Criteria

- Zero dead code references (`MemoryWindowSize`, `memory_window_size`, embedding fields) — verified via grep
- `go test -race ./...` passes (bridge race condition fixed)
- Agent routing works case-insensitively — verified via test with mixed-case agent name
- `telegram/` package contains no LLM classification calls, no token arithmetic, no session threshold logic — only I/O and delegation to injected services
- `bootstrapApp()` under 50 lines
- `processBridgeEventsAsync()` under 72 lines (83 → ~71 after Phase 2 token extraction net of replacement code; forced extraction beyond that would be over-engineering)
- `BotController` fields are all either I/O or injected interfaces
- Zero `_ = Send*()` and `_ = *.Close()` patterns — verified via grep
- Zero unnamed magic numbers in timeout/buffer contexts
- SQL index covers `ListDueJobs` query
- Cron execution + state update are atomic (wrapped in transaction)
- `agents` package remains a leaf with zero internal imports
- All tests green, all packages compile: `go build ./...` + `go test ./...` + `go vet ./...`
