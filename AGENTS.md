# AGENTS.md

This file defines how coding agents should operate within this repository.

The goal is to keep execution disciplined, documentation current, and changes easy to review.

---

# Workflow

Follow this workflow for any non-trivial task:

1. **Plan** — Draft a structured plan before touching code
2. **Review** — Present the plan and wait for human approval
3. **Execute** — Implement autonomously after approval
4. **Validate** — Run the relevant tests and checks
5. **Handoff** — Report completion and leave the commit to the human

For trivial tasks, implement directly and validate briefly.

Replan if the current approach becomes messy, uncertain, or incorrect.

---

# Autonomy Rules

After plan approval, execute autonomously.

Stop and ask the human when:

- a decision has architectural impact not covered in the approved plan
- requirements are ambiguous and multiple interpretations could break behavior
- the only valid fix would violate a project rule documented here or in `docs/`
- multiple valid approaches exist with significant tradeoffs

Do not ask for permission on obvious implementation details.

Do not over-communicate progress. Surface blockers, tradeoffs, and decisions that matter.

Always leave the final commit to the human.

---

# Subagents

Use subagents only when they help isolate work or protect the main context.

Typical use cases:

- repository exploration
- parallel research
- isolated debugging
- log analysis

Avoid unnecessary subagent usage.

---

# Tests And Validation

Never mark a task complete without running the relevant validation.

Validation loop:

1. run the relevant test suite or verification command
2. if it passes, report completion
3. if it fails, investigate the root cause
4. fix the root cause without violating repository rules
5. rerun validation until green
6. if the fix requires breaking a documented rule, stop and ask the human

Do not patch symptoms.
Fix root causes.

Do not skip validation because the change "should work".

---

# Engineering Principles

- simplicity first
- consistency over cleverness
- reliable code over fast guesses
- explicit rules over hidden behavior
- avoid unnecessary abstractions
- leave the codebase in a better state than you found it

---

# Bug Fixing

Investigate bugs independently when evidence is clear.

Use logs, failing tests, stack traces, and code paths to identify the root cause.

Fix autonomously when the scope is well understood.

If multiple risky interpretations exist, stop and ask for clarification.

---

# Learning Loop

When the human corrects an error:

- identify the pattern behind the mistake
- record the lesson in `docs/LEARNINGS.md` when it is likely to matter again

Record recurring mistakes, process failures, architectural misunderstandings, and operational traps.

---

# Project Docs

Keep project documentation updated as part of execution.

Canonical documents:

- workflow and execution rules: `AGENTS.md`
- architecture and technical boundaries: `docs/ARCHITECTURE.md`
- implementation patterns and coding rules: `docs/STYLE_GUIDE.md`
- operational lessons and recurring mistakes: `docs/LEARNINGS.md`

When a new pattern or decision is established:

- document it in the appropriate file before closing the task
- keep entries concise: context, decision, rationale

Do not leave decisions undocumented.

---

# Documentation Status

The documents above are the primary source of truth for ongoing work.

Older files in `docs/` may continue to exist during the transition to `Aurelia`, but they should be treated as secondary references until they are migrated, merged, or retired.

---

# Core Principles

Simplicity first.

Consistency over cleverness.

Reliable code over fast guesses.

Always leave the codebase in a better state than you found it.
