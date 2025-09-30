// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package llmclient

import (
	"encoding/json"
	"fmt"
)

type MessageArtifact struct {
	ID          string            `json:"id"`
	Name        string            `json:"name" binding:"required"`
	ContentType string            `json:"contenttype" binding:"required"`
	Description string            `json:"description"`
	Slug        string            `json:"slug" binding:"required"`
	Content     []byte            `json:"content"`
	Metadata    map[string]string `json:"metadata"`
}

// Message represents a communication unit with a specific role and content.
// It is used for constructing messages in agent-based communication.
type Message struct {
	ID        string             `json:"id"`
	Role      MessageRole        `json:"role"`
	Content   string             `json:"content"`
	Artifacts []*MessageArtifact `json:"artifacts"`
	ToolCall  *ToolCall          `json:"toolCall"`
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
type StreamModelResponse <-chan StreamChunk

type ModelResponse struct {
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
