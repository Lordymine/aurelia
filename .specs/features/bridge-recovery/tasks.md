# Bridge Recovery Automático — Tasks

**Design**: `.specs/features/bridge-recovery/design.md`
**Status**: Done

---

## Execution Plan

### Phase 1: Foundation (Sequential)

```
T1 → T2
```

### Phase 2: Core (Parallel)

```
     ┌→ T3 ─┐
T2 ──┤      ├──→ T5
     └→ T4 ─┘
```

### Phase 3: Integration (Sequential)

```
T5 → T6
```

---

## Task Breakdown

### T1: Add `DeactivateAll` to Session Store

**What**: Novo método que marca todas as sessions como `active: false` sem remover os session IDs.
**Where**: `internal/session/store.go`
**Depends on**: None
**Reuses**: Estrutura existente de `entry{sessionID, active}`

**Done when**:
- [ ] Método `DeactivateAll()` implementado
- [ ] Teste: store com 3 sessions ativas → `DeactivateAll()` → todas retornam `active: false` via `GetWithState()`
- [ ] Teste: session IDs preservados após `DeactivateAll()`
- [ ] Tests pass: `go test ./internal/session/...`

**Verify:**
```bash
go test ./internal/session/ -run TestDeactivateAll -v
```

---

### T2: Add `OnDeath` callback to Bridge

**What**: Campo `onDeath func()` no struct Bridge, setter `SetOnDeath()`, e chamada no `readLoop` quando processo morre inesperadamente (não durante `Stop()`).
**Where**: `internal/bridge/bridge.go`
**Depends on**: None (pode rodar em paralelo com T1, mas T2 é pré-req das tasks seguintes)
**Reuses**: Padrão de callback de `cron.NotifyingRuntime`

**Implementation details**:
- Novo campo `onDeath func()` no struct `Bridge`
- `SetOnDeath(fn func())` protegido por `b.mu`
- No `readLoop`, após `scanner.Scan()` retornar false: se `!b.stopping && b.onDeath != nil`, chamar `b.onDeath()` **antes** do cleanup de pending channels
- Chamar callback **antes** do cleanup garante que `DeactivateAll` roda antes dos channels fecharem e o retry começar

**Done when**:
- [ ] `SetOnDeath()` implementado
- [ ] Teste: bridge com processo que morre → callback chamado
- [ ] Teste: bridge com `Stop()` intencional → callback **não** chamado
- [ ] Teste: bridge sem callback configurado → sem panic
- [ ] Tests pass: `go test ./internal/bridge/...`

**Verify:**
```bash
go test ./internal/bridge/ -run TestOnDeath -v
```

---

### T3: Change `processBridgeEventsAsync` return type [P]

**What**: Mudar `processBridgeEventsAsync` de `void` pra retornar `bridgeOutcome` (success, llmError, processDeath). Mover a lógica de envio de mensagens pro caller onde necessário.
**Where**: `internal/telegram/input_pipeline.go`
**Depends on**: T2
**Reuses**: Toda a lógica existente de event processing

**Implementation details**:
- Novo tipo `bridgeOutcome int` com constantes `outcomeSuccess`, `outcomeLLMError`, `outcomeProcessDeath`
- `case "result"`: continua enviando resposta, retorna `outcomeSuccess`
- `case "error"`: continua enviando erro, retorna `outcomeLLMError`
- Channel closed (fim do `for range`): **não envia mais** mensagem de erro — retorna `outcomeProcessDeath` com o texto parcial. O caller decide se faz retry ou envia erro.
- Mudar assinatura pra retornar `(bridgeOutcome, string)` — segundo valor é texto parcial acumulado (usado no path de processDeath)

**Done when**:
- [ ] Tipo `bridgeOutcome` definido
- [ ] `processBridgeEventsAsync` retorna outcome + texto parcial
- [ ] Path de sucesso e erro do LLM inalterados (mensagens enviadas dentro da função)
- [ ] Path de process death não envia mais mensagem (retorna pro caller)
- [ ] `executeAsync` temporariamente trata `outcomeProcessDeath` como antes (envia erro) — retry vem em T4
- [ ] Tests pass: `go test ./internal/telegram/...`

**Verify:**
```bash
go test ./internal/telegram/ -run TestProcessBridgeEvents -v
```

---

### T4: Add retry logic to `executeAsync` [P]

**What**: Quando `processBridgeEventsAsync` retorna `outcomeProcessDeath`, refazer a request com `Resume` e tentar uma segunda vez.
**Where**: `internal/telegram/input_pipeline.go`
**Depends on**: T3
**Reuses**: `bridge.Execute()`, `sessions.Get()`

