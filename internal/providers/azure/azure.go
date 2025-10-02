package easyllm

import (
	"github.com/openai/openai-go/option"
)

type AzureOpenAIModel struct {
	*OpenAIModel
}

func NewAzureOpenAIModel(opts ...ModelOption) (*AzureOpenAIModel, error) {
	config := applyOptions(opts)

	if config.APIKey == "" {
		return nil, ErrAPIKeyEmpty
	}
	if config.BaseURL == "" {
		return nil, ErrBaseURLEmpty
	}
	if config.APIVersion == "" {
		return nil, ErrAPIVersionEmpty
	}

	// Build request options list
	requestOpts := []option.RequestOption{
		option.WithBaseURL(config.BaseURL),
		option.WithQuery("api-version", config.APIVersion),
	}

	// Append any custom options
	requestOpts = append(requestOpts, config.Options...)

	// Create base model with Azure OpenAI's API endpoint and required headers
	base, err := newOpenAIBaseModel(config.APIKey, requestOpts...)
	if err != nil {
		return nil, err
	}

	return &AzureOpenAIModel{
		OpenAIModel: &OpenAIModel{
			OpenAICompletionModel: &OpenAICompletionModel{OpenAIBaseModel: base},
			OpenAIEmbeddingModel:  &OpenAIEmbeddingModel{OpenAIBaseModel: base},
			OpenAIImageModel:      &OpenAIImageModel{OpenAIBaseModel: base},
		},
	}, nil
}

func (p *AzureOpenAIModel) Name() string {
	return "azure_openai"
}
