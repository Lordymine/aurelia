# Aurelia OS Rewrite — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rewrite Aurelia as an autonomous agent OS that uses Claude Code CLI as its brain via a TypeScript Bridge, with semantic memory and configurable markdown agents.

**Architecture:** Go process (daemon 24/7) spawns a TypeScript Bridge (~300 lines) that wraps `@anthropic-ai/claude-agent-sdk`. Communication is stdin/stdout NDJSON. Memory is sqlite-vec with embeddings. Agents are markdown files in `~/.aurelia/agents/`. Persona is 3 markdown files (IDENTITY/SOUL/USER).

**Tech Stack:** Go (main process), TypeScript (bridge), sqlite-vec (semantic memory), Claude Agent SDK, Telegram Bot API

**Spec:** `docs/superpowers/specs/2026-03-21-aurelia-os-spec.md`

---

## File Structure

### New files to create

```
bridge/
  index.ts                    # TypeScript Bridge — wraps Claude Agent SDK (~300 LOC)
  package.json                # Dependencies: @anthropic-ai/claude-agent-sdk
  tsconfig.json               # TypeScript config

internal/bridge/
  bridge.go                   # Go side — spawn TS process, stdin/stdout, parse events (~120 LOC)
  bridge_test.go              # Tests for bridge communication (~80 LOC)
  events.go                   # Event types (system, tool_use, tool_result, assistant, result) (~50 LOC)
  protocol.go                 # Request/response JSON protocol types (~30 LOC)

internal/memory/              # REWRITE (replaces current memory)
  store.go                    # sqlite-vec store, Save/Search/Inject (~200 LOC)
  store_test.go               # Tests (~150 LOC)
  embeddings.go               # Embedding generation via API (~80 LOC)
  embeddings_test.go          # Tests (~50 LOC)
  schema.go                   # DDL for memory_vectors table (~30 LOC)

internal/agents/
  registry.go                 # Load markdown agent definitions, route messages (~120 LOC)
  registry_test.go            # Tests (~80 LOC)
  types.go                    # Agent struct, frontmatter parsing (~50 LOC)
```

### Files to modify

```
cmd/aurelia/
  app.go                      # Rewrite bootstrap — remove old deps, wire Bridge + new memory + agents
  wiring.go                   # Remove old tool registry, replace with Bridge dispatch
  main.go                     # Minimal changes (keep CLI dispatcher)
  cron_support.go             # Adapt to use Bridge instead of agent loop

internal/persona/
  canonical_service.go        # Remove memory.MemoryManager dependency
  canonical_service_prompt.go # Simplify prompt building (remove retrieval ranking)
  canonical_service_retrieval.go      # Remove (retrieval ranking)
  canonical_service_retrieval_rank.go # Remove (retrieval ranking)
  conversation_memory.go      # Remove (replaced by semantic memory)
  loader.go                   # Simplify — keep IDENTITY/SOUL/USER loading

internal/cron/
  runtime.go                  # Replace AgentExecutor with Bridge call
  types.go                    # Simplify CronJob to match new agent model

internal/telegram/
  bot.go                      # Rewrite constructor — remove agent/skill/observability deps
  input_pipeline.go           # Replace Loop.Execute with Bridge dispatch
  bootstrap.go                # Simplify
  bootstrap_config.go         # Simplify
  bootstrap_memory.go         # Remove or rewrite for semantic memory
  memory_handlers.go          # Remove (semantic memory replaces)
  ops_handlers.go             # Remove (observability removed)
  cron_handlers.go            # Adapt to new cron model

internal/config/
  config.go                   # Refactor AppConfig for new schema (providers, embedding config)
  mcp.go                      # Simplify — MCPs configured per agent markdown
```

### Files/directories to delete

```
internal/agent/               # Entire directory (56 files, 7,871 LOC)
internal/tools/               # Entire directory (24 files, 3,098 LOC)
internal/observability/       # Entire directory (2 files, 386 LOC)
internal/mcp/                 # Entire directory (9 files, 768 LOC)
internal/skill/               # Entire directory (6 files, 686 LOC)
internal/memory/              # Entire directory (9 files, 1,209 LOC) — recreated from scratch
pkg/llm/                      # Entire directory (31 files, 3,498 LOC)
```

---

## Task 1: Delete modules that are being replaced

**Goal:** Remove all code that Claude Code CLI replaces. Codebase must compile after cleanup.

**Files:**
- Delete: `internal/agent/` (56 files)
- Delete: `internal/tools/` (24 files)
- Delete: `internal/observability/` (2 files)
- Delete: `internal/mcp/` (9 files)
- Delete: `internal/skill/` (6 files)
- Delete: `internal/memory/` (9 files)
- Delete: `pkg/llm/` (31 files)
- Delete: `internal/persona/canonical_service_retrieval.go`
- Delete: `internal/persona/canonical_service_retrieval_rank.go`
- Delete: `internal/persona/conversation_memory.go`
- Delete: `internal/telegram/memory_handlers.go`
- Delete: `internal/telegram/memory_handlers_test.go`
- Delete: `internal/telegram/ops_handlers.go`
- Delete: `internal/telegram/ops_handlers_test.go`
- Delete: `internal/telegram/context_policy_test.go`
- Delete: `internal/telegram/bootstrap_memory.go`
- Modify: `cmd/aurelia/app.go` — stub out bootstrap to compile
- Modify: `cmd/aurelia/wiring.go` — remove tool registry functions
- Modify: `cmd/aurelia/cron_support.go` — stub
- Modify: `internal/telegram/bot.go` — stub constructor
- Modify: `internal/telegram/input_pipeline.go` — stub
- Modify: `internal/telegram/bootstrap.go` — stub
- Modify: `internal/telegram/bootstrap_config.go` — stub
- Modify: `internal/cron/runtime.go` — stub
- Modify: `internal/persona/canonical_service.go` — remove memory dep
- Modify: `internal/persona/canonical_service_prompt.go` — remove retrieval calls
- Modify: `go.mod` — remove unused dependencies

