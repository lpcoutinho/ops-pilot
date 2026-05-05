# 🚀 Ops-Pilot

**Ops-Pilot** acts as a natural language co-pilot for Linux and macOS system administration. It translates your intent into safe system commands, provides hardware diagnostics, and monitors system health using state-of-the-art LLMs.

---

## 🛠 Features

- **Natural Language Interface:** Ask questions like "How is my system health?" or "What's consuming my memory?"
- **Multi-LLM Support:** Seamless integration with **Anthropic (Claude)**, **OpenAI (GPT-4)**, **Google Gemini**, and **Ollama** (local models).
- **Safety First:** Built-in validator to prevent destructive commands unless `--dangerous-mode` is explicitly used.
- **Native Diagnostics:** Direct system calls via `gopsutil` for high-performance health monitoring.
- **Global Installation:** Easy to install and use from anywhere in your terminal.

---

## 📥 Installation

### Homebrew (macOS & Linux)
```bash
brew tap lpcoutinho/tap
brew install ops-pilot
```

### From Source (using Makefile)
```bash
git clone https://github.com/lpcoutinho/ops-pilot.git
cd ops-pilot
make install
```

### Direct Download
Download the latest binary for your architecture from the [Releases](https://github.com/lpcoutinho/ops-pilot/releases) page.

---

## ⚙️ Configuration

Create a file named `.ops-pilot.yaml` in your home directory or the project root:

```yaml
llm_provider: gemini # options: gemini, openai, anthropic, ollama
llm_api_key: "your-api-key"
llm_model: "gemini-1.5-flash" # optional
```

You can also use environment variables:
```bash
export LLM_PROVIDER=gemini
export LLM_API_KEY="your-api-key"
```

---

## 🚀 Usage

```bash
# General health check
ops-pilot "How is my system doing?"

# List available models for your provider
ops-pilot models

# Enable debug logging
ops-pilot "Check my disk space" --debug
```

---

## 🛡 Security

Ops-Pilot includes a security sandbox. Commands like `sudo`, `rm`, `mv`, and network modifications are restricted by default. To allow potentially dangerous operations, use the flag:
`--dangerous-mode`

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
