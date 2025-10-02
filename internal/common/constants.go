// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package common

// Provider names
const (
	ProviderOpenAI      = "openai"
	ProviderAzureOpenAI = "azure-openai"
	ProviderClaude      = "claude"
	ProviderGemini      = "gemini"
	ProviderDeepSeek    = "deepseek"
	ProviderOpenRouter  = "openrouter"
)

// Message roles
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
	RoleTool      = "tool"
)

// Response formats
const (
	ResponseFormatText       = "text"
	ResponseFormatJSON       = "json"
	ResponseFormatJSONSchema = "json_schema"
)

// Reasoning effort levels
const (
	ReasoningEffortLow    = "low"
	ReasoningEffortMedium = "medium"
	ReasoningEffortHigh   = "high"
)

// Stream chunk types
const (
	ChunkTypeText      = "text"
	ChunkTypeReasoning = "reasoning"
	ChunkTypeUsage     = "usage"
)

// Embedding encoding formats
const (
	EncodingFormatFloat  = "float"
	EncodingFormatBase64 = "base64"
)

// Image sizes and quality
const (
	ImageSize256  = "256x256"
	ImageSize512  = "512x512"
	ImageSize1024 = "1024x1024"
	ImageSize1792 = "1024x1792"
	ImageSize1536 = "1536x1536"

	ImageQualityStandard = "standard"
	ImageQualityHD       = "hd"
)

// Validation constants
const (
	MinTemperature      = 0.0
	MaxTemperature      = 2.0
	MinTopP             = 0.0
	MaxTopP             = 1.0
	MinPresencePenalty  = -2.0
	MaxPresencePenalty  = 2.0
	MinFrequencyPenalty = -2.0
	MaxFrequencyPenalty = 2.0
	MinMaxTokens        = 1
	MaxMaxTokens        = 1000000
	MinAPIKeyLength     = 10
)
