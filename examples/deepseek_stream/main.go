// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"github.com/easymvp-ai/llm"
	"github.com/easymvp-ai/llm/internal/providers/deepseek"
	"log"
	"os"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		log.Fatal("DEEPSEEK_API_KEY environment variable is required")
	}

	// Create DeepSeek completion model
	model, err := deepseek.NewDeepSeekModel(
		llm.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create DeepSeek model: %v", err)
	}

	ctx := context.Background()

	// Example 1: Basic streaming
	fmt.Println("=== Example 1: Basic Streaming ===")
	req := &llm.CompletionRequest{
		Model:        "deepseek-chat",
		Instructions: "You are a helpful assistant.",
		Messages: []*llm.ModelMessage{
			{
				Role:    llm.RoleUser,
				Content: "Count from 1 to 5 and explain each number briefly.",
			},
		},
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
