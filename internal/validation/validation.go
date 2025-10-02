// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package validation

import "github.com/easymvp/easyllm"

import (
	"fmt"
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

// ValidateCompletionRequest performs comprehensive validation on completion requests
func ValidateCompletionRequest(req *easyllm.CompletionRequest) error {
	if req == nil {
		return easyllm.NewValidationError("request", "cannot be nil", nil)
	}

	// Validate model
	if req.Model == "" {
		return easyllm.NewValidationError("model", "cannot be empty", "")
	}

	// Validate messages
	if len(req.Messages) == 0 {
		return easyllm.NewValidationError("messages", "must contain at least one message", nil)
	}

	for i, msg := range req.Messages {
		if err := validateMessage(msg, i); err != nil {
			return err
		}
	}

	// Validate config if provided
	if req.Config != nil {
		if err := ValidateModelConfig(req.Config); err != nil {
			return err
		}
	}

	return nil
}

// validateMessage validates a single message
func validateMessage(msg *easyllm.ModelMessage, index int) error {
	if msg == nil {
		return easyllm.NewValidationError(fmt.Sprintf("messages[%d]", index), "cannot be nil", nil)
	}

	// Validate role
	if msg.Role == "" {
		return easyllm.NewValidationError(fmt.Sprintf("messages[%d].role", index), "cannot be empty", "")
	}

	validRoles := map[easyllm.MessageRole]bool{
		easyllm.MessageRoleUser:      true,
		easyllm.MessageRoleAssistant: true,
		easyllm.MessageRoleTool:      true,
	}

	if !validRoles[msg.Role] {
		return easyllm.NewValidationError(
			fmt.Sprintf("messages[%d].role", index),
			"must be one of: user, assistant, tool",
			string(msg.Role),
		)
	}

	// Content can be empty for tool calls, but check if both content and tool call are empty
	if msg.Content == "" && msg.ToolCall == nil && len(msg.Artifacts) == 0 {
		return easyllm.NewValidationError(
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
func validateArtifact(artifact *easyllm.ModelArtifact, msgIndex, artifactIndex int) error {
	if artifact == nil {
		return easyllm.NewValidationError(
			fmt.Sprintf("messages[%d].artifacts[%d]", msgIndex, artifactIndex),
			"cannot be nil",
			nil,
		)
	}

	if artifact.Name == "" {
		return easyllm.NewValidationError(
			fmt.Sprintf("messages[%d].artifacts[%d].name", msgIndex, artifactIndex),
			"cannot be empty",
			"",
		)
	}

	if artifact.ContentType == "" {
		return easyllm.NewValidationError(
			fmt.Sprintf("messages[%d].artifacts[%d].contentType", msgIndex, artifactIndex),
			"cannot be empty",
			"",
		)
	}

	return nil
}

// ValidateModelConfig validates model configuration parameters
func ValidateModelConfig(config *easyllm.ModelConfig) error {
	if config == nil {
		return nil // Config is optional
	}

	// Validate temperature
	if config.Temperature != 0 {
		if config.Temperature < MinTemperature || config.Temperature > MaxTemperature {
			return easyllm.NewValidationError(
				"temperature",
				fmt.Sprintf("must be between %.1f and %.1f", MinTemperature, MaxTemperature),
				config.Temperature,
			)
		}
	}

	// Validate top_p
	if config.TopP != 0 {
		if config.TopP < MinTopP || config.TopP > MaxTopP {
			return easyllm.NewValidationError(
				"topP",
				fmt.Sprintf("must be between %.1f and %.1f", MinTopP, MaxTopP),
				config.TopP,
			)
		}
	}

	// Validate max_tokens
	if config.MaxTokens != 0 {
		if config.MaxTokens < MinMaxTokens || config.MaxTokens > MaxMaxTokens {
			return easyllm.NewValidationError(
				"maxTokens",
				fmt.Sprintf("must be between %d and %d", MinMaxTokens, MaxMaxTokens),
				config.MaxTokens,
			)
		}
	}

	// Validate presence_penalty
	if config.PresencePenalty != 0 {
		if config.PresencePenalty < MinPresencePenalty || config.PresencePenalty > MaxPresencePenalty {
			return easyllm.NewValidationError(
				"presencePenalty",
				fmt.Sprintf("must be between %.1f and %.1f", MinPresencePenalty, MaxPresencePenalty),
				config.PresencePenalty,
			)
		}
	}

	// Validate frequency_penalty
	if config.FrequencyPenalty != 0 {
		if config.FrequencyPenalty < MinFrequencyPenalty || config.FrequencyPenalty > MaxFrequencyPenalty {
			return easyllm.NewValidationError(
				"frequencyPenalty",
				fmt.Sprintf("must be between %.1f and %.1f", MinFrequencyPenalty, MaxFrequencyPenalty),
				config.FrequencyPenalty,
			)
		}
	}

	// Validate reasoning effort
	if config.ReasoningEffort != "" {
		validEfforts := map[easyllm.ReasoningEffort]bool{
			easyllm.ReasoningEffortLow:    true,
			easyllm.ReasoningEffortMedium: true,
			easyllm.ReasoningEffortHigh:   true,
		}
		if !validEfforts[config.ReasoningEffort] {
			return easyllm.NewValidationError(
				"reasoningEffort",
				"must be one of: low, medium, high",
				string(config.ReasoningEffort),
			)
		}
	}

	// Validate response format
	if config.ResponseFormat != "" {
		validFormats := map[easyllm.ResponseFormat]bool{
			easyllm.ResponseFormatJson:       true,
			easyllm.ResponseFormatJsonSchema: true,
		}
		if !validFormats[config.ResponseFormat] {
			return easyllm.NewValidationError(
				"responseFormat",
				"must be one of: json, json_schema",
				string(config.ResponseFormat),
			)
		}

		// If json_schema is specified, JSONSchema must be provided
		if config.ResponseFormat == easyllm.ResponseFormatJsonSchema && config.JSONSchema == nil {
			return easyllm.NewValidationError(
				"jsonSchema",
				"must be provided when responseFormat is json_schema",
				nil,
			)
		}
	}

	return nil
}

// ValidateEmbeddingRequestWithDetails validates embedding request with detailed errors
func ValidateEmbeddingRequestWithDetails(req *easyllm.EmbeddingRequest) error {
	if req == nil {
		return easyllm.NewValidationError("request", "cannot be nil", nil)
	}

	if req.Model == "" {
		return easyllm.NewValidationError("model", "cannot be empty", "")
	}

	if len(req.Contents) == 0 {
		return easyllm.NewValidationError("contents", "must contain at least one item", nil)
	}

	// Validate each content item
	for i, content := range req.Contents {
		if strings.TrimSpace(content) == "" {
			return easyllm.NewValidationError(
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

// validateEmbeddingConfig validates embedding configuration
func validateEmbeddingConfig(config *easyllm.EmbeddingModelConfig) error {
	if config == nil {
		return nil
	}

	// Validate dimensions
	if config.Dimensions < 0 {
		return easyllm.NewValidationError(
			"dimensions",
			"must be a positive number or 0 for default",
			config.Dimensions,
		)
	}

	// Validate encoding format
	if config.EncodingFormat != "" {
		validFormats := map[easyllm.EmbeddingEncodingFormat]bool{
			easyllm.EmbeddingEncodingFormatFloat:  true,
			easyllm.EmbeddingEncodingFormatBase64: true,
		}
		if !validFormats[config.EncodingFormat] {
			return easyllm.NewValidationError(
				"encodingFormat",
				"must be one of: float, base64",
				string(config.EncodingFormat),
			)
		}
	}

	return nil
}

// ValidateImageRequestWithDetails validates image request with detailed errors
func ValidateImageRequestWithDetails(req *easyllm.ImageRequest) error {
	if req == nil {
		return easyllm.NewValidationError("request", "cannot be nil", nil)
	}

	if req.Model == "" {
		return easyllm.NewValidationError("model", "cannot be empty", "")
	}

	if strings.TrimSpace(req.Instructions) == "" {
		return easyllm.NewValidationError("instructions", "cannot be empty or whitespace only", req.Instructions)
	}

	// Validate config if provided
	if req.Config != nil {
		if err := validateImageConfig(req.Config); err != nil {
			return err
		}
	}

	return nil
}

// validateImageConfig validates image configuration
func validateImageConfig(config *easyllm.ImageModelConfig) error {
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
			return easyllm.NewValidationError(
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
			return easyllm.NewValidationError(
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
			return easyllm.NewValidationError(
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
		return easyllm.ErrAPIKeyEmpty
	}

	// Trim whitespace
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return easyllm.NewValidationError("apiKey", "cannot be whitespace only", apiKey)
	}

	// Check minimum length (most API keys are at least 20 characters)
	if len(apiKey) < 10 {
		return easyllm.NewValidationError("apiKey", "appears to be too short", len(apiKey))
	}

	return nil
}

// ValidateBaseURL validates base URL format
func ValidateBaseURL(baseURL string) error {
	if baseURL == "" {
		return easyllm.ErrBaseURLEmpty
	}

	// Parse URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return easyllm.NewValidationError("baseURL", "invalid URL format", baseURL)
	}

	// Check scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return easyllm.NewValidationError("baseURL", "must use http or https scheme", parsedURL.Scheme)
	}

	// Check host
	if parsedURL.Host == "" {
		return easyllm.NewValidationError("baseURL", "must have a valid host", baseURL)
	}

	return nil
}

// ValidateModelName validates that a model name is not empty and doesn't contain invalid characters
func ValidateModelName(modelName string) error {
	if modelName == "" {
		return easyllm.NewValidationError("model", "cannot be empty", "")
	}

	modelName = strings.TrimSpace(modelName)
	if modelName == "" {
		return easyllm.NewValidationError("model", "cannot be whitespace only", modelName)
	}

	// Model names should not contain certain characters
	invalidChars := []string{"\n", "\r", "\t", "\x00"}
	for _, char := range invalidChars {
		if strings.Contains(modelName, char) {
			return easyllm.NewValidationError("model", "contains invalid characters", modelName)
		}
	}

	return nil
}