- [ ] **Step 1: Delete the 7 module directories**

```bash
cd C:\Users\kocar\Documents\RafaClaw-aurelia-os
rm -rf internal/agent internal/tools internal/observability internal/mcp internal/skill internal/memory pkg/llm
```

- [ ] **Step 2: Delete persona files that won't be needed**

```bash
rm internal/persona/canonical_service_retrieval.go
rm internal/persona/canonical_service_retrieval_rank.go
rm internal/persona/conversation_memory.go
```

- [ ] **Step 3: Delete telegram files that won't be needed**

```bash
rm internal/telegram/memory_handlers.go internal/telegram/memory_handlers_test.go
rm internal/telegram/ops_handlers.go internal/telegram/ops_handlers_test.go
rm internal/telegram/context_policy_test.go
rm internal/telegram/bootstrap_memory.go
```

- [ ] **Step 4: Stub cmd/aurelia/app.go**

Rewrite `app.go` to a minimal struct that only holds what will remain:
- `resolver *runtime.PathResolver`
- `cronStore *cron.SQLiteCronStore`
- `bot *telegram.BotController`
- `cronScheduler *cron.Scheduler`

Remove all imports of deleted packages. `bootstrapApp()` returns a minimal app that doesn't wire deleted components. It's OK for the app to not function yet — just needs to compile.

- [ ] **Step 5: Stub cmd/aurelia/wiring.go**

Remove `buildToolRegistry()`, `registerScheduleTools()`, `maybeRegisterMCPTools()`, `registerSpawnAgentTool()`. File can be empty or deleted.

- [ ] **Step 6: Stub cmd/aurelia/cron_support.go**

Remove references to agent.ExecutionContext. Keep minimal cron setup function signature.

- [ ] **Step 7: Stub internal/telegram/bot.go constructor**

Simplify `NewBotController()` to only take:
- `cfg *config.AppConfig`
- `s stt.Transcriber`
- `canonical *persona.CanonicalIdentityService`
- `personasDir string`

Remove all agent/skill/observability parameters and fields.

- [ ] **Step 8: Stub internal/telegram/input_pipeline.go**

Remove references to agent.Loop. Leave the pipeline shell (Telegram message → text extraction) but comment out the agent execution part with a `// TODO: wire Bridge` comment.

- [ ] **Step 9: Stub internal/telegram/bootstrap.go and bootstrap_config.go**

Remove references to memory.MemoryManager. Keep basic config bootstrap.

- [ ] **Step 10: Stub internal/cron/runtime.go**

Remove AgentExecutor interface dependency. Replace with a simple `BridgeExecutor` interface:
```go
type BridgeExecutor interface {
    Execute(ctx context.Context, prompt string, agentName string) (string, error)
}
```

- [ ] **Step 11: Fix internal/persona/canonical_service.go**

Remove `memory.MemoryManager` from struct and constructor. Remove `conversation_memory.go` reference. Remove retrieval ranking calls from `canonical_service_prompt.go`.

- [ ] **Step 12: Run `go build ./...`**

```bash
go build ./...
```

Expected: compiles with no errors. Warnings about unused imports are OK to fix iteratively.

- [ ] **Step 13: Run `go vet ./...`**

```bash
go vet ./...
```

Expected: no errors.

- [ ] **Step 14: Clean go.mod**

```bash
go mod tidy
```

Remove dependencies that are no longer imported (google generative-ai-go, anthropic-sdk-go, modelcontextprotocol/go-sdk, etc.).

- [ ] **Step 15: Commit**

```bash
git add -A
git commit -m "chore: remove replaced modules (agent, tools, llm, mcp, skill, observability, memory)

Remove ~17.5k lines of code replaced by Claude Code CLI via Bridge.
Stub remaining modules to compile."
```

---

## Task 2: TypeScript Bridge

**Goal:** Create the TS Bridge that wraps Claude Agent SDK. Go sends JSON commands via stdin, Bridge returns NDJSON events via stdout.

**Files:**
- Create: `bridge/package.json`
- Create: `bridge/tsconfig.json`
- Create: `bridge/index.ts`

**Pre-requisite:** Check Claude Agent SDK docs via Context7 before implementation.

- [ ] **Step 1: Look up Claude Agent SDK documentation**

Use Context7 MCP to fetch `@anthropic-ai/claude-agent-sdk` docs. Understand:
- How to create a session
- How to send prompts
- How to receive streaming events
- How to configure model, system prompt, permissions, MCP servers
- How to resume sessions
- Cost tracking

- [ ] **Step 2: Create bridge/package.json**

```json
{
  "name": "aurelia-bridge",
  "version": "1.0.0",
  "private": true,
  "type": "module",
  "main": "index.ts",
  "scripts": {
    "start": "npx tsx index.ts"
  },
  "dependencies": {
    "@anthropic-ai/claude-agent-sdk": "latest"
  },
  "devDependencies": {
    "tsx": "^4.0.0",
    "typescript": "^5.0.0"
  }
}
```

