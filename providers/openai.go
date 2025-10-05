// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/openai"
)

// NewOpenAIModelProvider creates a new OpenAI model that supports llm, embedding, and image generation
func NewOpenAIModelProvider(opts ...llm.ModelOption) (llm.ModelProvider, error) {
	return openai.NewOpenAIModelProvider(opts...)
}
