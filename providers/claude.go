// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0
package providers

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/claude"
)

// NewClaudeModelProvider creates a new Claude model that supports llm
func NewClaudeModelProvider(opts ...llm.ModelOption) (llm.ModelProvider, error) {
	return claude.NewClaudeModelProvider(opts...)
}
