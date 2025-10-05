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
	req := &llm.ConversationRequest{
		Input: "You are a helpful assistant.",
	}
	model, err := provider.NewConversationModel("gpt-4o-mini", llm.WithReasoningSummary("concise"),
		llm.WithOptions(
			llm.WithReasoningEffort(llm.ReasoningEffortLow),
		))
	if err != nil {
		log.Fatalf("Failed to create DeepSeek model: %v", err)
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
