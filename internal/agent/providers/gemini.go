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

type geminiPart struct {
	Text         string              `json:"text,omitempty"`
	FunctionCall *geminiFunctionCall `json:"functionCall,omitempty"`
}

type geminiFunctionCall struct {
	Name string          `json:"name"`
	Args json.RawMessage `json:"args"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiFunctionDeclaration struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type geminiTool struct {
	FunctionDeclarations []geminiFunctionDeclaration `json:"function_declarations"`
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
	Tools    []geminiTool    `json:"tools,omitempty"`
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

func (p *GeminiProvider) GenerateResponse(ctx context.Context, prompt string, availableTools []tools.SystemTool) (*agent.LLMIntent, error) {
	slog.Debug("Calling Google Gemini API", "url", p.BaseURL, "model", p.Model)

	req := geminiRequest{
		Contents: []geminiContent{
			{
				Role: "user",
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
	}

	if len(availableTools) > 0 {
		tool := geminiTool{
			FunctionDeclarations: make([]geminiFunctionDeclaration, 0, len(availableTools)),
		}
		for _, t := range availableTools {
			tool.FunctionDeclarations = append(tool.FunctionDeclarations, geminiFunctionDeclaration{
				Name:        t.Name(),
				Description: t.Description(),
				Parameters:  t.Schema(),
			})
		}
		req.Tools = []geminiTool{tool}
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

	candidate := result.Candidates[0]
	intent := &agent.LLMIntent{}

	for _, part := range candidate.Content.Parts {
		if part.FunctionCall != nil {
			intent.ToolCall = &agent.ToolCallIntent{
				ToolName: part.FunctionCall.Name,
				Args:     part.FunctionCall.Args,
			}
			return intent, nil
		}
		if part.Text != "" {
			intent.Message = part.Text
		}
	}

	return intent, nil
}

func (p *GeminiProvider) ListModels(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/models?key=%s", p.BaseURL, p.APIKey)

	resp, err := p.Client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode(), resp.String())
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var models []string
	for _, m := range result.Models {
		models = append(models, m.Name)
	}
	return models, nil
}
