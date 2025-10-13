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
	key := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")
	if model == "":
		model = "gpt-4o-mini"
	}
	return &LLM{client: openai.NewClient(key), model: model}
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
