// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"fmt"
	"github.com/easyagent-dev/llm"
	"net/url"
	"strings"
)

// Validation constants
const (
	MinTemperature      = 0.0
	MaxTemperature      = 2.0
	MinTopP             = 0.0
	MaxTopP             = 1.0
	MinPresencePenalty  = -2.0
	MaxPresencePenalty  = 2.0
	MinFrequencyPenalty = -2.0
	MaxFrequencyPenalty = 2.0
	MinMaxTokens        = 1
	MaxMaxTokens        = 1000000
)

// ValidateCompletionRequest performs comprehensive validation on llm requests
func ValidateCompletionRequest(req *llm.CompletionRequest) error {
	if req == nil {
		return llm.NewValidationError("request", "cannot be nil", nil)
	}

	// Validate messages
	if len(req.Messages) == 0 {
		return llm.NewValidationError("messages", "must contain at least one message", nil)
	}

	for i, msg := range req.Messages {
		if err := validateMessage(msg, i); err != nil {
			return err
		}
	}

	return nil
}

// validateMessage validates a single message
func validateMessage(msg *llm.ModelMessage, index int) error {
	if msg == nil {
		return llm.NewValidationError(fmt.Sprintf("messages[%d]", index), "cannot be nil", nil)
	}

	// Validate role
	if msg.Role == "" {
		return llm.NewValidationError(fmt.Sprintf("messages[%d].role", index), "cannot be empty", "")
	}

	validRoles := map[llm.Role]bool{
		llm.RoleUser:      true,
		llm.RoleAssistant: true,
		llm.RoleTool:      true,
	}

	if !validRoles[msg.Role] {
		return llm.NewValidationError(
			fmt.Sprintf("messages[%d].role", index),
			"must be one of: user, assistant, tool",
			string(msg.Role),
		)
	}

	// Content can be empty for tool calls, but check if both content and tool call are empty
	if msg.Content == "" && msg.ToolCall == nil && len(msg.Artifacts) == 0 {
		return llm.NewValidationError(
			fmt.Sprintf("messages[%d]", index),
			"must have either content, tool call, or artifacts",
			nil,
		)
	}

	// Validate artifacts if present
	for j, artifact := range msg.Artifacts {
		if err := validateArtifact(artifact, index, j); err != nil {
			return err
		}
	}

	return nil
}

// validateArtifact validates a model artifact
func validateArtifact(artifact *llm.ModelArtifact, msgIndex, artifactIndex int) error {
	if artifact == nil {
		return llm.NewValidationError(
			fmt.Sprintf("messages[%d].artifacts[%d]", msgIndex, artifactIndex),
			"cannot be nil",
			nil,
		)
	}

	if artifact.Name == "" {
		return llm.NewValidationError(
			fmt.Sprintf("messages[%d].artifacts[%d].name", msgIndex, artifactIndex),
			"cannot be empty",
			"",
		)
	}

	if artifact.ContentType == "" {
		return llm.NewValidationError(
			fmt.Sprintf("messages[%d].artifacts[%d].contentType", msgIndex, artifactIndex),
			"cannot be empty",
			"",
		)
	}

	return nil
}

