# Aurelia OS — Spec Definitiva

**Data:** 2026-03-21
**Branch:** feat/aurelia-os
**Worktree:** RafaClaw-aurelia-os (a partir do main original)
**Conceito:** Sistema operacional orquestrador de agentes autonomos

---

## 1. Visao

Aurelia eh um OS de agentes autonomos acessivel via Telegram. O usuario conversa naturalmente — a Aurelia decide se responde direto, delega pra um agente, ou agenda execucao automatica. Agentes executam via Claude Code CLI (invisivel pro usuario). Qualquer provider/modelo compativel com tool calling.

O usuario final nao precisa ser programador.

---

## 2. Arquitetura

```
Go (processo principal — leve, 24/7)
├── Telegram I/O
├── Persona (IDENTITY.md, SOUL.md, USER.md)
├── Memoria Semantica (sqlite-vec + embeddings em Go)
├── Cron Scheduler
├── Config / Onboard TUI
├── Agent Registry (agentes configuraveis em markdown)
├── Provider Manager (env vars pra qualquer provider/modelo)
└── Bridge TypeScript (~300 linhas)
    └── @anthropic-ai/claude-agent-sdk (oficial)
        └── Claude Code CLI (embutido no SDK)
            ├── Herda MCPs da conta + MCPs configurados pela Aurelia
            ├── Herda skills da conta + skills configuradas pela Aurelia
            ├── Tools nativas (Read, Write, Edit, Bash, Glob, Grep, etc.)
            └── Subagents, worktrees, plugins
```

### Fluxo principal

```
1. Mensagem chega no Telegram
2. Aurelia (Go) busca contexto:
   - Persona (IDENTITY + SOUL + USER)
   - Memoria semantica (busca por similaridade)
   - Agent registry (qual agente responde?)
3. Monta system prompt = persona + memoria + instrucoes do agente
4. Envia pro Bridge TypeScript → SDK → Claude Code CLI
5. Bridge retorna eventos em stream (tool use, assistant text, result)
6. Aurelia envia progresso e resultado pro Telegram
7. Salva na memoria semantica
```

### Fluxo cron (autonomo)

```
1. Cron dispara no horario configurado
2. Aurelia carrega agente + system prompt
3. Envia pro Bridge → SDK → Claude Code
4. Resultado vai pro Telegram automaticamente
5. Salva na memoria
```

---

## 3. Componentes

### 3.1 Bridge TypeScript (~300 linhas)

Arquivo unico que wrappa o SDK oficial. Comunicacao com Go via stdin/stdout JSON.

**Protocolo Go → Bridge (stdin):**
```json
{
  "command": "query",
  "prompt": "analisa o projeto",
  "options": {
    "model": "kimi-k2-thinking",
    "cwd": "/path/to/project",
    "system_prompt": "Voce eh Lex, assistente pessoal...",
    "resume": "session-id-anterior",
    "max_turns": 50,
    "permission_mode": "bypassPermissions",
    "mcp_servers": {"github": {"command": "..."}},
    "allowed_tools": ["Read", "Write", "Edit", "Bash"]
  }
}
```

**Protocolo Bridge → Go (stdout, NDJSON):**
```json
{"event": "system", "session_id": "abc-123", "tools": ["Read","Write",...]}
{"event": "tool_use", "name": "Read", "input": {"file_path": "src/main.go"}}
{"event": "tool_result", "content": "package main..."}
{"event": "assistant", "text": "O projeto tem..."}
{"event": "result", "content": "...", "cost_usd": 0.12, "session_id": "abc-123", "duration_ms": 4500, "num_turns": 3}
```

**Features do SDK herdadas automaticamente:**
- Todas as tools nativas do Claude Code
- MCPs (da conta + configurados)
- Skills/Superpowers
- Subagents com worktrees
- Session resume
- Cost tracking
- File checkpointing
- Permission control

### 3.2 Persona (Go)

3 arquivos markdown em `~/.aurelia/personas/`:

- `IDENTITY.md` — nome, papel, regras de comportamento
- `SOUL.md` — personalidade, tom, estilo
- `USER.md` — informacoes sobre o usuario

Carregados e injetados no system_prompt antes de cada chamada ao Bridge.

Manter o loader existente. Simplificar: remover retrieval ranking, remover playbook/lessons (migrar pra memoria semantica).

