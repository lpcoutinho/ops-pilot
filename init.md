# Contexto: Senior Go Developer & Platform Engineer
Atue como um Engenheiro de Sistemas Staff. Vamos construir o "Ops-Pilot", uma ferramenta CLI de código aberto escrita em Go que funciona como um copiloto de linguagem natural para administração e auditoria de sistemas Linux.

# Objetivo
Criar um binário estático de alta performance que traduz intenções do usuário em comandos de sistema seguros, diagnósticos de rede e monitoramento de recursos, utilizando LLMs para o raciocínio.

# Requisitos Técnicos (Go Native)
1. **Arquitetura CLI**: Usar `spf13/cobra` para comandos e `spf13/viper` para gestão de configuração (YAML/Env).
2. **LLM Engine**: Implementar integração direta com a API da Anthropic ou OpenAI utilizando SDKs oficiais ou `go-resty` para chamadas resilientes.
3. **Typesafe Tool Calling**: 
   - Definir `structs` em Go para representar ferramentas do sistema (ex: `CheckDiskSpace`, `AnalyzeLogs`, `AuditNetwork`).
   - Usar tags de struct para gerar os esquemas JSON que a LLM utilizará para "Function Calling".
4. **Security Sandbox**:
   - Implementar um middleware de validação que impede a execução de comandos `sudo` ou destrutivos sem a flag `--dangerous-mode`.
   - Utilizar o pacote `os/exec` com contexto e timeout rigorosos.
5. **Observabilidade**: Log estruturado usando `slog` (nativo do Go) para auditoria de cada comando sugerido pela IA.

# Estrutura de Diretórios (Standard Go Layout)
- `/cmd`: Pontos de entrada da CLI.
- `/internal/agent`: Lógica do agente e orquestração da LLM.
- `/internal/tools`: Implementação das funções de sistema (net, disk, proc).
- `/pkg/validator`: Lógica de segurança e sanitização de comandos.

# Ação Inicial
1. Crie o arquivo `CLAUDE.md` com as instruções de `go build`, `go test` e padrões de código (interfaces, erro como valor).
2. Inicialize o `go.mod` e crie o esqueleto da CLI usando Cobra.
3. Planeje e implemente a interface `SystemTool` que todas as ferramentas de diagnóstico devem seguir.
4. Implemente uma ferramenta inicial `GetSystemHealth` que coleta CPU, Memória e Disco usando o pacote `shirou/gopsutil/v3`.

Mantenha o código idiomático, utilize interfaces para facilitar testes e garanta que o binário final seja estático.