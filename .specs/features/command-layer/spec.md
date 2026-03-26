# Command Layer — Specification

## Problem Statement

Toda mensagem no Aurelia OS passa pelo LLM classify, mesmo operações que o sistema pode resolver localmente sem inteligência artificial. Agendar um cron job, limpar sessão, ou consultar status são operações determinísticas que hoje gastam ~15s de latência e tokens desnecessários passando pelo bridge. O sistema precisa de uma camada de comandos do sistema que intercepta essas operações antes do LLM.

## Goals

- [x] Operações do sistema resolvidas localmente em Go, sem LLM
- [x] Latência de comandos locais < 500ms (vs ~15s atual do classify + bridge)
- [x] Zero tokens gastos em operações que o Go pode resolver sozinho
- [x] Fallback transparente pro LLM quando nenhum comando matcha

## Out of Scope

- Routing de agents (o LLM/SDK decide sozinho quem acionar)
- Troca manual de agent ("usa o prospector") — a Aurelia orquestra automaticamente
- Comandos com sintaxe rígida estilo CLI (o input é linguagem natural)
- Trocar modelo em runtime (modelos são definidos por agent/config)

---

## User Stories

### P1: CRUD de Cron Jobs via Comando Local ⭐ MVP

**User Story**: Como usuário, quero agendar, listar e cancelar tarefas recorrentes falando em linguagem natural, sem esperar o LLM classify e sem depender do cron do Claude SDK.

**Why P1**: Agendamento é a dor mais clara. O Go já tem todo o infra (`internal/cron/`), mas hoje a mensagem precisa passar pelo LLM classify + bridge pra chegar lá.

**Acceptance Criteria**:

1. WHEN mensagem contém intenção de agendar (ex: "agenda X todo dia às 9h", "cria um lembrete pra amanhã") THEN sistema SHALL criar o job no cron store local sem chamar LLM classify
2. WHEN mensagem pede listagem de agendamentos (ex: "meus agendamentos", "o que tá agendado?") THEN sistema SHALL consultar o SQLite e responder direto no Telegram
3. WHEN mensagem pede cancelamento (ex: "cancela o agendamento X", "remove o lembrete") THEN sistema SHALL desativar o job no cron store e confirmar no Telegram
4. WHEN mensagem não matcha nenhum padrão de comando THEN sistema SHALL seguir o fluxo normal (LLM → bridge)

**Independent Test**: Mandar "agenda um lembrete pra amanhã às 10h" no Telegram e verificar que o job aparece no SQLite em < 500ms, sem nenhuma chamada ao bridge.

---

### P1: Gerenciar Sessão via Comando Local ⭐ MVP

**User Story**: Como usuário, quero limpar o contexto da conversa dizendo "limpa o contexto" ou "nova conversa", sem gastar LLM.

**Why P1**: Reset de sessão é operação do sistema, não requer raciocínio.

**Acceptance Criteria**:

1. WHEN mensagem contém intenção de limpar sessão (ex: "nova conversa", "limpa o contexto", "reset") THEN sistema SHALL invalidar a sessão atual e confirmar no Telegram
2. WHEN sessão é limpa THEN a próxima mensagem SHALL iniciar uma sessão nova no bridge (sem continue/resume)

**Independent Test**: Ter uma sessão ativa, mandar "nova conversa", e verificar que a próxima mensagem cria uma nova session ID.

---

### P2: Status do Sistema

**User Story**: Como usuário, quero perguntar "status" e receber um diagnóstico rápido do sistema.

**Why P2**: Útil pra diagnóstico, mas não é dor diária.

**Acceptance Criteria**:

1. WHEN mensagem pede status (ex: "status", "tá funcionando?") THEN sistema SHALL responder com: estado do bridge (up/down/recovery), número de agents carregados, número de cron jobs ativos, modelo padrão configurado
2. WHEN bridge está em recovery THEN sistema SHALL indicar isso no status

**Independent Test**: Mandar "status" e verificar que a resposta contém informações do sistema em < 500ms.

---

### P2: Listar Agents Disponíveis

**User Story**: Como usuário, quero perguntar "quais agents eu tenho?" e receber a lista direto, sem LLM.

**Why P2**: Complementa o entendimento de quem a Aurelia pode escalar.

**Acceptance Criteria**:

1. WHEN mensagem pede lista de agents (ex: "quais agents?", "lista agents") THEN sistema SHALL responder com nome e descrição de cada agent registrado
2. WHEN não há agents registrados THEN sistema SHALL informar que nenhum agent está configurado

**Independent Test**: Mandar "quais agents?" e verificar que a resposta lista todos os agents do registry.

---

### P2: Listar Modelos e Provedores

**User Story**: Como usuário, quero perguntar "quais modelos eu tenho?" e ver uma tabela com modelos e provedores disponíveis.

**Why P2**: Precisa de um registry de modelos que ainda não existe.

**Acceptance Criteria**:

1. WHEN mensagem pede lista de modelos (ex: "quais modelos?", "lista provedores") THEN sistema SHALL responder com tabela de modelos disponíveis e seus provedores
2. WHEN um provedor está sem API key configurada THEN sistema SHALL indicar que o provedor não está ativo

**Independent Test**: Mandar "quais modelos?" e verificar que a resposta lista modelos com provedores.

---

## Edge Cases

- WHEN mensagem é ambígua entre comando e conversa (ex: "quero agendar" sem detalhes) THEN sistema SHALL pedir clarificação OU deixar cair pro LLM
- WHEN comando é enviado durante bootstrap (perfil/assistente pendente) THEN sistema SHALL ignorar o command layer e manter o fluxo de bootstrap
- WHEN mensagem contém keyword de comando dentro de contexto conversacional (ex: "ontem eu tentei agendar uma reunião") THEN sistema SHALL não interceptar e deixar pro LLM

---

## Success Criteria

- [x] Comandos P1 resolvidos localmente em < 500ms sem chamada ao bridge
- [x] Zero regressão no fluxo atual — mensagens que não são comandos continuam funcionando identicamente
- [x] Detecção de intenção com taxa de falso positivo < 5% (não interceptar conversa normal como comando)