- [ ] **Step 3: Create bridge/tsconfig.json**

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "esModuleInterop": true,
    "outDir": "dist",
    "rootDir": "."
  },
  "include": ["index.ts"]
}
```

- [ ] **Step 4: Write the failing test for bridge/index.ts**

Create a test script that sends a JSON command to the bridge via stdin and verifies NDJSON output. This will be a shell-based integration test:

```bash
# e2e/bridge_test.sh
echo '{"command":"query","prompt":"say hello","options":{"model":"claude-sonnet-4-6","system_prompt":"You are a test assistant. Reply with exactly: hello"}}' | npx tsx bridge/index.ts
```

Expected output: NDJSON lines with `event: system`, `event: assistant`, `event: result`.

- [ ] **Step 5: Implement bridge/index.ts**

Read JSON from stdin line by line. For each command:
1. Parse the JSON request
2. Create Claude Agent SDK session with options (model, system_prompt, cwd, permission_mode, mcp_servers, allowed_tools)
3. Stream events to stdout as NDJSON:
   - `{"event":"system","session_id":"...","tools":[...]}`
   - `{"event":"tool_use","name":"...","input":{...}}`
   - `{"event":"tool_result","content":"..."}`
   - `{"event":"assistant","text":"..."}`
   - `{"event":"result","content":"...","cost_usd":0.12,"session_id":"...","duration_ms":4500,"num_turns":3}`
4. Handle errors: `{"event":"error","message":"..."}`
5. Support `resume` option to continue a previous session

Key: stderr for bridge's own logs, stdout exclusively for NDJSON events.

- [ ] **Step 6: Install dependencies**

```bash
cd bridge && npm install
```

- [ ] **Step 7: Run integration test**

```bash
bash e2e/bridge_test.sh
```

Expected: NDJSON events appear on stdout, final event is `result`.

- [ ] **Step 8: Commit**

```bash
git add bridge/
git commit -m "feat(bridge): add TypeScript Bridge wrapping Claude Agent SDK

Stdin/stdout NDJSON protocol for Go ↔ Claude Code communication."
```

---

## Task 3: Go Bridge client (internal/bridge/)

**Goal:** Go package that spawns the TS Bridge process, sends commands via stdin, parses NDJSON events from stdout.

**Files:**
- Create: `internal/bridge/protocol.go`
- Create: `internal/bridge/events.go`
- Create: `internal/bridge/bridge.go`
- Create: `internal/bridge/bridge_test.go`

- [ ] **Step 1: Write protocol.go — request/response types**

```go
package bridge

// Request sent to Bridge via stdin
type Request struct {
    Command string         `json:"command"` // "query"
    Prompt  string         `json:"prompt"`
    Options RequestOptions `json:"options"`
}

type RequestOptions struct {
    Model          string            `json:"model,omitempty"`
    Cwd            string            `json:"cwd,omitempty"`
    SystemPrompt   string            `json:"system_prompt,omitempty"`
    Resume         string            `json:"resume,omitempty"`
    MaxTurns       int               `json:"max_turns,omitempty"`
    PermissionMode string            `json:"permission_mode,omitempty"`
    MCPServers     map[string]any    `json:"mcp_servers,omitempty"`
    AllowedTools   []string          `json:"allowed_tools,omitempty"`
}
```

- [ ] **Step 2: Write events.go — event types parsed from stdout**

```go
package bridge

type Event struct {
    Type      string  `json:"event"`
    // system event
    SessionID string  `json:"session_id,omitempty"`
    Tools     []string `json:"tools,omitempty"`
    // tool_use event
    Name      string  `json:"name,omitempty"`
    Input     any     `json:"input,omitempty"`
    // tool_result / assistant / result / error
    Content   string  `json:"content,omitempty"`
    Text      string  `json:"text,omitempty"`
    Message   string  `json:"message,omitempty"`
    // result event
    CostUSD    float64 `json:"cost_usd,omitempty"`
    DurationMs int64   `json:"duration_ms,omitempty"`
    NumTurns   int     `json:"num_turns,omitempty"`
}
```

- [ ] **Step 3: Write failing test for bridge.go**

```go
func TestBridge_Execute_ParsesEvents(t *testing.T) {
    // Test that Bridge can parse a stream of NDJSON events
    // Use a mock script that outputs known NDJSON instead of real Bridge
}
```

- [ ] **Step 4: Run test to verify it fails**

```bash
go test ./internal/bridge/ -v -run TestBridge_Execute_ParsesEvents
```

Expected: FAIL

- [ ] **Step 5: Implement bridge.go**

```go
package bridge

type Bridge struct {
    bridgePath string // path to bridge/index.ts
    nodePath   string // path to npx/tsx
}

func New(bridgePath string) *Bridge

// Execute sends a request to the Bridge and returns events via channel
func (b *Bridge) Execute(ctx context.Context, req Request) (<-chan Event, error)

