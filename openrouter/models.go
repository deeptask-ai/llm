// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package openrouter

import (
	"github.com/deeptask-ai/llm"
	"github.com/deeptask-ai/llm/internal/providers/openrouter"
)

// NewOpenRouterModel creates a new OpenRouter model that supports llm
func NewOpenRouterModel(opts ...llm.ModelOption) (llm.CompletionModel, error) {
	return openrouter.NewOpenRouterModel(opts...)
}
