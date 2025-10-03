// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package deeptask

import (
	"github.com/easymvp-ai/llm"
	"github.com/easymvp-ai/llm/internal/providers/deepseek"
)

// NewDeepSeekModel creates a new DeepSeek model that supports llm
func NewDeepSeekModel(opts ...llm.ModelOption) (llm.CompletionModel, error) {
	return deepseek.NewDeepSeekModel(opts...)
}
