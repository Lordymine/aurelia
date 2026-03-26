# Bridge Recovery Automático — Design

**Spec**: `.specs/features/bridge-recovery/spec.md`
**Status**: Done

---

## Architecture Overview

O retry vive no `executeAsync` (Telegram layer), não no Bridge. O Bridge é responsável por detectar a morte e notificar; o Telegram layer decide o que fazer (retry, feedback, abort).

```
executeAsync()
  │
  ├─ bridge.Execute(req) ──► ch
  │
  ├─ processBridgeEventsAsync(ch)
  │     │
  │     ├─ terminal "result" → resposta normal, return success
  │     ├─ terminal "error"  → erro do LLM, return llmError
  │     └─ channel closed    → process morreu, return processDeath
  │
  ├─ if processDeath:
  │     ├─ sessions.DeactivateAll()   ← já feito pelo OnDeath callback
  │     ├─ rebuild request com Resume (não Continue)
  │     ├─ bridge.Execute(req)        ← auto-restart via startLocked()
  │     └─ processBridgeEventsAsync(ch) ← segunda tentativa
  │
  └─ if retry also fails → SendError("Processador reiniciado mas não conseguiu completar.")
```

---

## Code Reuse Analysis

### Existing Components to Leverage

| Component | Location | How to Use |
|-----------|----------|------------|
| `Bridge.Execute()` | `internal/bridge/bridge.go:225` | Já faz auto-restart via `startLocked()` — retry é só chamar de novo |
| `session.Store` | `internal/session/store.go` | Estender com `DeactivateAll()` |
| `readLoop` cleanup | `bridge.go:164-175` | Hook point pra notificar morte do processo |
| `buildBridgeRequest` | `input_pipeline.go:102` | Reusar pra rebuild da request com Resume |
| `processBridgeEventsAsync` | `input_pipeline.go:188` | Modificar return type pra indicar causa do fim |

### Integration Points

| System | Integration Method |
|--------|--------------------|
| Bridge → Session Store | Callback `OnDeath` chamado do `readLoop` ao detectar saída do processo |
| Telegram → Bridge | `executeAsync` encapsula a lógica de retry, transparente pro resto |

---

## Components

### 1. Bridge Death Notification

- **Purpose**: Notificar listeners quando o processo bridge morre
- **Location**: `internal/bridge/bridge.go`
- **Changes**:
  - Novo campo `onDeath func()` no struct `Bridge`
  - Novo método `SetOnDeath(fn func())`
  - `readLoop()` chama `b.onDeath()` quando `scanner.Scan()` retorna false (mas não quando `stopping == true` — shutdown intencional não é crash)
- **Reuses**: Estrutura existente do `readLoop`, padrão de callback já usado em `cron.NotifyingRuntime`

```go
// bridge.go — novo campo
type Bridge struct {
    // ... existing fields ...
    onDeath func() // called when process exits unexpectedly
}

// bridge.go — setter
func (b *Bridge) SetOnDeath(fn func()) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.onDeath = fn
}

// readLoop — adicionar antes do cleanup de pending channels
func (b *Bridge) readLoop() {
    defer close(b.done)

    for b.scanner.Scan() {
        // ... existing event routing ...
    }

    // Process died — notify listener (only if not intentional shutdown)
    b.mu.Lock()
    stopping := b.stopping
    cb := b.onDeath
    b.mu.Unlock()
    if !stopping && cb != nil {
        cb()
    }

    // ... existing cleanup (close pending channels, reset state) ...
}
```

### 2. Session Store — DeactivateAll

- **Purpose**: Marcar todas as sessions como cold após morte do bridge
- **Location**: `internal/session/store.go`
- **Interface**: `DeactivateAll()` — itera sessions e seta `active = false`
- **Dependencies**: Nenhuma nova
- **Reuses**: Estrutura existente de `entry{sessionID, active}`

```go
// store.go — novo método
func (s *Store) DeactivateAll() {
    s.mu.Lock()
    defer s.mu.Unlock()
    for _, e := range s.sessions {
        e.active = false
    }
}
```

### 3. executeAsync — Retry Logic

- **Purpose**: Detectar morte do bridge e refazer a request automaticamente
- **Location**: `internal/telegram/input_pipeline.go`
- **Changes**:
  - `processBridgeEventsAsync` retorna um `outcome` (success, llmError, processDeath)
  - `executeAsync` verifica o outcome e faz retry se `processDeath`
  - Na retry, request usa `Resume` com session ID (se disponível), nunca `Continue`
  - Máximo 1 retry — se falhar de novo, SendError
- **Dependencies**: `bridge.Execute`, `session.Store`
- **Reuses**: Todo o flow existente — mesma função `processBridgeEventsAsync`

