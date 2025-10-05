// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package openai

import (
	"testing"

	"github.com/easyagent-dev/llm"
	"github.com/openai/openai-go/v3/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewOpenAIModelProvider_Success tests successful creation of OpenAI model provider
func TestNewOpenAIModelProvider_Success(t *testing.T) {
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
			wantName:    "openai",
			description: "Should create OpenAI model with basic configuration",
		},
		{
			name: "with_custom_base_url",
			opts: []llm.ModelOption{
				llm.WithAPIKey("test-api-key"),
				llm.WithBaseURL("https://custom.openai.com/"),
			},
			wantName:    "openai",
			description: "Should create OpenAI model with custom base URL",
		},
		{
			name: "with_custom_request_option",
			opts: []llm.ModelOption{
				llm.WithAPIKey("test-api-key"),
				llm.WithRequestOption(option.WithHeader("Custom-Header", "custom-value")),
			},
			wantName:    "openai",
			description: "Should create OpenAI model with custom request options",
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
			wantName:    "openai",
			description: "Should create OpenAI model with multiple custom request options",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewOpenAIModelProvider(tt.opts...)

			require.NoError(t, err, tt.description)
			require.NotNil(t, provider, "Provider should not be nil")
			assert.Equal(t, tt.wantName, provider.Name(), "Provider name should match expected value")
		})
	}
}

// TestOpenAIModelProvider_Name tests the Name() method
func TestOpenAIModelProvider_Name(t *testing.T) {
	provider, err := NewOpenAIModelProvider(llm.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, provider)

	assert.Equal(t, "openai", provider.Name(), "Name should return 'openai'")
}

// TestOpenAIModelProvider_SupportedModels tests the SupportedModels() method
func TestOpenAIModelProvider_SupportedModels(t *testing.T) {
	provider, err := NewOpenAIModelProvider(llm.WithAPIKey("test-api-key"))

	require.NoError(t, err)
	require.NotNil(t, provider)

	models := provider.SupportedModels()
	assert.NotNil(t, models, "SupportedModels should not return nil")
	assert.Greater(t, len(models), 0, "SupportedModels should return at least one model")

	// Verify each model has required fields
	for _, model := range models {
		assert.NotEmpty(t, model.Name, "Model name should not be empty")
		assert.NotEmpty(t, model.ID, "Model ID should not be empty")
	}
}

// TestOpenAIModelProvider_NewCompletionModel tests completion model creation
func TestOpenAIModelProvider_NewCompletionModel(t *testing.T) {
	provider, err := NewOpenAIModelProvider(llm.WithAPIKey("test-api-key"))
	require.NoError(t, err)
	require.NotNil(t, provider)

	tests := []struct {
		name        string
		model       string
		opts        []llm.CompletionOption
		expectError bool
		description string
	}{
		{
			name:        "valid_gpt4_model",
			model:       "gpt-4o",
			opts:        []llm.CompletionOption{},
			expectError: false,
			description: "Should create completion model for gpt-4o",
		},
		{
			name:        "valid_gpt4_mini_model",
			model:       "gpt-4o-mini",
			opts:        []llm.CompletionOption{},
			expectError: false,
			description: "Should create completion model for gpt-4o-mini",
		},
		{
			name:  "with_options",
			model: "gpt-4o",
			opts: []llm.CompletionOption{
				llm.WithTemperature(0.7),
				llm.WithMaxTokens(1000),
			},
			expectError: false,
			description: "Should create completion model with options",
		},
		{
			name:        "invalid_model",
			model:       "non-existent-model",
			opts:        []llm.CompletionOption{},
			expectError: true,
			description: "Should return error for non-existent model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := provider.NewCompletionModel(tt.model, tt.opts...)

			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, model, "Model should be nil when error occurs")
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, model, "Model should not be nil")
			}
		})
	}
}

