# Code Health Fixes — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix all important issues found in the code quality audit — bugs, dead code, missing timeouts, god files, DRY violations.

**Architecture:** Surgical fixes across existing packages. No new packages. No new abstractions unless strictly needed.

**Tech Stack:** Go, TypeScript

---

## File Structure

### Files to modify

```
internal/telegram/bot_middleware.go     — fix nil sender panic
internal/telegram/input_pipeline.go    — add context timeout for bridge/memory
internal/telegram/output.go            — remove SendAudio dead code
internal/telegram/bot.go               — extract interfaces for deps
internal/config/config.go              — split into 3 files, remove dead code
internal/config/mcp.go                 — delete (dead code)
internal/memory/embeddings.go          — remove unused NoopEmbedder
cmd/aurelia/app.go                     — fix duplicate cron jobs, remove dead app_test.go
cmd/aurelia/onboard_providers.go       — remove duplicated normalizeProvider, import from config
bridge/index.ts                        — add timeout to query stream
```

### Files to create

```
internal/config/config_editable.go     — extracted EditableConfig + conversions
internal/config/config_migration.go    — extracted legacy migration
```

### Files to delete

```
internal/config/mcp.go                 — entirely unused
cmd/aurelia/app_test.go                — empty file, no tests
```

---

## Task 1: Fix nil sender panic in middleware

**Files:**
- Modify: `internal/telegram/bot_middleware.go:11-12`

- [ ] **Step 1: Fix nil check**

```go
func (bc *BotController) whitelistMiddleware() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			sender := c.Sender()
			if sender == nil {
				return nil
			}
			if !bc.isAllowedUser(sender.ID) {
				log.Printf("blocked unauthorized user: %d\n", sender.ID)
				return nil
			}
			return next(c)
		}
	}
}
```

- [ ] **Step 2: Also fix processInput which calls c.Sender().ID**

In `input_pipeline.go:19`, add nil check:

```go
if state, ok := bc.popPendingBootstrap(c.Sender().ID); ok {
```

This is safe because the middleware already filters nil senders before reaching processInput. No change needed here.

- [ ] **Step 3: Run tests**

```bash
go test ./internal/telegram/ -v
```

- [ ] **Step 4: Commit**

```bash
git commit -m "fix(telegram): guard against nil sender in middleware"
```

---

## Task 2: Add context timeout to bridge and memory calls

**Files:**
- Modify: `internal/telegram/input_pipeline.go:78,110,194`

- [ ] **Step 1: Replace context.Background() with timeout in processInput**

```go
// Line 78: bridge execution
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
defer cancel()
ch, err := bc.bridge.Execute(ctx, req)
```

- [ ] **Step 2: Replace context.Background() in buildSystemPrompt**

```go
// Line 110: memory injection
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
memoryBlock, err := bc.memory.Inject(ctx, userText, bc.config.MemoryWindowSize)
```

