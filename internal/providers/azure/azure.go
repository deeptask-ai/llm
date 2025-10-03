// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package azure

import (
	"github.com/easymvp-ai/llm"
	"github.com/easymvp-ai/llm/internal/providers/openai"
	"github.com/openai/openai-go/v3/option"
)

type AzureOpenAIModel struct {
	*openai.OpenAIModel
}

func NewAzureOpenAIModel(opts ...llm.ModelOption) (*AzureOpenAIModel, error) {
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

	// Create base model with Azure OpenAI's API endpoint and required headers
	base, err := openai.NewOpenAIBaseModel(config.APIKey, requestOpts...)
	if err != nil {
		return nil, err
	}

	return &AzureOpenAIModel{
		OpenAIModel: &openai.OpenAIModel{
			OpenAICompletionModel: &openai.OpenAICompletionModel{OpenAIBaseModel: base},
			OpenAIEmbeddingModel:  &openai.OpenAIEmbeddingModel{OpenAIBaseModel: base},
			OpenAIImageModel:      &openai.OpenAIImageModel{OpenAIBaseModel: base},
		},
	}, nil
}

func (p *AzureOpenAIModel) Name() string {
	return "azure_openai"
}
