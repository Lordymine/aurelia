# Provider Expansion Roadmap

## Objective

Track the incremental expansion of LLM providers in `Aurelia` with explicit scope, validation, and documentation checkpoints.

This roadmap is the human-readable plan.

Machine-readable execution state lives in `docs/STATE_PROVIDER_EXPANSION.json`.

## State Machine

Allowed task states:

- `planned`
- `researching`
- `ready`
- `in_progress`
- `blocked`
- `validating`
- `done`

Valid transitions:

- `planned -> researching`
- `researching -> ready`
- `ready -> in_progress`
- `in_progress -> validating`
- `validating -> done`
- `validating -> in_progress`
- `* -> blocked`
- `blocked -> in_progress`

## Tasks

### T00 - Baseline Provider Infrastructure

- scope: add `llm_model`, remove hardcoded LLM model selection, introduce provider model catalog support, and adapt onboarding to select a model
- depends on: none
- official references:
  - https://opencode.ai/docs/providers
- acceptance criteria:
  - `app.json` persists `llm_provider` and `llm_model`
  - runtime uses `llm_model` instead of hardcoded defaults
  - onboarding can fetch or fallback to provider model options and save the selection
  - model discovery failures fall back to curated defaults without aborting onboarding
- validation:
  - `go test ./internal/config ./cmd/aurelia ./pkg/llm`
  - `go test ./...`

### T01 - Anthropic API Key

- scope: integrate Anthropic with API key auth, provider factory support, onboarding support, and model catalog
- depends on: `T00`
- official references:
  - https://github.com/anthropics/anthropic-sdk-go
  - https://docs.anthropic.com/
- acceptance criteria:
  - onboarding supports `anthropic`
  - onboarding stores `anthropic_api_key`
  - model selection lists Anthropic models or curated fallback
  - runtime can build Anthropic provider with selected model
  - docs describe setup and limitations
- validation:
  - provider unit tests
  - onboarding tests
  - `go test ./...`

### T02 - Google API Key

- scope: expose Google Gemini as a first-class provider using API key auth and model catalog
- depends on: `T00`
- official references:
  - https://ai.google.dev/api
  - https://ai.google.dev/gemini-api/docs/api-key
- acceptance criteria:
  - onboarding supports `google`
  - onboarding stores `google_api_key`
  - model selection lists Gemini models or curated fallback
  - runtime builds Gemini with selected model
  - docs describe setup and limitations
- validation:
  - provider unit tests
  - onboarding tests
  - `go test ./...`

### T03 - OpenRouter

- scope: integrate OpenRouter with model discovery and `auto/free` routing options
- depends on: `T00`, `T04`
- official references:
  - https://openrouter.ai/docs/docs/overview/models
  - https://openrouter.ai/docs/guides/routing/routers/free-models-router
- acceptance criteria:
  - onboarding supports `openrouter`
  - onboarding stores `openrouter_api_key`
  - model selection lists remote models plus `openrouter/auto` and `openrouter/free`
  - runtime builds OpenRouter provider with correct endpoint and headers
  - docs explain fixed model vs auto routing vs free routing
- validation:
  - provider unit tests
  - model catalog tests
  - onboarding tests
  - `go test ./...`

### T04 - OpenAI-Compatible Base

- scope: create a shared provider base for OpenAI-compatible chat completions APIs
- depends on: `T00`
- official references:
  - https://opencode.ai/docs/providers
- acceptance criteria:
  - shared request/response implementation handles model, base URL, API key, and optional headers
  - wrappers can configure OpenRouter, Z.ai, and Alibaba without duplicating transport logic
  - tool calling remains compatible with the existing loop
- validation:
  - provider base unit tests
  - `go test ./pkg/llm ./...`

### T05 - Z.ai Coding Plan

- scope: integrate Z.ai GLM Coding Plan using the OpenAI-compatible base and provider-specific model catalog strategy
- depends on: `T04`
- official references:
  - https://docs.z.ai/api-reference/introduction
  - https://docs.z.ai/devpack/tool/cursor
