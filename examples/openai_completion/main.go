package main

import (
	"context"
	"fmt"
	"github.com/easymvp/easyllm"
	"github.com/easymvp/easyllm/types/completion"
	"log"
	"os"

	"github.com/easymvp/easyllm/types"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create OpenAI model client
	model, err := easyllm.NewOpenAIModel(
		types.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create OpenAI model: %v", err)
	}

	ctx := context.Background()

	// Example 1: Basic completion
	fmt.Println("=== Example 1: Basic Completion ===")
	req := &completion.CompletionRequest{
		Model:        "gpt-4o-mini",
		Instructions: "You are a helpful assistant.",
		Messages: []*types.ModelMessage{
			{
				Role:    types.MessageRoleUser,
				Content: "What is the capital of France?",
			},
		},
	}

	resp, err := model.Complete(ctx, req, nil)
	if err != nil {
		log.Fatalf("Completion failed: %v", err)
	}

	fmt.Printf("Response: %s\n", resp.Output)
}
