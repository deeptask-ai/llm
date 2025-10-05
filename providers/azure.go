// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0
package providers

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/azure"
)

// NewAzureOpenAIModelProvider creates a new Azure OpenAI model that supports llm, embedding, and image generation
func NewAzureOpenAIModelProvider(opts ...llm.ModelOption) (llm.ModelProvider, error) {
	return azure.NewAzureOpenAIModelProvider(opts...)
}
