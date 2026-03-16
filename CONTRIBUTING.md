# CONTRIBUTING

## Objective

This repository accepts contributions only when they preserve the architectural, security, and validation rules of `Aurelia`.

Read before opening a PR:

- [AGENTS.md](AGENTS.md)
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- [docs/STYLE_GUIDE.md](docs/STYLE_GUIDE.md)
- [docs/LEARNINGS.md](docs/LEARNINGS.md)

## Contribution Rules

Every non-trivial change should follow this sequence:

1. understand the current architecture and scope
2. make the smallest defensible change
3. run relevant validation
4. update canonical docs if behavior or policy changed
5. open a PR with a clear explanation of scope and risk

## What Can Enter The Repository

Allowed:

- source code
- tests
- mocks without secrets
- canonical docs
- benchmark harness and sanitized results
- GitHub workflow and governance files

## What Must Not Enter The Repository

Forbidden:

- secrets
- `.env`
- local MCP config with real keys
- local runtime state
- local databases
- debug dumps
- personal paths
- machine-specific artifacts
- generated binaries

## Validation Minimum

Before opening a PR, run:

```bash
go test ./...
```

If the change affects performance claims, docs, or release positioning:

- update [docs/BENCHMARKS.md](docs/BENCHMARKS.md) when relevant
- keep the README aligned with measured data

## Documentation Rule

If the change introduces a new architectural rule, coding rule, or recurring operational lesson, update the canonical document that owns it:

- architecture -> [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- implementation rule -> [docs/STYLE_GUIDE.md](docs/STYLE_GUIDE.md)
- recurring lesson -> [docs/LEARNINGS.md](docs/LEARNINGS.md)

## Pull Request Expectations

A good PR should:

- explain what changed
- explain why the change was needed
- describe risk or behavioral impact
- list validation that was run
- mention any documentation updates

## Review Policy

Changes touching core architecture, runtime state, memory, agent orchestration, secrets, or CI should be reviewed carefully before merge.
