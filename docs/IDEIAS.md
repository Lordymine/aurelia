# IDEIAS

## Backlog

### TUI de modelos com sugestao por capability

Contexto:
- a TUI ja suporta busca por texto e filtros por `vision`, `tools` e `free`
- ainda falta aproveitar isso no runtime quando o usuario tenta usar uma capability que o modelo atual nao suporta

Proxima evolucao desejada:
- quando o usuario enviar imagem para um modelo sem vision, a Aurelia deve sugerir um modelo compativel do catalogo atual
- a sugestao pode priorizar o provider atual e depois providers proximos como `OpenRouter` e `Kilo`
- opcionalmente a TUI pode destacar `recommended for vision` quando houver metadado suficiente

Racional:
- reduz erro de configuracao
- melhora a experiencia de uso sem exigir que o usuario conheca todos os modelos
- reaproveita o catalogo e os badges de capability que ja existem

### Orcamento de uso, janela de contexto e fallback por limite

Contexto:
- hoje o runtime escolhe um provider/modelo principal, mas nao tem nocao de saude de contexto nem de orcamento de uso
- alguns modos de acesso, especialmente API keys e modos tipo Codex/ChatGPT plan, podem ter limites por janela de tempo ou consumo

Proxima evolucao desejada:
- expor uma visao de `context health` para evitar `context rot` e permitir reset/compactacao da janela
- rastrear uso e limites por provider quando houver sinal confiavel
- usar fallback automatico quando um modelo estiver perto do limite, indisponivel ou caro demais para a politica atual

Escopo inicial sugerido:
- `context window manager` com compactacao/reset previsivel
- `budget registry` por provider com sinais confiaveis de uso
- `fallback policy` entre modelos/providers configurados

Observacoes:
- para API key, parte disso pode vir de catalogo, usage local e APIs oficiais de credito/quota quando existirem
- para modos tipo OAuth/subscription com janela `5h` ou semanal, provavelmente sera preciso usar heuristica ou sinais indiretos, porque nem sempre ha API oficial publica para saldo restante
