# Architecture

**Pattern:** Modular Monolith — single binary with well-separated internal packages

## High-Level Structure

```
┌─────────────┐    Telegram     ┌──────────────┐
│  Telegram   │ ◄────────────► │   telebot.v3  │
│   Users     │   Bot API      │  Long Polling  │
└─────────────┘                └──────┬───────┘
                                      │
                         ┌────────────▼────────────┐
                         │     BotController       │
                         │  (internal/telegram)     │
                         │  Input Pipeline → Output │
                         └────────┬───────┬────────┘
                                  │       │
                    ┌─────────────▼─┐   ┌─▼──────────────┐
                    │    Persona    │   │  Agent Registry │
                    │  (identity,   │   │  (.md files,    │
                    │   soul, user) │   │   routing,      │
                    │               │   │   classification)│
                    └───────────────┘   └────────────────┘
                                  │
                         ┌────────▼────────┐
                         │     Bridge      │
                         │  (Go ↔ TS IPC)  │
                         │  NDJSON mux     │
                         └────────┬────────┘
                                  │ stdin/stdout
                         ┌────────▼────────┐
                         │  bridge/index.ts │
                         │  Claude Agent    │
                         │  SDK wrapper     │
                         └─────────────────┘
                                  │
            ┌─────────────────────┼──────────────────┐
            │                     │                  │
   ┌────────▼──────┐    ┌────────▼──────┐   ┌──────▼───────┐
   │  Cron Scheduler│   │ Session Store │   │  Cloud MCPs  │
   │  (SQLite, poll │   │ (in-memory,   │   │  (OAuth,     │
   │   every 15s)   │   │  per-chat)    │   │   5min cache)│
   └───────────────┘   └───────────────┘   └──────────────┘
```

## Identified Patterns

### NDJSON Request Multiplexing
**Location:** `internal/bridge/bridge.go`
**Purpose:** Multiple concurrent LLM requests over a single long-lived process
**Implementation:** Atomic request counter generates IDs. `readLoop()` goroutine routes events to per-request buffered channels (cap=16). Terminal events (`result`, `error`) close the channel.
**Example:** `Bridge.Execute()` → creates channel → sends JSON → returns `<-chan Event`

### Fire-and-Forget Async Execution
**Location:** `internal/telegram/input_pipeline.go`
**Purpose:** Non-blocking Telegram message handling — handler returns immediately, results sent asynchronously
**Implementation:** `processInput()` launches `go executeAsync()` which runs bridge request, processes streaming events, and sends Telegram reply on completion.
**Example:** `input_pipeline.go:56` → `go bc.executeAsync(...)`

### Constructor Injection with Interfaces
**Location:** All packages
**Purpose:** Testable, loosely coupled components
**Implementation:** Every struct receives dependencies via `New()` constructor. Key interfaces: `cron.Store`, `cron.Runtime`, `BridgeExecutor`, `ChatSender`, `PersonaBuilder`. Tests use hand-written fakes.
**Example:** `cron.NewScheduler(store, runtime, clock, config)`

### Persona-Based System Prompt Assembly
**Location:** `internal/persona/`, `internal/telegram/input_pipeline.go`
**Purpose:** Dynamic system prompt construction from identity files + agent config + context
**Implementation:** Layers: Persona (IDENTITY+SOUL+USER) → Agent instructions → Cron instructions → Telegram context. Each layer is optional.
**Example:** `buildSystemPrompt()` in `input_pipeline.go`

### Embedded Bridge Bundle
**Location:** `internal/bridge/embed.go`, `internal/bridge/setup.go`
**Purpose:** Self-contained binary with TypeScript bridge included
**Implementation:** `go:embed` bundles the TS code. On first run, writes to `~/.aurelia/bridge/`, installs npm deps. Auto-updates when embedded bundle changes.

## Data Flow

### Telegram Message → LLM Response

1. **Input:** Telegram long poller receives message → `handleText/Photo/Voice/Document`
2. **Bootstrap:** Check if user needs first-run persona setup (`popPendingBootstrap`)
3. **Routing:** `routeAgent()` — `@name` prefix match OR LLM classification (15s timeout)
4. **Prompt:** `buildSystemPrompt()` assembles persona + agent + context layers
5. **Execution:** `bridge.Execute()` sends NDJSON request, returns event channel
6. **Streaming:** `processBridgeEventsAsync()` accumulates assistant text, tracks tools, manages session
7. **Output:** `SendTextReply()` chunks at 3900 chars, converts MD→HTML, handles rate limits
8. **Session:** SessionID stored for context resumption on next message

### Cron Job Execution

1. **Poll:** Scheduler ticks every 15s, queries `ListDueJobs(now, limit=50)`
2. **Dedup:** `sync.Map.LoadOrStore(jobID)` prevents concurrent runs of same job
3. **Execute:** `BridgeCronRuntime.ExecuteJob()` builds persona+agent prompt, calls `bridge.ExecuteSync()`
4. **Record:** Atomic transaction: `RecordExecutionTx` + `UpdateJobTx`
5. **Deliver:** `TelegramDelivery.Deliver()` sends result to `target_chat_id`
6. **Schedule:** Compute `nextRunAt` (cron) or deactivate (once)

## Code Organization

**Approach:** Feature-based packages under `internal/`

**Module boundaries:**
- `cmd/aurelia/` — CLI entry points, dependency wiring, lifecycle
- `internal/telegram/` — All Telegram I/O, message processing, rendering
- `internal/bridge/` — TypeScript process management, NDJSON protocol
- `internal/cron/` — Scheduler, store, runtime, delivery (self-contained with SQLite)
- `internal/agents/` — Agent definition loading, routing, classification
- `internal/persona/` — Identity file parsing, system prompt building
- `internal/session/` — In-memory session and token tracking
- `internal/config/` — App configuration loading, provider management
- `internal/runtime/` — Path resolution, directory bootstrapping
- `pkg/stt/` — Speech-to-text client (Groq)
- `bridge/` — TypeScript source for Claude Agent SDK wrapper
