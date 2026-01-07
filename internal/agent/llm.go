package agent

import (
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type LLM struct {
	client *openai.Client
	model  string
}

func NewLLM() *LLM {
	key := os.Getenv("GROQ_API_KEY")
	model := os.Getenv("GROQ_MODEL")
	if model == "" {
		model = "llama-3.1-8b-instant"
	}

	baseURL := os.Getenv("GROQ_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.groq.com/openai/v1"
	}

	cfg := openai.DefaultConfig(key)
	cfg.BaseURL = baseURL
	return &LLM{client: openai.NewClientWithConfig(cfg), model: model}
}

func (l *LLM) Complete(ctx context.Context, system, user string) (string, error) {
	resp, err := l.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: l.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Temperature: 0.2,
	})
	if err != nil { return "", err }
	if len(resp.Choices) == 0 { return "", nil }
	return resp.Choices[0].Message.Content, nil
}
