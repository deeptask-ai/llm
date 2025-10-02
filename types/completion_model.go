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
	Config       *ModelConfig
	WithCost     bool
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