// ExecuteSync sends a request and waits for the final result event
func (b *Bridge) ExecuteSync(ctx context.Context, req Request) (*Event, error)
```

Implementation:
1. Spawn `npx tsx bridge/index.ts` as child process
2. Write JSON request to stdin, close stdin
3. Read stdout line by line, parse NDJSON events
4. Send events to channel
5. On `result` or `error` event, close channel
6. On context cancellation, kill process
7. Capture stderr for error diagnostics

- [ ] **Step 6: Run test to verify it passes**

```bash
go test ./internal/bridge/ -v -run TestBridge_Execute_ParsesEvents
```

Expected: PASS

- [ ] **Step 7: Write integration test with real Bridge**

```go
func TestBridge_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    // Requires ANTHROPIC_API_KEY
    b := New("../../bridge")
    result, err := b.ExecuteSync(context.Background(), Request{
        Command: "query",
        Prompt:  "reply with exactly: test_ok",
        Options: RequestOptions{
            Model:        "claude-sonnet-4-6",
            SystemPrompt: "Reply with exactly what is asked, nothing more.",
        },
    })
    require.NoError(t, err)
    require.Contains(t, result.Content, "test_ok")
}
```

- [ ] **Step 8: Run integration test**

```bash
go test ./internal/bridge/ -v -run TestBridge_Integration
```

Expected: PASS (requires API key)

- [ ] **Step 9: Commit**

```bash
git add internal/bridge/
git commit -m "feat(bridge): add Go client for TypeScript Bridge

Spawns Bridge process, sends JSON via stdin, parses NDJSON events from stdout."
```

---

## Task 4: Semantic Memory (internal/memory/)

**Goal:** sqlite-vec based semantic memory. Save content with embeddings, search by similarity, inject context into prompts.

**Files:**
- Create: `internal/memory/schema.go`
- Create: `internal/memory/embeddings.go`
- Create: `internal/memory/embeddings_test.go`
- Create: `internal/memory/store.go`
- Create: `internal/memory/store_test.go`

**Pre-requisite:** Research sqlite-vec Go bindings and embedding API options.

- [ ] **Step 1: Research sqlite-vec for Go**

Check if `github.com/asg017/sqlite-vec` has Go bindings or if we need to use CGo with the C extension. Check compatibility with `modernc.org/sqlite` (pure Go SQLite already in go.mod).

- [ ] **Step 2: Research embedding API**

Decide embedding provider. Options:
- Voyage AI (spec mentions voyage-3) — check API
- OpenAI embeddings — widely available
- Local embedding model — no external dependency but heavier

Decision should align with config schema in spec (embedding_provider, embedding_model fields).

- [ ] **Step 3: Write schema.go**

```go
package memory

const schemaSQL = `
CREATE TABLE IF NOT EXISTS memories (
    id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    category TEXT NOT NULL,
    agent TEXT DEFAULT '',
    created_at TEXT NOT NULL,
    embedding BLOB
);
CREATE INDEX IF NOT EXISTS idx_memories_category ON memories(category);
CREATE INDEX IF NOT EXISTS idx_memories_agent ON memories(agent);
`
```

Priority: use `CREATE VIRTUAL TABLE memory_vectors USING vec0(embedding float[1536], +content TEXT, +category TEXT, +agent TEXT, +created_at TEXT)` as the spec requires. If sqlite-vec doesn't work with modernc.org/sqlite (pure Go), fallback to regular table with BLOB embeddings and cosine similarity computed in Go. Research this in Step 1.

- [ ] **Step 4: Write embeddings.go — embedding generation interface**

```go
package memory

type Embedder interface {
    Embed(ctx context.Context, text string) ([]float32, error)
    Dimensions() int
}

// VoyageEmbedder calls Voyage AI API
type VoyageEmbedder struct {
    apiKey string
    model  string
    client *http.Client
}

func NewVoyageEmbedder(apiKey, model string) *VoyageEmbedder
func (e *VoyageEmbedder) Embed(ctx context.Context, text string) ([]float32, error)
func (e *VoyageEmbedder) Dimensions() int
```

- [ ] **Step 5: Write failing test for embeddings**

```go
func TestVoyageEmbedder_Embed(t *testing.T) {
    if testing.Short() { t.Skip("requires API key") }
    e := NewVoyageEmbedder(os.Getenv("VOYAGE_API_KEY"), "voyage-3")
    vec, err := e.Embed(context.Background(), "test text")
    require.NoError(t, err)
    require.Equal(t, e.Dimensions(), len(vec))
}
```

- [ ] **Step 6: Implement embeddings.go**

HTTP call to Voyage AI `/v1/embeddings` endpoint. Parse response, return float32 slice.

- [ ] **Step 7: Run embedding test**

```bash
go test ./internal/memory/ -v -run TestVoyageEmbedder -short=false
```

- [ ] **Step 8: Write failing test for store.go**

```go
func TestStore_SaveAndSearch(t *testing.T) {
    store := newTestStore(t) // in-memory sqlite
    // Use a mock embedder that returns deterministic vectors
    ctx := context.Background()
    err := store.Save(ctx, "Go is a compiled language", "fact", "")
    require.NoError(t, err)
    err = store.Save(ctx, "Python is interpreted", "fact", "")
    require.NoError(t, err)

    results, err := store.Search(ctx, "compiled programming", 1)
    require.NoError(t, err)
    require.Len(t, results, 1)
    require.Contains(t, results[0].Content, "Go")
}
```

- [ ] **Step 9: Run test to verify it fails**

```bash
go test ./internal/memory/ -v -run TestStore_SaveAndSearch
```

- [ ] **Step 10: Implement store.go**

```go
package memory

