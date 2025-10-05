// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/deepseek"
)

// NewDeepSeekModelProvider creates a new DeepSeek model that supports llm
func NewDeepSeekModelProvider(opts ...llm.ModelOption) (llm.ModelProvider, error) {
	return deepseek.NewDeepSeekModelProvider(opts...)
}
