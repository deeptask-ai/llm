package easyllm

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"text/template"

	"github.com/easymvp/easyllm/internal/conversion"
)

// Legacy errors - kept for backward compatibility
// New code should use errors from errors.go
var (
	ErrNoInstructions      = errors.New("no instructions provided")
	ErrNoCompletionChoices = errors.New("no completion choices returned")
	ErrNoImageData         = errors.New("no image data returned")
	ErrModelNotFound       = errors.New("model not found")
)

// Template cache for better performance
var (
	templateCache = make(map[string]*template.Template)
	templateMutex sync.RWMutex
)

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
// Delegates to internal/conversion package for consistency
func CalculateImageCost(modelInfo *ModelInfo, imageCount int) *float64 {
	if modelInfo == nil {
		return nil
	}
	return conversion.CalculateImageCost(modelInfo.Pricing.Image, imageCount)
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
func ValidateModelRequest(req *CompletionRequest) error {
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

// GetPrompts executes a template with caching for better performance
func GetPrompts(prompt string, params map[string]interface{}) (string, error) {
	// Try to get cached template first (read lock)
	templateMutex.RLock()
	tmpl, exists := templateCache[prompt]
	templateMutex.RUnlock()

	if !exists {
		// Parse and cache the template (write lock)
		templateMutex.Lock()
		// Double-check in case another goroutine added it
		if tmpl, exists = templateCache[prompt]; !exists {
			var err error
			tmpl, err = template.New("prompt").Parse(prompt)
			if err != nil {
				templateMutex.Unlock()
				return "", fmt.Errorf("failed to parse template: %w", err)
			}
			templateCache[prompt] = tmpl
		}
		templateMutex.Unlock()
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ClearTemplateCache clears the template cache to free memory
func ClearTemplateCache() {
	templateMutex.Lock()
	templateCache = make(map[string]*template.Template)
	templateMutex.Unlock()
}

// OptimizedCalculateCost is a more efficient version of cost calculation
func OptimizedCalculateCost(modelInfo *ModelInfo, usage *TokenUsage) *float64 {
	if modelInfo == nil || usage == nil {
		return nil
	}

	const tokensPerMillion = 1000000.0
	totalCost := 0.0

	// Parse prices once and cache them
	promptPrice, err := strconv.ParseFloat(modelInfo.Pricing.Prompt, 64)
	if err != nil {
		return nil
	}

	completionPrice, err := strconv.ParseFloat(modelInfo.Pricing.Completion, 64)
	if err != nil {
		return nil
	}

	// Calculate input token costs
	if cacheReadPrice, err := strconv.ParseFloat(modelInfo.Pricing.InputCacheRead, 64); err == nil && cacheReadPrice > 0.0 {
		totalInputTokens := usage.TotalInputTokens - usage.TotalCacheReadTokens
		totalCost += (float64(totalInputTokens) / tokensPerMillion) * promptPrice
		totalCost += (float64(usage.TotalCacheReadTokens) / tokensPerMillion) * cacheReadPrice
	} else {
		totalCost += (float64(usage.TotalInputTokens) / tokensPerMillion) * promptPrice
	}

	// Calculate reasoning token costs
	if internalReasoningPrice, err := strconv.ParseFloat(modelInfo.Pricing.InternalReasoning, 64); err == nil && internalReasoningPrice > 0.0 {
		totalCost += (float64(usage.TotalReasoningTokens) / tokensPerMillion) * internalReasoningPrice
	}

	// Calculate completion token costs
	totalCost += (float64(usage.TotalOutputTokens) / tokensPerMillion) * completionPrice

	return &totalCost
}