### 3.3 Memoria Semantica (Go)

SQLite com extensao sqlite-vec. Embeddings gerados por modelo (ex: via API de embeddings do provider ativo).

**Tabela:**
```sql
CREATE VIRTUAL TABLE memory_vectors USING vec0(
  id TEXT PRIMARY KEY,
  embedding FLOAT[1536],
  content TEXT,
  category TEXT,       -- "fact", "conversation", "decision", "preference"
  agent TEXT,          -- qual agente criou
  created_at TEXT
);
```

**Operacoes:**
- `Save(content, category)` — gera embedding, salva
- `Search(query, limit)` — busca por similaridade, retorna top N
- `Inject(query)` — busca relevantes e monta bloco de contexto pro system prompt

**Substituicao do sistema atual:**
- facts → category="fact" na memoria vetorial
- notes → category="conversation"
- archive → nao precisa (memoria vetorial eh persistente)
- context_policy → nao precisa (busca por similaridade ja seleciona o relevante)

### 3.4 Agent Registry (Go)

Agentes definidos em markdown em `~/.aurelia/agents/`:

```markdown
---
name: prospector
description: Busca leads e entra em contato
model: kimi-k2-thinking
schedule: "0 9 * * 1"
mcp_servers:
  google-places: { command: "npx google-places-mcp" }
allowed_tools: ["WebSearch", "WebFetch", "Bash"]
---

Voce eh um agente de prospeccao comercial.
Busque empresas no Google Places na regiao configurada.
Colete nome, telefone, email.
Entre em contato via WhatsApp ou email pra agendar reuniao.
```

**Funcionalidades:**
- Carregar agentes do diretorio
- Rotear mensagens pro agente certo (por nome ou classificacao LLM)
- Passar config do agente (model, mcp_servers, allowed_tools) pro Bridge
- Cron registra agentes com schedule automaticamente

### 3.5 Provider Manager (Go)

Reutiliza `env.go` existente. Monta env vars baseado no provider/modelo:

- Anthropic: `ANTHROPIC_API_KEY`
- Outros: `ANTHROPIC_BASE_URL` + `ANTHROPIC_AUTH_TOKEN`
- Modelo via `--model` flag

Cada agente pode ter seu proprio provider/modelo.

### 3.6 Cron Scheduler (Go)

Reutiliza `internal/cron/` existente. Adaptar runtime pra chamar o Bridge ao inves do ReAct loop.

### 3.7 Telegram I/O (Go)

Simplificar o existente. Manter:
- Text, photo, voice, document handling
- Album batching
- Markdown rendering
- Message chunking
- Typing indicators
- Progresso em tempo real (tool events do Bridge)

Remover:
- /ops command (substituido por observabilidade via memoria)
- /memory command (substituido por memoria semantica)
- Bootstrap flow complexo (simplificar)

### 3.8 Config / Onboard TUI (Go)

Manter TUI existente. Adaptar pra novo config:

```json
{
  "default_provider": "anthropic",
  "default_model": "claude-sonnet-4-6",
  "providers": {
    "anthropic": { "api_key": "sk-ant-..." },
    "kimi": { "api_key": "sk-kimi-...", "base_url": "..." },
    "groq": { "api_key": "gsk-..." }
  },
  "telegram_bot_token": "...",
  "telegram_allowed_user_ids": [123456],
  "embedding_provider": "anthropic",
  "embedding_model": "voyage-3"
}
```

---

## 4. O que reutilizar do codigo atual (main)

| Modulo atual | Acao | Destino |
|---|---|---|
| `internal/persona/` | Simplificar — manter loader de IDENTITY/SOUL/USER, remover retrieval ranking | `internal/persona/` |
| `internal/cron/` | Adaptar runtime pra chamar Bridge | `internal/cron/` |
| `internal/config/` | Refatorar pra novo schema | `internal/config/` |
| `internal/runtime/` | Manter (path resolver, bootstrap) | `internal/runtime/` |
| `internal/telegram/` | Simplificar — manter I/O, remover ops/memory commands | `internal/telegram/` |
| `pkg/stt/` | Manter (Groq STT) | `pkg/stt/` |
| `cmd/aurelia/onboard*.go` | Adaptar pra novo config | `cmd/aurelia/` |
| `cmd/aurelia/main.go` | Manter | `cmd/aurelia/` |

### O que criar do zero

