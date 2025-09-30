package llmclient

import (
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type AzureOpenAIModel struct {
	*OpenAIModel
}

type AzureOpenAIModelConfig struct {
	APIKey     string
	BaseURL    string
	APIVersion string
}

func NewAzureOpenAIModel(config AzureOpenAIModelConfig) (*AzureOpenAIModel, error) {
	if config.APIKey == "" {
		return nil, errors.New("API key cannot be empty")
	}
	if config.BaseURL == "" {
		return nil, errors.New("base URL cannot be empty")
	}
	if config.APIVersion == "" {
		return nil, errors.New("API version cannot be empty")
	}

	// Create the client with Azure OpenAI's API endpoint and required headers
	client := openai.NewClient(
		option.WithBaseURL(config.BaseURL),
		option.WithAPIKey(config.APIKey),
		option.WithQuery("api-version", config.APIVersion),
	)

	openaiProvider := &OpenAIModel{
		client: client,
		apiKey: config.APIKey,
	}

	provider := &AzureOpenAIModel{
		OpenAIModel: openaiProvider,
	}

	return provider, nil
}

func (p *AzureOpenAIModel) Name() string {
	return "azure_openai"
}

//go:embed data/openai.json
var azureOpenaiModels []byte

func (p *AzureOpenAIModel) SupportedModels() []*ModelInfo {
	var models []*ModelInfo
	if err := json.Unmarshal(azureOpenaiModels, &models); err != nil {
		return nil
	}
	return models
}
