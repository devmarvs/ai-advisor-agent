package rag

import (
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func EmbedTexts(ctx context.Context, texts []string) ([][]float32, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: texts,
		Model: openai.AdaEmbeddingV2, // swap to text-embedding-3-large in production client lib
	})
	if err != nil {
		return nil, err
	}
	out := make([][]float32, len(resp.Data))
	for i, d := range resp.Data {
		out[i] = make([]float32, len(d.Embedding))
		for j, v := range d.Embedding {
			out[i][j] = float32(v)
		}
	}
	return out, nil
}