// ValidateCompletionOptions validates llm configuration options
func ValidateCompletionOptions(config *llm.CompletionOptions) error {
	if config == nil {
		return nil // Config is optional
	}

	// Validate temperature
	if config.Temperature != nil {
		if *config.Temperature < MinTemperature || *config.Temperature > MaxTemperature {
			return llm.NewValidationError(
				"temperature",
				fmt.Sprintf("must be between %.1f and %.1f", MinTemperature, MaxTemperature),
				*config.Temperature,
			)
		}
	}

	// Validate top_p
	if config.TopP != nil {
		if *config.TopP < MinTopP || *config.TopP > MaxTopP {
			return llm.NewValidationError(
				"topP",
				fmt.Sprintf("must be between %.1f and %.1f", MinTopP, MaxTopP),
				*config.TopP,
			)
		}
	}

	// Validate max_tokens
	if config.MaxTokens != nil {
		if *config.MaxTokens < MinMaxTokens || *config.MaxTokens > MaxMaxTokens {
			return llm.NewValidationError(
				"maxTokens",
				fmt.Sprintf("must be between %d and %d", MinMaxTokens, MaxMaxTokens),
				*config.MaxTokens,
			)
		}
	}

	// Validate presence_penalty
	if config.PresencePenalty != nil {
		if *config.PresencePenalty < MinPresencePenalty || *config.PresencePenalty > MaxPresencePenalty {
			return llm.NewValidationError(
				"presencePenalty",
				fmt.Sprintf("must be between %.1f and %.1f", MinPresencePenalty, MaxPresencePenalty),
				*config.PresencePenalty,
			)
		}
	}

	// Validate frequency_penalty
	if config.FrequencyPenalty != nil {
		if *config.FrequencyPenalty < MinFrequencyPenalty || *config.FrequencyPenalty > MaxFrequencyPenalty {
			return llm.NewValidationError(
				"frequencyPenalty",
				fmt.Sprintf("must be between %.1f and %.1f", MinFrequencyPenalty, MaxFrequencyPenalty),
				*config.FrequencyPenalty,
			)
		}
	}

	// Validate reasoning effort
	if config.ReasoningEffort != nil {
		validEfforts := map[llm.ReasoningEffort]bool{
			llm.ReasoningEffortLow:    true,
			llm.ReasoningEffortMedium: true,
			llm.ReasoningEffortHigh:   true,
		}
		if !validEfforts[*config.ReasoningEffort] {
			return llm.NewValidationError(
				"reasoningEffort",
				"must be one of: low, medium, high",
				string(*config.ReasoningEffort),
			)
		}
	}

	// Validate response format
	if config.ResponseFormat != nil {
		validFormats := map[llm.ResponseFormat]bool{
			llm.ResponseFormatJson:       true,
			llm.ResponseFormatJsonSchema: true,
		}
		if !validFormats[*config.ResponseFormat] {
			return llm.NewValidationError(
				"responseFormat",
				"must be one of: json, json_schema",
				string(*config.ResponseFormat),
			)
		}

		// If json_schema is specified, JSONSchema must be provided
		if *config.ResponseFormat == llm.ResponseFormatJsonSchema && config.JSONSchema == nil {
			return llm.NewValidationError(
				"jsonSchema",
				"must be provided when responseFormat is json_schema",
				nil,
			)
		}
	}

	return nil
}

// ValidateEmbeddingRequestWithDetails validates llm request with detailed errors
func ValidateEmbeddingRequestWithDetails(req *llm.EmbeddingRequest) error {
	if req == nil {
		return llm.NewValidationError("request", "cannot be nil", nil)
	}

	if req.Model == "" {
		return llm.NewValidationError("model", "cannot be empty", "")
	}

	if len(req.Contents) == 0 {
		return llm.NewValidationError("contents", "must contain at least one item", nil)
	}

	// Validate each content item
	for i, content := range req.Contents {
		if strings.TrimSpace(content) == "" {
			return llm.NewValidationError(
				fmt.Sprintf("contents[%d]", i),
				"cannot be empty or whitespace only",
				content,
			)
		}
	}

	// Validate config if provided
	if req.Config != nil {
		if err := validateEmbeddingConfig(req.Config); err != nil {
			return err
		}
	}

	return nil
}

// validateEmbeddingConfig validates llm configuration
func validateEmbeddingConfig(config *llm.EmbeddingModelConfig) error {
	if config == nil {
		return nil
	}

	// Validate dimensions
	if config.Dimensions < 0 {
		return llm.NewValidationError(
			"dimensions",
			"must be a positive number or 0 for default",
			config.Dimensions,
		)
	}

	// Validate encoding format
	if config.EncodingFormat != "" {
		validFormats := map[llm.EmbeddingEncodingFormat]bool{
			llm.EmbeddingEncodingFormatFloat:  true,
			llm.EmbeddingEncodingFormatBase64: true,
		}
		if !validFormats[config.EncodingFormat] {
			return llm.NewValidationError(
				"encodingFormat",
				"must be one of: float, base64",
				string(config.EncodingFormat),
			)
		}
	}

	return nil
}

