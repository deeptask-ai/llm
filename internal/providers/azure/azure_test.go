package azure

import (
	"testing"

	"github.com/easymvp/easyllm/types"
	"github.com/openai/openai-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAzureOpenAIModel_Success(t *testing.T) {
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
				types.WithBaseURL("https://test.openai.azure.com"),
				types.WithAPIVersion("2024-02-15-preview"),
			},
			wantName:    "azure_openai",
			description: "Should create Azure OpenAI model with basic configuration",
		},
		{
			name: "with_custom_request_option",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithBaseURL("https://test.openai.azure.com"),
				types.WithAPIVersion("2024-02-15-preview"),
				types.WithRequestOption(option.WithHeader("Custom-Header", "custom-value")),
			},
			wantName:    "azure_openai",
			description: "Should create Azure OpenAI model with custom request options",
		},
		{
			name: "with_multiple_request_options",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithBaseURL("https://test.openai.azure.com"),
				types.WithAPIVersion("2024-02-15-preview"),
				types.WithRequestOptions(
					option.WithHeader("Custom-Header-1", "value1"),
					option.WithHeader("Custom-Header-2", "value2"),
				),
			},
			wantName:    "azure_openai",
			description: "Should create Azure OpenAI model with multiple custom request options",
		},
		{
			name: "different_api_version",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithBaseURL("https://test.openai.azure.com"),
				types.WithAPIVersion("2023-05-15"),
			},
			wantName:    "azure_openai",
			description: "Should create Azure OpenAI model with different API version",
		},
		{
			name: "custom_base_url",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithBaseURL("https://custom.openai.azure.com/openai/deployments/gpt-4"),
				types.WithAPIVersion("2024-02-15-preview"),
			},
			wantName:    "azure_openai",
			description: "Should create Azure OpenAI model with custom base URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewAzureOpenAIModel(tt.opts...)

			require.NoError(t, err, tt.description)
			require.NotNil(t, model, "Model should not be nil")
			assert.NotNil(t, model.OpenAIModel, "OpenAIModel should not be nil")
			assert.NotNil(t, model.OpenAICompletionModel, "OpenAICompletionModel should not be nil")
			assert.NotNil(t, model.OpenAIEmbeddingModel, "OpenAIEmbeddingModel should not be nil")
			assert.NotNil(t, model.OpenAIImageModel, "OpenAIImageModel should not be nil")
			assert.Equal(t, tt.wantName, model.Name(), "Model name should match expected value")
		})
	}
}

func TestNewAzureOpenAIModel_MissingAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		opts        []types.ModelOption
		description string
	}{
		{
			name: "empty_api_key",
			opts: []types.ModelOption{
				types.WithAPIKey(""),
				types.WithBaseURL("https://test.openai.azure.com"),
				types.WithAPIVersion("2024-02-15-preview"),
			},
			description: "Should return error when API key is empty string",
		},
		{
			name: "no_api_key",
			opts: []types.ModelOption{
				types.WithBaseURL("https://test.openai.azure.com"),
				types.WithAPIVersion("2024-02-15-preview"),
			},
			description: "Should return error when API key is not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewAzureOpenAIModel(tt.opts...)

			assert.Error(t, err, tt.description)
			assert.ErrorIs(t, err, types.ErrAPIKeyEmpty, "Error should be ErrAPIKeyEmpty")
			assert.Nil(t, model, "Model should be nil when error occurs")
		})
	}
}

func TestNewAzureOpenAIModel_MissingBaseURL(t *testing.T) {
	tests := []struct {
		name        string
		opts        []types.ModelOption
		description string
	}{
		{
			name: "empty_base_url",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithBaseURL(""),
				types.WithAPIVersion("2024-02-15-preview"),
			},
			description: "Should return error when base URL is empty string",
		},
		{
			name: "no_base_url",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithAPIVersion("2024-02-15-preview"),
			},
			description: "Should return error when base URL is not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewAzureOpenAIModel(tt.opts...)

			assert.Error(t, err, tt.description)
			assert.ErrorIs(t, err, types.ErrBaseURLEmpty, "Error should be ErrBaseURLEmpty")
			assert.Nil(t, model, "Model should be nil when error occurs")
		})
	}
}

func TestNewAzureOpenAIModel_MissingAPIVersion(t *testing.T) {
	tests := []struct {
		name        string
		opts        []types.ModelOption
		description string
	}{
		{
			name: "empty_api_version",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithBaseURL("https://test.openai.azure.com"),
				types.WithAPIVersion(""),
			},
			description: "Should return error when API version is empty string",
		},
		{
			name: "no_api_version",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithBaseURL("https://test.openai.azure.com"),
			},
			description: "Should return error when API version is not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewAzureOpenAIModel(tt.opts...)

			assert.Error(t, err, tt.description)
			assert.ErrorIs(t, err, types.ErrAPIVersionEmpty, "Error should be ErrAPIVersionEmpty")
			assert.Nil(t, model, "Model should be nil when error occurs")
		})
	}
}

