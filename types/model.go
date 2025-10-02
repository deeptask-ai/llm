// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package types

import (
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

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleTool      MessageRole = "tool"
)

type EmbeddingEncodingFormat string

const (
	EmbeddingEncodingFormatFloat  EmbeddingEncodingFormat = "float"
	EmbeddingEncodingFormatBase64 EmbeddingEncodingFormat = "base64"
)

type ModelArtifact struct {
	ID          string            `json:"id"`
	Name        string            `json:"name" binding:"required"`
	ContentType string            `json:"contentType" binding:"required"`
	Description string            `json:"description"`
	Content     []byte            `json:"content"`
	Metadata    map[string]string `json:"metadata"`
}

type ModelMessage struct {
	Role      MessageRole      `json:"role"`
	Content   string           `json:"content"`
	Artifacts []*ModelArtifact `json:"artifacts"`
	ToolCall  *ToolCall        `json:"toolCall"`
}

type TokenUsage struct {
	TotalInputTokens      int64 `json:"totalInputTokens"`
	TotalOutputTokens     int64 `json:"totalOutputTokens"`
	TotalReasoningTokens  int64 `json:"totalReasoningTokens"`
	TotalImages           int   `json:"totalImages"`
	TotalWebSearches      int   `json:"totalWebSearches"`
	TotalRequests         int   `json:"totalRequests"`
	TotalCacheReadTokens  int64 `json:"totalCacheReadTokens"`
	TotalCacheWriteTokens int64 `json:"totalCacheWriteTokens"`
}

func (s *TokenUsage) Append(usage *TokenUsage) {
	s.TotalInputTokens += usage.TotalInputTokens
	s.TotalOutputTokens += usage.TotalOutputTokens
	s.TotalReasoningTokens += usage.TotalReasoningTokens
	s.TotalImages += usage.TotalImages
	s.TotalWebSearches += usage.TotalWebSearches
	s.TotalRequests += usage.TotalRequests
	s.TotalCacheReadTokens += usage.TotalCacheReadTokens
	s.TotalCacheWriteTokens += usage.TotalCacheWriteTokens
}
