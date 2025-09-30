package llmclient

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"text/template"
)

// Standard error types for better error handling
var (
	ErrAPIKeyEmpty         = errors.New("API key cannot be empty")
	ErrBaseURLEmpty        = errors.New("base URL cannot be empty")
	ErrAPIVersionEmpty     = errors.New("API version cannot be empty")
	ErrNoInstructions      = errors.New("no instructions provided")
	ErrNoCompletionChoices = errors.New("no completion choices returned")
	ErrNoImageData         = errors.New("no image data returned")
	ErrModelNotFound       = errors.New("model not found")
)

// Capability errors for unsupported features
func NewUnsupportedCapabilityError(provider, capability string) error {
	if capability == "image generation" {
		return fmt.Errorf("%s is not supported by %s models", capability, provider)
	}
	return fmt.Errorf("%s are not supported by %s models", capability, provider)
}

// CalculateCost calculates the cost based on token usage and model pricing information
// This function is shared across all model implementations
func CalculateCost(modelInfo *ModelInfo, usage *TokenUsage) *float64 {
	if modelInfo == nil {
		return nil
	}

	totalCost := 0.0

	// Calculate input token costs
	cacheReadPrice, err := strconv.ParseFloat(modelInfo.Pricing.InputCacheRead, 64)
	if err != nil {
		cacheReadPrice = 0.0
	}
	promptPrice, err := strconv.ParseFloat(modelInfo.Pricing.Prompt, 64)
	if err != nil {
		return nil
	}

	if cacheReadPrice > 0.0 {
		totalInputTokens := usage.TotalInputTokens - usage.TotalCacheReadTokens
		totalCost += (float64(totalInputTokens) / 1000000.0) * promptPrice
		totalCost += (float64(usage.TotalCacheReadTokens) / 1000000.0) * cacheReadPrice
	} else {
		totalCost += (float64(usage.TotalInputTokens) / 1000000.0) * promptPrice
	}

	// Calculate internal reasoning token costs
	internalReasoningPrice, err := strconv.ParseFloat(modelInfo.Pricing.InternalReasoning, 64)
	if err != nil {
		internalReasoningPrice = 0.0
	}
	if internalReasoningPrice > 0.0 {
		totalCost += (float64(usage.TotalReasoningTokens) / 1000000.0) * internalReasoningPrice
	}

	// Calculate completion token costs
	completionPrice, err := strconv.ParseFloat(modelInfo.Pricing.Completion, 64)
	if err != nil {
		return nil
	}
	totalCost += (float64(usage.TotalOutputTokens) / 1000000.0) * completionPrice

	return &totalCost
}

// CalculateImageCost calculates cost for image generation
func CalculateImageCost(modelInfo *ModelInfo, imageCount int) *float64 {
	if modelInfo == nil {
		return nil
	}

	imagePrice, err := strconv.ParseFloat(modelInfo.Pricing.Image, 64)
	if err != nil {
		return nil
	}

	totalCost := float64(imageCount) * imagePrice
	return &totalCost
}

// CreateTokenUsage creates a standardized TokenUsage struct
func CreateTokenUsage(inputTokens, outputTokens, reasoningTokens int64, images, requests int, cacheReadTokens, cacheWriteTokens int64) *TokenUsage {
	return &TokenUsage{
		TotalInputTokens:      inputTokens,
		TotalOutputTokens:     outputTokens,
		TotalReasoningTokens:  reasoningTokens,
		TotalImages:           images,
		TotalWebSearches:      0,
		TotalRequests:         requests,
		TotalCacheReadTokens:  cacheReadTokens,
		TotalCacheWriteTokens: cacheWriteTokens,
	}
}

// ValidateModelRequest validates common model request fields
func ValidateModelRequest(req *ModelRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Model == "" {
		return errors.New("model cannot be empty")
	}
	return nil
}

// ValidateEmbeddingRequest validates embedding request fields
func ValidateEmbeddingRequest(req *EmbeddingRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Model == "" {
		return errors.New("model cannot be empty")
	}
	if len(req.Contents) == 0 {
		return errors.New("contents cannot be empty")
	}
	return nil
}

// ValidateImageRequest validates image request fields
func ValidateImageRequest(req *ImageRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.Instructions == "" {
		return ErrNoInstructions
	}
	return nil
}

func GetPrompts(prompt string, params map[string]interface{}) (string, error) {
	tmpl, err := template.New("prompt").Parse(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, params)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