type Memory struct {
    ID        string
    Content   string
    Category  string    // "fact", "conversation", "decision", "preference"
    Agent     string
    CreatedAt time.Time
}

type Store struct {
    db       *sql.DB
    embedder Embedder
}

func NewStore(dbPath string, embedder Embedder) (*Store, error)
func (s *Store) Save(ctx context.Context, content, category, agent string) error
func (s *Store) Search(ctx context.Context, query string, limit int) ([]Memory, error)
func (s *Store) Inject(ctx context.Context, query string, limit int) (string, error)
func (s *Store) Close() error
```

`Search`: generate embedding for query, compute cosine similarity against all stored embeddings, return top N.
`Inject`: calls Search, formats results as a markdown block for system prompt injection.

- [ ] **Step 11: Run test to verify it passes**

```bash
go test ./internal/memory/ -v -run TestStore_SaveAndSearch
```

- [ ] **Step 12: Write test for Inject**

```go
func TestStore_Inject(t *testing.T) {
    store := newTestStore(t)
    store.Save(ctx, "user prefers dark mode", "preference", "")
    block, err := store.Inject(ctx, "what theme does the user like", 5)
    require.NoError(t, err)
    require.Contains(t, block, "dark mode")
}
```

- [ ] **Step 13: Run all memory tests**

```bash
go test ./internal/memory/ -v
```

- [ ] **Step 14: Commit**

```bash
git add internal/memory/
git commit -m "feat(memory): add semantic memory with sqlite-vec and embeddings

Save/Search/Inject operations with cosine similarity search."
```

---

## Task 5: Agent Registry (internal/agents/)

**Goal:** Load agent definitions from markdown files, route messages to the right agent, provide agent config to Bridge.

**Files:**
- Create: `internal/agents/types.go`
- Create: `internal/agents/registry.go`
- Create: `internal/agents/registry_test.go`

- [ ] **Step 1: Write types.go — Agent struct and frontmatter**

```go
package agents

type Agent struct {
    Name         string            `yaml:"name"`
    Description  string            `yaml:"description"`
    Model        string            `yaml:"model,omitempty"`
    Schedule     string            `yaml:"schedule,omitempty"`
    MCPServers   map[string]any    `yaml:"mcp_servers,omitempty"`
    AllowedTools []string          `yaml:"allowed_tools,omitempty"`
    Prompt       string            `yaml:"-"` // body after frontmatter
}
```

- [ ] **Step 2: Write failing test for registry loading**

```go
func TestRegistry_LoadAgents(t *testing.T) {
    dir := t.TempDir()
    // Write a test agent markdown file
    os.WriteFile(filepath.Join(dir, "prospector.md"), []byte(`---
name: prospector
description: Busca leads
model: kimi-k2-thinking
schedule: "0 9 * * 1"
allowed_tools: ["WebSearch", "WebFetch"]
---

Voce eh um agente de prospeccao.
`), 0644)

    reg, err := Load(dir)
    require.NoError(t, err)
    require.Len(t, reg.Agents(), 1)

    a := reg.Get("prospector")
    require.NotNil(t, a)
    require.Equal(t, "kimi-k2-thinking", a.Model)
    require.Equal(t, "0 9 * * 1", a.Schedule)
    require.Contains(t, a.Prompt, "prospeccao")
}
```

- [ ] **Step 3: Run test to verify it fails**

```bash
go test ./internal/agents/ -v -run TestRegistry_LoadAgents
```

- [ ] **Step 4: Implement registry.go**

```go
package agents

type Registry struct {
    agents map[string]*Agent
}

func Load(dir string) (*Registry, error)  // reads all .md files from dir
func (r *Registry) Get(name string) *Agent
func (r *Registry) Agents() []*Agent
func (r *Registry) Scheduled() []*Agent   // agents with schedule != ""
func (r *Registry) Route(message string) *Agent  // simple keyword match, fallback to default
```

`Load`: read each `.md` file, split frontmatter (between `---` markers), parse YAML frontmatter into Agent struct, body goes to Agent.Prompt.

`Route`: for now, simple name-based matching (if message starts with `@agentname`). Default agent (no specific match) returns nil (handled by caller as general query).

- [ ] **Step 5: Run test to verify it passes**

```bash
go test ./internal/agents/ -v -run TestRegistry_LoadAgents
```

- [ ] **Step 6: Write test for Route**

```go
func TestRegistry_Route(t *testing.T) {
    reg := setupTestRegistry(t) // has "prospector" agent
    a := reg.Route("@prospector busque leads em SP")
    require.NotNil(t, a)
    require.Equal(t, "prospector", a.Name)

    a = reg.Route("oi, tudo bem?")
    require.Nil(t, a) // no specific agent, nil = default
}
```

- [ ] **Step 7: Run all agent tests**

```bash
go test ./internal/agents/ -v
```

- [ ] **Step 8: Commit**

```bash
git add internal/agents/
git commit -m "feat(agents): add agent registry with markdown definitions

