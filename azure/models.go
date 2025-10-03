// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package azure

import (
	"github.com/deeptask-ai/llm"
	"github.com/deeptask-ai/llm/internal/providers/azure"
)

// NewAzureOpenAIModel creates a new Azure OpenAI model that supports llm, embedding, and image generation
func NewAzureOpenAIModel(opts ...llm.ModelOption) (llm.CompletionModel, error) {
	return azure.NewAzureOpenAIModel(opts...)
}
