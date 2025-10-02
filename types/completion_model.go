package types

import (
	"context"
	"encoding/json"
	"fmt"
)

// CompletionModel defines the interface for text completion operations
type CompletionModel interface {
	BaseModel
	// Stream generates streaming content from the model
	Stream(ctx context.Context, req *CompletionRequest, tools []ModelTool) (StreamCompletionResponse, error)
	// Complete generates complete content from the model
	Complete(ctx context.Context, req *CompletionRequest, tools []ModelTool) (*CompletionResponse, error)
}

type CompletionRequest struct {
	Instructions string
	Model        string
	Messages     []*ModelMessage
	Options      []CompletionOption
}

// StreamModelResponse represents a stream of API chunks
type StreamCompletionResponse <-chan StreamChunk

type CompletionResponse struct {
	Output string `json:"output"`
	Usage  *TokenUsage
	Cost   *float64
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

// CompletionOption is a functional option for configuring completion requests
type CompletionOption func(*CompletionOptions)

// CompletionOptions contains configuration options for completion requests
type CompletionOptions struct {
	Temperature      *float64
	TopP             *float64
	MaxTokens        *int
	PresencePenalty  *float64
	FrequencyPenalty *float64
	Seed             *int64
	ReasoningEffort  *ReasoningEffort
	Stop             []string
	ResponseFormat   *ResponseFormat
	JSONSchema       any
	WithCost         *bool
	WithUsage        *bool
}

// WithTemperature sets the temperature for sampling
func WithTemperature(temperature float64) CompletionOption {
	return func(o *CompletionOptions) {
		o.Temperature = &temperature
	}
}

// WithTopP sets the top-p for nucleus sampling
func WithTopP(topP float64) CompletionOption {
	return func(o *CompletionOptions) {
		o.TopP = &topP
	}
}

// WithMaxTokens sets the maximum number of tokens to generate
func WithMaxTokens(maxTokens int) CompletionOption {
	return func(o *CompletionOptions) {
		o.MaxTokens = &maxTokens
	}
}

// WithPresencePenalty sets the presence penalty
func WithPresencePenalty(presencePenalty float64) CompletionOption {
	return func(o *CompletionOptions) {
		o.PresencePenalty = &presencePenalty
	}
}

// WithFrequencyPenalty sets the frequency penalty
func WithFrequencyPenalty(frequencyPenalty float64) CompletionOption {
	return func(o *CompletionOptions) {
		o.FrequencyPenalty = &frequencyPenalty
	}
}

// WithSeed sets the random seed for deterministic sampling
func WithSeed(seed int64) CompletionOption {
	return func(o *CompletionOptions) {
		o.Seed = &seed
	}
}

// WithReasoningEffort sets the reasoning effort level
func WithReasoningEffort(effort ReasoningEffort) CompletionOption {
	return func(o *CompletionOptions) {
		o.ReasoningEffort = &effort
	}
}

// WithStop sets the stop sequences
func WithStop(stop []string) CompletionOption {
	return func(o *CompletionOptions) {
		o.Stop = stop
	}
}

// WithResponseFormat sets the response format
func WithResponseFormat(format ResponseFormat) CompletionOption {
	return func(o *CompletionOptions) {
		o.ResponseFormat = &format
	}
}

// WithJSONSchema sets the JSON schema for structured output
func WithJSONSchema(schema any) CompletionOption {
	return func(o *CompletionOptions) {
		o.JSONSchema = schema
	}
}

// WithCost enables cost calculation in the response
func WithCost(enabled bool) CompletionOption {
	return func(o *CompletionOptions) {
		o.WithCost = &enabled
	}
}

// WithUsage enables usage information in the response
func WithUsage(enabled bool) CompletionOption {
	return func(o *CompletionOptions) {
		o.WithUsage = &enabled
	}
}

// ApplyCompletionOptions applies all options to create a CompletionOptions struct
func ApplyCompletionOptions(opts []CompletionOption) *CompletionOptions {
	options := &CompletionOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
