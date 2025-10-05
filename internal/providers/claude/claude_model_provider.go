// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package claude

import (
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/openai"
	"github.com/openai/openai-go/v3/option"
)

type ClaudeModelProvider struct {
	*openai.OpenAIModelProvider
}

//go:embed claude.json
var claudeModels []byte

var _ llm.ModelProvider = (*ClaudeModelProvider)(nil)

func NewClaudeModelProvider(opts ...llm.ModelOption) (*ClaudeModelProvider, error) {
	config := llm.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, llm.ErrAPIKeyEmpty
	}

	// Build request options list with defaults
	requestOpts := []option.RequestOption{
		option.WithHeader("anthropic-version", "2023-06-01"),
	}

	// Set base URL (use default if not provided)
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1/"
	}
	requestOpts = append(requestOpts, option.WithBaseURL(baseURL))

	// Append any custom options
	requestOpts = append(requestOpts, config.Options...)

	var models []*llm.ModelInfo
	if err := json.Unmarshal(claudeModels, &models); err != nil {
		return nil, errors.New("failed to read model info")
	}

	// Create the completion model with Claude's API endpoint and required headers
	provider, err := openai.NewBaseOpenAIModelProvider("claude", config.APIKey, models, requestOpts...)
	if err != nil {
		return nil, err
	}

	return &ClaudeModelProvider{
		OpenAIModelProvider: provider,
	}, nil
}
