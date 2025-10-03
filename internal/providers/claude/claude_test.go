// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package claude

import (
	"github.com/deeptask-ai/llm"
	"testing"

	"github.com/openai/openai-go/v3/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClaudeModel_Success(t *testing.T) {
	tests := []struct {
		name        string
		opts        []llm.ModelOption
		wantName    string
		description string
	}{
		{
			name: "basic_configuration",
			opts: []llm.ModelOption{
				llm.WithAPIKey("test-api-key"),
			},
			wantName:    "claude",
			description: "Should create Claude model with basic configuration",
		},
		{
			name: "with_custom_base_url",
			opts: []llm.ModelOption{
				llm.WithAPIKey("test-api-key"),
				llm.WithBaseURL("https://custom.anthropic.com/v1"),
			},
			wantName:    "claude",
			description: "Should create Claude model with custom base URL",
		},
		{
			name: "with_custom_request_option",
			opts: []llm.ModelOption{
				llm.WithAPIKey("test-api-key"),
				llm.WithRequestOption(option.WithHeader("Custom-Header", "custom-value")),
			},
			wantName:    "claude",
			description: "Should create Claude model with custom request options",
		},
		{
			name: "with_multiple_request_options",
			opts: []llm.ModelOption{
				llm.WithAPIKey("test-api-key"),
				llm.WithRequestOptions(
					option.WithHeader("Custom-Header-1", "value1"),
					option.WithHeader("Custom-Header-2", "value2"),
				),
			},
			wantName:    "claude",
			description: "Should create Claude model with multiple custom request options",
		},
		{
			name: "default_base_url",
			opts: []llm.ModelOption{
				llm.WithAPIKey("test-api-key"),
			},
			wantName:    "claude",
			description: "Should use default Anthropic base URL when not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewClaudeModel(tt.opts...)

			require.NoError(t, err, tt.description)
			require.NotNil(t, model, "Model should not be nil")
			assert.NotNil(t, model.OpenAICompletionModel, "OpenAICompletionModel should not be nil")
			assert.Equal(t, tt.wantName, model.Name(), "Model name should match expected value")
		})
	}
}

func TestNewClaudeModel_MissingAPIKey(t *testing.T) {
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
			model, err := NewClaudeModel(tt.opts...)

			assert.Error(t, err, tt.description)
			assert.ErrorIs(t, err, llm.ErrAPIKeyEmpty, "Error should be ErrAPIKeyEmpty")
			assert.Nil(t, model, "Model should be nil when error occurs")
		})
	}
}

func TestClaudeModel_Name(t *testing.T) {
	model, err := NewClaudeModel(llm.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	assert.Equal(t, "claude", model.Name(), "Name should return 'claude'")
}

func TestClaudeModel_SupportedModels(t *testing.T) {
	model, err := NewClaudeModel(llm.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	models := model.SupportedModels()
	assert.NotNil(t, models, "SupportedModels should not return nil")
	assert.Greater(t, len(models), 0, "SupportedModels should return at least one model")
}

func TestClaudeModel_ModelStructure(t *testing.T) {
	model, err := NewClaudeModel(llm.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	t.Run("has_completion_model", func(t *testing.T) {
		assert.NotNil(t, model.OpenAICompletionModel, "OpenAICompletionModel should be initialized")
		assert.NotNil(t, model.OpenAICompletionModel.OpenAIBaseModel, "OpenAIBaseModel should be initialized")
	})
}

func TestNewClaudeModel_MultipleInstances(t *testing.T) {
	// Create multiple instances with different configurations
	model1, err1 := NewClaudeModel(llm.WithAPIKey("test-api-key-1"))
	model2, err2 := NewClaudeModel(
		llm.WithAPIKey("test-api-key-2"),
		llm.WithBaseURL("https://custom.anthropic.com/v1"),
	)

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NotNil(t, model1)
	require.NotNil(t, model2)

	// Verify instances are independent
	assert.NotSame(t, model1, model2, "Different instances should be created")
	assert.NotSame(t, model1.OpenAICompletionModel, model2.OpenAICompletionModel, "Completion models should be independent")
}

// Benchmark tests
func BenchmarkNewClaudeModel_Success(b *testing.B) {
	opts := []llm.ModelOption{
		llm.WithAPIKey("test-api-key"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewClaudeModel(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewClaudeModel_WithOptions(b *testing.B) {
	opts := []llm.ModelOption{
		llm.WithAPIKey("test-api-key"),
		llm.WithBaseURL("https://custom.anthropic.com/v1"),
		llm.WithRequestOptions(
			option.WithHeader("Custom-Header-1", "value1"),
			option.WithHeader("Custom-Header-2", "value2"),
		),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewClaudeModel(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkClaudeModel_Name(b *testing.B) {
	model, err := NewClaudeModel(llm.WithAPIKey("test-api-key"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.Name()
	}
}

func BenchmarkClaudeModel_SupportedModels(b *testing.B) {
	model, err := NewClaudeModel(llm.WithAPIKey("test-api-key"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.SupportedModels()
	}
}
