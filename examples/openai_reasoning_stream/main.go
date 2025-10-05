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
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create OpenAI completion model
	provider, err := providers.NewOpenAIModelProvider(
		llm.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create OpenAI model: %v", err)
	}

	ctx := context.Background()

	// Example 1: Basic streaming
	fmt.Println("=== Example 1: Basic Streaming ===")
	req := &llm.CompletionRequest{
		Instructions: "You are a helpful assistant.",
		Messages: []*llm.ModelMessage{
			{
				Role:    llm.RoleUser,
				Content: "Count from 1 to 5 and explain each number briefly.",
			},
		},
	}
	model, err := provider.NewCompletionModel("o4-mini", llm.WithReasoningEffort(llm.ReasoningEffortLow))
	if err != nil {
		log.Fatalf("Failed to create DeepSeek model: %v", err)
	}
	stream, err := model.StreamComplete(ctx, req)
	if err != nil {
		log.Fatalf("Stream failed: %v", err)
	}

	fmt.Print("Response: ")
	for chunk := range stream {
		switch c := chunk.(type) {
		case llm.StreamTextChunk:
			fmt.Print(c.Text)
		}
	}
	fmt.Println()
}
