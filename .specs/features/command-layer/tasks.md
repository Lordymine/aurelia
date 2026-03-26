# Command Layer — Tasks

**Design**: `.specs/features/command-layer/design.md`
**Status**: Done

---

## Execution Plan

### Phase 1: Foundation (Sequential)

Command matcher e interceptação no pipeline.

```
T1 → T2 → T3
```

### Phase 2: P1 Handlers (Parallel)

Handlers dos comandos MVP — podem rodar em paralelo pois são independentes.

```
      ┌→ T4 (session reset) ─┐
T3 ──→┼→ T5 (cron list)      ┼──→ T7
      └→ T6 (cron cancel)    ┘
```

### Phase 3: Cron Create (Sequential)

Handler mais complexo — depende do cron list pra validar.

```
T7 → T8
```

### Phase 4: P2 Handlers (Parallel)

Comandos informativos — independentes entre si.

```
      ┌→ T9  (status)      ─┐
T8 ──→┼→ T10 (list agents)  ┼──→ T12
      └→ T11 (list models)  ┘
```

### Phase 5: Integração (Sequential)

Wiring final e teste E2E.

```
T12 → T13
```

---

## Task Breakdown

### T1: Criar struct MatchedCommand e função Match

**What**: Definir o tipo `MatchedCommand` e a função `Match(text) *MatchedCommand` que detecta comandos por keyword matching com heurística anti-falso-positivo.
**Where**: `internal/telegram/commands.go` (novo)
**Depends on**: Nenhuma
**Reuses**: Nada — lógica nova

**Done when**:
- [ ] Tipo `CommandType` com constantes: `CmdCronCreate`, `CmdCronList`, `CmdCronCancel`, `CmdSessionReset`, `CmdStatus`, `CmdListAgents`, `CmdListModels`
- [ ] Tipo `MatchedCommand` com `Type CommandType` e `Text string` (texto original)
- [ ] Função `Match(text string) *MatchedCommand` com keyword matching case-insensitive
- [ ] Heurística anti-falso-positivo: keywords no início da mensagem matcham, keywords após contexto narrativo não matcham
- [ ] Testes passam: `go test ./internal/telegram/ -run TestMatch -v`

**Verify**:
```
go test ./internal/telegram/ -run TestMatch -v
```
Espera: "agenda uma reunião" → `CmdCronCreate`, "ontem eu tentei agendar" → nil, "nova conversa" → `CmdSessionReset`, "bom dia" → nil

---

### T2: Criar função HandleCommand no BotController

**What**: Função dispatcher que recebe um `MatchedCommand` e roteia pro handler correto. Handlers retornam string de resposta.
**Where**: `internal/telegram/commands.go`
**Depends on**: T1
**Reuses**: `SendTextReply` de `internal/telegram/send.go`

**Done when**:
- [ ] Método `(bc *BotController) handleCommand(c telebot.Context, cmd *MatchedCommand) error`
- [ ] Switch/map de `cmd.Type` → handler function
- [ ] Handlers stub que retornam "não implementado" (implementados nas tasks seguintes)
- [ ] Teste unitário do dispatch

**Verify**:
```
go test ./internal/telegram/ -run TestHandleCommand -v
```

---

### T3: Interceptar no processInput

**What**: Inserir chamada a `Match()` + `handleCommand()` no `processInput()`, depois do bootstrap check e antes do `routeAgent()`.
**Where**: `internal/telegram/input_pipeline.go`
**Depends on**: T1, T2
**Reuses**: Fluxo existente do `processInput()`

**Done when**:
- [ ] `processInput()` chama `Match(text)` antes de `routeAgent()`
- [ ] Se match: chama `handleCommand()` e retorna (não prossegue pro LLM)
- [ ] Se não match: fluxo continua inalterado
- [ ] `go build ./...` compila sem erros
- [ ] Testes existentes continuam passando: `go test ./internal/telegram/ -v`

**Verify**:
```
go build ./... && go test ./internal/telegram/ -v
```

---

### T4: Handler session_reset [P]

**What**: Implementar handler que limpa sessão e tracker do chat.
**Where**: `internal/telegram/commands.go`
**Depends on**: T3
**Reuses**: `session.Store.Clear()`, `session.Tracker.Clear()`

