package providers

import (
	"fmt"

	"github.com/lpcoutinho/ops-pilot/internal/agent"
	"github.com/spf13/viper"
)

func NewProviderFromConfig() (agent.LLMProvider, error) {
	providerName := viper.GetString("llm_provider")
	apiKey := viper.GetString("llm_api_key")
	model := viper.GetString("llm_model")
	baseURL := viper.GetString("llm_base_url")

	switch providerName {
	case "openai":
		if model == "" {
			model = "gpt-4o"
		}
		return NewOpenAIProvider(apiKey, model, baseURL), nil
	case "anthropic":
		if model == "" {
			model = "claude-3-5-sonnet-20240620"
		}
		return NewAnthropicProvider(apiKey, model, baseURL), nil
	case "gemini":
		if model == "" {
			model = "gemini-1.5-pro"
		}
		return NewGeminiProvider(apiKey, model, baseURL), nil
	case "ollama":
		if model == "" {
			model = "llama3"
		}
		if baseURL == "" {
			baseURL = "http://localhost:11434/v1"
		}
		// Ollama is OpenAI-compatible for chat completions
		return NewOpenAIProvider(apiKey, model, baseURL), nil
	case "mock", "":
		return &agent.MockLLMProvider{
			Response: &agent.LLMIntent{Message: "Mock response: LLM integration is working!"},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported llm provider: %s", providerName)
	}
}
