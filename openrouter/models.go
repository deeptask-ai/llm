// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package openrouter

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/openrouter"
)

// NewOpenRouterModel creates a new OpenRouter model that supports llm
func NewOpenRouterModel(opts ...llm.ModelOption) (llm.CompletionModel, error) {
	return openrouter.NewOpenRouterModel(opts...)
}
