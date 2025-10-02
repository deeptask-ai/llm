package easyllm

import (
	"github.com/openai/openai-go/option"
)

type DeepSeekModel struct {
	*OpenAICompletionModel
}

type DeepSeekModelConfig struct {
	APIKey string
}

func NewDeepSeekModel(config DeepSeekModelConfig) (*DeepSeekModel, error) {
	if config.APIKey == "" {
		return nil, ErrAPIKeyEmpty
	}

	// Create the completion model with DeepSeek's API endpoint
	completionModel, err := NewOpenAICompletionModel(
		config.APIKey,
		option.WithBaseURL("https://api.deepseek.com/"),
	)
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