```go
// input_pipeline.go — novo tipo
type bridgeOutcome int

const (
    outcomeSuccess      bridgeOutcome = iota
    outcomeLLMError                         // terminal "error" event
    outcomeProcessDeath                     // channel closed without terminal
)

// processBridgeEventsAsync — muda return type
func (bc *BotController) processBridgeEventsAsync(...) bridgeOutcome {
    // ... existing logic ...
    // case "result": return outcomeSuccess
    // case "error":  return outcomeLLMError
    // channel closed: return outcomeProcessDeath
}

// executeAsync — adiciona retry
func (bc *BotController) executeAsync(chatID int64, messageID int, req bridge.Request, userText string) {
    // ... existing setup (typing, progress, ctx) ...

    ch, err := bc.bridge.Execute(ctx, req)
    if err != nil {
        // ... existing error handling ...
        return
    }

    outcome := bc.processBridgeEventsAsync(chat, ch, progress, userText, messageID)

    if outcome != outcomeProcessDeath {
        return // success or LLM error — done
    }

    // --- RETRY PATH ---
    log.Printf("bridge: process died mid-request, retrying for chat=%d", chatID)

    // Rebuild request with Resume (not Continue) using stored session ID
    retryReq := req
    retryReq.Options.Continue = false
    if sid := bc.sessions.Get(chatID); sid != "" {
        retryReq.Options.Resume = sid
    }
    retryReq.RequestID = "" // let bridge assign new ID

    ch, err = bc.bridge.Execute(ctx, retryReq)
    if err != nil {
        log.Printf("bridge: retry failed for chat=%d: %v", chatID, err)
        _ = SendError(bc.bot, chat, "Processador reiniciado mas não conseguiu completar. Tente novamente.")
        return
    }

    // Second attempt — no more retries
    outcome = bc.processBridgeEventsAsync(chat, ch, progress, userText, messageID)
    if outcome == outcomeProcessDeath {
        _ = SendError(bc.bot, chat, "Processador reiniciado mas não conseguiu completar. Tente novamente.")
    }
}
```

### 4. Wiring — app.go

- **Purpose**: Conectar o callback `OnDeath` ao `DeactivateAll` da session store
- **Location**: `cmd/aurelia/app.go`
- **Changes**: Após criar Bridge e SessionStore, registrar o callback

```go
// app.go — no bootstrapApp(), após criar bridge e sessions
br.SetOnDeath(func() {
    log.Printf("bridge: process died unexpectedly, deactivating all sessions")
    sessions.DeactivateAll()
})
```

---

## Data Models

Nenhum modelo novo. Apenas:
- Novo campo `onDeath func()` no struct `Bridge`
- Novo tipo `bridgeOutcome` (enum int) no package telegram

---

## Error Handling Strategy

| Error Scenario | Handling | User Impact |
|----------------|----------|-------------|
| Bridge morre mid-request | Retry 1x com Resume | Resposta chega com delay extra (~10-20s) |
| Retry também falha | SendError | "Processador reiniciado mas não conseguiu completar." |
| Bridge morre + sem session ID | Retry sem Resume (session nova) | Resposta sem contexto anterior |
| Stop() intencional | Sem callback, sem retry | N/A (shutdown normal) |
| Erro terminal do LLM | Sem retry, propaga erro | Usuário vê mensagem de erro do LLM |
| Contexto expira durante retry | Abort | Não estende timeout |

---

## Tech Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Retry no Telegram layer, não no Bridge | `executeAsync` | Bridge é genérico (cron também usa). Retry faz sentido só onde a mensagem do usuário está em jogo. Cron tem seu próprio error handling. |
| Callback ao invés de channel/event | `onDeath func()` | Mais simples, um único listener (session store). Não precisa de pub/sub. Padrão já usado em `cron.NotifyingRuntime`. |
| Max 1 retry | Hardcoded | Se o bridge morreu 2x seguidas, retry infinito piora. Backoff (P3) cobre falha persistente. |
| `DeactivateAll` ao invés de clear | Preserva session IDs | Sessions cold ainda permitem `Resume` — o SDK restaura do disco. Clear perderia os IDs. |

---

## Scope Summary

**Arquivos a modificar:**
1. `internal/bridge/bridge.go` — `onDeath` field, `SetOnDeath()`, callback no `readLoop`
2. `internal/session/store.go` — `DeactivateAll()`
3. `internal/telegram/input_pipeline.go` — `bridgeOutcome` type, return value de `processBridgeEventsAsync`, retry em `executeAsync`
4. `cmd/aurelia/app.go` — wiring do callback

**Arquivos de teste:**
1. `internal/bridge/bridge_test.go` — test OnDeath callback
2. `internal/session/store_test.go` — test DeactivateAll
3. `internal/telegram/input_pipeline_test.go` — test retry on process death, no retry on LLM error

**Estimativa de mudança:** ~80 linhas de produção, ~120 linhas de teste
