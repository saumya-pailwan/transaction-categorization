package categorization

import (
	"fmt"
	"log"
)

type Transaction struct {
	Description string
	Amount      float64
	Merchant    string
}

type CategorizationResult struct {
	Category   string
	Confidence float64
	Method     string // "rule", "vector", "llm"
}

func Categorize(tx Transaction) (CategorizationResult, error) {
	// 1. Try Rules
	ruleCategory, found := MatchRule(tx.Description)
	if found {
		return CategorizationResult{
			Category:   ruleCategory,
			Confidence: 1.0,
			Method:     "rule",
		}, nil
	}

	// 2. Try Vector Search
	embedding, err := GenerateEmbedding(tx.Description)
	if err != nil {
		log.Printf("Error generating embedding: %v", err)
		// Fallback to LLM if embedding fails
	} else {
		vectorCategory, confidence, found, err := SearchSimilar(embedding)
		if err == nil && found && confidence > 0.85 {
			return CategorizationResult{
				Category:   vectorCategory,
				Confidence: confidence,
				Method:     "vector",
			}, nil
		}
	}

	// 3. Fallback to LLM
	llmCategory, err := CategorizeWithLLM(tx)
	if err != nil {
		return CategorizationResult{}, fmt.Errorf("LLM categorization failed: %w", err)
	}

	return CategorizationResult{
		Category:   llmCategory,
		Confidence: 0.7,
		Method:     "llm",
	}, nil
}
