# CLAUDE.md

Instructions for coding agents working in this repository.

## Development Commands

```bash
go build ./...           # compile check
go test ./... -short     # fast tests
go test ./... -v         # full test suite
go vet ./...             # static analysis
```

Bridge rebuild (after modifying `bridge/index.ts`):

```bash
cd bridge && npx esbuild index.ts --bundle --platform=node --target=node18 --outfile=bundle.js --format=esm
cp bundle.js ../internal/bridge/bundle.js
```

## Workflow

1. **Plan** — Understand the problem, break into atomic tasks
2. **Review** — Question the plan before executing
3. **Execute** — One atomic task at a time, test-first
4. **Validate** — Run tests, verify completion criteria
5. **Commit** — Conventional Commits: `type(scope): description`

For trivial tasks, implement directly and validate.

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

## Key Packages

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

## Reference

- Architecture and codebase details: `.specs/codebase/`
- Project vision and roadmap: `.specs/project/`
