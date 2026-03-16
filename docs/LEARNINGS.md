# LEARNINGS

## Objective

This document records recurring mistakes, operational traps, and lessons that should remain visible across future tasks.

Each entry should stay concise and practical.

Recommended entry shape:

- context
- mistake, trap, or deviation
- impact
- decision or lesson
- prevention for future work

## Entries

### 2026-03-16 - Canonical Documentation For Aurelia

- context: start of the repository transition from the legacy project identity to `Aurelia`
- mistake, trap, or deviation: architectural guidance, coding rules, and project memory were spread across multiple overlapping files
- impact: the repository had no single concise operating baseline for future work, making migration planning and contribution governance harder
- decision or lesson: the canonical documentation set for ongoing work is now `AGENTS.md`, `docs/ARCHITECTURE.md`, `docs/STYLE_GUIDE.md`, and `docs/LEARNINGS.md`
- prevention for future work: new rules and decisions should be added to the canonical document that owns that concern instead of creating overlapping guidance elsewhere

### 2026-03-16 - Secrets And Local Runtime Artifacts Must Never Sit In Repo Root

- context: sanitization inventory for the transition to `Aurelia`
- mistake, trap, or deviation: local config, debug output, and database artifacts were present in the repository root, including a real secret in local MCP configuration
- impact: publication readiness was blocked and any tracked secret must be treated as compromised
- decision or lesson: local config files, debug traces, and runtime databases must be ignored by default and replaced with public examples when needed
- prevention for future work: only example config files may be committed; any real secret exposed in the repository must trigger rotation, removal from the working tree, and historical cleanup
