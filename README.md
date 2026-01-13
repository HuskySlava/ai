## AI CLI Tool

> **Note:** Ollama is a local AI model manager and inference engine that allows you to run large language models on your own machine. It must be **installed and configured separately**. See [Ollama installation guide](https://ollama.com/).  
> If no `--provider` flag is specified, Ollama will be used as the default provider.

### CLI interface

This CLI tool supports:

1. Rewriting text (spelling, grammar, etc...)
2. Translation
3. Summarization
4. General prompts

| Flag          | Shorthand | Description                                          |
|---------------|-----------|------------------------------------------------------|
| `--rewrite`   | `-r`      | Rewrite text                                         |
| `--translate` | `-t`      | Translate text                                       |
| `--summarize` | `-s`      | Summarize text                                       |
| `--language`  | `-l`      | Target language for translation                      |
| `--provider`  | `-p`      | AI provider (`ollama`, `openai`, `gemini`, `claude`) |
| `--input`     | `-i`      | Input text or prompt                                 |
| `--clipboard` | `-c`      | Copy result to clipboard automatically               |

> If --provider is not set → defaults to **Ollama**.

### Examples

* Rewrite
```bash
ai -r -p openai -i "A sentence to rewrite"
```
* Translate + Copy result to clipboard
```bash
ai -t -p gemini -c -i "翻訳する行"
```

* Summarize
```bash
ai -s -p gemini -i "A sentence to summarize"
```

* General Prompt
```bash
ai -p ollama -i "Summarize: Go concurrency"
```

### Required Environment Variables
```.env
OPENAI_API_KEY=sk-xxx
GEMINI_API_KEY=xyz
CLAUDE_API_KEY=xyz
```

### Configuration
Configuration is stored in ``config.yaml``

#### Global Settings
```yaml
# Sets HTTP request timeout in seconds for all AI providers.
httpTimeoutSeconds: 30

# Specifies the default model for each provider.
# Must match the provider’s supported models.
models:
  ollama: llama3:latest       # Local model, download required: https://ollama.com/
  gemini: gemini-2.5-flash
  openai: gpt-5-nano
  claude: claude-haiku-4-5-20251001

prompts:
  rewrite: |
    You are a professional editor. Improve the following text according to these instructions:

    Objectives:
    1. Correct all spelling and punctuation errors.
    2. Fix all grammatical and syntactical issues.
    3. Keep the meaning and tone as close as possible to the original text.

    Style Guidelines:
      - Use active voice and clear, simple sentences.
      - Contractions (e.g., “we’re,” “they’re”) are allowed.
      - Maintain a conversational, professional tone suitable for the workplace.
      - Avoid slang, emojis, and overly formal corporate language.

    Output Rules:
      - Return only the edited text, with no explanations, comments, or formatting other than plain text.

    Text to edit:

  translate: "Translate the following text to %s, return only the result:"
  summarize: "Summarize the following text in a clear, concise way:"

baseEndpoint:
  gemini: https://generativelanguage.googleapis.com/v1beta/models/
  ollama: http://localhost:11434/api/generate
  openai: https://api.openai.com/v1/chat/completions
  claude: https://api.anthropic.com/v1/messages
 ```