| Modulo | Descricao |
|---|---|
| `internal/bridge/` | Comunicacao Go ↔ TypeScript Bridge (spawn, stdin/stdout, parse eventos) |
| `internal/memory/` | Reescrever — sqlite-vec, embeddings, busca semantica |
| `internal/agents/` | Agent registry — carrega markdowns, roteia, configura |
| `bridge/index.ts` | Bridge TypeScript — wrappa SDK oficial |
| `bridge/package.json` | Dependencia do @anthropic-ai/claude-agent-sdk |

### O que remover

| Modulo | Razao |
|---|---|
| `internal/agent/` (inteiro) | ReAct loop, tool profiles, compaction — Claude Code faz |
| `internal/tools/` (inteiro) | Tool system proprio — Claude Code faz |
| `internal/observability/` | Substituido por memoria semantica + eventos do Bridge |
| `internal/mcp/` | Ja removido no orchestrator. MCPs no Claude Code |
| `internal/skill/` | Ja removido. Skills no Claude Code |
| `pkg/llm/` (inteiro) | 9 providers — substituido por env vars + Bridge |
| `internal/claude/` (inteiro) | Substituido por Bridge TypeScript |
| `internal/memory/` (atual) | Substituido por memoria semantica |

---

## 5. Estimativa

| Componente | Linhas estimadas |
|---|---|
| Bridge TypeScript | ~300 |
| internal/bridge/ (Go) | ~200 |
| internal/memory/ (semantica) | ~400 |
| internal/agents/ | ~200 |
| Persona simplificada | ~400 (vs 1143 atual) |
| Telegram simplificado | ~1500 (vs 2331 atual) |
| Cron adaptado | ~700 (reutiliza) |
| Config refatorado | ~300 |
| Onboard adaptado | ~1000 (reutiliza) |
| Runtime | ~120 (reutiliza) |
| STT | ~115 (reutiliza) |
| cmd/aurelia/ | ~400 |
| **Total Go** | **~5.300** |
| **Total TS** | **~300** |
| **Total** | **~5.600** |

vs 13.070 linhas atuais = **57% reducao**

---

## 6. Ordem de implementacao

1. Limpeza: remover modulos que saem (agent, tools, llm, claude, observability, memory atual, mcp, skill)
2. Bridge: criar bridge TS + internal/bridge/ Go
3. Memoria: implementar sqlite-vec + embeddings
4. Agents: agent registry com markdown
5. Persona: simplificar loader
6. Cron: adaptar runtime pra Bridge
7. Telegram: simplificar pipeline
8. Config: refatorar pra novo schema
9. Testes end-to-end
10. Validacao com usuario

---

## 7. Riscos

1. **sqlite-vec em Go** — precisa verificar maturidade da extensao pra Go
2. **API de embeddings** — precisa de um endpoint de embeddings (Anthropic nao tem proprio, usar OpenAI ou Voyage)
3. **Bridge TypeScript** — protocolo de comunicacao precisa ser robusto (erros, timeouts, crash recovery)
4. **Providers nao-Anthropic** — Claude Code CLI com ANTHROPIC_BASE_URL pode ter quirks com alguns providers
5. **Tamanho do SDK TS** — 12.4 MB do CLI embutido. Aceitavel pra servidor, ruim pra edge.

---

## 8. Decisoes consolidadas

| Decisao | Escolha | Razao |
|---|---|---|
| Linguagem principal | Go | Binario leve, baixo consumo, daemon 24/7 |
| Executor | Claude Code via SDK oficial (TS Bridge) | 100% features, Anthropic mantém |
| LLM proprio | Nao — Claude Code eh o cerebro | Elimina toda complexidade de orquestracao |
| Multi-provider | Sim — via env vars | Liberdade de escolha (Kimi, GLM, Minimax, etc.) |
| Memoria | Semantica com sqlite-vec | Busca por similaridade, melhor que sequential |
| Persona | IDENTITY + SOUL + USER (3 arquivos) | Simples, efetivo, diferencial da Aurelia |
| Agentes | Markdown configuravel | Nao-programadores podem criar agentes |
| Skills/MCPs | Herda da conta Claude + configuraveis por agente | Sem reimplementar |
| Interface | Telegram | Acessivel de qualquer lugar |
| Autonomia | Cron scheduler | Agentes rodam sozinhos |
