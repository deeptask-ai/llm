package llm

import (
	"context"
)

// ConversationModel defines the interface for conversation operations using the responses API
type ConversationModel interface {
	BaseModel
	// StreamResponse generates streaming content from the model using the responses API
	StreamResponse(ctx context.Context, req *ConversationRequest) (StreamConversationResponse, error)
	// Response generates complete content from the model using the responses API
	Response(ctx context.Context, req *ConversationRequest) (*ConversationResponse, error)
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
	Usage  *TokenUsage
	Cost   *float64
}

// StreamConversationResponse represents a stream of response chunks
type StreamConversationResponse <-chan StreamChunk

// ResponseOption is a functional option for configuring response requests
type ResponseOption func(*ResponseOptions)

// ResponseOptions contains configuration options for response requests
// It extends ResponseOptions with additional response-specific options
type ResponseOptions struct {
	ReasoningSummary   *string
	Store              *bool
	PreviousResponseId *string
	CompletionOptions  *CompletionOptions
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

func WithOptions(opts ...CompletionOption) ResponseOption {
	return func(o *ResponseOptions) {
		o.CompletionOptions = ApplyCompletionOptions(opts)
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
