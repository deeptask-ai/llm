// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"github.com/easymvp-ai/llm"
	"github.com/easymvp-ai/llm/openai"
	"log"
	"os"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create OpenAI model client
	model, err := openai.NewOpenAIModel(
		llm.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create OpenAI model: %v", err)
	}

	ctx := context.Background()

	// Example 1: Basic completion
	fmt.Println("=== Example 1: Basic Completion ===")
	req := &llm.CompletionRequest{
		Model:        "gpt-4o-mini",
		Instructions: "You are a helpful assistant.",
		Messages: []*llm.ModelMessage{
			{
				Role:    llm.RoleUser,
				Content: "What is the capital of France?",
			},
		},
	}

	resp, err := model.Complete(ctx, req)
	if err != nil {
		log.Fatalf("Completion failed: %v", err)
	}

	fmt.Printf("Response: %s\n", resp.Output)
}
