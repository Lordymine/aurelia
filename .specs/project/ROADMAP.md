# Roadmap

## Done

### 1. Bridge Recovery Automático ✓

**Shipped:** 2026-03-26 — commit `0484f08`
**Spec:** `.specs/features/bridge-recovery/`

Retry automático com session resume, feedback visual ("Reconectando..."), backoff com cooldown após 3 falhas consecutivas.

---

## Priority: High

---

### 2. Command Layer — Comandos Locais Antes do LLM

**Spec:** `.specs/features/command-layer/`

**Problem:** Toda mensagem passa pelo LLM classify, mesmo operações que o Go resolve sozinho. Gasta ~15s de latência e tokens desnecessariamente.

**Scope:**
- P1: CRUD de cron jobs local, reset de sessão
- P2: Status do sistema, listar agents, listar modelos/provedores
- Intercepta antes do LLM, fallback transparente

**Packages:** `internal/telegram/`, `internal/cron/`, `internal/agents/`

---

### 3. Orquestração de Agents via SDK

**Problem:** A Aurelia já passa agents pro Claude SDK, mas falta controle e visibilidade. O usuário não sabe quando um agent foi acionado, qual tá rodando, ou se travou. A delegação funciona mas é invisível.

**Scope:**
- Feedback visual no Telegram quando SDK delega pra um agent (quem tá rodando)
- Controle de profundidade máxima de delegação
- Timeout por agent delegado
- A Aurelia decide sozinha quem escalar — sem intervenção do usuário

**Packages:** `internal/bridge/`, `internal/telegram/`, `bridge/index.ts`

## Priority: Backlog

_(Vazio — features futuras serão adicionadas conforme necessidade)_
