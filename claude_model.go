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

type ClaudeModel struct {
	*OpenAIModel
}

type ClaudeModelConfig struct {
	APIKey string
}

func NewClaudeModel(config ClaudeModelConfig) (*ClaudeModel, error) {
	if config.APIKey == "" {
		return nil, errors.New("API key cannot be empty")
	}

	// Create the client with Claude's API endpoint and required headers
	client := openai.NewClient(
		option.WithBaseURL("https://api.anthropic.com/v1/"),
		option.WithAPIKey(config.APIKey),
		option.WithHeader("anthropic-version", "2023-06-01"),
	)

	openaiProvider := &OpenAIModel{
		client: client,
		apiKey: config.APIKey,
	}

	provider := &ClaudeModel{
		OpenAIModel: openaiProvider,
	}

	return provider, nil
}

//go:embed data/claude.json
var claudeModels []byte

func (p *ClaudeModel) SupportedModels() []*ModelInfo {
	var models []*ModelInfo
	if err := json.Unmarshal(claudeModels, &models); err != nil {
		return nil
	}
	return models
}

func (p *ClaudeModel) Name() string {
	return "claude"
}

func (p *ClaudeModel) GenerateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	return nil, fmt.Errorf("embeddings are not supported by Claude models")
}

func (p *ClaudeModel) GenerateImage(ctx context.Context, req *ImageRequest) (*ImageResponse, error) {
	return nil, fmt.Errorf("image generation is not supported by Claude models")
}
