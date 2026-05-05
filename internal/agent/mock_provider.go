package agent

import (
	"context"
	"github.com/lpcoutinho/ops-pilot/internal/tools"
)

// MockLLMProvider is a test implementation of LLMProvider.
type MockLLMProvider struct {
	Response *LLMIntent
	Err      error
}

func (m *MockLLMProvider) GenerateResponse(ctx context.Context, prompt string, tools []tools.SystemTool) (*LLMIntent, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Response, nil
}
