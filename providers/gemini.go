// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/gemini"
)

// NewGeminiModelProvider creates a new Gemini model that supports llm
func NewGeminiModelProvider(opts ...llm.ModelOption) (llm.ModelProvider, error) {
	return gemini.NewGeminiModelProvider(opts...)
}