// ValidateImageRequestWithDetails validates llm request with detailed errors
func ValidateImageRequestWithDetails(req *llm.ImageRequest) error {
	if req == nil {
		return llm.NewValidationError("request", "cannot be nil", nil)
	}

	if req.Model == "" {
		return llm.NewValidationError("model", "cannot be empty", "")
	}

	if strings.TrimSpace(req.Instructions) == "" {
		return llm.NewValidationError("instructions", "cannot be empty or whitespace only", req.Instructions)
	}

	// Validate config if provided
	if req.Config != nil {
		if err := validateImageConfig(req.Config); err != nil {
			return err
		}
	}

	return nil
}

// validateImageConfig validates llm configuration
func validateImageConfig(config *llm.ImageModelConfig) error {
	if config == nil {
		return nil
	}

	// Validate size format
	if config.Size != "" {
		validSizes := map[string]bool{
			"256x256":   true,
			"512x512":   true,
			"1024x1024": true,
			"1792x1024": true,
			"1024x1792": true,
		}
		if !validSizes[config.Size] {
			return llm.NewValidationError(
				"size",
				"must be one of: 256x256, 512x512, 1024x1024, 1792x1024, 1024x1792",
				config.Size,
			)
		}
	}

	// Validate quality
	if config.Quality != "" {
		validQualities := map[string]bool{
			"standard": true,
			"hd":       true,
		}
		if !validQualities[config.Quality] {
			return llm.NewValidationError(
				"quality",
				"must be one of: standard, hd",
				config.Quality,
			)
		}
	}

	// Validate style
	if config.Style != "" {
		validStyles := map[string]bool{
			"vivid":   true,
			"natural": true,
		}
		if !validStyles[config.Style] {
			return llm.NewValidationError(
				"style",
				"must be one of: vivid, natural",
				config.Style,
			)
		}
	}

	return nil
}

// ValidateAPIKey validates API key format
func ValidateAPIKey(apiKey string) error {
	if apiKey == "" {
		return llm.ErrAPIKeyEmpty
	}

	// Trim whitespace
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return llm.NewValidationError("apiKey", "cannot be whitespace only", apiKey)
	}

	// Check minimum length (most API keys are at least 20 characters)
	if len(apiKey) < 10 {
		return llm.NewValidationError("apiKey", "appears to be too short", len(apiKey))
	}

	return nil
}

// ValidateBaseURL validates base URL format
func ValidateBaseURL(baseURL string) error {
	if baseURL == "" {
		return llm.ErrBaseURLEmpty
	}

	// Parse URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return llm.NewValidationError("baseURL", "invalid URL format", baseURL)
	}

	// Check scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return llm.NewValidationError("baseURL", "must use http or https scheme", parsedURL.Scheme)
	}

	// Check host
	if parsedURL.Host == "" {
		return llm.NewValidationError("baseURL", "must have a valid host", baseURL)
	}

	return nil
}

// ValidateModelName validates that a model name is not empty and doesn't contain invalid characters
func ValidateModelName(modelName string) error {
	if modelName == "" {
		return llm.NewValidationError("model", "cannot be empty", "")
	}

	modelName = strings.TrimSpace(modelName)
	if modelName == "" {
		return llm.NewValidationError("model", "cannot be whitespace only", modelName)
	}

	// Model names should not contain certain characters
	invalidChars := []string{"\n", "\r", "\t", "\x00"}
	for _, char := range invalidChars {
		if strings.Contains(modelName, char) {
			return llm.NewValidationError("model", "contains invalid characters", modelName)
		}
	}

	return nil
}
