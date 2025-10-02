package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/easymvp/easyllm"
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

	// Example: Basic completion with cost tracking enabled
	fmt.Println("=== OpenAI Completion with Cost Tracking ===")
	req := &types.CompletionRequest{
		Model:        "gpt-4o-mini",
		Instructions: "You are a helpful assistant.",
		Messages: []*types.ModelMessage{
			{
				Role:    types.MessageRoleUser,
				Content: "What is the capital of France?",
			},
		},
		Options: []types.CompletionOption{
			types.WithCost(true),
			types.WithUsage(true),
		},
	}

	resp, err := model.Complete(ctx, req, nil)
	if err != nil {
		log.Fatalf("Completion failed: %v", err)
	}

	fmt.Printf("\nResponse: %s\n", resp.Output)

	// Print usage information
	if resp.Usage != nil {
		fmt.Printf("\n=== Usage Information ===\n")
		fmt.Printf("Input tokens: %d\n", resp.Usage.TotalInputTokens)
		fmt.Printf("Output tokens: %d\n", resp.Usage.TotalOutputTokens)
		fmt.Printf("Reasoning tokens: %d\n", resp.Usage.TotalReasoningTokens)
	}

	// Print cost information
	if resp.Cost != nil {
		fmt.Printf("\n=== Cost Information ===\n")
		fmt.Printf("Total cost: $%.6f\n", *resp.Cost)
	} else {
		fmt.Println("\n=== Cost Information ===")
		fmt.Println("Cost information not available")
	}
}
