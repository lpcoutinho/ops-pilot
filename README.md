# 🚀 Ops-Pilot

**Ops-Pilot** is an AI-powered CLI co-pilot for Linux and macOS system administration. It translates natural language into safe system diagnostics and monitoring actions.

---

## 🛠 Features

- **Natural Language Interface:** Ask "Why is my system slow?" or "Check my network health."
- **Multi-LLM Support:** Seamless integration with **Anthropic**, **OpenAI**, **Google Gemini**, and **Ollama**.
- **Comprehensive Diagnostic Tools:**
  - `get_system_health`: Real-time CPU, RAM, and Disk overview.
  - `get_top_processes`: Identify resource-heavy processes.
  - `audit_network`: Detailed interface and I/O traffic audit.
  - `analyze_logs`: AI-driven analysis of system logs (`syslog`/`journalctl`).
  - `get_hardware_info`: Kernel, Uptime, and CPU model details.
- **Safety First:** Built-in command validator to prevent accidental destruction.
- **Automated Distribution:** Native binaries via GitHub Releases and Homebrew.

---

## 📥 Installation

### Homebrew (Recommended)
```bash
brew tap lpcoutinho/tap
brew install ops-pilot
```

### From Source
```bash
git clone https://github.com/lpcoutinho/ops-pilot.git
cd ops-pilot
make install
```

---

## ⚙️ Configuration

Create a `.ops-pilot.yaml` in your home directory:

```yaml
llm_provider: gemini # options: gemini, openai, anthropic, ollama
llm_api_key: "your-api-key"
llm_model: "gemini-1.5-flash"
```

---

## 🚀 Usage

```bash
ops-pilot "How is my system health?"
ops-pilot "Which processes are consuming most RAM?"
ops-pilot "Check network statistics"
ops-pilot "Analyze recent system logs for errors"
```

---

## 🤝 Contributing

We love contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to add new tools or providers.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
