package deepseek

import (
	_ "embed"
	"encoding/json"
	"github.com/easymvp/easyllm/internal/providers/openai"
	"github.com/easymvp/easyllm/types"
	"github.com/openai/openai-go/option"
)

type DeepSeekModel struct {
	*openai.OpenAICompletionModel
}

func NewDeepSeekModel(opts ...types.ModelOption) (*DeepSeekModel, error) {
	config := types.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, types.ErrAPIKeyEmpty
	}

	// Build request options list
	requestOpts := []option.RequestOption{}

	// Set base URL (use default if not provided)
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/"
	}
	requestOpts = append(requestOpts, option.WithBaseURL(baseURL))

	// Append any custom options
	requestOpts = append(requestOpts, config.Options...)

	// Create the completion model with DeepSeek's API endpoint
	completionModel, err := openai.NewOpenAICompletionModel(config.APIKey, requestOpts...)
	if err != nil {
		return nil, err
	}

	return &DeepSeekModel{
		OpenAICompletionModel: completionModel,
	}, nil
}

//go:embed deepseek.json
var deepSeekModels []byte

func (p *DeepSeekModel) SupportedModels() []*types.ModelInfo {
	var models []*types.ModelInfo
	if err := json.Unmarshal(deepSeekModels, &models); err != nil {
		return nil
	}
	return models
}

func (p *DeepSeekModel) Name() string {
	return "deepseek"
}
