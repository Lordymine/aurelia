# Code Conventions

## Naming Conventions

**Files:**
`lowercase_with_underscores.go` — grouped by feature within package
Examples: `bot_middleware.go`, `cron_handlers.go`, `store_executions.go`, `markdown_renderer_blocks.go`

**Functions/Methods:**
Exported: `CamelCase` with verb-first. Private: `camelCase`.
Examples: `NewBotController()`, `RunDueJobs()`, `buildSystemPrompt()`, `computeNextRun()`, `isAllowedUser()`

**Variables:**
Short, lowercase camelCase. Abbreviations for common patterns.
Examples: `mu sync.Mutex`, `bc *BotController`, `cfg *AppConfig`, `ctx context.Context`

**Receiver names:** Single letter or short abbreviation matching type.
Examples: `(s *Store)`, `(b *Bridge)`, `(bc *BotController)`, `(r *Registry)`

**Constants:**
Unexported `camelCase` with descriptive names.
Examples: `defaultMaxIterations = 500`, `defaultLLMProvider = "kimi"`, `defaultSTTProvider = "groq"`

## Code Organization

**Import ordering:** Three groups separated by blank lines
1. Standard library (alphabetical)
2. Third-party (alphabetical)
3. Internal `github.com/kocar/aurelia/...` (alphabetical)

```go
import (
    "context"
    "fmt"
    "log"

    "gopkg.in/telebot.v3"

    "github.com/kocar/aurelia/internal/bridge"
    "github.com/kocar/aurelia/internal/config"
)
```

**Aliases:** Used sparingly for disambiguation: `robfigcron "github.com/robfig/cron/v3"`

**File structure:** Types → constructors → public methods → private methods. Two blank lines between top-level declarations.

**Struct initialization:** Always explicit field names, never positional.
```go
bc := &BotController{
    bot:    b,
    config: cfg,
    bridge: br,
}
```

## Error Handling

**Pattern:** Wrapped returns with operation context prefix
```go
if err != nil {
    return nil, fmt.Errorf("reading agents dir: %w", err)
}
```

**Key practices:**
- Early returns on error — no nested else blocks
- Context prefix describes the operation: `"open cron sqlite store: %w"`
- Blank assignment for intentionally ignored errors: `_ = db.Close()`
- Logging via `log.Printf("package.component: message %s", val)` for non-fatal issues
- Graceful degradation: nil-safe checks before use, optional components allowed

## Comments/Documentation

**Style:** Pragmatic — doc comments on exported APIs, inline comments on non-obvious logic.

```go
// Bridge manages a long-lived TypeScript bridge process and communicates via
// stdin/stdout using NDJSON.
type Bridge struct { ... }

// safeClose closes a channel, recovering from panic if already closed.
func safeClose(ch chan Event) { ... }
```

**Coverage:** ~60% exported functions documented, ~40% complex private logic. No comments on self-explanatory code.

## Concurrency

- `sync.Mutex` / `sync.RWMutex` for shared state protection
- `sync.Map` for concurrent deduplication (cron jobs)
- `atomic.Uint64` for counters
- Buffered channels (cap=16) for event streaming
- `context.Context` propagation for cancellation
- `select{}` for non-blocking sends with backpressure handling
