// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package models

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/deepseek"
)

// NewDeepSeekModel creates a new DeepSeek model that supports llm
func NewDeepSeekModel(opts ...llm.ModelOption) (llm.CompletionModel, error) {
	return deepseek.NewDeepSeekModel(opts...)
}
