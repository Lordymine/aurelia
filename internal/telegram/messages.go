package telegram

const (
	unsupportedDocumentMessage = "⚠️ **Formato nao suportado**\n\n" +
		"No momento eu consigo processar:\n" +
		"- arquivos `.md`\n" +
		"- arquivos `.pdf`\n" +
		"- audio e voz\n"

	downloadFailureMessage = "❌ **Falha no download**\n\n" +
		"Nao consegui baixar o arquivo enviado pelo Telegram. Tente novamente."

	audioNotConfiguredMessage = "⚠️ **Audio indisponivel**\n\n" +
		"Meu modulo de transcricao nao esta configurado.\n\n" +
		"Configure a variavel `GROQ_API_KEY` no arquivo `.env`."

	audioProcessingFailureMessage = "❌ **Falha na transcricao**\n\n" +
		"Nao consegui compreender o audio. Tente falar mais claro ou mais perto do microfone."

	emptyAudioMessage = "⚠️ **Audio vazio**\n\n" +
		"Nao captei conteudo util. Pode reenviar?"

	alreadyConfiguredMessage = "✅ **Aurelia online**\n\n" +
		"Ja estou configurado e pronto. Como posso ajudar?"

	bootstrapWelcomeMessage = "# Boas-vindas\n\n" +
		"Eu sou o **Aurelia** recem-iniciado.\n\n" +
		"Escolha como voce quer que eu atue primariamente hoje."

	bootstrapFailureMessage = "❌ **Falha no bootstrap**\n\n" +
		"Nao consegui criar os arquivos base de persona."

	bootstrapProfileMessage = "✅ **Modo inicial configurado**\n\n" +
		"Agora me diga seu nome e como prefere que eu trabalhe com voce.\n\n" +
		"Exemplo: `Me chamo Rafael e quero respostas diretas, sem floreios.`"

	bootstrapSuccessMessage = "✅ **Personas criadas**\n\n" +
		"Suas configuracoes base foram salvas em `~/.aurelia/memory/personas/`.\n\n" +
		"Voce ja pode conversar comigo ou editar:\n" +
		"- `IDENTITY.md`\n" +
		"- `SOUL.md`\n" +
		"- `USER.md`\n\n" +
		"para refinar nosso comportamento."
)


