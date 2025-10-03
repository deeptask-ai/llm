// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package openrouter

import (
	"github.com/easymvp-ai/llm"
	"testing"

	"github.com/openai/openai-go/v3/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenRouterModel_MissingAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		opts        []llm.ModelOption
		description string
	}{
		{
			name: "empty_api_key",
			opts: []llm.ModelOption{
				llm.WithAPIKey(""),
			},
			description: "Should return error when API key is empty string",
		},
		{
			name:        "no_api_key",
			opts:        []llm.ModelOption{},
			description: "Should return error when API key is not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewOpenRouterModel(tt.opts...)

			assert.Error(t, err, tt.description)
			assert.ErrorIs(t, err, llm.ErrAPIKeyEmpty, "Error should be ErrAPIKeyEmpty")
			assert.Nil(t, model, "Model should be nil when error occurs")
		})
	}
}

func TestOpenRouterModel_Name(t *testing.T) {
	// Note: This test will make an actual API call to OpenRouter
	// In a real-world scenario, you might want to mock this
	t.Skip("Skipping test that requires actual API call - enable with valid API key")

	model, err := NewOpenRouterModel(llm.WithAPIKey("test-api-key"))
	if err != nil {
		// If the API call fails (expected with invalid key), skip the rest
		t.Skipf("API call failed as expected with test key: %v", err)
		return
	}

	require.NotNil(t, model)
	assert.Equal(t, "openrouter", model.Name(), "Name should return 'openrouter'")
}

func TestOpenRouterModel_ModelStructure(t *testing.T) {
	// Note: This test will make an actual API call to OpenRouter
	// In a real-world scenario, you might want to mock this
	t.Skip("Skipping test that requires actual API call - enable with valid API key")

	model, err := NewOpenRouterModel(llm.WithAPIKey("test-api-key"))
	if err != nil {
		t.Skipf("API call failed as expected with test key: %v", err)
		return
	}

	require.NotNil(t, model)

	t.Run("has_completion_model", func(t *testing.T) {
		assert.NotNil(t, model.OpenAICompletionModel, "OpenAICompletionModel should be initialized")
		assert.NotNil(t, model.OpenAICompletionModel.OpenAIBaseModel, "OpenAIBaseModel should be initialized")
	})

	t.Run("has_models_map", func(t *testing.T) {
		assert.NotNil(t, model.models, "Models map should be initialized")
	})
}

func TestOpenRouterModel_SupportedModels(t *testing.T) {
	// Note: This test will make an actual API call to OpenRouter
	// In a real-world scenario, you might want to mock this
	t.Skip("Skipping test that requires actual API call - enable with valid API key")

	model, err := NewOpenRouterModel(llm.WithAPIKey("test-api-key"))
	if err != nil {
		t.Skipf("API call failed as expected with test key: %v", err)
		return
	}

	require.NotNil(t, model)

	models := model.SupportedModels()
	assert.NotNil(t, models, "SupportedModels should not return nil")
	// OpenRouter should have many models available
	assert.Greater(t, len(models), 0, "SupportedModels should return at least one model")
}

func TestNewOpenRouterModel_WithCustomBaseURL(t *testing.T) {
	// This test verifies the constructor accepts custom base URL option
	// but will still fail due to missing/invalid API key
	model, err := NewOpenRouterModel(
		llm.WithAPIKey(""),
		llm.WithBaseURL("https://custom.openrouter.ai/api/v1/"),
	)

	// Should fail due to empty API key, not base URL
	assert.Error(t, err)
	assert.ErrorIs(t, err, llm.ErrAPIKeyEmpty)
	assert.Nil(t, model)
}

func TestNewOpenRouterModel_WithCustomRequestOptions(t *testing.T) {
	// This test verifies the constructor accepts custom request options
	// but will still fail due to missing/invalid API key
	model, err := NewOpenRouterModel(
		llm.WithAPIKey(""),
		llm.WithRequestOptions(
			option.WithHeader("Custom-Header-1", "value1"),
			option.WithHeader("Custom-Header-2", "value2"),
		),
	)

	// Should fail due to empty API key
	assert.Error(t, err)
	assert.ErrorIs(t, err, llm.ErrAPIKeyEmpty)
	assert.Nil(t, model)
}

// Note: Benchmark tests are skipped for OpenRouter since they require actual API calls
// In a production environment, you would want to implement mocking to test these

func BenchmarkOpenRouterModel_Name(b *testing.B) {
	b.Skip("Skipping benchmark that requires actual API call")

	model, err := NewOpenRouterModel(llm.WithAPIKey("test-api-key"))
	if err != nil {
		b.Skipf("API call failed: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.Name()
	}
}

func BenchmarkOpenRouterModel_SupportedModels(b *testing.B) {
	b.Skip("Skipping benchmark that requires actual API call")

	model, err := NewOpenRouterModel(llm.WithAPIKey("test-api-key"))
	if err != nil {
		b.Skipf("API call failed: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.SupportedModels()
	}
}
