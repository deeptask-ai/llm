package conversation

import (
	"context"
	"github.com/easymvp/easyllm/types"
)

// ConversationModel defines the interface for conversation operations using the responses API
type ConversationModel interface {
	types.BaseModel
	// StreamResponse generates streaming content from the model using the responses API
	StreamResponse(ctx context.Context, req *ConversationRequest, tools []types.ModelTool) (StreamConversationResponse, error)
	// Response generates complete content from the model using the responses API
	Response(ctx context.Context, req *ConversationRequest, tools []types.ModelTool) (*ConversationResponse, error)
}

// ConversationRequest represents a request to the conversation/responses API
type ConversationRequest struct {
	Input   string
	Model   string
	Options []ResponseOption
}

// ConversationResponse represents a complete response from the conversation/responses API
type ConversationResponse struct {
	Output string `json:"output"`
	Usage  *types.TokenUsage
	Cost   *float64
}

// StreamConversationResponse represents a stream of response chunks
type StreamConversationResponse <-chan types.StreamChunk

// ResponseOption is a functional option for configuring response requests
type ResponseOption func(*ResponseOptions)

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

// ResponseOptions contains configuration options for response requests
// It extends ResponseOptions with additional response-specific options
type ResponseOptions struct {
	Temperature        *float64
	TopP               *float64
	MaxTokens          *int
	PresencePenalty    *float64
	FrequencyPenalty   *float64
	Seed               *int64
	ReasoningEffort    *ReasoningEffort
	Stop               []string
	ResponseFormat     *ResponseFormat
	JSONSchema         any
	WithCost           *bool
	WithUsage          *bool
	MaxOutputTokens    *int
	ParallelToolCalls  *bool
	TopLogprobs        *int
	ReasoningSummary   *string
	Store              *bool
	PreviousResponseId *string
}

// WithReasoningSummary sets the reasoning summary type (auto, concise, or detailed)
func WithReasoningSummary(summary string) ResponseOption {
	return func(o *ResponseOptions) {
		o.ReasoningSummary = &summary
	}
}

func WithPreviousResponseId(previousResponseId string) ResponseOption {
	return func(o *ResponseOptions) {
		o.PreviousResponseId = &previousResponseId
	}
}

// WithStore sets whether to store the response
func WithStore(enabled bool) ResponseOption {
	return func(o *ResponseOptions) {
		o.Store = &enabled
	}
}

// WithTemperature sets the temperature for sampling
func WithTemperature(temperature float64) ResponseOption {
	return func(o *ResponseOptions) {
		o.Temperature = &temperature
	}
}

// WithTopP sets the top-p for nucleus sampling
func WithTopP(topP float64) ResponseOption {
	return func(o *ResponseOptions) {
		o.TopP = &topP
	}
}

// WithMaxTokens sets the maximum number of tokens to generate
func WithMaxTokens(maxTokens int) ResponseOption {
	return func(o *ResponseOptions) {
		o.MaxTokens = &maxTokens
	}
}

// WithPresencePenalty sets the presence penalty
func WithPresencePenalty(presencePenalty float64) ResponseOption {
	return func(o *ResponseOptions) {
		o.PresencePenalty = &presencePenalty
	}
}

// WithFrequencyPenalty sets the frequency penalty
func WithFrequencyPenalty(frequencyPenalty float64) ResponseOption {
	return func(o *ResponseOptions) {
		o.FrequencyPenalty = &frequencyPenalty
	}
}

// WithSeed sets the random seed for deterministic sampling
func WithSeed(seed int64) ResponseOption {
	return func(o *ResponseOptions) {
		o.Seed = &seed
	}
}

// WithReasoningEffort sets the reasoning effort level
func WithReasoningEffort(effort ReasoningEffort) ResponseOption {
	return func(o *ResponseOptions) {
		o.ReasoningEffort = &effort
	}
}

// WithStop sets the stop sequences
func WithStop(stop []string) ResponseOption {
	return func(o *ResponseOptions) {
		o.Stop = stop
	}
}

// WithResponseFormat sets the response format
func WithResponseFormat(format ResponseFormat) ResponseOption {
	return func(o *ResponseOptions) {
		o.ResponseFormat = &format
	}
}

// WithJSONSchema sets the JSON schema for structured output
func WithJSONSchema(schema any) ResponseOption {
	return func(o *ResponseOptions) {
		o.JSONSchema = schema
	}
}

// WithCost enables cost calculation in the response
func WithCost(enabled bool) ResponseOption {
	return func(o *ResponseOptions) {
		o.WithCost = &enabled
	}
}

// WithUsage enables usage information in the response
func WithUsage(enabled bool) ResponseOption {
	return func(o *ResponseOptions) {
		o.WithUsage = &enabled
	}
}

// WithMaxOutputTokens sets the maximum number of output tokens to generate
func WithMaxOutputTokens(maxTokens int) ResponseOption {
	return func(o *ResponseOptions) {
		o.MaxOutputTokens = &maxTokens
	}
}

func WithParallelToolCalls(enabled bool) ResponseOption {
	return func(o *ResponseOptions) {
		o.ParallelToolCalls = &enabled
	}
}

func WithTopLogprobs(topLogprobs int) ResponseOption {
	return func(o *ResponseOptions) {
		o.TopLogprobs = &topLogprobs
	}
}

// ApplyResponseOptions applies all options to create a ResponseOptions struct
func ApplyResponseOptions(opts []ResponseOption) *ResponseOptions {
	options := &ResponseOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
