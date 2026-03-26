# External Integrations

## Telegram Bot API

**Service:** Telegram Bot API via telebot.v3
**Purpose:** Primary user interface — receive messages, send responses
**Implementation:** `internal/telegram/` (25 files)
**Configuration:** `TelegramBotToken` + `TelegramAllowedUserIDs` in app config
**Authentication:** Bot token (long polling, 10s timeout)

**Capabilities:**
- Text, photo, document, voice/audio message handling
- Album buffering (900ms wait for grouped photos)
- Chunked output (3900 char limit) with Markdown→HTML conversion
- Typing indicators during LLM processing
- Emoji reactions
- Rate limit handling (FloodError detection + retry)
- User whitelist middleware

**Commands:** `/start`, `/help`, `/cwd`, `/reset`, `/cron`, `/agents`

## LLM Bridge (Claude Agent SDK)

**Service:** Anthropic Claude (and compatible providers)
**Purpose:** LLM inference, tool use, agentic loops
**Implementation:** `internal/bridge/` (Go client) + `bridge/index.ts` (TS wrapper)
**Configuration:** Provider API key + base URL in app config
**Authentication:** API key via `ANTHROPIC_API_KEY` env var, or OAuth for subscription mode

**Protocol:** NDJSON over stdin/stdout
- Request: `{command, prompt, request_id, options}`
- Events: `system` → `assistant`/`tool_use` → `result`/`error`
- Timeout: 10 minutes per query
- Multiplexed: Multiple concurrent requests via request_id

**Supported Providers:**
| Provider | Base URL | Default Model |
|----------|----------|---------------|
| Anthropic | (default) | `claude-sonnet-4-6` |
| Kimi | `api.kimi.com/coding/` | `kimi-k2-thinking` |
| OpenRouter | `openrouter.ai/api/v1` | (user-selected) |
| Zai | `api.z.ai/api/anthropic` | (user-selected) |
| Alibaba | `dashscope-intl.aliyuncs.com/apps/anthropic` | (user-selected) |

## Cloud MCP Servers

**Service:** Anthropic MCP Proxy
**Purpose:** Discover and connect to cloud-hosted MCP servers (tools)
**Implementation:** `bridge/index.ts` (lines 54-122)
**Configuration:** OAuth token from `~/.claude/.credentials.json`
**Authentication:** Bearer token + `anthropic-beta: mcp-servers-2025-12-04`

**Flow:**
1. Read OAuth from `~/.claude/.credentials.json` → `claudeAiOauth.accessToken`
2. Fetch `GET api.anthropic.com/v1/mcp_servers?limit=1000`
3. Convert to `claudeai-proxy` format with MCP proxy URL
4. Cache for 5 minutes
5. Merge with agent-defined MCP servers (agent config wins)

## Speech-to-Text (Groq)

**Service:** Groq Whisper API
**Purpose:** Transcribe voice messages and audio files
**Implementation:** `pkg/stt/groq.go`
**Configuration:** `groq` API key in provider config
**Authentication:** Bearer token

**Endpoint:** `POST https://api.groq.com/openai/v1/audio/transcriptions`
**Model:** `whisper-large-v3`
**Format:** Multipart form-data (file + model + response_format=json)

## SQLite Database

**Service:** Embedded SQLite (modernc.org/sqlite, pure Go)
**Purpose:** Persistent storage for cron jobs and execution history
**Implementation:** `internal/cron/store*.go`
**Configuration:** `~/.aurelia/data/cron.db` (configurable via `DBPath`)
**Authentication:** N/A (local file)

**Tables:**
- `cron_jobs` — Job definitions, schedule, status
- `cron_executions` — Execution history, cost tracking

**Features:** WAL mode, transactions via `WithTx()`, indexed queries

## File System (Persona & Config)

**Purpose:** Persistent identity, configuration, and agent definitions

**Runtime directory:** `~/.aurelia/` (overridable via `$AURELIA_HOME`)
```
~/.aurelia/
├── config/app.json              # Main configuration
├── config/mcp_servers.json      # MCP server definitions
├── data/cron.db                 # SQLite database
├── memory/personas/             # IDENTITY.md, SOUL.md, USER.md
├── memory/OWNER_PLAYBOOK.md     # Optional owner instructions
├── agents/*.md                  # Agent definitions (YAML frontmatter)
└── bridge/                      # TypeScript runtime files
```

**Temporary files:** Media downloads to `os.TempDir()` (photos, documents, audio)

## Background Jobs

**System:** Custom polling scheduler (no external queue)
**Location:** `internal/cron/scheduler.go`
**Interval:** 15 seconds
**Capacity:** Up to 50 jobs per tick
**Deduplication:** `sync.Map` prevents concurrent runs of same job

**Job types:**
- Recurring: Cron expressions (e.g., `"0 9 * * MON"`)
- One-time: Absolute timestamp (`run_at`)
- Agent-scheduled: Auto-registered from agents with `schedule` field
