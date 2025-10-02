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
	"github.com/easymvp/easyllm/types/completion"
	"github.com/easymvp/easyllm/types/conversation"
	"github.com/easymvp/easyllm/types/embedding"
	"github.com/easymvp/easyllm/types/image"
)

// NewOpenAIModel creates a new OpenAI model that supports completion, embedding, and image generation
func NewOpenAIModel(opts ...types.ModelOption) (completion.CompletionModel, error) {
	return openai.NewOpenAIModel(opts...)
}

func NewOpenAIConversationModel(opts ...types.ModelOption) (conversation.ConversationModel, error) {
	return openai.NewOpenAIModel(opts...)
}

func NewOpenAIEmbeddingModel(opts ...types.ModelOption) (embedding.EmbeddingModel, error) {
	return openai.NewOpenAIModel(opts...)
}

func NewOpenAIImageModel(opts ...types.ModelOption) (image.ImageModel, error) {
	return openai.NewOpenAIModel(opts...)
}

// NewAzureOpenAIModel creates a new Azure OpenAI model that supports completion, embedding, and image generation
func NewAzureOpenAIModel(opts ...types.ModelOption) (completion.CompletionModel, error) {
	return azure.NewAzureOpenAIModel(opts...)
}

// NewClaudeModel creates a new Claude model that supports completion
func NewClaudeModel(opts ...types.ModelOption) (completion.CompletionModel, error) {
	return claude.NewClaudeModel(opts...)
}

// NewDeepSeekModel creates a new DeepSeek model that supports completion
func NewDeepSeekModel(opts ...types.ModelOption) (completion.CompletionModel, error) {
	return deepseek.NewDeepSeekModel(opts...)
}

// NewGeminiModel creates a new Gemini model that supports completion
func NewGeminiModel(opts ...types.ModelOption) (completion.CompletionModel, error) {
	return gemini.NewGeminiModel(opts...)
}

// NewOpenRouterModel creates a new OpenRouter model that supports completion
func NewOpenRouterModel(opts ...types.ModelOption) (completion.CompletionModel, error) {
	return openrouter.NewOpenRouterModel(opts...)
}
