package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/lpcoutinho/ops-pilot/internal/agent"
	"github.com/lpcoutinho/ops-pilot/internal/tools"
)

type OpenAIProvider struct {
	Client  *resty.Client
	APIKey  string
	Model   string
	BaseURL string
}

func NewOpenAIProvider(apiKey, model, baseURL string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIProvider{
		Client:  resty.New(),
		APIKey:  apiKey,
		Model:   model,
		BaseURL: baseURL,
	}
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (p *OpenAIProvider) GenerateResponse(ctx context.Context, prompt string, tools []tools.SystemTool) (*agent.LLMIntent, error) {
	slog.Debug("Calling OpenAI compatible API", "url", p.BaseURL, "model", p.Model)

	// Note: In a full implementation, we would convert tools to OpenAI's tool format.
	// For now, we'll implement the basic chat completion to establish connectivity.
	
	req := openAIRequest{
		Model: p.Model,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
	}

	resp, err := p.Client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+p.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(p.BaseURL + "/chat/completions")

	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode(), resp.String())
	}

	var result openAIResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error.Message != "" {
		return nil, fmt.Errorf("openai error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("empty response from openai")
	}

	return &agent.LLMIntent{
		Message: result.Choices[0].Message.Content,
	}, nil
}