**Done when**:
- [ ] Chama `bc.sessions.Clear(chatID)` e `bc.tracker.Clear(chatID)`
- [ ] Responde: "Contexto limpo. Próxima mensagem começa uma conversa nova."
- [ ] Teste unitário: verifica que sessão e tracker foram limpos
- [ ] `go test ./internal/telegram/ -run TestSessionReset -v`

**Verify**:
```
go test ./internal/telegram/ -run TestSessionReset -v
```

---

### T5: Handler cron_list [P]

**What**: Implementar handler que lista cron jobs ativos do chat.
**Where**: `internal/telegram/commands.go`
**Depends on**: T3
**Reuses**: `cron.Service.ListJobs()`, formatação existente

**Done when**:
- [ ] Chama `cronService.ListJobs(ctx, chatID)`
- [ ] Formata resposta: ID (8 chars), prompt (truncado), schedule, status (ativo/pausado)
- [ ] Se lista vazia: "Nenhum agendamento encontrado."
- [ ] Teste unitário com mock do cron service
- [ ] `go test ./internal/telegram/ -run TestCronList -v`

**Verify**:
```
go test ./internal/telegram/ -run TestCronList -v
```

---

### T6: Handler cron_cancel [P]

**What**: Implementar handler que desativa um cron job por prefixo de ID.
**Where**: `internal/telegram/commands.go`
**Depends on**: T3
**Reuses**: `cron.Store.ResolveJobID()`, `cron.Service.DeleteJob()`

**Done when**:
- [ ] Extrai identificador da mensagem (últimas palavras ou ID)
- [ ] Resolve ID completo via `ResolveJobID(ctx, prefix)`
- [ ] Se não encontra: responde com erro amigável + sugere "meus agendamentos"
- [ ] Se encontra: deleta e confirma com ID + prompt do job
- [ ] Teste unitário com mock
- [ ] `go test ./internal/telegram/ -run TestCronCancel -v`

**Verify**:
```
go test ./internal/telegram/ -run TestCronCancel -v
```

---

### T7: Validar P1 parcial (session + cron list/cancel)

**What**: Rodar suite completa e validar que os 3 handlers locais funcionam e o fluxo LLM não foi quebrado.
**Where**: N/A — validação
**Depends on**: T4, T5, T6

**Done when**:
- [ ] `go test ./... -v` — todos os testes passam
- [ ] `go vet ./...` — sem warnings
- [ ] `go build ./...` — compila

**Verify**:
```
go test ./... -v && go vet ./... && go build ./...
```

---

### T8: Handler cron_create

**What**: Implementar handler que cria cron job usando LLM focado pra extração de parâmetros.
**Where**: `internal/telegram/commands.go`
**Depends on**: T7
**Reuses**: `bridge.ExecuteSync()`, `cron.Service.AddRecurringJob()`, `cron.Service.AddOnceJob()`

**Done when**:
- [ ] Chama `bridge.ExecuteSync()` com prompt focado de extração (system prompt: extrair schedule type, cron_expr ou run_at, e prompt)
- [ ] Timeout de 5s no context
- [ ] Parse do JSON retornado pelo LLM
- [ ] Se "cron": chama `AddRecurringJob(ctx, userID, chatID, cronExpr, prompt)`
- [ ] Se "once": chama `AddOnceJob(ctx, userID, chatID, runAt, prompt)`
- [ ] Se parse falha: responde pedindo clarificação
- [ ] Confirma no Telegram: tipo de schedule + próxima execução
- [ ] Teste unitário com mock do bridge (retorna JSON fixo)
- [ ] `go test ./internal/telegram/ -run TestCronCreate -v`

**Verify**:
```
go test ./internal/telegram/ -run TestCronCreate -v
```

---

### T9: Handler status [P]

**What**: Implementar handler que retorna status do sistema.
**Where**: `internal/telegram/commands.go`
**Depends on**: T8
**Reuses**: `bridge` (alive check), `agents.Registry`, `cron.Service`, `config.AppConfig`

**Done when**:
- [ ] Responde com: bridge status, número de agents, número de cron jobs ativos, modelo padrão
- [ ] Teste unitário
- [ ] `go test ./internal/telegram/ -run TestStatus -v`