func TestNewAzureOpenAIModel_ErrorPriority(t *testing.T) {
	tests := []struct {
		name        string
		opts        []types.ModelOption
		expectedErr error
		description string
	}{
		{
			name: "missing_api_key_first",
			opts: []types.ModelOption{
				types.WithBaseURL("https://test.openai.azure.com"),
				types.WithAPIVersion("2024-02-15-preview"),
			},
			expectedErr: types.ErrAPIKeyEmpty,
			description: "Should check API key first",
		},
		{
			name: "missing_base_url_second",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithAPIVersion("2024-02-15-preview"),
			},
			expectedErr: types.ErrBaseURLEmpty,
			description: "Should check base URL after API key",
		},
		{
			name: "missing_api_version_third",
			opts: []types.ModelOption{
				types.WithAPIKey("test-api-key"),
				types.WithBaseURL("https://test.openai.azure.com"),
			},
			expectedErr: types.ErrAPIVersionEmpty,
			description: "Should check API version after base URL",
		},
		{
			name:        "missing_all_required",
			opts:        []types.ModelOption{},
			expectedErr: types.ErrAPIKeyEmpty,
			description: "Should return API key error first when all are missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewAzureOpenAIModel(tt.opts...)

			assert.Error(t, err, tt.description)
			assert.ErrorIs(t, err, tt.expectedErr, "Error should match expected error type")
			assert.Nil(t, model, "Model should be nil when error occurs")
		})
	}
}

func TestAzureOpenAIModel_Name(t *testing.T) {
	model, err := NewAzureOpenAIModel(
		types.WithAPIKey("test-api-key"),
		types.WithBaseURL("https://test.openai.azure.com"),
		types.WithAPIVersion("2024-02-15-preview"),
	)

	require.NoError(t, err)
	require.NotNil(t, model)

	assert.Equal(t, "azure_openai", model.Name(), "Name should return 'azure_openai'")
}

func TestAzureOpenAIModel_ModelStructure(t *testing.T) {
	model, err := NewAzureOpenAIModel(
		types.WithAPIKey("test-api-key"),
		types.WithBaseURL("https://test.openai.azure.com"),
		types.WithAPIVersion("2024-02-15-preview"),
	)

	require.NoError(t, err)
	require.NotNil(t, model)

	t.Run("has_openai_model", func(t *testing.T) {
		assert.NotNil(t, model.OpenAIModel, "OpenAIModel should be initialized")
	})

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

func TestNewAzureOpenAIModel_MultipleInstances(t *testing.T) {
	// Create multiple instances with different configurations
	model1, err1 := NewAzureOpenAIModel(
		types.WithAPIKey("test-api-key-1"),
		types.WithBaseURL("https://test1.openai.azure.com"),
		types.WithAPIVersion("2024-02-15-preview"),
	)

	model2, err2 := NewAzureOpenAIModel(
		types.WithAPIKey("test-api-key-2"),
		types.WithBaseURL("https://test2.openai.azure.com"),
		types.WithAPIVersion("2023-05-15"),
	)

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NotNil(t, model1)
	require.NotNil(t, model2)

	// Verify instances are independent
	assert.NotSame(t, model1, model2, "Different instances should be created")
	assert.NotSame(t, model1.OpenAIModel, model2.OpenAIModel, "OpenAI models should be independent")
}

// Benchmark tests
func BenchmarkNewAzureOpenAIModel_Success(b *testing.B) {
	opts := []types.ModelOption{
		types.WithAPIKey("test-api-key"),
		types.WithBaseURL("https://test.openai.azure.com"),
		types.WithAPIVersion("2024-02-15-preview"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewAzureOpenAIModel(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewAzureOpenAIModel_WithOptions(b *testing.B) {
	opts := []types.ModelOption{
		types.WithAPIKey("test-api-key"),
		types.WithBaseURL("https://test.openai.azure.com"),
		types.WithAPIVersion("2024-02-15-preview"),
		types.WithRequestOptions(
			option.WithHeader("Custom-Header-1", "value1"),
			option.WithHeader("Custom-Header-2", "value2"),
		),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewAzureOpenAIModel(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAzureOpenAIModel_Name(b *testing.B) {
	model, err := NewAzureOpenAIModel(
		types.WithAPIKey("test-api-key"),
		types.WithBaseURL("https://test.openai.azure.com"),
		types.WithAPIVersion("2024-02-15-preview"),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.Name()
	}
}
