## AI CLI Tool

### CLI interface for

1. Rewriting text (spelling, grammar, etc...)
2. Translation
3. General prompts

| Flag          | Shorthand | Description                                |
| ------------- | --------- | ------------------------------------------ |
| `--rewrite`   | `-r`      | Rewrite text                               |
| `--translate` | `-t`      | Translate text                             |
| `--model`     | `-m`      | AI provider (`ollama`, `openai`, `gemini`) |
| `--prompt`    | `-p`      | Input text or prompt                       |

> If --model is not set → defaults to Ollama.

### Examples

* Rewrite
```bash
ai -r -m openai -p "A sentence to rewrite"
```
* Translate
```bash
ai -t -m gemini -p "翻訳する行"
```
* Test
```bash
ai -m ollama -p="Summarize: Go concurrency"
```

### Required Environment Variables
```.env
OPENAI_API_KEY=sk-xxx
GEMINI_API_KEY=xyz
```

### Configuration
``config.yaml``

#### Global Settings
  ```yaml
# Sets HTTP request timeout in seconds for all AI providers and requests.
httpTimeoutSeconds: 30

# Specifies the default model for each provider.
# must match the provider’s supported models.
models:
  gemini: gemini-2.5-flash
  openai: gpt-5-nano
  ollama: llama3:latest # Local model, download required https://ollama.com/

# Defines the prompts sent to the AI for rewriting, translating, or testing.
prompts:
  test: ""
  rewrite: |
    Act as a professional editor. Your goal is to improve the following text based on these three criteria:

    1. Correct all spelling errors.
    2. Fix all grammatical mistakes.
    3. The content should be close to the provided text

    Guidelines
      - Use active voice and simple sentence structures.
      - It is okay to use contractions (e.g., "we're" instead of "we are").
      - Be conversational and approachable, but remain respectful and appropriate for a workplace environment.
      - Avoid stiff corporate jargon, but do not use slang or emojis.

    As a response, provide only the result

    Here is the text to rewrite:

  translate: "Translate the following text to english, return only the result:"
  
# Specifies API base URLs for each AI provider.
# Endpoints should point to the live API or local test servers.
baseEndpoint:
  gemini: https://generativelanguage.googleapis.com/v1beta/models/
  ollama: http://localhost:11434/api/generate
  openai: https://api.openai.com/v1/chat/completions
  ```







