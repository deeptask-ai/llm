// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/openai"
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
	model, err := openai.NewOpenAIConversationModel(
		llm.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create OpenAI model: %v", err)
	}

	ctx := context.Background()

	// Example 1: Basic streaming
	fmt.Println("=== Example 1: Basic Streaming ===")
	req := &llm.ConversationRequest{
		Model: "o4-mini",
		Input: "You are a helpful assistant.",
		Options: []llm.ResponseOption{
			llm.WithReasoningSummary("concise"),
			llm.WithOptions(
				llm.WithReasoningEffort(llm.ReasoningEffortLow),
			),
		},
	}

	stream, err := model.StreamResponse(ctx, req)
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