Load agents from ~/.aurelia/agents/, route by @name, parse frontmatter config."
```

---

## Task 6: Simplify Persona (internal/persona/)

**Goal:** Keep IDENTITY/SOUL/USER loading. Remove retrieval ranking, remove conversation memory dependency. Persona just loads files and builds a system prompt string.

**Files:**
- Modify: `internal/persona/canonical_service.go`
- Modify: `internal/persona/canonical_service_prompt.go`
- Modify: `internal/persona/canonical_service_files.go`
- Modify: `internal/persona/canonical_service_test.go`
- Delete: `internal/persona/canonical_service_archive.go` (if unused after cleanup)
- Keep: `internal/persona/loader.go`, `loader_files.go`, `loader_prompt.go`
- Keep: `internal/persona/optional_file.go`

- [ ] **Step 1: Read all persona files to understand current interfaces**

Read every file in `internal/persona/` to map what's used and what depends on removed packages.

- [ ] **Step 2: Remove memory dependency from canonical_service.go**

Remove `memory.MemoryManager` from `CanonicalIdentityService` struct and constructor. Update `NewCanonicalIdentityService` signature.

- [ ] **Step 3: Simplify canonical_service_prompt.go**

Remove calls to retrieval ranking. The prompt builder should:
1. Load IDENTITY.md content
2. Load SOUL.md content
3. Load USER.md content
4. Concatenate them with section headers
5. Return the combined system prompt string

No memory injection here — that happens at the caller level (app.go will inject semantic memory separately).

- [ ] **Step 4: Delete canonical_service_archive.go if unused**

Check if any function in this file is called after removing retrieval. If not, delete.

- [ ] **Step 5: Update tests**

Fix `canonical_service_test.go` to remove memory mocks. Tests should verify:
- Loading persona files
- Building system prompt from 3 files
- Handling missing optional files gracefully

- [ ] **Step 6: Run persona tests**

```bash
go test ./internal/persona/ -v
```

Expected: all pass.

- [ ] **Step 7: Run `go build ./...`**

Verify full project still compiles.

- [ ] **Step 8: Commit**

```bash
git add internal/persona/
git commit -m "refactor(persona): simplify loader, remove memory and retrieval ranking

Persona now loads IDENTITY/SOUL/USER and builds a prompt string. No memory dependency."
```

---

## Task 7: Adapt Cron for Bridge (internal/cron/)

**Goal:** Cron scheduler calls Bridge instead of agent loop. Scheduled agents (from registry) run automatically.

**Files:**
- Modify: `internal/cron/runtime.go`
- Modify: `internal/cron/runtime_test.go`
- Modify: `internal/cron/types.go`

- [ ] **Step 1: Read current cron files**

Understand `AgentCronRuntime`, `CronJob`, `Scheduler`, and how they wire together.

- [ ] **Step 2: Redefine Runtime interface**

```go
// Runtime executes a cron job via the Bridge
type Runtime interface {
    ExecuteJob(ctx context.Context, job CronJob) (string, error)
}
```

- [ ] **Step 3: Write failing test for new BridgeCronRuntime**

```go
func TestBridgeCronRuntime_ExecuteJob(t *testing.T) {
    mockBridge := &mockBridge{
        result: &bridge.Event{Type: "result", Content: "job done"},
    }
    rt := NewBridgeCronRuntime(mockBridge)
    result, err := rt.ExecuteJob(ctx, CronJob{
        AgentName: "prospector",
        Prompt:    "run your weekly task",
    })
    require.NoError(t, err)
    require.Equal(t, "job done", result)
}
```

- [ ] **Step 4: Implement BridgeCronRuntime**

```go
type BridgeCronRuntime struct {
    bridge *bridge.Bridge
    agents *agents.Registry
    persona *persona.CanonicalIdentityService
    memory  *memory.Store
}

func (r *BridgeCronRuntime) ExecuteJob(ctx context.Context, job CronJob) (string, error) {
    agent := r.agents.Get(job.AgentName)
    systemPrompt := r.persona.BuildPrompt() + "\n\n" + agent.Prompt
    memCtx, _ := r.memory.Inject(ctx, job.Prompt, 10)
    systemPrompt += "\n\n" + memCtx

    result, err := r.bridge.ExecuteSync(ctx, bridge.Request{
        Command: "query",
        Prompt:  job.Prompt,
        Options: bridge.RequestOptions{
            Model:        agent.Model,
            SystemPrompt: systemPrompt,
            AllowedTools: agent.AllowedTools,
            MCPServers:   agent.MCPServers,
            PermissionMode: "bypassPermissions",
        },
    })
    if err != nil { return "", fmt.Errorf("bridge execute: %w", err) }
    r.memory.Save(ctx, result.Content, "conversation", agent.Name)
    return result.Content, nil
}
```

- [ ] **Step 5: Run cron tests**

```bash
go test ./internal/cron/ -v
```

- [ ] **Step 6: Commit**

```bash
git add internal/cron/
git commit -m "feat(cron): adapt scheduler to use Bridge instead of agent loop

CronJob now executes via Bridge with agent config from registry."
```

---

## Task 8: Simplify Telegram (internal/telegram/)

**Goal:** Telegram receives messages, builds context (persona + memory + agent), sends to Bridge, streams progress back. Remove all agent/skill/observability dependencies.

**Files:**
- Modify: `internal/telegram/bot.go`
- Modify: `internal/telegram/input_pipeline.go`
- Modify: `internal/telegram/input.go`
- Modify: `internal/telegram/bootstrap.go`
- Modify: `internal/telegram/bootstrap_config.go`
- Modify: `internal/telegram/cron_handlers.go`
- Modify: `internal/telegram/cron_handlers_test.go`
- Keep: `internal/telegram/format.go`, `markdown*.go`, `output.go`, `send.go`, `messages.go`, `activity.go`
- Delete tests for removed handlers (already done in Task 1)

- [ ] **Step 1: Read remaining telegram files**

Understand what's left after Task 1 deletions. Map the flow: message in → processing → response out.

- [ ] **Step 2: Rewrite bot.go constructor**

```go
type BotController struct {
    cfg         *config.AppConfig
    bridge      *bridge.Bridge
    agents      *agents.Registry
    memory      *memory.Store
    persona     *persona.CanonicalIdentityService
    stt         stt.Transcriber
    personasDir string
    bot         *telebot.Bot
}

