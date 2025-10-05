// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package gemini

import (
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/openai"
	"github.com/openai/openai-go/v3/option"
)

type GeminiModelProvider struct {
	*openai.OpenAIModelProvider
}

var _ llm.ModelProvider = (*GeminiModelProvider)(nil)

//go:embed gemini.json
var geminiModels []byte

func NewGeminiModelProvider(opts ...llm.ModelOption) (*GeminiModelProvider, error) {
	config := llm.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, llm.ErrAPIKeyEmpty
	}

	// Build request options list
	requestOpts := []option.RequestOption{}
	requestOpts = append(requestOpts, option.WithAPIKey(config.APIKey))

	// Set base URL (use default if not provided)
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta/openai/"
	}
	requestOpts = append(requestOpts, option.WithBaseURL(baseURL))

	// Append any custom options
	requestOpts = append(requestOpts, config.Options...)

	var models []*llm.ModelInfo
	if err := json.Unmarshal(geminiModels, &models); err != nil {
		return nil, errors.New("failed to read model info")
	}

	// Create the completion model with DeepSeek's API endpoint
	provider, err := openai.NewBaseOpenAIModelProvider("gemini", models, requestOpts)
	if err != nil {
		return nil, err
	}

	return &GeminiModelProvider{
		OpenAIModelProvider: provider,
	}, nil
}
