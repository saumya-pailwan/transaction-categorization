package categorization

import (
	"log"

	openai "github.com/sashabaranov/go-openai"
)

var (
	// LLMClient is the shared OpenAI client instance
	LLMClient *openai.Client
)

// Init initializes the categorization package:
// 1. Sets up the shared OpenAI client
// 2. Loads and compiles rules into memory
func Init(apiKey string) {
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is required")
	}

	// 1. Setup Shared Client
	LLMClient = openai.NewClient(apiKey)

	// 2. Load Rules
	if err := LoadRules(); err != nil {
		log.Printf("Warning: Failed to load rules on startup: %v", err)
	} else {
		log.Println("Categorization rules loaded successfully.")
	}
}