// TestOpenAIModelProvider_NewImageModel tests image model creation
func TestOpenAIModelProvider_NewImageModel(t *testing.T) {
	provider, err := NewOpenAIModelProvider(llm.WithAPIKey("test-api-key"))
	require.NoError(t, err)
	require.NotNil(t, provider)

	tests := []struct {
		name        string
		model       string
		expectError bool
		description string
	}{
		{
			name:        "valid_image_model",
			model:       "gpt-image-1",
			expectError: false,
			description: "Should create image model for gpt-image-1",
		},
		{
			name:        "invalid_model",
			model:       "non-existent-image-model",
			expectError: true,
			description: "Should return error for non-existent model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := provider.NewImageModel(tt.model)

			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, model, "Model should be nil when error occurs")
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, model, "Model should not be nil")
			}
		})
	}
}

// TestOpenAIModelProvider_NewConversationModel tests conversation model creation
func TestOpenAIModelProvider_NewConversationModel(t *testing.T) {
	provider, err := NewOpenAIModelProvider(llm.WithAPIKey("test-api-key"))
	require.NoError(t, err)
	require.NotNil(t, provider)

	tests := []struct {
		name        string
		model       string
		opts        []llm.ResponseOption
		expectError bool
		description string
	}{
		{
			name:        "valid_gpt4_model",
			model:       "gpt-4o",
			opts:        []llm.ResponseOption{},
			expectError: false,
			description: "Should create conversation model for gpt-4o",
		},
		{
			name:  "with_options",
			model: "gpt-4o",
			opts: []llm.ResponseOption{
				llm.WithOptions(llm.WithTemperature(0.7)),
			},
			expectError: false,
			description: "Should create conversation model with options",
		},
		{
			name:        "invalid_model",
			model:       "non-existent-model",
			opts:        []llm.ResponseOption{},
			expectError: true,
			description: "Should return error for non-existent model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := provider.NewConversationModel(tt.model, tt.opts...)

			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, model, "Model should be nil when error occurs")
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, model, "Model should not be nil")
			}
		})
	}
}

// TestToChatCompletionMessage tests message conversion
func TestToChatCompletionMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     *llm.ModelMessage
		expectError bool
		description string
	}{
		{
			name: "user_message",
			message: &llm.ModelMessage{
				Role:    llm.RoleUser,
				Content: "Hello, world!",
			},
			expectError: false,
			description: "Should convert user message",
		},
		{
			name: "assistant_message",
			message: &llm.ModelMessage{
				Role:    llm.RoleAssistant,
				Content: "Hello! How can I help you?",
			},
			expectError: false,
			description: "Should convert assistant message",
		},
		{
			name: "assistant_with_tool_call",
			message: &llm.ModelMessage{
				Role:    llm.RoleAssistant,
				Content: "Calling tool",
				ToolCall: &llm.ToolCall{
					Name:  "test_tool",
					Input: map[string]any{"arg1": "value1"},
				},
			},
			expectError: false,
			description: "Should convert assistant message with tool call",
		},
		{
			name: "tool_message",
			message: &llm.ModelMessage{
				Role: llm.RoleTool,
				ToolCall: &llm.ToolCall{
					Name:   "test_tool",
					Output: "success",
				},
			},
			expectError: false,
			description: "Should convert tool message",
		},
		{
			name:        "nil_message",
			message:     nil,
			expectError: true,
			description: "Should return error for nil message",
		},
		{
			name: "invalid_role",
			message: &llm.ModelMessage{
				Role:    "invalid",
				Content: "test",
			},
			expectError: true,
			description: "Should return error for invalid role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToChatCompletionMessage(tt.message)

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, result, "Result should not be nil")
			}
		})
	}
}

