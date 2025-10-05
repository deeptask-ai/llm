// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package azure

import (
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/openai"
	"github.com/openai/openai-go/v3/option"
)

//go:embed openai.json
var openaiModels []byte

type AzureOpenAIModelProvider struct {
	*openai.OpenAIModelProvider
}

func NewAzureOpenAIModelProvider(opts ...llm.ModelOption) (*AzureOpenAIModelProvider, error) {
	config := llm.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, llm.ErrAPIKeyEmpty
	}
	if config.BaseURL == "" {
		return nil, llm.ErrBaseURLEmpty
	}
	if config.APIVersion == "" {
		return nil, llm.ErrAPIVersionEmpty
	}

	// Build request options list
	requestOpts := []option.RequestOption{
		option.WithBaseURL(config.BaseURL),
		option.WithQuery("api-version", config.APIVersion),
	}

	// Append any custom options
	requestOpts = append(requestOpts, config.Options...)
	var models []*llm.ModelInfo
	if err := json.Unmarshal(openaiModels, &models); err != nil {
		return nil, errors.New("failed to read model info")
	}
	// Create baseProvider model with Azure OpenAI's API endpoint and required headers
	provider, err := openai.NewBaseOpenAIModelProvider("azure_openai", models, requestOpts)
	if err != nil {
		return nil, err
	}

	return &AzureOpenAIModelProvider{
		OpenAIModelProvider: provider,
	}, nil
}

func (p *AzureOpenAIModelProvider) Name() string {
	return "azure_openai"
}
