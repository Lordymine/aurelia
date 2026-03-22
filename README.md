# Aurelia OS

Autonomous agent operating system accessible via Telegram. Talk naturally — Aurelia decides whether to respond directly, delegate to a specialist agent, or schedule automated execution.

## Architecture

Go daemon (24/7) ↔ TypeScript Bridge ↔ Claude Code SDK

- **Bridge**: Long-lived TypeScript process wrapping `@anthropic-ai/claude-agent-sdk`. Communicates via stdin/stdout NDJSON with request multiplexing.
- **Semantic Memory**: SQLite + local ONNX embeddings (all-MiniLM-L6-v2) with cosine similarity search.
- **Agents**: Configurable in markdown (`~/.aurelia/agents/`). Each agent has its own model, tools, MCPs, and cwd.
- **Persona**: IDENTITY.md + SOUL.md + USER.md define personality and behavior.
- **Cron**: Persistent scheduler with Telegram delivery. Create schedules via natural conversation.
- **Multi-provider**: Anthropic (API key or Max subscription), Kimi, OpenRouter, Z.ai, Alibaba.

## Quick Start

1. Install dependencies:
   ```bash
   cd bridge && npm install
   ```

2. Run onboarding:
   ```bash
   go run ./cmd/aurelia/ onboard
   ```

3. Start:
   ```bash
   go run ./cmd/aurelia/
   ```

   Or with hot reload:
   ```bash
   air
   ```

4. Send `/start` to your bot on Telegram.

## Telegram Commands

| Command | Description |
|---------|-------------|
| `/start` | Setup persona (first run) |
| `/help` | List commands |
| `/cwd <path>` | Set working directory |
| `/reset` | New session |
| `/cron` | Manage schedules |
| `/agents` | List agents |

## CLI

```bash
aurelia onboard                    # Configure providers and Telegram
aurelia cron add "0 9 * * *" "..." # Create schedule
aurelia cron list                  # List schedules
aurelia telegram react <chat> <msg> <emoji>  # React to message
```

## Project Structure

```
cmd/aurelia/          CLI entry point + onboarding
internal/bridge/      Go ↔ TypeScript Bridge (long-lived, multiplexed)
internal/telegram/    Telegram I/O, pipeline, progress, sessions
internal/memory/      Semantic memory (SQLite + local ONNX embeddings)
internal/agents/      Agent registry (markdown definitions)
internal/persona/     Persona loader (IDENTITY/SOUL/USER)
internal/cron/        Persistent cron scheduler
internal/config/      App configuration
internal/runtime/     Path resolver + bootstrap
pkg/stt/              Speech-to-text (Groq)
bridge/               TypeScript Bridge (Claude Agent SDK wrapper)
```

## Development

```bash
go build ./...        # Build
go test ./... -short  # Test
go vet ./...          # Lint
air                   # Hot reload
```
