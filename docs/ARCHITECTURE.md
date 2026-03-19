# ARCHITECTURE

## Context

`Aurelia` is a local-first autonomous coding agent built in Go.

The system operates through Telegram, persists operational state in SQLite, and executes work through an explicit runtime made of:

- a ReAct loop
- a live tool registry
- Agent Teams orchestration
- layered memory
- controlled local execution
- optional MCP integration

This document is the architectural source of truth for the current codebase.

Use it together with:

- `AGENTS.md`
- `docs/STYLE_GUIDE.md`
- `docs/LEARNINGS.md`

## Architectural Shape

The project is a modular monolith in Go.

```text
/cmd/aurelia/
  main.go          [thin entrypoint]
  app.go           [composition root]
  wiring.go        [dependency and tool registration]

/internal/
  agent/           [loop, tool registry, team manager, task graph, recovery]
  config/          [configuration loading]
  cron/            [store, service, scheduler, runtime]
  mcp/             [manager, discovery, transport]
  memory/          [messages, facts, notes, archive]
  persona/         [canonical identity, prompt building, retrieval, file sync]
  runtime/         [instance and project path resolution]
  skill/           [skill loading and routing]
  telegram/        [handlers, bootstrap, input pipeline, output]
  tools/           [tool definitions, handlers, MCP adapter]

/pkg/
  llm/             [LLM providers and model catalogs]
  stt/             [speech-to-text]
```

Current code still contains legacy repository names and paths.
Those are implementation details under migration, not the product identity to preserve.

## Runtime Scope Separation

The runtime distinguishes three scopes:

### Repository

Contains:

- source code
- tests
- project documentation
- default assets shipped with the product

### Local Instance

Lives outside the repository in a per-user hidden directory.

Contains:

- local config
- canonical app config in `config/app.json`
- SQLite state
- logs
- learned notes
- runtime skills
- canonical persona files

### Target Project

The external codebase the agent is acting on.

It may define project-specific docs and working conventions.

Project-specific rules must not leak into global defaults unless explicitly promoted.

## Layer Boundaries

### Entry And Wiring

`cmd/aurelia` is responsible for:

- loading configuration
- building services
- registering tools
- starting and stopping runtimes

It must stay thin.

Runtime configuration is instance-local and file-backed.

Rule:

- app runtime config is loaded from `~/.aurelia/config/app.json`
- LLM selection is persisted as `llm_provider` plus `llm_model`
- provider-specific auth modes are explicit config fields when needed, such as `openai_auth_mode`
- repository-local `.env` files are not part of the supported runtime config model

Current auth variants:

- `google`: `api_key`
- `kilo`: `api_key`
- `zai`: `coding_plan_api_key`
- `alibaba`: `coding_plan_api_key`
- `openai`: `api_key` or experimental local `codex` CLI mode

Vision input:

- Telegram photos, image documents, and photo albums enter the runtime as multimodal user messages
- providers decide support at runtime based on the selected model
- when the active model does not support image input, Aurelia responds with a natural limitation message instead of silently ignoring the image

### Interface Layer

`internal/telegram` is the interface boundary.

Responsibilities:

- receiving Telegram events
- adapting text, files, and audio into internal input
- sending output back to Telegram

It must not become the source of truth for business rules.

### Domain And Orchestration

`internal/agent` is the center of execution behavior.

Responsibilities:

- ReAct loop
- tool execution contracts
- Agent Teams orchestration
- task graph
- recovery

### Memory And Identity

`internal/memory` persists:

- recent messages
- durable facts
- compact notes
- raw archive

`internal/persona` resolves:

- canonical identity
- prompt assembly
- file synchronization
- owner and project context injection
- retrieval helpers

Important rule:

- identity and operating rules come from canonical persona and project context
- actual runtime capabilities come from the live tool registry

### Tooling

`internal/tools` owns:

- native tool definitions
- handler contracts
- schedule tools
- MCP registration adapters
- team mailbox and control tools

Tool schemas belong here, not in the composition root.

Provider and model metadata should have a single source of truth in `pkg/llm`.

That includes:

- provider identity and labels
- default model selection
- model catalog fallbacks
- model normalization rules
- context window lookup tables

Interface layers and wiring may consume this metadata, but should not redefine it locally.

## Core Runtime Model

### ReAct Execution

The main loop is tool-driven.

The runtime injects:

- available tools
- execution guidance
- workdir guidance
- runtime capabilities

The agent must reason from real capabilities, not assumed capabilities.

### Agent Teams

The multi-agent model is master-led.

Rules:

- `master` is the only agent that answers the end user
- workers operate on explicit tasks
- tasks carry status, ownership, dependencies, result, and canonical `workdir`
- workers may coordinate through internal mailbox tools
- operational team state should persist in SQLite whenever possible

### Memory

The memory model is deterministic and layered.

Priority order:

1. canonical identity
2. stable facts
3. episodic notes
4. recent conversation window

The architecture explicitly avoids treating larger prompts or vector search as the default fix for identity and continuity problems.

Conversation continuity follows a simple local context-budget policy:

- context capacity is resolved from a small provider/model table
- prompt size is estimated locally before execution
- when the estimated prompt crosses the soft threshold, Aurelia compacts older chat history into a deterministic summary note and trims the working window
- when it crosses the hard threshold, Aurelia rotates the working window more aggressively while preserving the latest compact summary plus durable facts and notes

This is intentionally not a provider-driven session-introspection subsystem.
The runtime keeps continuity by deterministic local policy over SQLite-backed memory.
The source of truth for this policy lives in `internal/memory`; interface layers may trigger it, but should not own the rule.

### Local Execution

The runtime is expected to observe the environment, not only describe it.

That includes:

- reading files
- writing files
- listing directories
- running controlled local commands
- acting on a canonical project `workdir`

Operational observability is local-first and intentionally minimal:

- correlated logs carry stable execution identifiers such as `run_id`, `team_id`, and `task_id` when available
- critical runtime boundaries record short operational events in SQLite for recent inspection
- inspection stays inside the product surface through lightweight debug commands instead of external APM dependencies by default

## Architectural Rules

1. Telegram is an interface layer, not a domain layer.
2. Identity and memory rules belong in `persona` and `memory`.
3. Multi-agent orchestration belongs in `agent`.
4. Long-lived operational state should persist in SQLite when practical.
5. Tools are explicit runtime capabilities, not hidden side channels.
6. New code should preserve the modular monolith shape instead of adding ad hoc service sprawl.
7. Workers operating on external projects must preserve canonical `workdir` across local tools.
8. Architecture changes must be reflected here before the task is considered complete.

## Current Capabilities

Implemented in the current codebase:

- tool-driven ReAct loop
- Agent Teams orchestration with task graph, mailbox, recovery, and final synthesis
- SQLite-backed memory and operational persistence
- Telegram text, markdown, image, and audio input flow
- cron scheduling subsystem
- MCP discovery and registration
- project-aware execution through propagated `workdir`

## Current Constraints

Known constraints in the current codebase:

- repository naming still reflects the old product identity
- some documentation is still transitional and overlaps with future canonical docs
- benchmark evidence for memory and CPU footprint is not yet documented
- contribution governance and GitHub policy gates are not fully established yet

## Documentation Policy

This file captures architecture and boundaries.

Implementation conventions belong in `docs/STYLE_GUIDE.md`.

Operational mistakes, traps, and recurring lessons belong in `docs/LEARNINGS.md`.