**Implementation details**:
- Após `outcomeProcessDeath`:
  1. Log: `"bridge: process died mid-request, retrying for chat=%d"`
  2. Rebuild request: `Continue = false`, `Resume = sessions.Get(chatID)` se existir, `RequestID = ""`
  3. `bridge.Execute(ctx, retryReq)` — auto-restart via `startLocked()`
  4. Segunda chamada a `processBridgeEventsAsync`
  5. Se segunda tentativa retorna `outcomeProcessDeath` → `SendError("Processador reiniciado mas não conseguiu completar.")`
  6. Se segunda tentativa retorna `outcomeSuccess` ou `outcomeLLMError` → já tratado dentro da função
- O context original (10min timeout) é compartilhado — não estende o timeout
- Se texto parcial da primeira tentativa existir, descarta — a retry refaz do zero com Resume

**Done when**:
- [ ] Retry executado quando outcome é `processDeath`
- [ ] Request reconstruída com `Resume` (não `Continue`)
- [ ] Sem retry quando outcome é `llmError` ou `success`
- [ ] Máximo 1 retry
- [ ] Teste: simular process death → retry → success
- [ ] Teste: simular process death → retry → death novamente → SendError
- [ ] Teste: simular LLM error → sem retry
- [ ] Tests pass: `go test ./internal/telegram/...`

**Verify:**
```bash
go test ./internal/telegram/ -run TestRetry -v
```

---

### T5: Wire OnDeath callback in app.go

**What**: Conectar `bridge.SetOnDeath` ao `sessions.DeactivateAll` no bootstrap da aplicação.
**Where**: `cmd/aurelia/app.go`
**Depends on**: T1, T2 (precisa de ambos)
**Reuses**: Padrão existente de wiring em `bootstrapApp()`

**Implementation details**:
- Após criar `br` (Bridge) e `sessions` (session.Store), adicionar:
  ```go
  br.SetOnDeath(func() {
      log.Printf("bridge: process died, deactivating all sessions")
      sessions.DeactivateAll()
  })
  ```
- Posicionar após a criação de ambos e antes de `telegram.NewBotController()`

**Done when**:
- [ ] Callback registrado no bootstrap
- [ ] Log de morte do bridge visível nos logs da aplicação
- [ ] Build compila: `go build ./cmd/aurelia/`

**Verify:**
```bash
go build ./cmd/aurelia/
```

---

### T6: Integration test — full retry path

**What**: Teste end-to-end que valida o fluxo completo: bridge morre → sessions desativadas → retry com Resume → resposta entregue.
**Where**: `internal/telegram/input_pipeline_test.go`
**Depends on**: T3, T4, T5
**Reuses**: Fake bridge pattern existente nos testes

**Implementation details**:
- Fake bridge que aceita primeira request normalmente, mata o "processo" (fecha channel sem terminal), e na segunda request responde com sucesso
- Verificar que a session foi marcada como cold (mock de session store)
- Verificar que a segunda request usa `Resume` (não `Continue`)
- Verificar que o usuário recebe a resposta final

**Done when**:
- [ ] Teste do fluxo completo: death → deactivate → retry → success
- [ ] Teste do fluxo de falha: death → retry → death → error message
- [ ] Tests pass: `go test ./internal/telegram/ -run TestBridgeRecovery -v`

**Verify:**
```bash
go test ./internal/telegram/ -run TestBridgeRecovery -v
go test ./... -short
```

---

## Parallel Execution Map

```
Phase 1 (Sequential):
  T1 ──→ ─┐
  T2 ──→ ─┤  (T1 e T2 podem rodar em paralelo entre si)
           │
Phase 2 (Parallel, após T2):
           ├── T3 [P] (change return type)
           └── T4 [P] (retry logic — depende de T3 na prática)
               Note: T4 depende de T3, então são sequenciais

Phase 3 (Sequential, após T1+T2+T3+T4):
  T5 ──→ T6
```

**Ordem real de execução:**
```
T1 ─────────────────────────┐
T2 → T3 → T4 ──────────────┼→ T5 → T6
                            │
```

---

## Task Granularity Check

| Task | Scope | Status |
|------|-------|--------|
| T1: DeactivateAll | 1 método + teste | Granular |
| T2: OnDeath callback | 1 campo + 1 setter + hook no readLoop | Granular |
| T3: Return type change | 1 tipo + refactor de 1 função | Granular |
| T4: Retry logic | ~30 linhas em 1 função | Granular |
| T5: Wiring | 4 linhas em app.go | Granular |
| T6: Integration test | 1 test file | Granular |
