# ARCHITECTURE

## Context

Aurelia OS is a local-first agent operating system built in Go.

The system delegates LLM reasoning to a TypeScript Bridge process that wraps the Claude SDK, while Go owns orchestration, session management, scheduling, identity, and interfaces.

User interaction happens through Telegram. Operational state persists in SQLite.

Use this document together with:

- `AGENTS.md`
- `docs/STYLE_GUIDE.md`
- `docs/LEARNINGS.md`

## Architectural Shape

```text
cmd/aurelia/
  main.go          [thin entrypoint]
  app.go           [composition root]

bridge/
  index.ts         [TypeScript Bridge — Claude SDK wrapper]
  package.json     [Bridge dependencies]

internal/
  agents/          [agent registry — markdown-defined agents with YAML frontmatter]
  bridge/          [Go client for the TS Bridge process]
  config/          [configuration loading and validation]
  cron/            [schedule store, scheduler, bridge-backed runtime, delivery]
  session/         [session store, token tracking, auto-reset]
  persona/         [identity files, prompt assembly]
  runtime/         [instance and project path resolution]
  telegram/        [Telegram bot handlers]

pkg/
  stt/             [speech-to-text]
```

## Bridge Protocol

The Bridge is the boundary between Go and Claude. Go starts a long-lived multiplexed TypeScript process, communicates via stdin/stdout using NDJSON with request multiplexing via `request_id`.

Flow:

1. Go serializes a `Request` (command, prompt, request_id, options) as JSON to stdin
2. Bridge process reads the request, calls the Claude SDK
3. Bridge streams NDJSON events back on stdout (system, tool_use, assistant, result, error, pong)
4. Go reads events, correlates by request_id, and acts on them (forwarding text to Telegram, storing results, etc.)

The Bridge is long-lived — multiple concurrent requests are multiplexed over the same process.

## Agent Registry

Agents are defined as markdown files with YAML frontmatter. The registry loads all `.md` files from a configurable directory.

Frontmatter fields: `name`, `description`, `model` (optional override), `schedule` (cron expression), `mcp_servers`, `allowed_tools`.

The markdown body becomes the system prompt for that agent.

Scheduled agents are registered with the cron scheduler at startup.

## Session Management

`internal/session` provides session store and token tracking, extracted from the telegram package for reuse across channels.

Capabilities:

- Session ID management per chat (warm continue / cold resume)
- Token usage accumulation and cost tracking
- Auto-reset when configurable token threshold is exceeded (`max_session_tokens`)

## Runtime Scope Separation

### Repository

Source code, tests, project documentation, default assets.

### Local Instance

Lives outside the repository in `~/.aurelia/`.

Contains: config, SQLite state, logs, persona files, runtime artifacts.

### Target Project

External codebase the agent acts on. Project-specific rules stay local.

## Layer Boundaries

### Entry And Wiring

`cmd/aurelia` loads configuration, builds services, and starts runtimes. Must stay thin.

### Interface Layer

`internal/telegram` receives Telegram events, adapts input, sends output. Not a domain layer.

### Identity

`internal/persona` resolves canonical identity files and assembles system prompts.

### Session

`internal/session` provides session store and token tracking with auto-reset logic.

### Scheduling

`internal/cron` persists schedules in SQLite and executes them through the Bridge.

## Architectural Rules

1. Telegram is an interface layer, not a domain layer.
2. Identity rules belong in `persona`.
3. Session management belongs in `session`.
4. Agent definitions are declarative markdown, not code.
5. The Bridge is the only path to LLM reasoning — Go never calls LLM APIs directly.
6. Long-lived state persists in SQLite.
7. New code should preserve the modular shape.
8. Architecture changes must be reflected here before the task is complete.

## Current Capabilities

- Bridge-based LLM execution via Claude SDK
- Agent registry with markdown-defined agents and LLM classification
- Session management with token tracking and auto-reset
- Telegram text, photo, and audio input
- Cron scheduling with bridge-backed execution and Telegram delivery
- Configurable multi-provider support
- Persona-driven identity and prompt assembly

## Current Constraints

- Bridge requires Node.js runtime available on PATH
- No multi-agent orchestration yet (single agent per execution)
