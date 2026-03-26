# Project

## Vision

Aurelia OS is an autonomous agent operating system accessible via Telegram. The goal is not to reimplement what Claude Code already does — it's to **orchestrate it**, adding persistence, scheduling, multi-project support, and a natural Telegram interface on top.

One persistent Go daemon, many projects, many agents.

## Goals

- **Natural interface** — Talk to an AI assistant via Telegram with text, photos, voice, documents. No CLI required for daily use.
- **Agent orchestration** — Route messages to specialist agents, schedule autonomous execution, deliver results back to Telegram.
- **Local-first** — Single binary, SQLite, no cloud dependencies beyond LLM providers. Runs on your machine, owns your data.
- **Stay light** — Don't rebuild what Claude Code SDK already provides. Wrap it, orchestrate it, extend it.
- **Multi-provider** — Not locked to Anthropic. Support Kimi, OpenRouter, Zai, Alibaba, and whatever comes next.

## Constraints

- **Single user** — Personal assistant, not a multi-tenant platform
- **Telegram-only interface** — No web UI, no other chat platforms (for now)
- **Bridge dependency** — LLM reasoning requires Node.js runtime for Claude Agent SDK
- **Windows primary** — CI runs on Windows, developed on Windows
- **No Docker** — Single binary deployment, no container orchestration

## Current State

- Core loop working: Telegram → Agent routing → Bridge → Claude SDK → Response
- Persona system: IDENTITY.md + SOUL.md + USER.md assembled into system prompts
- Cron scheduler: SQLite-backed, recurring and one-time jobs, Telegram delivery
- Multi-modal input: text, photos (albums), voice (Groq STT), documents
- Session continuity: resume via session ID, auto-reset on token threshold
- Agent registry: markdown-defined agents with model/tool/MCP overrides
- Onboarding CLI: interactive setup for providers, tokens, and configuration
- ~6.8K Go LOC + ~400 TS LOC, comprehensive test coverage
