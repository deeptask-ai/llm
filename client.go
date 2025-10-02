// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package easyllm

import (
	azureprovider "github.com/easymvp/easyllm/internal/providers/azure"
	claudeprovider "github.com/easymvp/easyllm/internal/providers/claude"
	deepseekprovider "github.com/easymvp/easyllm/internal/providers/deepseek"
	geminiprovider "github.com/easymvp/easyllm/internal/providers/gemini"
	openaiprovider "github.com/easymvp/easyllm/internal/providers/openai"
	openrouterprovider "github.com/easymvp/easyllm/internal/providers/openrouter"
)

// NewOpenAIModel creates a new OpenAI model instance supporting completion, embedding, and image generation
// Note: This will work once provider packages are updated to use their own package names
func NewOpenAIModel(opts ...ModelOption) (CompletionModel, error) {
	return openaiprovider.NewOpenAIModel(opts...)
}

// NewOpenAICompletionModel creates a new OpenAI model instance for completions only
func NewOpenAICompletionModel(opts ...ModelOption) (CompletionModel, error) {
	return openaiprovider.NewOpenAIModel(opts...)
}

// NewOpenAIEmbeddingModel creates a new OpenAI model instance for embeddings only
func NewOpenAIEmbeddingModel(opts ...ModelOption) (EmbeddingModel, error) {
	return openaiprovider.NewOpenAIModel(opts...)
}

// NewOpenAIImageModel creates a new OpenAI model instance for image generation only
func NewOpenAIImageModel(opts ...ModelOption) (ImageModel, error) {
	return openaiprovider.NewOpenAIModel(opts...)
}

// NewClaudeModel creates a new Claude model instance
func NewClaudeModel(opts ...ModelOption) (CompletionModel, error) {
	return claudeprovider.NewClaudeModel(opts...)
}

// NewClaudeCompletionModel creates a new Claude model instance for completions
func NewClaudeCompletionModel(opts ...ModelOption) (CompletionModel, error) {
	return claudeprovider.NewClaudeModel(opts...)
}

// NewGeminiModel creates a new Gemini model instance
func NewGeminiModel(opts ...ModelOption) (CompletionModel, error) {
	return geminiprovider.NewGeminiModel(opts...)
}

// NewGeminiCompletionModel creates a new Gemini model instance for completions
func NewGeminiCompletionModel(opts ...ModelOption) (CompletionModel, error) {
	return geminiprovider.NewGeminiModel(opts...)
}

// NewDeepSeekModel creates a new DeepSeek model instance
func NewDeepSeekModel(opts ...ModelOption) (CompletionModel, error) {
	return deepseekprovider.NewDeepSeekModel(opts...)
}

// NewDeepSeekCompletionModel creates a new DeepSeek model instance for completions
func NewDeepSeekCompletionModel(opts ...ModelOption) (CompletionModel, error) {
	return deepseekprovider.NewDeepSeekModel(opts...)
}

// NewAzureOpenAIModel creates a new Azure OpenAI model instance
func NewAzureOpenAIModel(opts ...ModelOption) (CompletionModel, error) {
	return azureprovider.NewAzureOpenAIModel(opts...)
}

// NewAzureOpenAICompletionModel creates a new Azure OpenAI model instance for completions
func NewAzureOpenAICompletionModel(opts ...ModelOption) (CompletionModel, error) {
	return azureprovider.NewAzureOpenAIModel(opts...)
}

// NewOpenRouterModel creates a new OpenRouter model instance
func NewOpenRouterModel(opts ...ModelOption) (CompletionModel, error) {
	return openrouterprovider.NewOpenRouterModel(opts...)
}

// NewOpenRouterCompletionModel creates a new OpenRouter model instance for completions
func NewOpenRouterCompletionModel(opts ...ModelOption) (CompletionModel, error) {
	return openrouterprovider.NewOpenRouterModel(opts...)
}
