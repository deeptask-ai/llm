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

type OpenRouterModel struct {
	*openai.OpenAICompletionModel
	models map[string]OpenRouterModelInfo
	apiKey string
}

func NewOpenRouterModel(opts ...llm.ModelOption) (*OpenRouterModel, error) {
	config := llm.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, llm.ErrAPIKeyEmpty
	}

	// Build request options list
	requestOpts := []option.RequestOption{}

	// Set base URL (use default if not provided)
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1/"
	}
	requestOpts = append(requestOpts, option.WithBaseURL(baseURL))

	// Append any custom options
	requestOpts = append(requestOpts, config.Options...)

	// Create the completion model with OpenRouter's API endpoint
	completionModel, err := openai.NewOpenAICompletionModel(config.APIKey, requestOpts...)
	if err != nil {
		return nil, err
	}

	provider := &OpenRouterModel{
		OpenAICompletionModel: completionModel,
		models:                make(map[string]OpenRouterModelInfo),
		apiKey:                config.APIKey,
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
func (p *OpenRouterModel) SupportedModels() []*llm.ModelInfo {
	var models []*llm.ModelInfo

	for _, model := range p.models {
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

	return models
}

// getModelInfo returns the ModelInfo for a given model from OpenRouter's model list
func (p *OpenRouterModel) getModelInfo(modelID string) *llm.ModelInfo {
	openRouterModel, exists := p.models[modelID]
	if !exists {
		return nil
	}

	return &llm.ModelInfo{
		ID:   openRouterModel.ID,
		Name: openRouterModel.Name,
		Pricing: llm.ModelPricing{
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
