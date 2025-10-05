// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"errors"
	"github.com/easyagent-dev/llm"
)

// CalculateCost calculates the cost based on token usage and model pricing information
// This function is shared across all model implementations
// With float64 pricing, no parsing or caching needed - direct calculation!
func CalculateCost(modelInfo *llm.ModelInfo, usage *llm.TokenUsage) *float64 {
	if modelInfo == nil || usage == nil {
		return nil
	}

	const tokensPerMillion = 1000000.0
	totalCost := 0.0

	// Calculate input token costs
	if modelInfo.Pricing.InputCacheRead > 0.0 {
		totalInputTokens := usage.TotalInputTokens - usage.TotalCacheReadTokens
		totalCost += (float64(totalInputTokens) / tokensPerMillion) * modelInfo.Pricing.Prompt
		totalCost += (float64(usage.TotalCacheReadTokens) / tokensPerMillion) * modelInfo.Pricing.InputCacheRead
	} else {
		totalCost += (float64(usage.TotalInputTokens) / tokensPerMillion) * modelInfo.Pricing.Prompt
	}

	// Calculate internal reasoning token costs
	if modelInfo.Pricing.InternalReasoning > 0.0 {
		totalCost += (float64(usage.TotalReasoningTokens) / tokensPerMillion) * modelInfo.Pricing.InternalReasoning
	}

	// Calculate completion token costs
	totalCost += (float64(usage.TotalOutputTokens) / tokensPerMillion) * modelInfo.Pricing.Completion

	return &totalCost
}

// ValidateEmbeddingRequest validates llm request fields
func ValidateEmbeddingRequest(req *llm.EmbeddingRequest) error {
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
