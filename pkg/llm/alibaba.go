package llm

const alibabaChatCompletionsURL = "https://coding-intl.dashscope.aliyuncs.com/v1/chat/completions"

func NewAlibabaProvider(apiKey string, model string) *OpenAICompatibleProvider {
	return NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		Provider:  "alibaba",
		APIKey:    apiKey,
		BaseURL:   alibabaChatCompletionsURL,
		Model:     model,
		UserAgent: "Aurelia/1.0",
	})
}
