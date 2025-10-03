package deepseek

import (
	"testing"

	"github.com/deeptask-ai/llm/types"
	"github.com/openai/openai-go/v3/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeepSeekModel_Success(t *testing.T) {
	tests := []struct {
		name        string
		opts        []types.ModelOption
		wantName    string
		description string
	}{
		{
			name: "basic_configuration",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
			},
			wantName:    "deepseek",
			description: "Should create DeepSeek model with basic configuration",
		},
		{
			name: "with_custom_base_url",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithBaseURL("https://custom.deepseek.com/"),
			},
			wantName:    "deepseek",
			description: "Should create DeepSeek model with custom base URL",
		},
		{
			name: "with_custom_request_option",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithRequestOption(option.WithHeader("Custom-Header", "custom-value")),
			},
			wantName:    "deepseek",
			description: "Should create DeepSeek model with custom request options",
		},
		{
			name: "with_multiple_request_options",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithRequestOptions(
					option.WithHeader("Custom-Header-1", "value1"),
					option.WithHeader("Custom-Header-2", "value2"),
				),
			},
			wantName:    "deepseek",
			description: "Should create DeepSeek model with multiple custom request options",
		},
		{
			name: "default_base_url",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
			},
			wantName:    "deepseek",
			description: "Should use default DeepSeek base URL when not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewDeepSeekModel(tt.opts...)

			require.NoError(t, err, tt.description)
			require.NotNil(t, model, "Model should not be nil")
			assert.NotNil(t, model.OpenAICompletionModel, "OpenAICompletionModel should not be nil")
			assert.Equal(t, tt.wantName, model.Name(), "Model name should match expected value")
		})
	}
}

func TestNewDeepSeekModel_MissingAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		opts        []types.ModelOption
		description string
	}{
		{
			name: "empty_api_key",
			opts: []types.ModelOption{
				types.WithAPIKey(""),
			},
			description: "Should return error when API key is empty string",
		},
		{
			name:        "no_api_key",
			opts:        []types.ModelOption{},
			description: "Should return error when API key is not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewDeepSeekModel(tt.opts...)

			assert.Error(t, err, tt.description)
			assert.ErrorIs(t, err, types.ErrAPIKeyEmpty, "Error should be ErrAPIKeyEmpty")
			assert.Nil(t, model, "Model should be nil when error occurs")
		})
	}
}

func TestDeepSeekModel_Name(t *testing.T) {
	model, err := NewDeepSeekModel(types.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	assert.Equal(t, "deepseek", model.Name(), "Name should return 'deepseek'")
}

func TestDeepSeekModel_SupportedModels(t *testing.T) {
	model, err := NewDeepSeekModel(types.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	models := model.SupportedModels()
	assert.NotNil(t, models, "SupportedModels should not return nil")
	assert.Greater(t, len(models), 0, "SupportedModels should return at least one model")
}

func TestDeepSeekModel_ModelStructure(t *testing.T) {
	model, err := NewDeepSeekModel(types.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	t.Run("has_completion_model", func(t *testing.T) {
		assert.NotNil(t, model.OpenAICompletionModel, "OpenAICompletionModel should be initialized")
		assert.NotNil(t, model.OpenAICompletionModel.OpenAIBaseModel, "OpenAIBaseModel should be initialized")
	})
}

func TestNewDeepSeekModel_MultipleInstances(t *testing.T) {
	// Create multiple instances with different configurations
	model1, err1 := NewDeepSeekModel(types.WithAPIKey("test-api-key-1"))
	model2, err2 := NewDeepSeekModel(
		types.WithAPIKey("test-api-key-2"),
		types.WithBaseURL("https://custom.deepseek.com/"),
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
func BenchmarkNewDeepSeekModel_Success(b *testing.B) {
	opts := []types.ModelOption{
		types.WithAPIKey("test-api-key"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewDeepSeekModel(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewDeepSeekModel_WithOptions(b *testing.B) {
	opts := []types.ModelOption{
		types.WithAPIKey("test-api-key"),
		types.WithBaseURL("https://custom.deepseek.com/"),
		types.WithRequestOptions(
			option.WithHeader("Custom-Header-1", "value1"),
			option.WithHeader("Custom-Header-2", "value2"),
		),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewDeepSeekModel(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeepSeekModel_Name(b *testing.B) {
	model, err := NewDeepSeekModel(types.WithAPIKey("test-api-key"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.Name()
	}
}

func BenchmarkDeepSeekModel_SupportedModels(b *testing.B) {
	model, err := NewDeepSeekModel(types.WithAPIKey("test-api-key"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.SupportedModels()
	}
}
