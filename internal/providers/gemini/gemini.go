// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package gemini

import (
	_ "embed"
	"encoding/json"
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/openai"
	"github.com/openai/openai-go/v3/option"
)

type GeminiModel struct {
	*openai.OpenAICompletionModel
}

func NewGeminiModel(opts ...llm.ModelOption) (*GeminiModel, error) {
	config := llm.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, llm.ErrAPIKeyEmpty
	}

	// Build request options list
	requestOpts := []option.RequestOption{}

	// Set base URL (use default if not provided)
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta/openai/"
	}
	requestOpts = append(requestOpts, option.WithBaseURL(baseURL))

	// Append any custom options
	requestOpts = append(requestOpts, config.Options...)

	// Create the completion model with Gemini's OpenAI-compatible API endpoint
	completionModel, err := openai.NewOpenAICompletionModel(config.APIKey, requestOpts...)
	if err != nil {
		return nil, err
	}

	return &GeminiModel{
		OpenAICompletionModel: completionModel,
	}, nil
}

//go:embed gemini.json
var geminiModels []byte

func (p *GeminiModel) SupportedModels() []*llm.ModelInfo {
	var models []*llm.ModelInfo
	if err := json.Unmarshal(geminiModels, &models); err != nil {
		return nil
	}
	return models
}

func (p *GeminiModel) Name() string {
	return "gemini"
}
