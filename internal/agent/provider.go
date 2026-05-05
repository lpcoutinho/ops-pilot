package agent

import (
	"context"
	"github.com/lpcoutinho/ops-pilot/internal/tools"
)

// LLMProvider defines the interface for different AI models.
type LLMProvider interface {
	// GenerateResponse sends a prompt and available tools to the LLM and returns the intent.
	GenerateResponse(ctx context.Context, prompt string, tools []tools.SystemTool) (*LLMIntent, error)
	// ListModels returns a list of available model identifiers for the provider.
	ListModels(ctx context.Context) ([]string, error)
}

// LLMIntent represents the output from the LLM, which could be a tool call or a direct message.
type LLMIntent struct {
	Message  string
	ToolCall *ToolCallIntent
}

type ToolCallIntent struct {
	ToolName string
	Args     []byte // JSON raw message
}
