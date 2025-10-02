package openai

import (
	"testing"

	"github.com/easymvp/easyllm/types"
	"github.com/openai/openai-go/v3/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOpenAIModel_Success(t *testing.T) {
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
			wantName:    "openai",
			description: "Should create OpenAI model with basic configuration",
		},
		{
			name: "with_custom_base_url",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithBaseURL("https://custom.openai.com/v1"),
			},
			wantName:    "openai",
			description: "Should create OpenAI model with custom base URL",
		},
		{
			name: "with_custom_request_option",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithRequestOption(option.WithHeader("Custom-Header", "custom-value")),
			},
			wantName:    "openai",
			description: "Should create OpenAI model with custom request options",
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
			wantName:    "openai",
			description: "Should create OpenAI model with multiple custom request options",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewOpenAIModel(tt.opts...)

			require.NoError(t, err, tt.description)
			require.NotNil(t, model, "Model should not be nil")
			assert.NotNil(t, model.OpenAICompletionModel, "OpenAICompletionModel should not be nil")
			assert.NotNil(t, model.OpenAIEmbeddingModel, "OpenAIEmbeddingModel should not be nil")
			assert.NotNil(t, model.OpenAIImageModel, "OpenAIImageModel should not be nil")
			assert.Equal(t, tt.wantName, model.Name(), "Model name should match expected value")
		})
	}
}

func TestNewOpenAIModel_MissingAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		opts        []types.ModelOption
		description string
	}{
		{
			name:        "empty_api_key",
			opts:        []types.ModelOption{types.WithAPIKey("")},
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
			model, err := NewOpenAIModel(tt.opts...)

			assert.Error(t, err, tt.description)
			assert.ErrorIs(t, err, types.ErrAPIKeyEmpty, "Error should be ErrAPIKeyEmpty")
			assert.Nil(t, model, "Model should be nil when error occurs")
		})
	}
}

