package gemini

import (
	_ "embed"
	"encoding/json"
	"github.com/easymvp/easyllm/internal/providers/openai"
	"github.com/easymvp/easyllm/types"
	"github.com/openai/openai-go/v3/option"
)

type GeminiModel struct {
	*openai.OpenAICompletionModel
}

func NewGeminiModel(opts ...types.ModelOption) (*GeminiModel, error) {
	config := types.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, types.ErrAPIKeyEmpty
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

func (p *GeminiModel) SupportedModels() []*types.ModelInfo {
	var models []*types.ModelInfo
	if err := json.Unmarshal(geminiModels, &models); err != nil {
		return nil
	}
	return models
}

func (p *GeminiModel) Name() string {
	return "gemini"
}
