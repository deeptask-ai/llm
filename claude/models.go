// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package claude

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/claude"
)

// NewClaudeModel creates a new Claude model that supports llm
func NewClaudeModel(opts ...llm.ModelOption) (llm.CompletionModel, error) {
	return claude.NewClaudeModel(opts...)
}