func TestOpenAIModel_Name(t *testing.T) {
	model, err := NewOpenAIModel(types.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	assert.Equal(t, "openai", model.Name(), "Name should return 'openai'")
}

func TestOpenAIModel_ModelStructure(t *testing.T) {
	model, err := NewOpenAIModel(types.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	t.Run("has_completion_model", func(t *testing.T) {
		assert.NotNil(t, model.OpenAICompletionModel, "OpenAICompletionModel should be initialized")
		assert.NotNil(t, model.OpenAICompletionModel.OpenAIBaseModel, "OpenAIBaseModel should be initialized")
	})

	t.Run("has_embedding_model", func(t *testing.T) {
		assert.NotNil(t, model.OpenAIEmbeddingModel, "OpenAIEmbeddingModel should be initialized")
		assert.NotNil(t, model.OpenAIEmbeddingModel.OpenAIBaseModel, "OpenAIBaseModel should be initialized")
	})

	t.Run("has_image_model", func(t *testing.T) {
		assert.NotNil(t, model.OpenAIImageModel, "OpenAIImageModel should be initialized")
		assert.NotNil(t, model.OpenAIImageModel.OpenAIBaseModel, "OpenAIBaseModel should be initialized")
	})
}

func TestOpenAIModel_SupportedModels(t *testing.T) {
	model, err := NewOpenAIModel(types.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, model)

	models := model.SupportedModels()
	assert.NotNil(t, models, "SupportedModels should not return nil")
	assert.Greater(t, len(models), 0, "SupportedModels should return at least one model")
}

func TestOpenAIModel_MultipleInstances(t *testing.T) {
	// Create multiple instances with different configurations
	model1, err1 := NewOpenAIModel(types.WithAPIKey("test-api-key-1"))
	model2, err2 := NewOpenAIModel(types.WithAPIKey("test-api-key-2"))

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NotNil(t, model1)
	require.NotNil(t, model2)

	// Verify instances are independent
	assert.NotSame(t, model1, model2, "Different instances should be created")
	assert.NotSame(t, model1.OpenAICompletionModel, model2.OpenAICompletionModel, "Completion models should be independent")
	assert.NotSame(t, model1.OpenAIEmbeddingModel, model2.OpenAIEmbeddingModel, "Embedding models should be independent")
	assert.NotSame(t, model1.OpenAIImageModel, model2.OpenAIImageModel, "Image models should be independent")
}

func TestNewOpenAIBaseModel_Success(t *testing.T) {
	model, err := NewOpenAIBaseModel("test-api-key")

	require.NoError(t, err)
	require.NotNil(t, model)
	assert.NotNil(t, model.modelCache, "Model cache should be initialized")
}

func TestNewOpenAIBaseModel_MissingAPIKey(t *testing.T) {
	model, err := NewOpenAIBaseModel("")

	assert.Error(t, err)
	assert.ErrorIs(t, err, types.ErrAPIKeyEmpty)
	assert.Nil(t, model)
}

func TestOpenAIBaseModel_ClearModelCache(t *testing.T) {
	model, err := NewOpenAIBaseModel("test-api-key")
	require.NoError(t, err)

	// Add something to cache by calling getModelInfo
	_ = model.getModelInfo("gpt-4")

	// Clear cache
	model.ClearModelCache()

	// Verify cache is empty
	assert.Equal(t, 0, len(model.modelCache), "Cache should be empty after clearing")
}

func TestNewOpenAICompletionModel_Success(t *testing.T) {
	model, err := NewOpenAICompletionModel("test-api-key")

	require.NoError(t, err)
	require.NotNil(t, model)
	assert.NotNil(t, model.OpenAIBaseModel, "OpenAIBaseModel should be initialized")
}

func TestNewOpenAICompletionModel_MissingAPIKey(t *testing.T) {
	model, err := NewOpenAICompletionModel("")

	assert.Error(t, err)
	assert.ErrorIs(t, err, types.ErrAPIKeyEmpty)
	assert.Nil(t, model)
}

func TestNewOpenAIEmbeddingModel_Success(t *testing.T) {
	model, err := NewOpenAIEmbeddingModel("test-api-key")

	require.NoError(t, err)
	require.NotNil(t, model)
	assert.NotNil(t, model.OpenAIBaseModel, "OpenAIBaseModel should be initialized")
}

func TestNewOpenAIEmbeddingModel_MissingAPIKey(t *testing.T) {
	model, err := NewOpenAIEmbeddingModel("")

	assert.Error(t, err)
	assert.ErrorIs(t, err, types.ErrAPIKeyEmpty)
	assert.Nil(t, model)
}

func TestNewOpenAIImageModel_Success(t *testing.T) {
	model, err := NewOpenAIImageModel("test-api-key")

	require.NoError(t, err)
	require.NotNil(t, model)
	assert.NotNil(t, model.OpenAIBaseModel, "OpenAIBaseModel should be initialized")
}

func TestNewOpenAIImageModel_MissingAPIKey(t *testing.T) {
	model, err := NewOpenAIImageModel("")

	assert.Error(t, err)
	assert.ErrorIs(t, err, types.ErrAPIKeyEmpty)
	assert.Nil(t, model)
}

// Benchmark tests
func BenchmarkNewOpenAIModel_Success(b *testing.B) {
	opts := []types.ModelOption{
		types.WithAPIKey("test-api-key"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewOpenAIModel(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewOpenAIModel_WithOptions(b *testing.B) {
	opts := []types.ModelOption{
		types.WithAPIKey("test-api-key"),
		types.WithBaseURL("https://custom.openai.com/v1"),
		types.WithRequestOptions(
			option.WithHeader("Custom-Header-1", "value1"),
			option.WithHeader("Custom-Header-2", "value2"),
		),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewOpenAIModel(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOpenAIModel_Name(b *testing.B) {
	model, err := NewOpenAIModel(types.WithAPIKey("test-api-key"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.Name()
	}
}

func BenchmarkOpenAIModel_SupportedModels(b *testing.B) {
	model, err := NewOpenAIModel(types.WithAPIKey("test-api-key"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.SupportedModels()
	}
}
