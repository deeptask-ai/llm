package easyllm

import (
	"github.com/openai/openai-go/option"
)

type DeepSeekModel struct {
	*OpenAICompletionModel
}

func NewDeepSeekModel(opts ...ModelOption) (*DeepSeekModel, error) {
	config := applyOptions(opts)

	if config.APIKey == "" {
		return nil, ErrAPIKeyEmpty
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
	completionModel, err := NewOpenAICompletionModel(config.APIKey, requestOpts...)
	if err != nil {
		return nil, err
	}

	return &DeepSeekModel{
		OpenAICompletionModel: completionModel,
	}, nil
}

func (p *DeepSeekModel) SupportedModels() []*ModelInfo {
	return []*ModelInfo{
		{
			ID:   "deepseek-chat",
			Name: "DeepSeek Chat",
			Pricing: ModelPricing{
				Prompt:            "0.28",
				Completion:        "0.42",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "0.028",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "deepseek-reasoner",
			Name: "DeepSeek Reasoner",
			Pricing: ModelPricing{
				Prompt:            "0.28",
				Completion:        "0.42",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "0.028",
				InputCacheWrite:   "0",
			},
		},
	}
}

func (p *DeepSeekModel) Name() string {
	return "deepseek"
}