// TestToChatCompletionParams tests parameter conversion
func TestToChatCompletionParams(t *testing.T) {
	tests := []struct {
		name         string
		model        string
		instructions string
		messages     []*llm.ModelMessage
		opts         *llm.CompletionOptions
		description  string
	}{
		{
			name:         "basic_params",
			model:        "gpt-4o",
			instructions: "You are a helpful assistant",
			messages: []*llm.ModelMessage{
				{Role: llm.RoleUser, Content: "Hello"},
			},
			opts:        nil,
			description: "Should create basic chat completion params",
		},
		{
			name:         "with_temperature",
			model:        "gpt-4o",
			instructions: "You are a helpful assistant",
			messages: []*llm.ModelMessage{
				{Role: llm.RoleUser, Content: "Hello"},
			},
			opts: &llm.CompletionOptions{
				Temperature: func() *float64 { v := 0.7; return &v }(),
			},
			description: "Should create params with temperature",
		},
		{
			name:         "with_max_tokens",
			model:        "gpt-4o",
			instructions: "",
			messages: []*llm.ModelMessage{
				{Role: llm.RoleUser, Content: "Hello"},
			},
			opts: &llm.CompletionOptions{
				MaxTokens: func() *int { v := 1000; return &v }(),
			},
			description: "Should create params with max tokens",
		},
		{
			name:         "with_json_response_format",
			model:        "gpt-4o",
			instructions: "You are a helpful assistant",
			messages: []*llm.ModelMessage{
				{Role: llm.RoleUser, Content: "Hello"},
			},
			opts: &llm.CompletionOptions{
				ResponseFormat: func() *llm.ResponseFormat { v := llm.ResponseFormatJson; return &v }(),
			},
			description: "Should create params with JSON response format",
		},
		{
			name:         "empty_messages",
			model:        "gpt-4o",
			instructions: "You are a helpful assistant",
			messages:     []*llm.ModelMessage{},
			opts:         nil,
			description:  "Should handle empty messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := ToChatCompletionParams(tt.model, tt.instructions, tt.messages, tt.opts)

			assert.NoError(t, err, tt.description)
			assert.Equal(t, tt.model, string(params.Model), "Model should match")
			assert.NotNil(t, params.Messages, "Messages should not be nil")
		})
	}
}

// TestNewOpenAIModelProvider_MultipleInstances tests creating multiple instances
func TestNewOpenAIModelProvider_MultipleInstances(t *testing.T) {
	provider1, err1 := NewOpenAIModelProvider(llm.WithAPIKey("test-api-key-1"))
	provider2, err2 := NewOpenAIModelProvider(
		llm.WithAPIKey("test-api-key-2"),
		llm.WithBaseURL("https://custom.openai.com/"),
	)

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NotNil(t, provider1)
	require.NotNil(t, provider2)

	// Verify instances are independent
	assert.NotSame(t, provider1, provider2, "Different instances should be created")
}

// Benchmark tests
func BenchmarkNewOpenAIModelProvider_Success(b *testing.B) {
	opts := []llm.ModelOption{
		llm.WithAPIKey("test-api-key"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewOpenAIModelProvider(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewOpenAIModelProvider_WithOptions(b *testing.B) {
	opts := []llm.ModelOption{
		llm.WithAPIKey("test-api-key"),
		llm.WithBaseURL("https://custom.openai.com/"),
		llm.WithRequestOptions(
			option.WithHeader("Custom-Header-1", "value1"),
			option.WithHeader("Custom-Header-2", "value2"),
		),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewOpenAIModelProvider(opts...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOpenAIModelProvider_Name(b *testing.B) {
	provider, err := NewOpenAIModelProvider(llm.WithAPIKey("test-api-key"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.Name()
	}
}

func BenchmarkOpenAIModelProvider_SupportedModels(b *testing.B) {
	provider, err := NewOpenAIModelProvider(llm.WithAPIKey("test-api-key"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.SupportedModels()
	}
}

func BenchmarkToChatCompletionMessage(b *testing.B) {
	msg := &llm.ModelMessage{
		Role:    llm.RoleUser,
		Content: "Hello, world!",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ToChatCompletionMessage(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkToChatCompletionParams(b *testing.B) {
	messages := []*llm.ModelMessage{
		{Role: llm.RoleUser, Content: "Hello"},
		{Role: llm.RoleAssistant, Content: "Hi there!"},
		{Role: llm.RoleUser, Content: "How are you?"},
	}
	opts := &llm.CompletionOptions{
		Temperature: func() *float64 { v := 0.7; return &v }(),
		MaxTokens:   func() *int { v := 1000; return &v }(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ToChatCompletionParams("gpt-4o", "You are a helpful assistant", messages, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}
