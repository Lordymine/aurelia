# Testing Infrastructure

## Test Frameworks

**Unit/Integration:** Standard Go `testing` package
**E2E:** Standard `testing` with bridge helpers
**Assertions:** Direct `if` comparisons + `t.Fatalf()` — no testify
**Coverage:** Not configured (no explicit targets or enforcement)

## Test Organization

**Location:** Co-located with source files (`*_test.go` in same package)
**Naming:** `{source_file}_test.go` (e.g., `scheduler.go` → `scheduler_test.go`)
**E2E:** Separate `e2e/` directory at project root

**Package distribution:**
- `cmd/aurelia/` — 1 test file (onboarding)
- `internal/agents/` — 2 test files (registry, SDK)
- `internal/bridge/` — 1 test file
- `internal/config/` — 1 test file
- `internal/cron/` — 5 test files (scheduler, service, store, delivery, runtime)
- `internal/persona/` — 3 test files (canonical, loader, optional file)
- `internal/runtime/` — 2 test files (bootstrap, resolver)
- `internal/session/` — 2 test files (store, tracker)
- `internal/telegram/` — 8 test files (activity, bootstrap, cron handlers, input pipeline, input, markdown, output, send)
- `pkg/stt/` — 1 test file
- `e2e/` — 2 test files

## Testing Patterns

### Unit Tests

**Approach:** Test public behavior via direct function calls. Use hand-written fakes for dependencies.

**Fake pattern example:**
```go
type fakeCronStore struct {
    jobs       []CronJob
    executions []CronExecution
    createErr  error
}
func (f *fakeCronStore) CreateJob(ctx context.Context, job CronJob) error {
    if f.createErr != nil { return f.createErr }
    f.jobs = append(f.jobs, job)
    return nil
}
```

**Test naming:** `Test{Function}_{Scenario}` or `Test{Function}_{Behavior}`
```go
func TestRunOnboard_SavesInteractiveConfig(t *testing.T) { ... }
func TestScheduler_RunDueJobs_ExecutesAndRecords(t *testing.T) { ... }
func TestClassifyPrompt_ContainsAgentNames(t *testing.T) { ... }
```

**Setup pattern:** `t.TempDir()` + `t.Setenv()` for isolated file system tests
```go
tmpDir := t.TempDir()
t.Setenv("AURELIA_HOME", tmpDir)
```

### Integration Tests

**Approach:** Wire real components together (e.g., SQLite store with scheduler). No external service mocking.

### E2E Tests

**Location:** `e2e/e2e_test.go`, `e2e/bridge_helper_test.go`
**Approach:** Test full wiring — bootstrap, routing, scheduling, bridge protocol
**Helper:** Bridge protocol test helpers for NDJSON communication

## Test Execution

**Commands:**
```bash
go test ./...                    # All tests
go test ./internal/cron/...      # Package-specific
go test -run TestScheduler ./... # Specific test
```

**CI:** GitHub Actions on Windows Latest with Go 1.25.8
- Caching: `~/.cache/go-build` + `~/go/pkg/mod`
- Single step: `go test ./...`

## Patterns & Observations

- **No external mocking framework** — all fakes are hand-written structs implementing interfaces
- **No assertion library** — uses `if val != expected { t.Fatalf(...) }` pattern
- **Subtests:** Used occasionally with `t.Run()` but not the primary pattern
- **Table-driven tests:** Not widely adopted — most tests are individual functions
- **Parallel tests:** Not explicitly used (`t.Parallel()` not observed)
- **Build tags:** None for test separation
