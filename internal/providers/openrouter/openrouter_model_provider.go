// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package openrouter

import (
	"encoding/json"
	"fmt"
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/openai"
	"net/http"

	"github.com/openai/openai-go/v3/option"
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

type OpenRouterModelProvider struct {
	*openai.OpenAIModelProvider
	models map[string]OpenRouterModelInfo
	apiKey string
}

var _ llm.ModelProvider = (*OpenRouterModelProvider)(nil)

func NewOpenRouterModelProvider(opts ...llm.ModelOption) (*OpenRouterModelProvider, error) {
	config := llm.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, llm.ErrAPIKeyEmpty
	}

	// Build request options list
	requestOpts := []option.RequestOption{}
	requestOpts = append(requestOpts, option.WithAPIKey(config.APIKey))

	// Set base URL (use default if not provided)
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1/"
	}
	requestOpts = append(requestOpts, option.WithBaseURL(baseURL))

	// Append any custom options
	requestOpts = append(requestOpts, config.Options...)

	models, err := loadModels(config.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load models: %w", err)
	}
	// Create the completion model with OpenRouter's API endpoint
	openAIModelProvider, err := openai.NewBaseOpenAIModelProvider("openrouter", models, requestOpts)
	if err != nil {
		return nil, err
	}

	provider := &OpenRouterModelProvider{
		OpenAIModelProvider: openAIModelProvider,
		apiKey:              config.APIKey,
	}

	return provider, nil
}

// loadModels fetches all available models from OpenRouter API
func loadModels(apiKey string) ([]*llm.ModelInfo, error) {
	req, err := http.NewRequest("GET", "https://openrouter.ai/api/v1/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var modelsResponse OpenRouterModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var models []*llm.ModelInfo

	for _, model := range modelsResponse.Data {
		modelInfo := &llm.ModelInfo{
			ID:   model.ID,
			Name: model.Name,
			Pricing: llm.ModelPricing{
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

	return models, nil
}
