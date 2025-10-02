package azure

import (
	"github.com/easymvp/easyllm/internal/providers/openai"
	"github.com/easymvp/easyllm/types"
	"github.com/openai/openai-go/v3/option"
)

type AzureOpenAIModel struct {
	*openai.OpenAIModel
}

func NewAzureOpenAIModel(opts ...types.ModelOption) (*AzureOpenAIModel, error) {
	config := types.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, types.ErrAPIKeyEmpty
	}
	if config.BaseURL == "" {
		return nil, types.ErrBaseURLEmpty
	}
	if config.APIVersion == "" {
		return nil, types.ErrAPIVersionEmpty
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
