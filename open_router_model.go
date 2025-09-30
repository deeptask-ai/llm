package llmclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
		return nil, errors.New("API key cannot be empty")
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

func (p *OpenRouterModel) calculateCost(model string, usage *TokenUsage) *float64 {
	modelInfo, exists := p.models[model]
	if !exists {
		return nil
	}

	totalCost := 0.0

	cacheReadPrice, err := strconv.ParseFloat(modelInfo.Pricing.InputCacheRead, 64)
	if err != nil {
		cacheReadPrice = 0.0
	}
	promptPrice, err := strconv.ParseFloat(modelInfo.Pricing.Prompt, 64)
	if err != nil {
		return nil
	}
	if cacheReadPrice > 0.0 {
		totalInputTokens := usage.TotalInputTokens - usage.TotalCacheReadTokens
		totalCost += (float64(totalInputTokens) / 1000000.0) * promptPrice
		totalCost += (float64(usage.TotalCacheReadTokens) / 1000000.0) * cacheReadPrice
	} else {
		totalCost += (float64(usage.TotalInputTokens) / 1000000.0) * promptPrice
	}

	internalReasoningPrice, err := strconv.ParseFloat(modelInfo.Pricing.InternalReasoning, 64)
	if err != nil {
		internalReasoningPrice = 0.0
	}
	if internalReasoningPrice > 0.0 {
		totalCost += (float64(usage.TotalReasoningTokens) / 1000000.0) * internalReasoningPrice
	}

	completionPrice, err := strconv.ParseFloat(modelInfo.Pricing.Completion, 64)
	if err != nil {
		return nil
	}
	totalCost += (float64(usage.TotalOutputTokens) / 1000000.0) * completionPrice

	return &totalCost
}
