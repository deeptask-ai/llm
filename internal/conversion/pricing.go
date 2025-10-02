// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package conversion

import (
	"strconv"
)

// CalculateCost calculates the total cost based on model pricing and token usage
// Returns nil if pricing information is not available
func CalculateCost(inputPrice, completionPrice, reasoningPrice string, inputTokens, outputTokens, reasoningTokens int64) *float64 {
	if inputPrice == "" || completionPrice == "" {
		return nil
	}

	inputPriceFloat, err := strconv.ParseFloat(inputPrice, 64)
	if err != nil {
		return nil
	}

	completionPriceFloat, err := strconv.ParseFloat(completionPrice, 64)
	if err != nil {
		return nil
	}

	// Calculate basic input and output costs
	cost := (float64(inputTokens) / 1_000_000.0 * inputPriceFloat) +
		(float64(outputTokens) / 1_000_000.0 * completionPriceFloat)

	// Add reasoning cost if applicable
	if reasoningTokens > 0 && reasoningPrice != "" {
		reasoningPriceFloat, err := strconv.ParseFloat(reasoningPrice, 64)
		if err == nil {
			cost += float64(reasoningTokens) / 1_000_000.0 * reasoningPriceFloat
		}
	}

	return &cost
}

// CalculateImageCost calculates the cost for image generation
// Returns nil if pricing information is not available
func CalculateImageCost(imagePrice string, imageCount int) *float64 {
	if imagePrice == "" || imageCount <= 0 {
		return nil
	}

	priceFloat, err := strconv.ParseFloat(imagePrice, 64)
	if err != nil {
		return nil
	}

	cost := float64(imageCount) * priceFloat
	return &cost
}

// CalculateCostWithCache calculates cost including cache read/write tokens
func CalculateCostWithCache(inputPrice, completionPrice, cacheReadPrice, cacheWritePrice string,
	inputTokens, outputTokens, cacheReadTokens, cacheWriteTokens int64) *float64 {

	cost := CalculateCost(inputPrice, completionPrice, "", inputTokens, outputTokens, 0)
	if cost == nil {
		return nil
	}

	// Add cache read cost
	if cacheReadTokens > 0 && cacheReadPrice != "" {
		if readPrice, err := strconv.ParseFloat(cacheReadPrice, 64); err == nil {
			*cost += float64(cacheReadTokens) / 1_000_000.0 * readPrice
		}
	}

	// Add cache write cost
	if cacheWriteTokens > 0 && cacheWritePrice != "" {
		if writePrice, err := strconv.ParseFloat(cacheWritePrice, 64); err == nil {
			*cost += float64(cacheWriteTokens) / 1_000_000.0 * writePrice
		}
	}

	return cost
}
