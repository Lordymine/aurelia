# Tech Stack

**Analyzed:** 2026-03-26

## Core

- Language: Go 1.25.0
- Module: `github.com/kocar/aurelia`
- Package manager: Go modules
- Build: `go build -trimpath -ldflags "-s -w" -o ./build/aurelia ./cmd/aurelia`
- Hot reload: Air (`.air.toml` configured)

## Backend

- API Style: Telegram Bot API via `gopkg.in/telebot.v3` v3.3.8
- LLM Bridge: TypeScript process using `@anthropic-ai/claude-agent-sdk` (NDJSON IPC)
- Database: SQLite via `modernc.org/sqlite` v1.46.1 (pure Go, WAL mode)
- Cron: `github.com/robfig/cron/v3` v3.0.1 (expression parsing)
- Markdown: `github.com/yuin/goldmark` v1.7.8 (Telegram HTML rendering)
- UUID: `github.com/google/uuid` v1.6.0
- Terminal UI: `golang.org/x/term` v0.41.0 (onboarding CLI)
- YAML: `gopkg.in/yaml.v3` v3.0.1 (agent frontmatter parsing)

## Bridge (TypeScript)

- Runtime: Node.js with `tsx` or `--experimental-strip-types`
- SDK: `@anthropic-ai/claude-agent-sdk` (latest)
- TypeScript: ^5.7.0
- Target: ES2022, ESNext modules

## Testing

- Unit/Integration: Standard `testing` package (no testify assertions)
- Mocking: Hand-written fakes (e.g., `fakeCronStore`)
- E2E: `e2e/` directory with bridge integration tests
- CI: `go test ./...` on Windows Latest, Go 1.25.8

## External Services

- LLM: Anthropic, Kimi, OpenRouter, Zai, Alibaba (multi-provider)
- Telegram: Bot API (long polling)
- STT: Groq Whisper API (`whisper-large-v3`)
- Cloud MCPs: Anthropic MCP proxy (OAuth-based discovery)

## Development Tools

- Linting: golangci-lint v2.10.1 (CI workflow, no local config)
- Security: gitleaks (secret scanning), govulncheck (vulnerability check)
- CI: GitHub Actions (Windows)
- Formatting: `gofmt` (implicit standard)
