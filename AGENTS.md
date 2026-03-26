# AGENTS.md

Instructions for coding agents working in this repository.

---

## Architecture Overview

Aurelia OS is a local-first agent operating system written in Go. It delegates LLM reasoning to a TypeScript Bridge process that wraps the Claude SDK, keeping Go responsible for orchestration, memory, scheduling, and interfaces.

```text
User ──► Telegram ──► Go runtime ──► Bridge (TS) ──► Claude SDK
                          │
                          ├── Agent Registry (markdown-defined agents)
                          ├── Session Management (token tracking, auto-reset)
                          ├── Cron Scheduler (bridge-backed)
                          └── Persona (identity + context)
```

### Bridge Protocol

The Bridge (`bridge/index.ts`) is a long-lived TypeScript process with request multiplexing via `request_id`. Go starts it once and communicates via stdin/stdout NDJSON.

Request shape: `{ command, prompt, request_id, options: { model, cwd, system_prompt, resume, max_turns, permission_mode, mcp_servers, allowed_tools, disabled_tools } }`

Events: `{ type: "system" | "tool_use" | "assistant" | "result" | "error" | "pong", ... }`

Source: `internal/bridge/` (Go client), `bridge/index.ts` (TS process).

### Agent Markdown Format

Agents are defined as `.md` files with YAML frontmatter:

```markdown
---
name: agent-name
description: What this agent does
model: claude-sonnet-4-6  # optional override
schedule: "0 9 * * *"       # optional cron expression
mcp_servers:                 # optional
  server-name:
    command: ...
allowed_tools:               # optional whitelist
  - tool_name
---

System prompt / instructions for this agent go here as the body.
```

Loaded by `internal/agents/Registry` from a configurable directory.

### Config Schema

Runtime config lives in `~/.aurelia/config/app.json`:

```json
{
  "llm_provider": "anthropic",
  "llm_model": "claude-sonnet-4-6",
  "telegram_bot_token": "...",
  "telegram_allowed_user_ids": [123],
  "anthropic_api_key": "...",
  "stt_provider": "groq",
  "groq_api_key": "...",
  "max_iterations": 500,
  "max_session_tokens": 100000
}
```

Source: `internal/config/`.

### Key Packages

| Package | Responsibility |
|---------|---------------|
| `cmd/aurelia/` | Entrypoint, wiring, onboarding |
| `internal/bridge/` | Go client for the TS Bridge process |
| `internal/agents/` | Agent registry (load markdown definitions) |
| `internal/session/` | Session store and token tracking |
| `internal/persona/` | Identity files, prompt assembly |
| `internal/cron/` | Schedule store, scheduler, bridge-backed runtime |
| `internal/telegram/` | Telegram bot handlers |
| `internal/config/` | Config loading and validation |
| `internal/runtime/` | Instance and project path resolution |
| `bridge/` | TypeScript Bridge (Claude SDK wrapper) |
| `pkg/stt/` | Speech-to-text |

---

## Workflow

1. **Plan** — Understand the problem, break into atomic tasks
2. **Review** — Question the plan before executing
3. **Execute** — One atomic task at a time, test-first
4. **Validate** — Run tests, verify completion criteria
5. **Commit** — Conventional Commits: `type(scope): description`

For trivial tasks, implement directly and validate.

---

## Development Commands

```bash
go build ./...           # compile check
go test ./... -short     # fast tests
go test ./... -v         # full test suite
go vet ./...             # static analysis
```

Bridge is embedded in the Go binary via `go:embed`. To rebuild after modifying `bridge/index.ts`:

```bash
cd bridge && npx esbuild index.ts --bundle --platform=node --target=node18 --outfile=bundle.js --format=esm
cp bundle.js ../internal/bridge/bundle.js
```

---

## Rules

- Service layer for business logic — never in handlers or entrypoints
- Errors treated explicitly — no silent swallowing
- `context.Context` with timeout on external operations
- Secrets never in repository — use `~/.aurelia/config/app.json`
- Tests required before marking work complete
- No new dependencies without justification
- Prefer editing over rewriting
- Keep interfaces small
- Update docs when behavior changes

---

## Canonical Documentation

| Document | Scope |
|----------|-------|
| `AGENTS.md` | Agent instructions, architecture overview |
| `docs/ARCHITECTURE.md` | Detailed architecture and boundaries |
| `docs/STYLE_GUIDE.md` | Coding conventions and patterns |
| `docs/LEARNINGS.md` | Operational lessons and recurring mistakes |
