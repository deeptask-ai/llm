package llmclient

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type GeminiModel struct {
	*OpenAIModel
}

type GeminiModelConfig struct {
	APIKey string
}

func NewGeminiModel(config GeminiModelConfig) (*GeminiModel, error) {
	if config.APIKey == "" {
		return nil, ErrAPIKeyEmpty
	}

	// Create the client with Gemini's OpenAI-compatible API endpoint
	client := openai.NewClient(
		option.WithBaseURL("https://generativelanguage.googleapis.com/v1beta/openai/"),
		option.WithAPIKey(config.APIKey),
	)

	openaiProvider := &OpenAIModel{
		client: client,
		apiKey: config.APIKey,
	}

	provider := &GeminiModel{
		OpenAIModel: openaiProvider,
	}

	return provider, nil
}

//go:embed data/gemini.json
var geminiModels []byte

func (p *GeminiModel) SupportedModels() []*ModelInfo {
	var models []*ModelInfo
	if err := json.Unmarshal(geminiModels, &models); err != nil {
		return nil
	}
	return models
}

func (p *GeminiModel) Name() string {
	return "gemini"
}

func (p *GeminiModel) GenerateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	return nil, NewUnsupportedCapabilityError("Gemini", "embeddings")
}

func (p *GeminiModel) GenerateImage(ctx context.Context, req *ImageRequest) (*ImageResponse, error) {
	return nil, NewUnsupportedCapabilityError("Gemini", "image generation")
}

// Override getModelInfo to use Gemini-specific models
func (p *GeminiModel) getModelInfo(modelID string) *ModelInfo {
	models := p.SupportedModels()
	for _, model := range models {
		if model.ID == modelID {
			return model
		}
	}
	return nil
}
