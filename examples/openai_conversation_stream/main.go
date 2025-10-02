package main

import (
	"context"
	"fmt"
	"github.com/easymvp/easyllm"
	"github.com/easymvp/easyllm/types/conversation"
	"log"
	"os"

	"github.com/easymvp/easyllm/types"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create OpenAI completion model
	model, err := easyllm.NewOpenAIConversationModel(
		types.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create OpenAI model: %v", err)
	}

	ctx := context.Background()

	// Example 1: Basic streaming
	fmt.Println("=== Example 1: Basic Streaming ===")
	req := &conversation.ConversationRequest{
		Model: "o4-mini",
		Input: "You are a helpful assistant.",
		Options: []conversation.ResponseOption{
			conversation.WithReasoningEffort(conversation.ReasoningEffortLow),
			conversation.WithReasoningSummary("concise"),
		},
	}

	stream, err := model.StreamResponse(ctx, req, nil)
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