- [ ] **Step 3: Replace context.Background() in saveToMemory**

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/telegram/ -v
```

- [ ] **Step 5: Commit**

```bash
git commit -m "fix(telegram): add timeouts to bridge and memory operations"
```

---

## Task 3: Fix duplicate cron jobs on restart

**Files:**
- Modify: `cmd/aurelia/app.go:262-279`

- [ ] **Step 1: Make registerScheduledAgents idempotent**

Before creating a job, check if one already exists for this agent name:

```go
func registerScheduledAgents(store *cron.SQLiteCronStore, reg *agents.Registry) {
	if reg == nil {
		return
	}
	svc := cron.NewService(store, nil)
	existing, _ := svc.ListJobs(context.Background())

	// Build set of agents that already have jobs
	hasJob := make(map[string]bool)
	for _, job := range existing {
		if job.AgentName != "" {
			hasJob[job.AgentName] = true
		}
	}

	for _, a := range reg.Scheduled() {
		if hasJob[a.Name] {
			continue
		}
		_, err := svc.CreateJob(context.Background(), cron.CronJob{
			AgentName:    a.Name,
			ScheduleType: "cron",
			CronExpr:     a.Schedule,
			Prompt:       a.Prompt,
		})
		if err != nil {
			log.Printf("Warning: failed to register scheduled agent %q: %v", a.Name, err)
		}
	}
}
```

Check if `svc.ListJobs` exists. If not, query the store directly. Read the cron service/store to find the right method.

- [ ] **Step 2: Run tests**

```bash
go test ./cmd/aurelia/ -v
go test ./internal/cron/ -v
```

- [ ] **Step 3: Commit**

```bash
git commit -m "fix(cron): prevent duplicate jobs on restart"
```

---

## Task 4: Remove dead code

**Files:**
- Delete: `internal/config/mcp.go` (entire file, unused)
- Delete: `cmd/aurelia/app_test.go` (empty, just a comment)
- Modify: `internal/telegram/output.go` — remove `SendAudio` function
- Modify: `internal/memory/embeddings.go` — remove `NoopEmbedder`
- Modify: `internal/config/config.go` — remove `defaultLLMModelForProvider` and `DefaultFileConfig`
- Modify: `internal/config/config_test.go` — replace `defaultLLMModelForProvider` calls with `defaultModelForProvider`

- [ ] **Step 1: Delete mcp.go**

```bash
rm internal/config/mcp.go
```

- [ ] **Step 2: Delete empty app_test.go**

```bash
rm cmd/aurelia/app_test.go
```

- [ ] **Step 3: Remove SendAudio from output.go**

Delete lines 113-117.

- [ ] **Step 4: Remove NoopEmbedder from embeddings.go**

Delete the NoopEmbedder struct and its methods (lines 91-107 approximately).

- [ ] **Step 5: Remove defaultLLMModelForProvider from config.go**

Delete the function. Update config_test.go to call `defaultModelForProvider` directly.

- [ ] **Step 6: Remove DefaultFileConfig from config.go**

Delete the exported function (line 333-349 approximately). It's only used to create the unexported `defaultFileConfig` which is the one that matters.

Actually — check if DefaultFileConfig is called anywhere. If it's only used internally, unexport it. If it's not used at all, delete it.

- [ ] **Step 7: Run full test suite**

```bash
go build ./...
go test ./... -short
```

- [ ] **Step 8: Commit**

```bash
git commit -m "chore: remove dead code (mcp.go, SendAudio, NoopEmbedder, unused exports)"
```

---

## Task 5: Fix DRY — deduplicate normalizeProvider

**Files:**
- Modify: `internal/config/config.go` — export `NormalizeProvider`
- Modify: `cmd/aurelia/onboard_providers.go` — delete local `normalizeProvider`, import from config
- Modify: `cmd/aurelia/onboard_helpers.go` — update calls to use `config.NormalizeProvider`

- [ ] **Step 1: Export normalizeProvider in config.go**

Rename `normalizeProvider` to `NormalizeProvider` in `internal/config/config.go`. Update all internal calls.

- [ ] **Step 2: Delete normalizeProvider from onboard_providers.go**

Remove the duplicate function (lines 127-141).

- [ ] **Step 3: Update all callers in cmd/aurelia/**

Replace `normalizeProvider(...)` with `config.NormalizeProvider(...)` in:
- `onboard_helpers.go`
- `onboard_providers.go`

- [ ] **Step 4: Run tests**

```bash
go build ./...
go test ./... -short
```

- [ ] **Step 5: Commit**

```bash
git commit -m "refactor: deduplicate normalizeProvider into config package"
```

---

## Task 6: Split config.go god file

**Files:**
- Modify: `internal/config/config.go` (keep core AppConfig, Load, Save)
- Create: `internal/config/config_editable.go` (EditableConfig, conversions, comparison)
- Create: `internal/config/config_migration.go` (legacyFileConfig, migrateLegacy)

- [ ] **Step 1: Read config.go fully to identify extraction boundaries**

Identify:
- `EditableConfig` struct + `LoadEditable` + `SaveEditable` + `appConfigToEditable` + `editableToFileConfig` + `sameFileConfig` → `config_editable.go`
- `legacyFileConfig` struct + `migrateLegacy` → `config_migration.go`

- [ ] **Step 2: Create config_editable.go**

Move the `EditableConfig` struct, its conversion functions, and the `sameFileConfig` comparison into this file. Keep the same package.

- [ ] **Step 3: Create config_migration.go**

Move `legacyFileConfig` struct and `migrateLegacy()` function into this file.

- [ ] **Step 4: Verify config.go is now focused**

After extraction, `config.go` should contain:
- `AppConfig` struct + methods
- `ProviderConfig` struct
- `fileConfig` struct
- `Load()`, `Save()`
- `NormalizeProvider()`, `defaultModelForProvider()`

Target: under 300 lines.

- [ ] **Step 5: Run tests**

```bash
go test ./internal/config/ -v
go build ./...
```

- [ ] **Step 6: Commit**

```bash
git commit -m "refactor(config): split god file into editable and migration modules"
```

---

## Task 7: Add timeout to Bridge TypeScript query stream

**Files:**
- Modify: `bridge/index.ts`

- [ ] **Step 1: Add timeout to handleQuery**

Wrap the query stream with a timeout. If no events arrive for 10 minutes, emit an error and break:

```typescript
async function handleQuery(req: Request): Promise<void> {
  const sdkOptions = buildSDKOptions(req.options);
  const timeoutMs = 10 * 60 * 1000; // 10 minutes

  log(`query start — model=${sdkOptions.model ?? "default"}`);

  try {
    const stream = query({
      prompt: req.prompt,
      options: sdkOptions as Parameters<typeof query>[0]["options"],
    });

    const timeout = setTimeout(() => {
      log("query timeout — no result after 10 minutes");
      emit({ event: "error", message: "query timeout: no result after 10 minutes" });
      // Force exit to clean up — Go will handle the process termination
      process.exit(1);
    }, timeoutMs);

    try {
      for await (const message of stream) {
        // ... existing event handling
      }
    } finally {
      clearTimeout(timeout);
    }
  } catch (err: unknown) {
    // ... existing error handling
  }
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd bridge && npx tsc --noEmit
```

- [ ] **Step 3: Commit**

```bash
git commit -m "fix(bridge): add 10-minute timeout to query stream"
```

---

## Task 8: Final verification

- [ ] **Step 1: Full build**

```bash
go build ./...
```

- [ ] **Step 2: Full vet**

```bash
go vet ./...
```

- [ ] **Step 3: Full test suite**

```bash
go test ./... -v -short
```

- [ ] **Step 4: Bridge compilation**

```bash
cd bridge && npx tsc --noEmit
```

- [ ] **Step 5: Verify no remaining dead code**

```bash
grep -rn "SendAudio\|LoadMCPConfig\|DefaultFileConfig\|defaultLLMModelForProvider\|NoopEmbedder" --include="*.go" .
```

Expected: no matches (or only in test assertions if applicable).
