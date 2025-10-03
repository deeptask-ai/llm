package gemini

import (
	"github.com/deeptask-ai/llm"
	"testing"

	"github.com/openai/openai-go/v3/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGeminiModel_Success(t *testing.T) {
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
			wantName:    "gemini",
			description: "Should create Gemini model with basic configuration",
		},
		{
			name: "with_custom_base_url",
			opts: []llm.ModelOption{
				llm.WithAPIKey("test-api-key"),
				llm.WithBaseURL("https://custom.googleapis.com/v1beta/openai/"),
			},
			wantName:    "gemini",
			description: "Should create Gemini model with custom base URL",
		},
		{
			name: "with_custom_request_option",
			opts: []llm.ModelOption{
				llm.WithAPIKey("test-api-key"),
				llm.WithRequestOption(option.WithHeader("Custom-Header", "custom-value")),
			},
			wantName:    "gemini",
			description: "Should create Gemini model with custom request options",
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
			wantName:    "gemini",
			description: "Should create Gemini model with multiple custom request options",
		},
		{
			name: "default_base_url",
			opts: []llm.ModelOption{
				llm.WithAPIKey("test-api-key"),
			},
			wantName:    "gemini",
			description: "Should use default Google AI base URL when not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewGeminiModel(tt.opts...)

			require.NoError(t, err, tt.description)
			require.NotNil(t, model, "Model should not be nil")
			assert.NotNil(t, model.OpenAICompletionModel, "OpenAICompletionModel should not be nil")
			assert.Equal(t, tt.wantName, model.Name(), "Model name should match expected value")
		})
	}
}

func TestNewGeminiModel_MissingAPIKey(t *testing.T) {
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
			model, err := NewGeminiModel(tt.opts...)

			assert.Error(t, err, tt.description)
			assert.ErrorIs(t, err, llm.ErrAPIKeyEmpty, "Error should be ErrAPIKeyEmpty")
			assert.Nil(t, model, "Model should be nil when error occurs")
		})
	}
}

func TestGeminiModel_Name(t *testing.T) {
	model, err := NewGeminiModel(llm.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	assert.Equal(t, "gemini", model.Name(), "Name should return 'gemini'")
}

func TestGeminiModel_SupportedModels(t *testing.T) {
	model, err := NewGeminiModel(llm.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	models := model.SupportedModels()
	assert.NotNil(t, models, "SupportedModels should not return nil")
	assert.Greater(t, len(models), 0, "SupportedModels should return at least one model")
}

func TestGeminiModel_ModelStructure(t *testing.T) {
	model, err := NewGeminiModel(llm.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	t.Run("has_completion_model", func(t *testing.T) {
		assert.NotNil(t, model.OpenAICompletionModel, "OpenAICompletionModel should be initialized")
		assert.NotNil(t, model.OpenAICompletionModel.OpenAIBaseModel, "OpenAIBaseModel should be initialized")
	})
}

func TestNewGeminiModel_MultipleInstances(t *testing.T) {
	// Create multiple instances with different configurations
	model1, err1 := NewGeminiModel(llm.WithAPIKey("test-api-key-1"))
	model2, err2 := NewGeminiModel(
		llm.WithAPIKey("test-api-key-2"),
		llm.WithBaseURL("https://custom.googleapis.com/v1beta/openai/"),
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
func BenchmarkNewGeminiModel_Success(b *testing.B) {
	opts := []llm.ModelOption{
		llm.WithAPIKey("test-api-key"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewGeminiModel(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewGeminiModel_WithOptions(b *testing.B) {
	opts := []llm.ModelOption{
		llm.WithAPIKey("test-api-key"),
		llm.WithBaseURL("https://custom.googleapis.com/v1beta/openai/"),
		llm.WithRequestOptions(
			option.WithHeader("Custom-Header-1", "value1"),
			option.WithHeader("Custom-Header-2", "value2"),
		),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewGeminiModel(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGeminiModel_Name(b *testing.B) {
	model, err := NewGeminiModel(llm.WithAPIKey("test-api-key"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.Name()
	}
}

func BenchmarkGeminiModel_SupportedModels(b *testing.B) {
	model, err := NewGeminiModel(llm.WithAPIKey("test-api-key"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.SupportedModels()
	}
}
