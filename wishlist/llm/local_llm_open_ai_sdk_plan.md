# Local LLM (Ollama) + OpenAI SDK (Go)

## Goal
Use a local LLM during development to avoid costs, while using OpenAI (e.g. `gpt-5-nano`) in production — **without changing application code**.

---

## Architecture
- Single OpenAI-compatible client
- Switch provider via environment variables
- Ollama provides a local OpenAI-compatible API

```
Go App (OpenAI SDK)
        |
        |  /v1/chat/completions
        |
   Base URL (env)
   ├─ http://localhost:11434/v1  → Ollama (local)
   └─ https://api.openai.com/v1  → OpenAI (prod)
```

---

## Step 1: Install & Run Ollama

```bash
brew install ollama   # or apt / windows installer
ollama pull llama3.1
ollama serve
```

API available at:
```
http://localhost:11434/v1
```

---

## Step 2: Environment Configuration

### Local (dev)
```env
OPENAI_BASE_URL=http://localhost:11434/v1
OPENAI_API_KEY=ollama
OPENAI_MODEL=llama3.1
```

### Production
```env
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-5-nano
```

---

## Step 3: Go Client Setup (OpenAI SDK)

```go
client := openai.NewClient(
  openai.WithBaseURL(os.Getenv("OPENAI_BASE_URL")),
  openai.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
)
```

---

## Step 4: Usage (Same for Local & Prod)

```go
resp, err := client.Chat.Completions.Create(ctx, openai.ChatCompletionRequest{
  Model: os.Getenv("OPENAI_MODEL"),
  Messages: []openai.ChatCompletionMessage{
    {Role: "system", Content: "You are a helpful recipe assistant."},
    {Role: "user", Content: "Suggest a vegetarian pasta recipe."},
  },
})
```

---

## Step 5: Prompt Guidelines
- Use explicit system instructions
- Avoid relying on advanced tool/function calling locally
- Keep output format simple (markdown or plain text)

---

## Step 6: Optional Improvements
- Cache responses (Redis / in-memory)
- Add provider health check (OpenAI ↔ Ollama)
- Use RAG before LLM to reduce calls

---

## Result
- One SDK
- One code path
- Zero local inference cost
- Production-grade LLM in prod

