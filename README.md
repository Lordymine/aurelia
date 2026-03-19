<div align="center">

# Aurelia

<img src="assets/aurelia_cover.png" alt="Aurelia cover" width="720" />

**A local-first autonomous coding agent in Go.**

Telegram-native. Tool-driven. SQLite-backed. Built to stay light.

One persistent agent, many target projects.

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Runtime](https://img.shields.io/badge/Runtime-Local--First-0F172A)](#runtime-model)
[![Architecture](https://img.shields.io/badge/Architecture-Modular_Monolith-1F2937)](docs/ARCHITECTURE.md)
[![Memory](https://img.shields.io/badge/Memory-Layered-0E7490)](#why-aurelia)
[![Storage](https://img.shields.io/badge/Storage-SQLite-003B57?logo=sqlite&logoColor=white)](https://sqlite.org/)
[![Telegram](https://img.shields.io/badge/Interface-Telegram-26A5E4?logo=telegram)](https://core.telegram.org/bots/api)
[![MCP](https://img.shields.io/badge/MCP-Enabled-6E56CF)](https://modelcontextprotocol.io/)
[![Bench](https://img.shields.io/badge/Idle_RAM-25.66_MB-15803D)](#lightweight-baseline)
[![Bench](https://img.shields.io/badge/Binary-23.22_MB-15803D)](#lightweight-baseline)

</div>

## Why Aurelia

`Aurelia` is an autonomous coding agent designed to run locally with a small operational footprint and explicit runtime behavior.

It is built around a practical execution model:

- Go runtime
- SQLite persistence
- tool-driven ReAct loop
- master-led Agent Teams
- deterministic layered memory
- controlled local command execution
- optional MCP expansion

The goal is not to fake autonomy through giant prompts.
The goal is to observe the local environment, act with real tools, persist useful state, and remain understandable under load.

## Lightweight Baseline

Current measured baseline is documented in [docs/BENCHMARKS.md](docs/BENCHMARKS.md).

Release build:

- binary size: `23.22 MB`

Idle runtime baseline on the current Windows machine:

- startup average: `15.75 ms`
- idle working set: `25.66 MB`
- idle private memory: `53.39 MB`
- idle CPU average: `0.00%`

Method summary:

- release build with `-trimpath -ldflags "-s -w"`
- isolated `AURELIA_HOME`
- temporary DB and MCP config paths
- three-run sample

These numbers are local baseline measurements, not universal guarantees across all hardware.

## Core Capabilities

- Telegram-native interaction with text, markdown documents, images, markdown documents, and audio input
- tool-driven ReAct execution
- Agent Teams with master-only user response, task graph, mailbox, and recovery
- SQLite-backed recent memory, facts, notes, archive, and operational task state
- local file operations and controlled command execution
- project-aware `workdir` propagation
- cron scheduling subsystem
- optional MCP tool discovery and registration

## Runtime Model

Aurelia separates three scopes:

1. Repository
2. Local instance
3. Target project

This keeps product source, user runtime state, and external project work clearly separated.

High-level flow:

```mermaid
flowchart LR
    U[User] --> T[Telegram]
    T --> LOOP[ReAct Loop]
    LOOP --> REG[Tool Registry]
    LOOP --> MEM[Layered Memory]
    REG --> TEAM[Agent Teams]
    REG --> CORE[Native Tools]
    REG --> MCP[MCP Tools]
    TEAM --> TASKS[Task Graph + Mailbox]
    LOOP --> OUT[Final Response]
    OUT --> U
```

## Architecture

Aurelia follows a modular monolith layout:

```text
/cmd/aurelia/
  main.go
  app.go
  wiring.go

/internal/
  agent/
  config/
  cron/
  mcp/
  memory/
  persona/
  runtime/
  skill/
  telegram/
  tools/

/pkg/
  llm/
  stt/
```

Architectural source of truth:

- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

Implementation rules:

- [docs/STYLE_GUIDE.md](docs/STYLE_GUIDE.md)

Operational lessons:

- [docs/LEARNINGS.md](docs/LEARNINGS.md)

## Setup

Requirements:

- Go `1.25+`
- Telegram bot token
- one LLM provider key:
  - Kimi API key
  - Anthropic API key
- Google API key
- Kilo API key
- OpenRouter API key
- Z.ai Coding Plan API key
- Alibaba Coding Plan API key
- OpenAI API key
- or a local Codex CLI login for experimental OpenAI Codex mode
- Groq API key for STT

Main runtime config lives in `~/.aurelia/config/app.json`.

Example config:

```json
{
  "llm_provider": "kimi",
  "llm_model": "kimi-k2-thinking",
  "stt_provider": "groq",
  "telegram_bot_token": "your-telegram-bot-token",
  "telegram_allowed_user_ids": [123456789],
  "anthropic_api_key": "",
  "google_api_key": "",
  "kilo_api_key": "",
  "kimi_api_key": "your-kimi-api-key",
  "openrouter_api_key": "",
  "zai_api_key": "",
  "alibaba_api_key": "",
  "openai_auth_mode": "api_key",
  "openai_api_key": "",
  "groq_api_key": "your-groq-api-key",
  "max_iterations": 500,
  "db_path": "C:/Users/you/.aurelia/data/aurelia.db",
  "memory_window_size": 20,
  "mcp_servers_config_path": "C:/Users/you/.aurelia/config/mcp_servers.json"
}
```

If the file does not exist, Aurelia creates it automatically with defaults on first start.

Interactive onboarding:

```bash
go run ./cmd/aurelia onboard
```

The onboarding runs in guided steps:

- select the LLM provider (`Kimi`, `Anthropic`, `Google`, `Kilo Code`, `OpenRouter`, `Z.ai`, `Alibaba`, or `OpenAI`)
- select auth mode when the provider supports multiple paths:
  - `OpenAI`: `api_key` or `codex`
- select the LLM model
- select the STT provider
- configure Telegram
- review runtime settings and save

Model selection notes:

- large catalogs such as `OpenRouter` and `Kilo Code` support live search by model or provider
- the onboarding shows capability badges such as `[vision]`, `[tools]`, and `[free]` when the catalog exposes that metadata
- use `→` in the model step to cycle between `all`, `vision`, `tools`, and `free`

Provider notes:

- `zai` expects a `GLM Coding Plan` key and uses the dedicated coding endpoint
- `alibaba` expects a Coding Plan key and uses the dedicated coding endpoint
- Telegram photos, image documents, and photo albums are forwarded as multimodal input when the current provider/model supports vision
- if the current model does not support image input, Aurelia replies with a natural fallback explaining that limitation
- `openai` in `codex` mode is experimental and should be treated as text-first; use API key mode or another provider for image analysis

OpenAI Codex login:

```bash
go run ./cmd/aurelia auth openai
```

When `llm_provider` is `openai` and `openai_auth_mode` is `codex`, Aurelia uses the local `codex mcp-server` as an experimental runtime backend. This mode requires:

- `codex` installed and available on `PATH`
- a successful local `codex login`
- choosing an OpenAI model in onboarding or `app.json`

Run:

```bash
go run ./cmd/aurelia
```

Release build:

```bash
go build -trimpath -ldflags "-s -w" -o ./build/aurelia.exe ./cmd/aurelia
```

## Documentation

Canonical documents for ongoing work:

- [AGENTS.md](AGENTS.md)
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- [docs/STYLE_GUIDE.md](docs/STYLE_GUIDE.md)
- [docs/LEARNINGS.md](docs/LEARNINGS.md)
- [docs/BENCHMARKS.md](docs/BENCHMARKS.md)

Recommended reading order:

1. `README.md`
2. `AGENTS.md`
3. `docs/ARCHITECTURE.md`
4. `docs/STYLE_GUIDE.md`
5. `docs/LEARNINGS.md`
6. `docs/BENCHMARKS.md`

## Current State

- the active codebase runs under the `Aurelia` identity
- the repository has benchmark-backed footprint data
- the Go test suite is green
- contribution policy and CI workflows are present in the repository
