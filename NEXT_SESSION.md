# Prompt para próxima sessão

Cole isso no início da conversa:

---

Estou trabalhando no projeto Aurelia OS. Leia sua memória e a spec em:

C:\Users\kocar\Documents\RafaClaw-aurelia-os\docs\superpowers\specs\2026-03-21-aurelia-os-spec.md

Contexto: Aurelia é um OS de agentes autônomos em Go que usa Claude Code como cérebro via Bridge TypeScript (SDK oficial). Worktree: RafaClaw-aurelia-os, branch: feat/aurelia-os.

O main original tem o código base pra reutilizar (persona, cron, telegram, config, TUI, STT). Precisa limpar o que não serve (agent loop, tools, pkg/llm, observability) e criar o que falta (bridge TS, memória semântica com sqlite-vec, agent registry).

Implemente a spec seguindo a ordem: limpeza → bridge → memória → agents → persona → cron → telegram → config → testes.
