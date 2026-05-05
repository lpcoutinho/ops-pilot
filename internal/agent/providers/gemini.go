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

type GeminiProvider struct {
	Client  *resty.Client
	APIKey  string
	Model   string
	BaseURL string
}

func NewGeminiProvider(apiKey, model, baseURL string) *GeminiProvider {
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	return &GeminiProvider{
		Client:  resty.New(),
		APIKey:  apiKey,
		Model:   model,
		BaseURL: baseURL,
	}
}

type geminiContent struct {
	Parts []struct {
		Text string `json:"text"`
	} `json:"parts"`
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

func (p *GeminiProvider) GenerateResponse(ctx context.Context, prompt string, tools []tools.SystemTool) (*agent.LLMIntent, error) {
	slog.Debug("Calling Google Gemini API", "url", p.BaseURL, "model", p.Model)

	req := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.BaseURL, p.Model, p.APIKey)

	resp, err := p.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode(), resp.String())
	}

	var result geminiResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error.Message != "" {
		return nil, fmt.Errorf("gemini error (code %d): %s", result.Error.Code, result.Error.Message)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from gemini")
	}

	return &agent.LLMIntent{
		Message: result.Candidates[0].Content.Parts[0].Text,
	}, nil
}
