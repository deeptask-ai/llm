package easyllm

import (
	"github.com/openai/openai-go/option"
)

type GeminiModel struct {
	*OpenAICompletionModel
}

type GeminiModelConfig struct {
	APIKey string
}

func NewGeminiModel(config GeminiModelConfig) (*GeminiModel, error) {
	if config.APIKey == "" {
		return nil, ErrAPIKeyEmpty
	}

	// Create the completion model with Gemini's OpenAI-compatible API endpoint
	completionModel, err := NewOpenAICompletionModel(
		config.APIKey,
		option.WithBaseURL("https://generativelanguage.googleapis.com/v1beta/openai/"),
	)
	if err != nil {
		return nil, err
	}

	return &GeminiModel{
		OpenAICompletionModel: completionModel,
	}, nil
}

func (p *GeminiModel) SupportedModels() []*ModelInfo {
	return []*ModelInfo{
		{
			ID:   "gemini-2.0-flash-exp",
			Name: "Gemini 2.0 Flash Experimental",
			Pricing: ModelPricing{
				Prompt:            "0.075",
				Completion:        "0.3",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "0.0375",
				InputCacheWrite:   "1.125",
			},
		},
		{
			ID:   "gemini-2.5-flash",
			Name: "Gemini 2.5 Flash",
			Pricing: ModelPricing{
				Prompt:            "0.075",
				Completion:        "0.3",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "0.0375",
				InputCacheWrite:   "1.125",
			},
		},
		{
			ID:   "gemini-1.5-pro",
			Name: "Gemini 1.5 Pro",
			Pricing: ModelPricing{
				Prompt:            "1.25",
				Completion:        "5",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "0.625",
				InputCacheWrite:   "18.75",
			},
		},
		{
			ID:   "gemini-1.5-flash",
			Name: "Gemini 1.5 Flash",
			Pricing: ModelPricing{
				Prompt:            "0.075",
				Completion:        "0.3",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "0.0375",
				InputCacheWrite:   "1.125",
			},
		},
		{
			ID:   "gemini-1.5-flash-8b",
			Name: "Gemini 1.5 Flash-8B",
			Pricing: ModelPricing{
				Prompt:            "0.0375",
				Completion:        "0.15",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "0.01875",
				InputCacheWrite:   "0.5625",
			},
		},
		{
			ID:   "gemini-exp-1206",
			Name: "Gemini Experimental 1206",
			Pricing: ModelPricing{
				Prompt:            "0.075",
				Completion:        "0.3",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "0.0375",
				InputCacheWrite:   "1.125",
			},
		},
	}
}

func (p *GeminiModel) Name() string {
	return "gemini"
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
