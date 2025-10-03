// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package gemini

import (
	"github.com/deeptask-ai/llm"
	"github.com/deeptask-ai/llm/internal/providers/gemini"
)

// NewGeminiModel creates a new Gemini model that supports llm
func NewGeminiModel(opts ...llm.ModelOption) (llm.CompletionModel, error) {
	return gemini.NewGeminiModel(opts...)
}
