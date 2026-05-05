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

type AnthropicProvider struct {
	Client  *resty.Client
	APIKey  string
	Model   string
	BaseURL string
}

func NewAnthropicProvider(apiKey, model, baseURL string) *AnthropicProvider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	return &AnthropicProvider{
		Client:  resty.New(),
		APIKey:  apiKey,
		Model:   model,
		BaseURL: baseURL,
	}
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	Messages  []anthropicMessage `json:"messages"`
	MaxTokens int                `json:"max_tokens"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"content"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func (p *AnthropicProvider) GenerateResponse(ctx context.Context, prompt string, tools []tools.SystemTool) (*agent.LLMIntent, error) {
	slog.Debug("Calling Anthropic API", "url", p.BaseURL, "model", p.Model)

	req := anthropicRequest{
		Model:     p.Model,
		MaxTokens: 1024,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
	}

	resp, err := p.Client.R().
		SetContext(ctx).
		SetHeader("x-api-key", p.APIKey).
		SetHeader("anthropic-version", "2023-06-01").
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(p.BaseURL + "/messages")

	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode(), resp.String())
	}

	var result anthropicResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error.Message != "" {
		return nil, fmt.Errorf("anthropic error (%s): %s", result.Error.Type, result.Error.Message)
	}

	if len(result.Content) == 0 {
		return nil, fmt.Errorf("empty response from anthropic")
	}

	return &agent.LLMIntent{
		Message: result.Content[0].Text,
	}, nil
}
