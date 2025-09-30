package easyllm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenRouterModelInfo struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Pricing OpenRouterModelPricing `json:"pricing"`
}

type OpenRouterModelPricing struct {
	Prompt            string `json:"prompt"`
	Completion        string `json:"completion"`
	Request           string `json:"request"`
	Image             string `json:"image"`
	WebSearch         string `json:"web_search"`
	InternalReasoning string `json:"internal_reasoning"`
	InputCacheRead    string `json:"input_cache_read"`
	InputCacheWrite   string `json:"input_cache_write"`
}

type OpenRouterModelsResponse struct {
	Data []OpenRouterModelInfo `json:"data"`
}

type OpenRouterModel struct {
	*OpenAIModel
	models map[string]OpenRouterModelInfo
}

type OpenRouterModelConfig struct {
	APIKey string
}

func NewOpenRouterModel(config OpenRouterModelConfig) (*OpenRouterModel, error) {
	if config.APIKey == "" {
		return nil, ErrAPIKeyEmpty
	}
	client := openai.NewClient(
		option.WithBaseURL("https://openrouter.ai/api/v1/"),
		option.WithAPIKey(config.APIKey),
	)
	openaiProvider := &OpenAIModel{
		client: client,
		apiKey: config.APIKey,
	}

	provider := &OpenRouterModel{
		OpenAIModel: openaiProvider,
		models:      make(map[string]OpenRouterModelInfo),
	}

	if err := provider.loadModels(); err != nil {
		return nil, fmt.Errorf("failed to load models: %w", err)
	}
	return provider, nil
}

// loadModels fetches all available models from OpenRouter API
func (p *OpenRouterModel) loadModels() error {
	req, err := http.NewRequest("GET", "https://openrouter.ai/api/v1/models", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var modelsResponse OpenRouterModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResponse); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	for _, model := range modelsResponse.Data {
		p.models[model.ID] = model
	}

	return nil
}

// SupportedModels returns all available models from OpenRouter
func (p *OpenRouterModel) SupportedModels() []*ModelInfo {
	var models []*ModelInfo

	for _, model := range p.models {
		modelInfo := &ModelInfo{
			ID:   model.ID,
			Name: model.Name,
			Pricing: ModelPricing{
				Prompt:            model.Pricing.Prompt,
				Completion:        model.Pricing.Completion,
				Request:           model.Pricing.Request,
				Image:             model.Pricing.Image,
				WebSearch:         model.Pricing.WebSearch,
				InternalReasoning: model.Pricing.InternalReasoning,
				InputCacheRead:    model.Pricing.InputCacheRead,
				InputCacheWrite:   model.Pricing.InputCacheWrite,
			},
		}
		models = append(models, modelInfo)
	}

	return models
}

// getModelInfo returns the ModelInfo for a given model from OpenRouter's model list
func (p *OpenRouterModel) getModelInfo(modelID string) *ModelInfo {
	openRouterModel, exists := p.models[modelID]
	if !exists {
		return nil
	}

	return &ModelInfo{
		ID:   openRouterModel.ID,
		Name: openRouterModel.Name,
		Pricing: ModelPricing{
			Prompt:            openRouterModel.Pricing.Prompt,
			Completion:        openRouterModel.Pricing.Completion,
			Request:           openRouterModel.Pricing.Request,
			Image:             openRouterModel.Pricing.Image,
			WebSearch:         openRouterModel.Pricing.WebSearch,
			InternalReasoning: openRouterModel.Pricing.InternalReasoning,
			InputCacheRead:    openRouterModel.Pricing.InputCacheRead,
			InputCacheWrite:   openRouterModel.Pricing.InputCacheWrite,
		},
	}
}

func (p *OpenRouterModel) Name() string {
	return "openrouter"
}

func (p *OpenRouterModel) GenerateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	return nil, NewUnsupportedCapabilityError("OpenRouter", "embeddings")
}

func (p *OpenRouterModel) GenerateImage(ctx context.Context, req *ImageRequest) (*ImageResponse, error) {
	return nil, NewUnsupportedCapabilityError("OpenRouter", "image generation")
}
