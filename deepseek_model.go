package llmclient

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type DeepSeekModel struct {
	*OpenAIModel
}

type DeepSeekModelConfig struct {
	APIKey string
}

func NewDeepSeekModel(config DeepSeekModelConfig) (*DeepSeekModel, error) {
	if config.APIKey == "" {
		return nil, errors.New("API key cannot be empty")
	}

	// Create the client with DeepSeek's API endpoint
	client := openai.NewClient(
		option.WithBaseURL("https://api.deepseek.com/"),
		option.WithAPIKey(config.APIKey),
	)

	openaiProvider := &OpenAIModel{
		client: client,
		apiKey: config.APIKey,
	}

	provider := &DeepSeekModel{
		OpenAIModel: openaiProvider,
	}

	return provider, nil
}

//go:embed data/deepseek.json
var deepseekModels []byte

func (p *DeepSeekModel) SupportedModels() []*ModelInfo {
	var models []*ModelInfo
	if err := json.Unmarshal(deepseekModels, &models); err != nil {
		return nil
	}
	return models
}

func (p *DeepSeekModel) Name() string {
	return "deepseek"
}

func (p *DeepSeekModel) GenerateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	return nil, fmt.Errorf("embeddings are not supported by DeepSeek models")
}