**Verify**:
```
go test ./internal/telegram/ -run TestStatus -v
```

---

### T10: Handler list_agents [P]

**What**: Implementar handler que lista agents disponíveis.
**Where**: `internal/telegram/commands.go`
**Depends on**: T8
**Reuses**: `agents.Registry.Agents()`

**Done when**:
- [ ] Itera `bc.agents.Agents()` e formata: nome + descrição
- [ ] Se vazio: "Nenhum agent configurado."
- [ ] Teste unitário
- [ ] `go test ./internal/telegram/ -run TestListAgents -v`

**Verify**:
```
go test ./internal/telegram/ -run TestListAgents -v
```

---

### T11: Handler list_models [P]

**What**: Implementar handler que lista modelos e provedores disponíveis.
**Where**: `internal/telegram/commands.go` + `internal/config/models.go` (novo)
**Depends on**: T8
**Reuses**: `config.AppConfig.Providers`

**Done when**:
- [ ] Mapa `ProviderModels` definido em `internal/config/models.go`
- [ ] Handler itera providers do config, cruza com `ProviderModels`, marca ativos (tem API key) vs inativos
- [ ] Formata tabela: provedor | modelos | status
- [ ] Teste unitário
- [ ] `go test ./internal/config/ -run TestProviderModels -v`
- [ ] `go test ./internal/telegram/ -run TestListModels -v`

**Verify**:
```
go test ./internal/config/ -run TestProviderModels -v && go test ./internal/telegram/ -run TestListModels -v
```

---

### T12: Validar suite completa

**What**: Rodar todos os testes, vet, e build pra garantir que tudo funciona junto.
**Where**: N/A — validação
**Depends on**: T9, T10, T11

**Done when**:
- [ ] `go test ./... -v` — todos os testes passam
- [ ] `go vet ./...` — sem warnings
- [ ] `go build ./...` — compila

**Verify**:
```
go test ./... -v && go vet ./... && go build ./...
```

---

### T13: Teste E2E

**What**: Adicionar teste E2E que verifica o command layer no fluxo completo (mensagem Telegram → comando → resposta).
**Where**: `e2e/e2e_test.go` (modificar)
**Depends on**: T12
**Reuses**: Padrões de teste E2E existentes

**Done when**:
- [ ] Teste E2E: enviar "nova conversa" → verificar resposta de reset
- [ ] Teste E2E: enviar "meus agendamentos" → verificar resposta de lista
- [ ] Teste E2E: enviar "mensagem normal" → verificar que cai no fluxo LLM (não interceptada)
- [ ] `go test ./e2e/ -v`

**Verify**:
```
go test ./e2e/ -v
```

---

## Parallel Execution Map

```
Phase 1 (Sequential):
  T1 ──→ T2 ──→ T3

Phase 2 (Parallel):
  T3 complete, then:
    ├── T4 [P] session_reset
    ├── T5 [P] cron_list
    └── T6 [P] cron_cancel

Phase 3 (Sequential):
  T4+T5+T6 ──→ T7 (validate) ──→ T8 (cron_create)

Phase 4 (Parallel):
  T8 complete, then:
    ├── T9  [P] status
    ├── T10 [P] list_agents
    └── T11 [P] list_models

Phase 5 (Sequential):
  T9+T10+T11 ──→ T12 (validate) ──→ T13 (E2E)
```

---

## Task Granularity Check

| Task | Scope | Status |
|------|-------|--------|
| T1: Match function + types | 1 função + tipos | OK |
| T2: HandleCommand dispatcher | 1 função | OK |
| T3: Interceptar processInput | 1 modificação | OK |
| T4: session_reset handler | 1 handler | OK |
| T5: cron_list handler | 1 handler | OK |
| T6: cron_cancel handler | 1 handler | OK |
| T7: Validação P1 parcial | validação | OK |
| T8: cron_create handler | 1 handler (mais complexo — usa bridge) | OK |
| T9: status handler | 1 handler | OK |
| T10: list_agents handler | 1 handler | OK |
| T11: list_models handler | 1 handler + 1 constante | OK |
| T12: Validação completa | validação | OK |
| T13: Teste E2E | testes | OK |
