package categorization

import (
	"context"
	"fmt"

	"autonomoustx/internal/db"

	"github.com/pgvector/pgvector-go"
	openai "github.com/sashabaranov/go-openai"
)

func GenerateEmbedding(text string) (pgvector.Vector, error) {
	// Use shared client
	client := LLMClient

	resp, err := client.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: []string{text},
			Model: openai.AdaEmbeddingV2,
		},
	)

	if err != nil {
		return pgvector.Vector{}, err
	}

	if len(resp.Data) == 0 {
		return pgvector.Vector{}, fmt.Errorf("no embedding returned")
	}

	// Convert []float32 to pgvector.Vector
	// pgvector-go handles []float32 directly usually, but let's check the type
	return pgvector.NewVector(resp.Data[0].Embedding), nil
}

func SearchSimilar(embedding pgvector.Vector) (string, float64, bool, error) {
	ctx := context.Background()

	// Find the most similar transaction that has a category
	// We use cosine distance (<=>) operator.
	// Note: 1 - cosine_distance = cosine_similarity

	var category string
	var distance float64

	err := db.Pool.QueryRow(ctx, `
		SELECT category, embedding <=> $1 as distance
		FROM transactions
		WHERE category IS NOT NULL AND category != ''
		ORDER BY distance ASC
		LIMIT 1
	`, embedding).Scan(&category, &distance)

	if err != nil {
		return "", 0, false, err
	}

	similarity := 1 - distance
	return category, similarity, true, nil
}
