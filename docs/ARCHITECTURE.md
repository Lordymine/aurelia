# ARCHITECTURE

## Context

Aurelia OS is a local-first agent operating system built in Go.

The system delegates LLM reasoning to a TypeScript Bridge process that wraps the Claude SDK, while Go owns orchestration, memory, scheduling, identity, and interfaces.

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
  cron/            [schedule store, scheduler, bridge-backed runtime]
  memory/          [semantic memory with sqlite-vec embeddings]
  persona/         [identity files, prompt assembly]
  runtime/         [instance and project path resolution]
  telegram/        [Telegram bot handlers]

pkg/
  stt/             [speech-to-text]
```

## Bridge Protocol

The Bridge is the boundary between Go and Claude. Go spawns a short-lived TypeScript process per execution, communicates via stdin/stdout using JSON/NDJSON.

Flow:

1. Go serializes a `Request` (command, prompt, options) as JSON to stdin
2. Bridge process reads the request, calls the Claude SDK
3. Bridge streams NDJSON events back on stdout (text, tool_use, tool_result, result, error)
4. Go reads events and acts on them (forwarding text to Telegram, storing results, etc.)

The Bridge is stateless from Go's perspective — each call is independent.

## Agent Registry

Agents are defined as markdown files with YAML frontmatter. The registry loads all `.md` files from a configurable directory.

Frontmatter fields: `name`, `description`, `model` (optional override), `schedule` (cron expression), `mcp_servers`, `allowed_tools`.

The markdown body becomes the system prompt for that agent.

Scheduled agents are registered with the cron scheduler at startup.

## Memory

The memory system uses SQLite with sqlite-vec for semantic search.

Capabilities:

- Store and retrieve text with vector embeddings
- Semantic similarity search
- Deterministic recent-message window

Embeddings are generated via configurable provider (default: Gemini text-embedding-004).

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

### Memory

`internal/memory` provides semantic storage and retrieval backed by SQLite + sqlite-vec.

### Scheduling

`internal/cron` persists schedules in SQLite and executes them through the Bridge.

## Architectural Rules

1. Telegram is an interface layer, not a domain layer.
2. Identity rules belong in `persona`.
3. Memory and embeddings belong in `memory`.
4. Agent definitions are declarative markdown, not code.
5. The Bridge is the only path to LLM reasoning — Go never calls LLM APIs directly.
6. Long-lived state persists in SQLite.
7. New code should preserve the modular shape.
8. Architecture changes must be reflected here before the task is complete.

## Current Capabilities

- Bridge-based LLM execution via Claude SDK
- Agent registry with markdown-defined agents
- Semantic memory with sqlite-vec
- Telegram text and audio input
- Cron scheduling with bridge-backed execution
- Configurable multi-provider support
- Persona-driven identity and prompt assembly

## Current Constraints

- Bridge requires Node.js runtime and `npx` available on PATH
- Embedding provider must be configured separately from LLM provider
- No multi-agent orchestration yet (single agent per execution)
