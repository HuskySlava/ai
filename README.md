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

#### Rewrite
```bash
ai -r -m openai -p "A sentence to rewrite"
```
#### Translate
```bash
ai -t -m gemini -p "翻訳する行"
```
#### Test
```bash
ai -m ollama -p="Summarize: Go concurrency"
```

Required Environment Variables
```.env
OPENAI_API_KEY=sk-xxx
GEMINI_API_KEY=xyz
```

