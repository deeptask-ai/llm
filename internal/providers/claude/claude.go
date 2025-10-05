// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package claude

import (
	_ "embed"
	"encoding/json"
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/openai"
	"github.com/openai/openai-go/v3/option"
)

type ClaudeModel struct {
	*openai.OpenAICompletionModel
}

func NewClaudeModel(opts ...llm.ModelOption) (*ClaudeModel, error) {
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

	// Create the completion model with Claude's API endpoint and required headers
	completionModel, err := openai.NewOpenAICompletionModel(config.APIKey, requestOpts...)
	if err != nil {
		return nil, err
	}

	return &ClaudeModel{
		OpenAICompletionModel: completionModel,
	}, nil
}

//go:embed claude.json
var claudeModels []byte

func (p *ClaudeModel) SupportedModels() []*llm.ModelInfo {
	var models []*llm.ModelInfo
	if err := json.Unmarshal(claudeModels, &models); err != nil {
		return nil
	}
	return models
}

func (p *ClaudeModel) Name() string {
	return "claude"
}
