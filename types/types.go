// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"context"
	"encoding/json"
	"fmt"
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

type CompletionRequest struct {
	Instructions string
	Model        string
	Messages     []*ModelMessage
	Config       *ModelConfig
	WithCost     bool
}

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleTool      MessageRole = "tool"
)

type ResponseFormat string

const (
	ResponseFormatJson       ResponseFormat = "json"
	ResponseFormatJsonSchema ResponseFormat = "json_schema"
)

type ReasoningEffort string

const (
	ReasoningEffortLow    ReasoningEffort = "low"
	ReasoningEffortMedium ReasoningEffort = "medium"
	ReasoningEffortHigh   ReasoningEffort = "high"
)

// ModelConfig contains all OpenAI completion options
type ModelConfig struct {
	Temperature      float64         `json:"temperature,omitempty"`
	TopP             float64         `json:"topP,omitempty"`
	MaxTokens        int             `json:"maxTokens,omitempty"`
	PresencePenalty  float64         `json:"presencePenalty,omitempty"`
	FrequencyPenalty float64         `json:"frequencyPenalty,omitempty"`
	Seed             int64           `json:"seed,omitempty"`
	ReasoningEffort  ReasoningEffort `json:"reasoningEffort,omitempty"`
	Stop             []string        `json:"stop,omitempty"`
	ResponseFormat   ResponseFormat  `json:"responseFormat,omitempty"`
	JSONSchema       any             `json:"jsonSchema,omitempty"`
}

type EmbeddingEncodingFormat string

const (
	EmbeddingEncodingFormatFloat  EmbeddingEncodingFormat = "float"
	EmbeddingEncodingFormatBase64 EmbeddingEncodingFormat = "base64"
)

type EmbeddingRequest struct {
	Model    string                `json:"model"`
	Contents []string              `json:"contents"`
	Config   *EmbeddingModelConfig `json:"config,omitempty"`
}

type EmbeddingModelConfig struct {
	Dimensions     int64                   `json:"dimensions,omitempty"`
	EncodingFormat EmbeddingEncodingFormat `json:"encoding_format,omitempty"`
}

type ImageRequest struct {
	Model        string            `json:"model"`
	Instructions string            `json:"instructions"`
	Artifacts    []*ModelArtifact  `json:"artifacts"`
	Config       *ImageModelConfig `json:"config,omitempty"`
}

type ImageModelConfig struct {
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	Style          string `json:"style,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
}

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

// StreamChunkType defines the type of chunk in the API stream
type StreamChunkType string

const (
	TextChunkType      StreamChunkType = "text"
	ReasoningChunkType StreamChunkType = "reasoning"
	UsageChunkType     StreamChunkType = "usage"
)

// StreamChunk is the interface for all types of chunks in the API stream
type StreamChunk interface {
	Type() StreamChunkType

	String() string
}

// StreamTextChunk represents a text chunk in the API stream
type StreamTextChunk struct {
	// Text contains the actual text content
	Text string `json:"text"`
}

// Type returns the type of the chunk
func (c StreamTextChunk) Type() StreamChunkType {
	return TextChunkType
}
func (c StreamTextChunk) String() string {
	return c.Text
}

// StreamReasoningChunk represents a reasoning chunk in the API stream
type StreamReasoningChunk struct {
	// Reasoning contains the reasoning text
	Reasoning string `json:"reasoning"`
}

// Type returns the type of the chunk
func (c StreamReasoningChunk) Type() StreamChunkType {
	return ReasoningChunkType
}
func (c StreamReasoningChunk) String() string {
	return c.Reasoning
}

// StreamUsageChunk represents a outputExample information chunk in the API stream
type StreamUsageChunk struct {
	Usage *TokenUsage
	Cost  *float64
}

// Type returns the type of the chunk
func (c StreamUsageChunk) Type() StreamChunkType {
	return UsageChunkType
}
func (c StreamUsageChunk) String() string {
	jsonBytes, err := json.Marshal(c.Usage)
	if err != nil {
		return "usage: {}"
	}
	return fmt.Sprintf("usage: %s", string(jsonBytes))
}

// StreamModelResponse represents a stream of API chunks
type StreamCompletionResponse <-chan StreamChunk

type CompletionResponse struct {
	Output string `json:"output"`
	Usage  *TokenUsage
	Cost   *float64
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

type Embedding struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
	Object    string    `json:"object"`
}

type EmbeddingResponse struct {
	Embeddings []Embedding `json:"embeddings"`
	Usage      *TokenUsage `json:"usage,omitempty"`
	Cost       *float64    `json:"cost,omitempty"`
}

type ImageResponse struct {
	Output []byte      `json:"output"`
	Usage  *TokenUsage `json:"usage,omitempty"`
	Cost   *float64    `json:"cost,omitempty"`
}

type ModelTool interface {
	Name() string

	Description() string

	InputSchema() any

	OutputSchema() any

	Run(ctx context.Context, input any) (any, error)

	Usage() string
}

type ToolCall struct {
	ID           string
	Name         string
	Input        any
	Output       any
	ErrorMessage *string
}
