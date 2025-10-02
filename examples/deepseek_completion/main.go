package main

import (
	"context"
	"fmt"
	"github.com/easymvp/easyllm"
	"log"
	"os"

	"github.com/easymvp/easyllm/types"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable is required")
	}

	// Create DeepSeek model client
	model, err := easyllm.NewDeepSeekModel(
		types.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create DeepSeek model: %v", err)
	}

	ctx := context.Background()

	// Example 1: Basic completion
	fmt.Println("=== Example 1: Basic Completion ===")
	req := &types.CompletionRequest{
		Model:        "deepseek-chat",
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
