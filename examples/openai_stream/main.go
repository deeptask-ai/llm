package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/easymvp/easyllm/internal/providers/openai"
	"github.com/easymvp/easyllm/types"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create OpenAI completion model
	model, err := openai.NewOpenAICompletionModel(apiKey)
	if err != nil {
		log.Fatalf("Failed to create OpenAI model: %v", err)
	}

	ctx := context.Background()

	// Example 1: Basic streaming
	fmt.Println("=== Example 1: Basic Streaming ===")
	req := &types.CompletionRequest{
		Model:        "gpt-4o-mini",
		Instructions: "You are a helpful assistant.",
		Messages: []*types.ModelMessage{
			{
				Role:    types.MessageRoleUser,
				Content: "Count from 1 to 5 and explain each number briefly.",
			},
		},
	}

	stream, err := model.Stream(ctx, req, nil)
	if err != nil {
		log.Fatalf("Stream failed: %v", err)
	}

	fmt.Print("Response: ")
	for chunk := range stream {
		switch c := chunk.(type) {
		case types.StreamTextChunk:
			fmt.Print(c.Text)
		}
	}
	fmt.Println()
}
