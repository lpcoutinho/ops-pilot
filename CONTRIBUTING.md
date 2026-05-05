# Contributing to Ops-Pilot 🚀

First off, thank you for considering contributing to Ops-Pilot! It's people like you that make the open-source community such an amazing place.

## 🛠 Tech Stack
- **Language:** Go 1.22+
- **CLI Framework:** Cobra & Viper
- **System Metrics:** gopsutil
- **LLM Integration:** Agnostic Provider Interface (Function Calling support required)

## 🏗 Project Structure
- `/cmd/ops-pilot`: Main entry point and CLI commands.
- `/internal/agent`: Agent orchestration and LLM provider implementations.
- `/internal/tools`: System tools implementation (This is where you add new capabilities!).
- `/pkg/validator`: Security logic and command sanitization.

## 🚀 How to add a new Tool
1. Create a new file in `internal/tools/your_tool.go`.
2. Implement the `SystemTool` interface:
   ```go
   type SystemTool interface {
       Name() string
       Description() string
       Schema() map[string]interface{}
       Execute(ctx context.Context, args json.RawMessage) (interface{}, error)
   }
   ```
3. Register your tool in `cmd/ops-pilot/main.go`.
4. Add a test case in `internal/agent/agent_test.go` or a new test file.

## 🧪 Running Tests
```bash
go test ./...
```

## 📮 Pull Request Process
1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. Ensure the test suite passes.
4. Update the README.md if you added a new feature.
5. Submit your PR!

## 📜 Code of Conduct
Be respectful, professional, and helpful. We are building this for the community.
