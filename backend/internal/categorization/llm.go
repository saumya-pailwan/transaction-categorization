package categorization

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func CategorizeWithLLM(tx Transaction) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	prompt := fmt.Sprintf(`
You are an expert accountant. Categorize the following transaction into one of these categories: 
[Groceries, Dining Out, Transport, Utilities, Rent, Entertainment, Health, Shopping, Income, Transfer, Other].

Transaction:
Description: %s
Merchant: %s
Amount: %.2f

Return ONLY the category name.
`, tx.Description, tx.Merchant, tx.Amount)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
