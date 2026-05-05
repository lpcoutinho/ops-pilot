package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lpcoutinho/ops-pilot/internal/tools"
	"github.com/lpcoutinho/ops-pilot/pkg/validator"
)

// Agent orchestrates the interaction between the user, LLM, and system tools.
type Agent struct {
	Provider  LLMProvider
	Tools     map[string]tools.SystemTool
	Validator *validator.CommandValidator
}

func NewAgent(provider LLMProvider, v *validator.CommandValidator) *Agent {
	return &Agent{
		Provider:  provider,
		Tools:     make(map[string]tools.SystemTool),
		Validator: v,
	}
}

func (a *Agent) RegisterTool(t tools.SystemTool) {
	a.Tools[t.Name()] = t
}

func (a *Agent) Process(ctx context.Context, input string) (string, error) {
	slog.Info("Agent processing input", "input", input)

	availableTools := make([]tools.SystemTool, 0, len(a.Tools))
	for _, t := range a.Tools {
		availableTools = append(availableTools, t)
	}

	intent, err := a.Provider.GenerateResponse(ctx, input, availableTools)
	if err != nil {
		return "", fmt.Errorf("llm provider error: %w", err)
	}

	if intent.ToolCall != nil {
		tool, ok := a.Tools[intent.ToolCall.ToolName]
		if !ok {
			return "", fmt.Errorf("llm requested unknown tool: %s", intent.ToolCall.ToolName)
		}

		slog.Info("Executing tool", "tool", tool.Name())
		result, err := tool.Execute(ctx, intent.ToolCall.Args)
		if err != nil {
			return "", fmt.Errorf("tool execution error: %w", err)
		}

		// Re-send to LLM to interpret the result
		resultJSON, _ := json.Marshal(result)
		followUpPrompt := fmt.Sprintf("System Tool %s result: %s\nUser original intent: %s\nSummarize the findings.", 
			tool.Name(), string(resultJSON), input)
		
		finalIntent, err := a.Provider.GenerateResponse(ctx, followUpPrompt, nil)
		if err != nil {
			return "", fmt.Errorf("llm follow-up error: %w", err)
		}
		return finalIntent.Message, nil
	}

	return intent.Message, nil
}
