// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package gemini

import (
	"github.com/easymvp-ai/llm"
	"github.com/easymvp-ai/llm/internal/providers/gemini"
)

// NewGeminiModel creates a new Gemini model that supports llm
func NewGeminiModel(opts ...llm.ModelOption) (llm.CompletionModel, error) {
	return gemini.NewGeminiModel(opts...)
}
