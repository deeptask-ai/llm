// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package easyllm

import (
	"context"
)

// ModelInfo contains metadata about a specific model including its ID, name, and pricing information
type ModelInfo struct {
	ID      string       `json:"id"`      // Unique identifier for the model
	Name    string       `json:"name"`    // Human-readable name for the model
	Pricing ModelPricing `json:"pricing"` // Pricing information for different operations
}

// ModelPricing contains pricing information for various model operations
type ModelPricing struct {
	Prompt            string `json:"prompt"`            // Price per million input tokens
	Completion        string `json:"completion"`        // Price per million output tokens
	Request           string `json:"request"`           // Price per request (if applicable)
	Image             string `json:"image"`             // Price per image generation
	WebSearch         string `json:"webSearch"`         // Price per web search operation
	InternalReasoning string `json:"internalReasoning"` // Price per million reasoning tokens
	InputCacheRead    string `json:"inputCacheRead"`    // Price per million cached input tokens read
	InputCacheWrite   string `json:"inputCacheWrite"`   // Price per million cached input tokens written
}

// Model defines the interface that all LLM providers must implement
type Model interface {
	// Name returns the provider name (e.g., "openai", "claude", "gemini")
	Name() string

	// SupportedModels returns a list of all models supported by this provider
	SupportedModels() []*ModelInfo

	// GenerateContentStream generates streaming content from the model
	GenerateContentStream(ctx context.Context, req *ModelRequest) (StreamModelResponse, error)

	// GenerateContent generates complete content from the model
	GenerateContent(ctx context.Context, req *ModelRequest) (*ModelResponse, error)

	// GenerateEmbeddings generates text embeddings (may not be supported by all providers)
	GenerateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)

	// GenerateImage generates images from text prompts (may not be supported by all providers)
	GenerateImage(ctx context.Context, req *ImageRequest) (*ImageResponse, error)
}
