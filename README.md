## AI CLI Tool

### CLI interface for

1. Rewriting text (spelling, grammar, etc...)
2. Translation
3. General prompts

| Flag          | Shorthand | Description                                |
|---------------|-----------|--------------------------------------------|
| `--rewrite`   | `-r`      | Rewrite text                               |
| `--translate` | `-t`      | Translate text                             |
| `--language`  | `-l`      | Translate target language                  |
| `--provider`  | `-p`      | AI provider (`ollama`, `openai`, `gemini`) |
| `--input`     | `-i`      | Input text or prompt                       |

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
httpTimeoutSeconds: 30

models:
  gemini: gemini-2.5-flash
  openai: gpt-5-nano
  ollama: llama3:latest

prompts:
  test: ""
  rewrite: |
    You are a professional editor. Improve the following text according to these instructions:

    Objectives:
    1. Correct all spelling and punctuation errors.
    2. Fix all grammatical and syntactical issues.
    3. Keep the meaning and tone as close as possible to the original text.

    Style Guidelines
      - Use active voice and clear, simple sentences.
      - Contractions (e.g., “we’re,” “they’re”) are allowed.
      - Maintain a conversational, professional tone suitable for the workplace.
      - Avoid slang, emojis, and overly formal corporate language.

    Output Rules
     - Return only the edited text, with no explanations, comments, or formatting other than plain text.

    Text to edit:

  translate: "Translate the following text to %s, return only the result: "

baseEndpoint:
  gemini: https://generativelanguage.googleapis.com/v1beta/models/
  ollama: http://localhost:11434/api/generate
  openai: https://api.openai.com/v1/chat/completions
  ```