func NewBotController(
    cfg *config.AppConfig,
    bridge *bridge.Bridge,
    agents *agents.Registry,
    memory *memory.Store,
    persona *persona.CanonicalIdentityService,
    stt stt.Transcriber,
    personasDir string,
) (*BotController, error)
```

- [ ] **Step 3: Rewrite input_pipeline.go**

The pipeline should:
1. Extract text from message (text, photo caption, voice transcription, document)
2. Route to agent via registry (`agents.Route(text)`)
3. Build system prompt: persona + agent prompt + memory context
4. Send to Bridge via `bridge.Execute(ctx, request)` (streaming)
5. Forward events to Telegram:
   - `tool_use` → typing indicator or progress message
   - `assistant` → send text chunk
   - `result` → send final message, save to memory

- [ ] **Step 4: Write test for new pipeline flow**

```go
func TestInputPipeline_SendsToBridge(t *testing.T) {
    mockBridge := &mockBridge{...}
    mockMemory := &mockMemory{...}
    ctrl := newTestBotController(t, mockBridge, mockMemory)
    // Simulate incoming message
    result := ctrl.processMessage(ctx, "analisa esse projeto")
    require.NotEmpty(t, result)
    require.True(t, mockBridge.called)
    require.True(t, mockMemory.saveCalled)
}
```

- [ ] **Step 5: Run test**

```bash
go test ./internal/telegram/ -v -run TestInputPipeline_SendsToBridge
```

- [ ] **Step 6: Simplify bootstrap.go**

Remove conversation bootstrap complexity. Keep:
- Bot token setup
- Allowed user IDs check
- Handler registration

- [ ] **Step 7: Adapt cron_handlers.go**

Update cron handlers to work with new agent registry (create/list/delete scheduled agents).

- [ ] **Step 8: Run all telegram tests**

```bash
go test ./internal/telegram/ -v
```

- [ ] **Step 9: Commit**

```bash
git add internal/telegram/
git commit -m "refactor(telegram): simplify pipeline to use Bridge

Messages route through agent registry → Bridge → response.
Remove agent/skill/observability dependencies."
```

---

## Task 9: Refactor Config (internal/config/)

**Goal:** Update AppConfig to match the new schema from spec (providers, embedding config, simplified structure).

**Files:**
- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`
- Modify: `internal/config/mcp.go` (or delete if MCPs are per-agent only)

- [ ] **Step 1: Read current config.go**

Understand the current `AppConfig` struct and what fields are used by remaining code.

- [ ] **Step 2: Write failing test for new config schema**

```go
func TestLoadConfig_NewSchema(t *testing.T) {
    cfg := `{
        "default_provider": "anthropic",
        "default_model": "claude-sonnet-4-6",
        "providers": {
            "anthropic": {"api_key": "sk-ant-test"},
            "kimi": {"api_key": "sk-kimi-test", "base_url": "https://api.kimi.ai"}
        },
        "telegram_bot_token": "test-token",
        "telegram_allowed_user_ids": [123456],
        "embedding_provider": "voyage",
        "embedding_model": "voyage-3"
    }`
    appCfg, err := LoadFromReader(strings.NewReader(cfg))
    require.NoError(t, err)
    require.Equal(t, "anthropic", appCfg.DefaultProvider)
    require.Equal(t, "claude-sonnet-4-6", appCfg.DefaultModel)
    require.Equal(t, "sk-ant-test", appCfg.Providers["anthropic"].APIKey)
    require.Equal(t, "voyage", appCfg.EmbeddingProvider)
    require.Equal(t, "voyage-3", appCfg.EmbeddingModel)
}
```

Config also supports `embedding_api_key` (optional, falls back to default provider key):

```json
{
    "embedding_api_key": "pa-voyage-test"
}
```

- [ ] **Step 3: Run test to verify it fails**

```bash
go test ./internal/config/ -v -run TestLoadConfig_NewSchema
```

- [ ] **Step 4: Implement new AppConfig struct**

```go
type ProviderConfig struct {
    APIKey  string `json:"api_key"`
    BaseURL string `json:"base_url,omitempty"`
}

type AppConfig struct {
    DefaultProvider string                    `json:"default_provider"`
    DefaultModel    string                    `json:"default_model"`
    Providers       map[string]ProviderConfig `json:"providers"`

    TelegramBotToken      string  `json:"telegram_bot_token"`
    TelegramAllowedUserIDs []int64 `json:"telegram_allowed_user_ids"`

    EmbeddingProvider string `json:"embedding_provider"`
    EmbeddingModel    string `json:"embedding_model"`
    EmbeddingAPIKey   string `json:"embedding_api_key,omitempty"` // falls back to provider's key if empty
}
```

Keep backward-compatible loading if possible, but prioritize new schema.

- [ ] **Step 5: Run all config tests**

```bash
go test ./internal/config/ -v
```

- [ ] **Step 6: Update cmd/aurelia/ to use new config**

