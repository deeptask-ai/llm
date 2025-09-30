// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package llmclient

type ModelRequest struct {
	Instructions string
	Model        string
	Messages     []*Message
	Config       *ModelConfig
	Tools        []Tool
	Metadata     map[string]string
	Cost         bool
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
	Metadata map[string]string     `json:"metadata,omitempty"`
}

type EmbeddingModelConfig struct {
	Dimensions     int64                   `json:"dimensions,omitempty"`
	EncodingFormat EmbeddingEncodingFormat `json:"encoding_format,omitempty"`
}
