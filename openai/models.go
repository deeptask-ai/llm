// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package openai

import (
	"github.com/easymvp-ai/llm"
	"github.com/easymvp-ai/llm/internal/providers/openai"
)

// NewOpenAIModel creates a new OpenAI model that supports llm, embedding, and image generation
func NewOpenAIModel(opts ...llm.ModelOption) (llm.CompletionModel, error) {
	return openai.NewOpenAIModel(opts...)
}

func NewOpenAIConversationModel(opts ...llm.ModelOption) (llm.ConversationModel, error) {
	return openai.NewOpenAIModel(opts...)
}

func NewOpenAIEmbeddingModel(opts ...llm.ModelOption) (llm.EmbeddingModel, error) {
	return openai.NewOpenAIModel(opts...)
}

func NewOpenAIImageModel(opts ...llm.ModelOption) (llm.ImageModel, error) {
	return openai.NewOpenAIModel(opts...)
}