Wire the new config fields into app bootstrap (providers → env vars for Bridge, embedding config → memory store).

- [ ] **Step 7: Commit**

```bash
git add internal/config/ cmd/aurelia/
git commit -m "refactor(config): update schema for providers, embedding, simplified structure"
```

---

## Task 10: Wire Everything in cmd/aurelia/

**Goal:** Rewrite app.go bootstrap to wire all new components: Bridge, Memory, Agents, Persona, Cron, Telegram.

**Files:**
- Modify: `cmd/aurelia/app.go`
- Modify: `cmd/aurelia/main.go` (if needed)
- Delete: `cmd/aurelia/wiring.go` (replaced by app.go)
- Delete: `cmd/aurelia/auth_openai.go` (no longer needed — providers via config)
- Delete: `cmd/aurelia/onboard_catalog.go` (simplify onboard)
- Modify: `cmd/aurelia/onboard.go` — adapt for new config schema
- Modify: `cmd/aurelia/onboard_ui.go` — adapt for new config schema
- Modify: `cmd/aurelia/onboard_helpers.go` — adapt for new config schema

- [ ] **Step 1: Rewrite app.go bootstrap**

```go
type app struct {
    resolver  *runtime.PathResolver
    bridge    *bridge.Bridge
    memory    *memory.Store
    agents    *agents.Registry
    persona   *persona.CanonicalIdentityService
    cronStore *cron.SQLiteCronStore
    bot       *telegram.BotController
    scheduler *cron.Scheduler
}

func bootstrapApp() (*app, error) {
    // 1. Resolve paths
    // 2. Load config
    // 3. Set provider env vars (ANTHROPIC_API_KEY, ANTHROPIC_BASE_URL, etc.)
    // 4. Create Bridge (points to bridge/index.ts)
    // 5. Create Embedder (from config.EmbeddingProvider)
    // 6. Create Memory Store
    // 7. Load Agent Registry (~/.aurelia/agents/)
    // 8. Load Persona
    // 9. Create Cron Store + Scheduler with BridgeCronRuntime
    // 10. Create Telegram BotController
    // 11. Register scheduled agents in cron
    // Return app
}
```

- [ ] **Step 2: Set up provider env vars**

```go
func setProviderEnv(cfg *config.AppConfig) {
    p := cfg.Providers[cfg.DefaultProvider]
    os.Setenv("ANTHROPIC_API_KEY", p.APIKey)
    if p.BaseURL != "" {
        os.Setenv("ANTHROPIC_BASE_URL", p.BaseURL)
    }
}
```

- [ ] **Step 3: Register scheduled agents**

```go
func registerScheduledAgents(svc *cron.Service, agents *agents.Registry) error {
    for _, a := range agents.Scheduled() {
        svc.CreateOrUpdate(cron.CronJob{
            Name:      a.Name,
            Schedule:  a.Schedule,
            AgentName: a.Name,
            Prompt:    a.Prompt,
        })
    }
    return nil
}
```

- [ ] **Step 4: Adapt onboard for new config**

Update onboard TUI to collect:
- Default provider + API key
- Telegram bot token + user ID
- Embedding provider + key (optional, can default to same as main provider)

Remove provider catalog complexity (was for 9 providers, now just key/url pairs).

- [ ] **Step 5: Run `go build ./...`**

```bash
go build ./...
```

Expected: compiles successfully.

- [ ] **Step 6: Run all tests**

```bash
go test ./... -short
```

Expected: all pass (integration tests skipped with -short).

- [ ] **Step 7: Commit**

```bash
git add cmd/aurelia/ internal/
git commit -m "feat: wire all components in app bootstrap

Bridge, Memory, Agents, Persona, Cron, Telegram fully connected."
```

---

## Task 11: End-to-End Integration Test

**Goal:** Verify the full flow works: message in → persona + memory + agent → Bridge → response out.

**Files:**
- Modify: `e2e/e2e_test.go`

- [ ] **Step 1: Write E2E test**

```go
func TestE2E_FullFlow(t *testing.T) {
    if testing.Short() { t.Skip("requires API key") }

    // 1. Create temp instance dir with persona files and an agent
    // 2. Create config with test API key
    // 3. Bootstrap app
    // 4. Simulate sending a message through the pipeline
    // 5. Verify response comes back
    // 6. Verify memory was saved
    // 7. Send another message, verify memory context is injected
}
```

- [ ] **Step 2: Run E2E test**

```bash
go test ./e2e/ -v -short=false
```

- [ ] **Step 3: Fix any issues**

Iterate until E2E passes.

- [ ] **Step 4: Commit**

```bash
git add e2e/
git commit -m "test: add end-to-end integration test for full message flow"
```

---

## Task 12: Cleanup and Documentation

**Goal:** Remove stale docs, update AGENTS.md, verify final line count.

**Files:**
- Modify: `AGENTS.md`
- Modify: `README.md`
- Delete: stale docs in `docs/` that reference removed systems

- [ ] **Step 1: Update AGENTS.md for new architecture**

Document: Bridge protocol, agent markdown format, new config schema.

- [ ] **Step 2: Verify line count target**

```bash
find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1
find bridge/ -name "*.ts" | xargs wc -l | tail -1
```

Target: ~5,600 lines total (Go + TS).

- [ ] **Step 3: Run full test suite one final time**

```bash
go test ./... -v
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "docs: update documentation for Aurelia OS architecture"
```
