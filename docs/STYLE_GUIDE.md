# STYLE GUIDE

## Objective

This guide defines the implementation standards for `Aurelia`.

The goal is to keep the codebase simple, explicit, testable, and easy to evolve without architectural drift.

## Principles

- simplicity first
- consistency over cleverness
- explicit behavior over implicit magic
- small focused components
- fix root causes, not symptoms
- prefer real validation over assumption

## Repository Organization

Expected high-level structure:

- `cmd/` for entrypoints and wiring
- `internal/` for product modules
- `pkg/` for reusable package-level integrations
- `docs/` for canonical project documentation

Rules:

- entrypoints stay thin
- domain logic does not belong in wiring
- tool schemas and handlers stay near the tool layer
- documentation should reflect the code that actually exists

## Go Conventions

### Errors

Rules:

- do not use `panic` outside startup boundaries unless there is no valid recovery path
- return explicit errors from internal functions
- wrap errors with enough context to diagnose the failing operation
- do not swallow errors silently

### Context

Rules:

- pass `context.Context` as the first argument in long-running or external operations
- propagate cancellation and deadlines when calling providers, MCP, web, command execution, or storage layers
- do not create background contexts deep inside business logic unless there is a documented reason

### Dependency Injection

Rules:

- inject dependencies through constructors
- do not instantiate network clients or heavy dependencies deep inside handlers
- keep interfaces small and behavior-oriented
- avoid interface extraction with only one speculative implementation
- when provider or model metadata is needed, consume the canonical helpers in `pkg/llm` instead of re-declaring local switch tables

### Functions And Files

Rules:

- prefer small focused functions
- split files when they become hard to scan or hold multiple responsibilities
- keep modules responsibility-focused
- extract helpers only when they improve readability or testability

### Naming

Rules:

- use descriptive names
- avoid abbreviations without clear payoff
- names should describe behavior, not implementation accident
- product-facing naming must use `Aurelia`, not legacy names

### Comments

Rules:

- comment only when intent is not obvious from the code
- avoid redundant comments
- prefer concise comments ahead of non-trivial logic

## Testing

### Validation Rule

No task is complete without relevant validation.

### Test Strategy

Use the simplest validation that proves the change:

- unit tests for small deterministic rules
- integration tests for persistence, orchestration, and provider boundaries
- end-to-end validation for critical user flows when relevant

### TDD Guidance

For domain rules and regressions:

- prefer test-first
- add a regression test when fixing a real bug

For infrastructure or wiring:

- minimal implementation followed immediately by real validation is acceptable

### Benchmarking

Performance claims in docs must come from measured data.

When claiming that `Aurelia` is lightweight, measure:

- binary size
- startup time
- idle memory
- idle CPU
- focused benchmark results when relevant

Do not publish guessed numbers.

## Tooling Rules

### Local Execution

Rules:

- prefer real local execution when the task asks for running, testing, or observing behavior
- do not claim environment limitations without observing an actual failure
- preserve canonical `workdir` when acting on a target project

### Tools And Runtime Capabilities

Rules:

- treat the live runtime tool list as the source of truth
- do not assume hidden capabilities
- document new runtime capabilities when they become part of the supported architecture

## Documentation Rules

Rules:

- update canonical docs as part of the task when behavior or policy changes
- architecture decisions go to `docs/ARCHITECTURE.md`
- coding and implementation patterns go to `docs/STYLE_GUIDE.md`
- recurring operational traps and mistakes go to `docs/LEARNINGS.md`

Do not leave important decisions undocumented.

## Security And Hygiene

Rules:

- do not commit secrets
- do not commit local databases, logs, or runtime memory artifacts
- do not log secrets in plain text
- prefer example config files over real local config
- sanitize benchmark and documentation outputs before publishing

## Decisions That Should Not Regress

For this project:

- do not solve continuity problems by inflating prompt size alone
- do not replace deterministic memory with vector-first behavior as a default
- do not move domain logic into Telegram handlers
- do not rebuild tool selection logic in interface layers; use canonical runtime helpers and expose the minimum necessary tool surface
- do not expose the full MCP surface to ordinary chat; when MCP is needed, resolve only the explicitly requested registered server or tool subset
- do not hardcode provider- or server-specific MCP routing in interface code when the live registry can infer the relevant server and preserve short conversation continuity
- do not send every hidden tool schema to the model when a dynamic capability catalog plus explicit runtime expansion can preserve behavior with lower token cost
- do not append oversized raw tool output to active history when a compact deterministic summary plus local artifact is enough
- do not bypass the team model by letting workers speak directly to the user
- do not introduce architectural sprawl without a documented reason