- acceptance criteria:
  - onboarding supports `zai`
  - onboarding stores `zai_api_key`
  - runtime builds Z.ai provider against the Coding Plan endpoint
  - model selection uses curated Coding Plan models
  - docs describe Coding Plan behavior and limitations
- validation:
  - provider unit tests
  - onboarding tests
  - `go test ./...`

### T06 - Alibaba Coding Plan

- scope: integrate Alibaba Coding Plan using the compatible coding endpoint
- depends on: `T04`
- official references:
  - https://www.alibabacloud.com/help/en/model-studio/compatibility-of-openai-with-dashscope
- acceptance criteria:
  - onboarding supports `alibaba`
  - onboarding stores `alibaba_api_key`
  - runtime builds Alibaba provider with the Coding Plan endpoint
  - model selection uses curated coding-focused models
  - docs describe Coding Plan behavior and limitations
- validation:
  - provider unit tests
  - onboarding tests
  - `go test ./...`

### T07 - OpenAI API Key

- scope: integrate OpenAI API key auth and model discovery
- depends on: `T00`, `T04`
- official references:
  - https://platform.openai.com/docs/api-reference
  - https://github.com/openai/openai-go
- acceptance criteria:
  - onboarding supports `openai`
  - onboarding stores `openai_api_key`
  - runtime builds OpenAI provider
  - model selection lists official models or curated fallback
  - docs describe setup and limitations
- validation:
  - provider unit tests
  - onboarding tests
  - `go test ./...`

### T08 - Google OAuth

- scope: add Google OAuth user auth in addition to API key auth
- depends on: `T02`
- official references:
  - https://ai.google.dev/gemini-api/docs/oauth
- acceptance criteria:
  - onboarding supports Google auth mode selection
  - OAuth tokens are persisted and refreshed correctly
  - runtime can build Google provider from OAuth credentials
  - docs clearly explain auth flow and local persistence
- validation:
  - auth persistence tests
  - onboarding tests
  - provider wiring tests
  - `go test ./...`

Status note:

- retired from the active product surface; Aurelia now keeps Google on `api_key` only
- preserved in roadmap only as historical context, not as an active feature

### T09 - OpenAI Codex via MCP

- scope: add experimental OpenAI Codex provider through the local Codex MCP server while keeping OpenAI API key mode in parallel
- depends on: `T07`
- official references:
  - current Codex CLI behavior and current OpenAI model docs
- acceptance criteria:
  - onboarding marks the flow as experimental and local-runtime dependent
  - config persists Codex CLI mode separately from API-key OpenAI mode
  - runtime fails clearly when `codex` is missing or unauthenticated
  - docs clearly state the local installation prerequisite and experimental status
- validation:
  - config tests
  - onboarding tests
  - provider wiring tests
  - `go test ./...`

### T10 - Documentation Finalization

- scope: align `README.md`, `docs/ARCHITECTURE.md`, and any recurring lessons after the provider expansion work
- depends on: `T01`, `T02`, `T03`, `T04`, `T05`, `T06`, `T07`, `T08`, `T09`
- official references:
  - provider-specific official docs used in the completed tasks
- acceptance criteria:
  - canonical docs reflect supported providers, auth modes, and onboarding behavior
  - docs do not mention unsupported providers as available
  - recurring operational traps are captured in `docs/LEARNINGS.md` when relevant
- validation:
  - manual doc review
  - `go test ./...`

### T11 - Kilo Gateway

- scope: integrate the Kilo Code gateway as an OpenAI-compatible provider with automatic model discovery
- depends on: `T04`
- official references:
  - https://kilo.ai/docs/gateway/api-reference
  - https://kilo.ai/docs/gateway/sdks-and-frameworks
  - https://kilo.ai/docs/ai-providers/kilocode
- acceptance criteria:
  - onboarding supports `kilo`
  - onboarding stores `kilo_api_key`
  - model selection lists models from the official Kilo `/models` endpoint or curated fallback
  - runtime builds Kilo provider with the gateway chat completions endpoint
  - docs describe the API-key setup and provider behavior
- validation:
  - provider unit tests
  - model catalog tests
  - onboarding tests
  - `go test ./...`
