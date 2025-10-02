// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package easyllm

import (
	"context"
	"time"
)

type ModelMediaType string

const (
	ModelMediaTypeText  ModelMediaType = "text"
	ModelMediaTypeImage ModelMediaType = "image"
	ModelMediaTypeAudio ModelMediaType = "audio"
	ModelMediaTypeVideo ModelMediaType = "video"
)

// ModelInfo contains metadata about a specific model including its ID, name, and pricing information
type ModelInfo struct {
	ID              string           `json:"id"`              // Unique identifier for the model
	Name            string           `json:"name"`            // Human-readable name for the model
	Pricing         ModelPricing     `json:"pricing"`         // Pricing information for different operations
	Reasoning       bool             `json:"reasoning"`       // Whether the model supports reasoning
	Embedding       bool             `json:"embedding"`       // Whether the model supports embeddings
	Input           []ModelMediaType `json:"input"`           // Input type (e.g., "text", "image")
	Output          []ModelMediaType `json:"output"`          // Output type (e.g., "text", "image")
	ContextWindow   int              `json:"contextWindow"`   // Maximum context window size in tokens
	MaxOutputTokens int              `json:"maxOutputTokens"` // Maximum output tokens
	UpdatedAt       time.Time        `json:"updatedAt"`       // Last updated time
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

// BaseModel defines the base interface that all model providers must implement
type BaseModel interface {
	// Name returns the provider name (e.g., "openai", "claude", "gemini")
	Name() string
	// SupportedModels returns a list of all models supported by this provider
	SupportedModels() []*ModelInfo
}

// CompletionModel defines the interface for text completion operations
type CompletionModel interface {
	BaseModel
	// Stream generates streaming content from the model
	Stream(ctx context.Context, req *CompletionRequest, tools []ModelTool) (StreamCompletionResponse, error)
	// Complete generates complete content from the model
	Complete(ctx context.Context, req *CompletionRequest, tools []ModelTool) (*CompletionResponse, error)
}

// EmbeddingModel defines the interface for generating text embeddings
type EmbeddingModel interface {
	BaseModel
	// GenerateEmbeddings generates text embeddings
	GenerateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
}

// ImageModel defines the interface for image generation operations
type ImageModel interface {
	BaseModel
	// GenerateImage generates images from text prompts
	GenerateImage(ctx context.Context, req *ImageRequest) (*ImageResponse, error)
}
