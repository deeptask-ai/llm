package easyllm

import (
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
		return nil, ErrAPIKeyEmpty
	}
	if config.BaseURL == "" {
		return nil, ErrBaseURLEmpty
	}
	if config.APIVersion == "" {
		return nil, ErrAPIVersionEmpty
	}

	// Create base model with Azure OpenAI's API endpoint and required headers
	base, err := newOpenAIBaseModel(
		config.APIKey,
		option.WithBaseURL(config.BaseURL),
		option.WithQuery("api-version", config.APIVersion),
	)
	if err != nil {
		return nil, err
	}

	return &AzureOpenAIModel{
		OpenAIModel: &OpenAIModel{
			OpenAICompletionModel: &OpenAICompletionModel{OpenAIBaseModel: base},
			OpenAIEmbeddingModel:  &OpenAIEmbeddingModel{OpenAIBaseModel: base},
			OpenAIImageModel:      &OpenAIImageModel{OpenAIBaseModel: base},
		},
	}, nil
}

func (p *AzureOpenAIModel) Name() string {
	return "azure_openai"
}

func (p *AzureOpenAIModel) SupportedModels() []*ModelInfo {
	return []*ModelInfo{
		{
			ID:   "gpt-5",
			Name: "GPT-5",
			Pricing: ModelPricing{
				Prompt:            "1.25",
				Completion:        "10",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.01",
				InternalReasoning: "0",
				InputCacheRead:    "0.125",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-5-mini",
			Name: "GPT-5 Mini",
			Pricing: ModelPricing{
				Prompt:            "0.25",
				Completion:        "2",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.01",
				InternalReasoning: "0",
				InputCacheRead:    "0.025",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-5-nano",
			Name: "GPT-5 Nano",
			Pricing: ModelPricing{
				Prompt:            "0.05",
				Completion:        "0.4",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.01",
				InternalReasoning: "0",
				InputCacheRead:    "0.005",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-4.1",
			Name: "GPT-4.1 (fine-tuning)",
			Pricing: ModelPricing{
				Prompt:            "3",
				Completion:        "12",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.025",
				InternalReasoning: "0",
				InputCacheRead:    "0.75",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-4.1-mini",
			Name: "GPT-4.1 Mini (fine-tuning)",
			Pricing: ModelPricing{
				Prompt:            "0.8",
				Completion:        "3.2",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.025",
				InternalReasoning: "0",
				InputCacheRead:    "0.2",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-4.1-nano",
			Name: "GPT-4.1 Nano (fine-tuning)",
			Pricing: ModelPricing{
				Prompt:            "0.2",
				Completion:        "0.8",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.025",
				InternalReasoning: "0",
				InputCacheRead:    "0.05",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "o4-mini",
			Name: "o4-mini (reinforcement fine-tuning)",
			Pricing: ModelPricing{
				Prompt:            "4",
				Completion:        "16",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.01",
				InternalReasoning: "0",
				InputCacheRead:    "1",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-realtime-text",
			Name: "gpt-realtime (text)",
			Pricing: ModelPricing{
				Prompt:            "4",
				Completion:        "16",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.01",
				InternalReasoning: "0",
				InputCacheRead:    "0.4",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-realtime-audio",
			Name: "gpt-realtime (audio)",
			Pricing: ModelPricing{
				Prompt:            "32",
				Completion:        "64",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.01",
				InternalReasoning: "0",
				InputCacheRead:    "0.4",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-realtime-image",
			Name: "gpt-realtime (image)",
			Pricing: ModelPricing{
				Prompt:            "5",
				Completion:        "0",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.01",
				InternalReasoning: "0",
				InputCacheRead:    "0.5",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-4o-mini-text",
			Name: "gpt-4o-mini (text)",
			Pricing: ModelPricing{
				Prompt:            "0.6",
				Completion:        "2.4",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.025",
				InternalReasoning: "0",
				InputCacheRead:    "0.3",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-4o-mini-audio",
			Name: "gpt-4o-mini (audio)",
			Pricing: ModelPricing{
				Prompt:            "10",
				Completion:        "20",
				Request:           "0",
				Image:             "0",
				WebSearch:         "0.025",
				InternalReasoning: "0",
				InputCacheRead:    "0.3",
				InputCacheWrite:   "0",
			},
		},
		{
			ID:   "gpt-image-1",
			Name: "GPT-image-1",
			Pricing: ModelPricing{
				Prompt:            "5",
				Completion:        "40",
				Request:           "0",
				Image:             "0.17",
				WebSearch:         "0",
				InternalReasoning: "0",
				InputCacheRead:    "10",
				InputCacheWrite:   "0",
			},
		},
	}
}
