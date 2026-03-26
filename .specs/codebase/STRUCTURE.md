# Project Structure

**Root:** `RafaClaw-aurelia-os/`

## Directory Tree

```
.
├── .air.toml                  # Hot reload config
├── .github/workflows/         # CI: test, lint, vulncheck, gitleaks
├── AGENTS.md                  # Architecture docs for AI assistants
├── README.md                  # Project overview & setup
├── go.mod / go.sum            # Go dependencies
├── mcp_servers.example.json   # MCP config template
│
├── cmd/aurelia/               # CLI entry points
│   ├── main.go                # Subcommand dispatcher
│   ├── app.go                 # Dependency wiring & lifecycle
│   ├── onboard*.go            # Interactive setup wizard (5 files)
│   ├── cron_cli.go            # Cron job management CLI
│   └── telegram_cli.go        # Telegram message CLI
│
├── bridge/                    # TypeScript bridge source
│   ├── index.ts               # Claude Agent SDK wrapper (~400 LOC)
│   ├── bundle.js              # Compiled JS (embedded in Go binary)
│   ├── package.json           # SDK dependency
│   └── tsconfig.json
│
├── internal/                  # Core application packages
│   ├── agents/                # Agent registry & routing (5 files)
│   ├── bridge/                # Go↔TS IPC client (7 files)
│   ├── config/                # App configuration (4 files)
│   ├── cron/                  # Scheduled jobs (14 files)
│   ├── persona/               # Identity & prompt assembly (10 files)
│   ├── runtime/               # Path resolution (5 files)
│   ├── session/               # Session & token tracking (4 files)
│   └── telegram/              # Telegram bot I/O (25 files)
│
├── pkg/stt/                   # Public: speech-to-text client (2 files)
│
└── e2e/                       # End-to-end tests (2 files)
```

## Module Organization

### cmd/aurelia — Application Entry
**Purpose:** CLI commands, dependency wiring, process lifecycle
**Key files:** `main.go` (dispatch), `app.go` (bootstrap + start/stop)

### internal/telegram — Telegram Bot Interface
**Purpose:** All Telegram I/O — input handling, output formatting, markdown rendering, commands
**Key files:** `bot.go` (controller), `input_pipeline.go` (message flow), `output.go` (event processing), `send.go` (chunked delivery)

### internal/bridge — LLM Bridge Client
**Purpose:** Manages TypeScript process, NDJSON protocol, request multiplexing
**Key files:** `bridge.go` (process + IPC), `protocol.go` (types), `events.go` (event model), `embed.go` (bundle embedding)

### internal/cron — Job Scheduler
**Purpose:** Persistent scheduled job execution with SQLite storage
**Key files:** `scheduler.go` (polling loop), `store.go` + `store_*.go` (SQLite CRUD), `runtime.go` (bridge execution), `delivery.go` (Telegram delivery)

### internal/persona — Identity Management
**Purpose:** Load persona files and assemble system prompts
**Key files:** `canonical_service.go` (prompt builder), `loader.go` (file parser)

### internal/agents — Agent Registry
**Purpose:** Load agent definitions from markdown, route messages to agents
**Key files:** `registry.go` (loading + routing), `types.go` (Agent struct)

## Where Things Live

**Telegram message processing:**
- Input handlers: `internal/telegram/input.go`, `input_pipeline.go`
- Output processing: `internal/telegram/output.go`
- Message sending: `internal/telegram/send.go`
- Markdown→HTML: `internal/telegram/markdown*.go`

**LLM communication:**
- Go client: `internal/bridge/bridge.go`
- TS wrapper: `bridge/index.ts`
- Protocol: `internal/bridge/protocol.go`, `events.go`

**Persistence:**
- Database: `internal/cron/store*.go` (SQLite)
- Sessions: `internal/session/store.go` (in-memory)
- Config: `internal/config/config.go` (JSON file)
- Personas: `internal/persona/loader*.go` (markdown files)

**Runtime configuration:**
- App config: `~/.aurelia/config/app.json`
- Agent defs: `~/.aurelia/agents/*.md`
- Personas: `~/.aurelia/memory/personas/`
- Database: `~/.aurelia/data/cron.db`
