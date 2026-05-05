package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/lpcoutinho/ops-pilot/internal/tools"
	"github.com/lpcoutinho/ops-pilot/pkg/validator"
)

func TestAgentProcess_ToolCallSimulation(t *testing.T) {
	// 1. Configurar o validador de segurança
	val := &validator.CommandValidator{DangerousMode: false}

	// 2. Configurar o MockLLMProvider para simular a intenção de usar a ferramenta
	mockProvider := &MockLLMProvider{
		Response: &LLMIntent{
			// Simulamos que a LLM decidiu chamar a ferramenta 'get_system_health'
			ToolCall: &ToolCallIntent{
				ToolName: "get_system_health",
				Args:     []byte(`{}`), // Sem argumentos necessários para esta ferramenta
			},
			// Esta mensagem simula a resposta final do Agente
			Message: "The system health metrics have been collected successfully.",
		},
	}

	// 3. Inicializar o Agente
	agent := NewAgent(mockProvider, val)

	// 4. Registrar a ferramenta real que queremos testar
	agent.RegisterTool(&tools.GetSystemHealthTool{})

	// 5. Executar o fluxo do Agente
	ctx := context.Background()
	response, err := agent.Process(ctx, "how is my system doing?")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !strings.Contains(response, "metrics have been collected") {
		t.Errorf("Unexpected response from agent: %s", response)
	}
}
