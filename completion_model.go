package llm

import (
	"context"
)

// CompletionModel defines the interface for text completion operations
type CompletionModel interface {
	// StreamComplete generates streaming content from the model
	StreamComplete(ctx context.Context, req *CompletionRequest) (StreamCompletionResponse, error)
	// Complete generates complete content from the model
	Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
}

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

type CompletionRequest struct {
	Instructions string
	Messages     []*ModelMessage
}

// StreamModelResponse represents a stream of API chunks
type StreamCompletionResponse <-chan StreamChunk

type CompletionResponse struct {
	Output string `json:"output"`
	Usage  *TokenUsage
	Cost   *float64
}

// CompletionOption is a functional option for configuring completion requests
type CompletionOption func(*CompletionOptions)

// CompletionOptions contains configuration options for completion requests
type CompletionOptions struct {
	Temperature       *float64
	TopP              *float64
	MaxTokens         *int
	PresencePenalty   *float64
	FrequencyPenalty  *float64
	Seed              *int64
	ReasoningEffort   *ReasoningEffort
	Stop              []string
	ResponseFormat    *ResponseFormat
	JSONSchema        any
	WithCost          *bool
	WithUsage         *bool
	MaxOutputTokens   *int
	ParallelToolCalls *bool
	TopLogprobs       *int
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

// WithMaxOutputTokens sets the maximum number of output tokens to generate
func WithMaxOutputTokens(maxTokens int) CompletionOption {
	return func(o *CompletionOptions) {
		o.MaxOutputTokens = &maxTokens
	}
}

func WithParallelToolCalls(enabled bool) CompletionOption {
	return func(o *CompletionOptions) {
		o.ParallelToolCalls = &enabled
	}
}

func WithTopLogprobs(topLogprobs int) CompletionOption {
	return func(o *CompletionOptions) {
		o.TopLogprobs = &topLogprobs
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
