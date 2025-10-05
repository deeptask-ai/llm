// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package openai

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/openai"
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
