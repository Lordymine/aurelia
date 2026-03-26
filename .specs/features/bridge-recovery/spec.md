# Bridge Recovery Automático — Specification

## Problem Statement

Quando o processo bridge (TypeScript) morre durante uma execução, a request em voo é perdida. O usuário no Telegram recebe "O processador encerrou sem resposta" e precisa reenviar a mensagem manualmente. Além disso, todas as sessions marcadas como `active` ficam stale — apontando pra um processo que não existe mais — e a próxima mensagem tenta `Continue` (warm session) num processo novo, o que falha silenciosamente.

Hoje o bridge nunca morreu em produção, mas quando acontecer o impacto é perda total da request + estado de sessão inconsistente.

## Goals

- [x] Request em voo sobrevive a um crash do bridge — retry automático com session resume
- [x] Sessions invalidadas corretamente após restart do processo
- [x] Usuário informado do que está acontecendo (não silêncio)
- [x] Falhas permanentes (config, credencial) não causam retry loop infinito

## Out of Scope

- Health check proativo / heartbeat periódico (pode vir depois)
- Restart automático do bridge quando idle (sem request ativa)
- Persistência de sessions em disco (já existe via Claude SDK, usamos o que tem)
- Mudanças no lado TypeScript do bridge

---

## User Stories

### P1: Retry Transparente da Request — MVP

**User Story**: Como usuário do Telegram, quero que minha mensagem seja reprocessada automaticamente se o bridge morrer durante a execução, pra não precisar reenviar manualmente.

**Why P1**: Sem isso, qualquer crash do bridge perde a mensagem do usuário. É o cenário mais visível e frustrante.

**Acceptance Criteria**:

1. WHEN o bridge morre durante uma `Execute()` (channel fecha sem evento terminal) THEN o sistema SHALL reiniciar o processo bridge e reenviar a mesma request com `Resume` (session ID) se disponível
2. WHEN o retry também falha THEN o sistema SHALL enviar mensagem de erro ao Telegram: "Processador reiniciado mas não conseguiu completar. Tente novamente."
3. WHEN o bridge retorna erro terminal (evento `error`) THEN o sistema SHALL NOT fazer retry — erro do LLM não é crash do processo
4. WHEN o retry é bem-sucedido THEN o sistema SHALL entregar a resposta normalmente ao usuário, como se nada tivesse acontecido

**Independent Test**: Matar o processo bridge (kill PID) durante uma execução e verificar que a resposta chega ao Telegram sem reenvio manual.

---

### P1: Invalidação de Sessions Após Restart — MVP

**User Story**: Como sistema, preciso que todas as sessions ativas sejam marcadas como cold após um restart do bridge, pra que a próxima mensagem use `Resume` (restore de disco) em vez de `Continue` (warm session inexistente).

**Why P1**: Sem isso, a primeira mensagem após um restart tenta `Continue` num processo novo que não tem a session em memória, resultando em comportamento imprevisível.

**Acceptance Criteria**:

1. WHEN o processo bridge é reiniciado (por crash ou Stop/Start) THEN o sistema SHALL marcar todas as sessions como `active: false` (cold)
2. WHEN uma session é marcada como cold THEN a próxima request para aquele chat SHALL usar `Resume` com o session ID em vez de `Continue`
3. WHEN o bridge reinicia e a session original não pode ser restaurada THEN o sistema SHALL iniciar uma session nova sem erro

**Independent Test**: Reiniciar o bridge, enviar mensagem, e verificar nos logs que o modo é `resume` (não `continue`).

---

### P2: Feedback Visual Durante Recovery

**User Story**: Como usuário do Telegram, quero saber que o sistema está se recuperando de um problema, pra não achar que minha mensagem foi ignorada.

**Why P2**: Melhora a experiência mas não é blocker — sem isso o retry ainda funciona, só demora uns segundos em silêncio.

**Acceptance Criteria**:

1. WHEN o bridge morre e o retry é iniciado THEN o sistema SHALL enviar uma reação (emoji) ou mensagem indicando reconexão (ex: "Reconectando...")
2. WHEN o retry é bem-sucedido THEN o sistema SHALL remover/editar a mensagem de reconexão ou simplesmente enviar a resposta final
3. WHEN o retry falha THEN o sistema SHALL atualizar a mensagem de reconexão com o erro final

**Independent Test**: Matar o bridge durante execução e verificar que o Telegram mostra indicação de reconexão antes da resposta.

---

### P3: Backoff em Falhas Consecutivas

**User Story**: Como sistema, quero evitar restart loops quando o bridge está falhando repetidamente (ex: dependência quebrada, Node.js corrompido).

**Why P3**: Protege contra cenários raros de falha persistente. Sem isso, o sistema tenta restart infinito em intervalos de milissegundos.

**Acceptance Criteria**:

1. WHEN o bridge falha 3 vezes consecutivas em menos de 1 minuto THEN o sistema SHALL entrar em modo cooldown e parar de restartar por 30 segundos
2. WHEN em cooldown THEN requests novas SHALL receber erro imediato: "Processador temporariamente indisponível."
3. WHEN uma execução é bem-sucedida THEN o contador de falhas SHALL ser resetado

**Independent Test**: Configurar bridge pra falhar no start (comando errado) e verificar que após 3 tentativas o sistema para e reporta cooldown.

---

## Edge Cases

- WHEN o bridge morre e não existe session ID armazenado (primeira mensagem do chat) THEN o sistema SHALL fazer retry sem Resume — session nova
- WHEN o bridge morre durante uma `ExecuteSync()` (usado no classify) THEN o sistema SHALL propagar o erro pro caller sem retry — classify é idempotente, o caller pode refazer
- WHEN o contexto (timeout de 10min) expira durante o retry THEN o sistema SHALL abortar e enviar erro — não estender o timeout
- WHEN múltiplas requests estão em voo e o bridge morre THEN o sistema SHALL fazer retry de todas as requests pendentes, não apenas uma
- WHEN o bridge morre durante execução de cron job THEN o sistema SHALL seguir o mesmo fluxo de retry — cron usa `ExecuteSync`, mesmo path

---

## Success Criteria

- [x] Bridge crash durante execução Telegram → resposta chega ao usuário sem reenvio manual
- [x] Sessions consistentes após restart — nenhum `Continue` pra processo inexistente
- [x] Zero retry loops infinitos — backoff previne restart storm
- [x] Testes unitários cobrem: retry com sucesso, retry com falha, session invalidation, backoff
