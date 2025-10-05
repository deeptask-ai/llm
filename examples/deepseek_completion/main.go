// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/providers"
	"log"
	"os"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable is required")
	}

	// Create DeepSeek model client
	provider, err := providers.NewDeepSeekModelProvider(
		llm.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create DeepSeek model: %v", err)
	}

	ctx := context.Background()

	// Example 1: Basic completion
	fmt.Println("=== Example 1: Basic Completion ===")
	req := &llm.CompletionRequest{
		Instructions: "You are a helpful assistant.",
		Messages: []*llm.ModelMessage{
			{
				Role:    llm.RoleUser,
				Content: "What is the capital of France?",
			},
		},
	}
	model, err := provider.NewCompletionModel("deepseek-chat")
	if err != nil {
		log.Fatalf("Failed to create DeepSeek model: %v", err)
	}
	resp, err := model.Complete(ctx, req)
	if err != nil {
		log.Fatalf("Completion failed: %v", err)
	}

	fmt.Printf("Response: %s\n", resp.Output)
}
