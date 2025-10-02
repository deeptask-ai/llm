package easyllm

import (
	"github.com/openai/openai-go/option"
)

type ClaudeModel struct {
	*OpenAICompletionModel
}

func NewClaudeModel(opts ...ModelOption) (*ClaudeModel, error) {
	config := applyOptions(opts)

	if config.APIKey == "" {
		return nil, ErrAPIKeyEmpty
	}

	// Build request options list with defaults
	requestOpts := []option.RequestOption{
		option.WithHeader("anthropic-version", "2023-06-01"),
	}

	// Set base URL (use default if not provided)
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1/"
	}
	requestOpts = append(requestOpts, option.WithBaseURL(baseURL))

	// Append any custom options
	requestOpts = append(requestOpts, config.Options...)

	// Create the completion model with Claude's API endpoint and required headers
	completionModel, err := NewOpenAICompletionModel(config.APIKey, requestOpts...)
	if err != nil {
		return nil, err
	}

	return &ClaudeModel{
		OpenAICompletionModel: completionModel,
	}, nil
}

func (p *ClaudeModel) SupportedModels() []*ModelInfo {
	return []*ModelInfo{
		{
			ID:   "opus-4.1",
			Name: "Opus 4.1",
			Pricing: ModelPricing{
				Prompt:            "15",
				Completion:        "75",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "1.5",
				InputCacheWrite:   "18.75",
			},
		},
		{
			ID:   "sonnet-4.5",
			Name: "Sonnet 4.5",
			Pricing: ModelPricing{
				Prompt:            "3",
				Completion:        "15",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "0.3",
				InputCacheWrite:   "3.75",
			},
		},
		{
			ID:   "haiku-3.5",
			Name: "Haiku 3.5",
			Pricing: ModelPricing{
				Prompt:            "0.8",
				Completion:        "4",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "0.08",
				InputCacheWrite:   "1",
			},
		},
	}
}

func (p *ClaudeModel) Name() string {
	return "claude"
}
