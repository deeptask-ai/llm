// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package easyllm

import (
	"github.com/easymvp/easyllm/internal/providers/azure"
	"github.com/easymvp/easyllm/internal/providers/claude"
	"github.com/easymvp/easyllm/internal/providers/deepseek"
	"github.com/easymvp/easyllm/internal/providers/gemini"
	"github.com/easymvp/easyllm/internal/providers/openai"
	"github.com/easymvp/easyllm/internal/providers/openrouter"
	"github.com/easymvp/easyllm/types"
)

// NewOpenAIModel creates a new OpenAI model that supports completion, embedding, and image generation
func NewOpenAIModel(opts ...types.ModelOption) (types.CompletionModel, error) {
	return openai.NewOpenAIModel(opts...)
}

// NewAzureOpenAIModel creates a new Azure OpenAI model that supports completion, embedding, and image generation
func NewAzureOpenAIModel(opts ...types.ModelOption) (types.CompletionModel, error) {
	return azure.NewAzureOpenAIModel(opts...)
}

// NewClaudeModel creates a new Claude model that supports completion
func NewClaudeModel(opts ...types.ModelOption) (types.CompletionModel, error) {
	return claude.NewClaudeModel(opts...)
}

// NewDeepSeekModel creates a new DeepSeek model that supports completion
func NewDeepSeekModel(opts ...types.ModelOption) (types.CompletionModel, error) {
	return deepseek.NewDeepSeekModel(opts...)
}

// NewGeminiModel creates a new Gemini model that supports completion
func NewGeminiModel(opts ...types.ModelOption) (types.CompletionModel, error) {
	return gemini.NewGeminiModel(opts...)
}

// NewOpenRouterModel creates a new OpenRouter model that supports completion
func NewOpenRouterModel(opts ...types.ModelOption) (types.CompletionModel, error) {
	return openrouter.NewOpenRouterModel(opts...)
}
