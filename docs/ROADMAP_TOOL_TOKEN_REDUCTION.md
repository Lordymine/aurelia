# Tool Token Reduction Roadmap

## Objective

Track the incremental reduction of token cost caused by tool exposure and tool outputs in `Aurelia` while preserving the current local-first architecture and avoiding overengineering.

This roadmap is the human-readable plan.

Machine-readable execution state lives in `docs/STATE_TOOL_TOKEN_REDUCTION.json`.

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

### T00 - Baseline Measurement

- scope: measure the current tool surface and establish a local baseline for tool count, serialized schema size, and large tool outputs
- depends on: none
- official references:
  - https://platform.openai.com/docs/guides/prompt-caching/prompt-caching%3F.pdf
  - https://docs.anthropic.com/id/docs/build-with-claude/prompt-caching
- acceptance criteria:
  - the runtime can report how many tools are exposed in a normal chat execution
  - the runtime can estimate the serialized size of the tool block before provider calls
  - the runtime can identify oversized tool outputs in local observability
- validation:
  - unit tests for measurement helpers
  - `go test ./...`

### T01 - Canonical Tool Profiles

- scope: define small canonical tool profiles so the runtime stops defaulting to the full tool surface for normal executions
- depends on: `T00`
- official references:
  - https://github.com/pro-vi/mcp-filter
- acceptance criteria:
  - canonical profiles exist in a domain/runtime layer instead of the Telegram boundary
  - profiles cover the common execution modes such as plain chat, local files, local execution, web research, scheduling, team operations, and MCP access
  - profiles are explicit and testable
- validation:
  - unit tests for profile definitions
  - `go test ./...`

### T02 - Main Chat Tool Selection

- scope: apply profile-based tool selection in the main chat path so normal conversations do not receive every available tool
- depends on: `T01`
- official references:
  - https://platform.openai.com/docs/guides/function-calling
  - https://docs.anthropic.com/en/docs/agents-and-tools/tool-use/implement-tool-use
- acceptance criteria:
  - the main chat execution chooses a minimal allowed-tool set before entering the loop
  - normal chat can run with a small profile instead of full registry exposure
  - filesystem, command execution, web, scheduling, and team intents can still access their required tools
- validation:
  - unit tests for intent-to-profile selection
  - integration tests for the Telegram execution pipeline
  - `go test ./...`

### T03 - Cron And Team Alignment

- scope: align cron and team execution with the same minimal-tool discipline already expected in the main chat path
- depends on: `T02`
- official references:
  - https://docs.anthropic.com/en/docs/agents-and-tools/tool-use/implement-tool-use
- acceptance criteria:
  - cron executions no longer default to an unnecessarily large tool surface
  - team tasks preserve explicit allowlists without falling back to broad exposure unless intentionally configured
  - execution behavior remains consistent across chat, cron, and team runtimes
- validation:
  - runtime tests for cron allowed-tools behavior
  - team execution tests
  - `go test ./...`

### T04 - Safe Schema Pruning

- scope: reduce the serialized weight of tool schemas without changing the actual argument contract
- depends on: `T02`
- official references:
  - https://github.com/pro-vi/mcp-filter
- acceptance criteria:
  - tool serialization keeps only contract-relevant schema fields such as `type`, `properties`, and `required`
  - redundant metadata like `title`, `examples`, and defaults can be removed safely
  - pruning is applied centrally and covered by tests
- validation:
  - schema pruning unit tests
  - provider request-building tests
  - `go test ./...`

### T05 - Short Tool Descriptions

- scope: shorten verbose tool descriptions, especially for MCP-discovered tools, so tools stay understandable without bloating the prompt
- depends on: `T04`
- official references:
  - https://github.com/pro-vi/mcp-filter
- acceptance criteria:
  - MCP tool descriptions are reduced to short functional summaries
  - local tool descriptions stay concise and contract-oriented
  - no runtime path reintroduces long repeated provider or server boilerplate into descriptions
- validation:
  - discovery and description unit tests
  - `go test ./...`

### T06 - MCP Surface Reduction

- scope: ensure MCP tools are only exposed when a profile actually needs them instead of contaminating every execution
- depends on: `T01`, `T05`
- official references:
  - https://github.com/pro-vi/mcp-filter
- acceptance criteria:
  - MCP tools are grouped behind explicit profile or allowlist decisions
  - a normal chat execution does not automatically receive the full discovered MCP surface
  - configured MCP allowlists remain respected
- validation:
  - MCP exposure tests
  - execution-path tests
  - `go test ./...`

### T07 - Tool Output Compaction

- scope: keep oversized tool outputs out of the active conversation history while preserving local traceability
- depends on: `T00`
- official references:
  - https://docs.anthropic.com/en/docs/claude-code/mcp
- acceptance criteria:
  - the runtime detects oversized tool results before appending them to history
  - raw output can remain available locally while a compact summary is injected into the conversation
  - the history stays deterministic and debuggable
- validation:
  - unit tests for compaction thresholds and summaries
  - loop tests covering oversized tool outputs
  - `go test ./...`

### T08 - Type-Specific Output Policies

- scope: apply simple deterministic compaction policies for common heavy tools like `run_command`, `read_file`, and MCP calls
- depends on: `T07`
- official references:
  - https://docs.anthropic.com/en/docs/claude-code/mcp
- acceptance criteria:
  - command output can be truncated or summarized deterministically
  - file reads can preserve the useful head/body shape without dumping the full file into history by default
  - MCP outputs follow the same compacting rules instead of bypassing them
- validation:
  - tool-specific compaction tests
  - `go test ./...`

### T09 - Lightweight Observability

- scope: expose local evidence that the reduction work is helping without introducing external APM or heavy metrics stacks
- depends on: `T00`, `T02`, `T07`
- official references:
  - https://platform.openai.com/docs/guides/prompt-caching/prompt-caching%3F.pdf
  - https://docs.anthropic.com/id/docs/build-with-claude/prompt-caching
- acceptance criteria:
  - local observability records the tool count for an execution
  - local observability records the approximate serialized size of the exposed tool block
  - local observability records raw-versus-compacted output size for oversized tool results
- validation:
  - observability store tests
  - command or handler tests for inspection output when relevant
  - `go test ./...`

### T10 - Provider-Friendly Prompt Stability

- scope: keep tool blocks stable and cache-friendly where possible so provider-side prompt caching can help later without adding provider-specific complexity now
- depends on: `T02`, `T04`, `T05`
- official references:
  - https://platform.openai.com/docs/guides/prompt-caching/prompt-caching%3F.pdf
  - https://docs.anthropic.com/id/docs/build-with-claude/prompt-caching
- acceptance criteria:
  - tool ordering and serialization are stable for the same profile
  - the implementation does not depend on provider-specific caching APIs to be correct
  - documentation states that prompt stability is intentional but secondary to correctness
- validation:
  - deterministic serialization tests
  - `go test ./...`

### T11 - Documentation Finalization

- scope: align the canonical docs after the token-reduction work is in place
- depends on: `T01`, `T02`, `T04`, `T06`, `T07`, `T08`, `T09`, `T10`
- official references:
  - task references used in the completed tasks
- acceptance criteria:
  - canonical docs explain that tools are exposed by minimum necessary profile
  - docs explain that oversized tool output should not pollute active history
  - recurring implementation traps are captured in `docs/LEARNINGS.md` when relevant
- validation:
  - manual doc review
  - `go test ./...`
