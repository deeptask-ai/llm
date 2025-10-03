// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"github.com/deeptask-ai/llm"
	"github.com/deeptask-ai/llm/internal/providers/deepseek"
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
	model, err := deepseek.NewDeepSeekModel(
		llm.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create DeepSeek model: %v", err)
	}

	ctx := context.Background()

	// Example 1: Basic completion
	fmt.Println("=== Example 1: Basic Completion ===")
	req := &llm.CompletionRequest{
		Model:        "deepseek-chat",
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
