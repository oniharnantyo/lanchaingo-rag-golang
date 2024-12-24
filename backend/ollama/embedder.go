package ollama

import (
	"context"
	"github.com/tmc/langchaingo/embeddings"
	"strings"
)

type Ollama struct {
	client embeddings.EmbedderClient
}

func NewOllama(client embeddings.EmbedderClient) embeddings.Embedder {
	return &Ollama{client}
}

func (o *Ollama) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	return o.client.CreateEmbedding(ctx, texts)
}

func (o *Ollama) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	text = strings.ReplaceAll(text, "\n", " ")

	emb, err := o.client.CreateEmbedding(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	return emb[0], nil
}
