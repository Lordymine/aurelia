package agent

import "context"

type MockLLMProvider struct {
	response *ModelResponse
	err      error
}

func (m *MockLLMProvider) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	return m.response, m.err
}
