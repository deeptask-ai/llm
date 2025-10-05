// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package replicate

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/easyagent-dev/llm"
	"github.com/replicate/replicate-go"
)

// ReplicateImageModel implements ImageModel interface for Replicate
type ReplicateImageModel struct {
	apiKey string
	client *replicate.Client
}

var _ llm.ImageModel = (*ReplicateImageModel)(nil)

// NewReplicateImageModel creates a new Replicate image model
func NewReplicateImageModel(apiKey string, opts ...llm.ModelOption) (*ReplicateImageModel, error) {
	if apiKey == "" {
		return nil, llm.ErrAPIKeyEmpty
	}

	// Create Replicate client
	r8, err := replicate.NewClient(replicate.WithToken(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create replicate client: %w", err)
	}

	return &ReplicateImageModel{
		apiKey: apiKey,
		client: r8,
	}, nil
}

// Name returns the provider name
func (m *ReplicateImageModel) Name() string {
	return "replicate"
}

// SupportedModels returns a list of supported models
func (m *ReplicateImageModel) SupportedModels() []*llm.ModelInfo {
	ctx := context.Background()

	// List models from Replicate API
	page, err := m.client.ListModels(ctx)
	if err != nil {
		// Return empty list on error
		return []*llm.ModelInfo{}
	}

	// Convert replicate models to llm.ModelInfo
	var models []*llm.ModelInfo
	for _, model := range page.Results {
		// Create model ID from owner and name
		modelID := fmt.Sprintf("%s/%s", model.Owner, model.Name)

		// Get version ID if available
		if model.LatestVersion != nil {
			modelID = model.LatestVersion.ID
		}

		modelInfo := &llm.ModelInfo{
			ID:     modelID,
			Name:   model.Name,
			Input:  []llm.ModelMediaType{llm.ModelMediaTypeText},
			Output: []llm.ModelMediaType{llm.ModelMediaTypeImage},
		}

		models = append(models, modelInfo)
	}

	return models
}

// GenerateImage generates an image from a text prompt
func (m *ReplicateImageModel) GenerateImage(ctx context.Context, req *llm.ImageRequest) (*llm.ImageResponse, error) {
	if req.Instructions == "" {
		return nil, llm.ErrEmptyInstructions
	}

	// Build input parameters
	input := replicate.PredictionInput{
		"prompt": req.Instructions,
	}

	// Apply config if provided
	if req.Config != nil {
		if req.Config.Size != "" {
			input["size"] = req.Config.Size
		}
		if req.Config.Quality != "" {
			input["quality"] = req.Config.Quality
		}
		if req.Config.Style != "" {
			input["style"] = req.Config.Style
		}
	}

	// Get model from request
	model := req.Model
	if model == "" {
		return nil, fmt.Errorf("model must be specified in request")
	}

	// Create prediction
	prediction, err := m.client.CreatePrediction(ctx, model, input, nil, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create prediction: %w", err)
	}

	// Wait for completion
	err = m.client.Wait(ctx, prediction)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for prediction: %w", err)
	}

	// Check for errors in the prediction
	if prediction.Error != nil {
		return nil, fmt.Errorf("prediction failed: %v", prediction.Error)
	}

	// Check if we have output
	if prediction.Output == nil {
		return nil, llm.ErrEmptyContent
	}

	// Extract URL from output
	var url string
	switch output := prediction.Output.(type) {
	case string:
		url = output
	case []interface{}:
		if len(output) == 0 {
			return nil, fmt.Errorf("empty output array")
		}
		var ok bool
		url, ok = output[0].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected output format: expected string in array")
		}
	default:
		return nil, fmt.Errorf("unexpected output format: %T", output)
	}

	// Download the image
	imageData, err := downloadURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	// Create usage information
	usage := &llm.TokenUsage{
		TotalImages:   1,
		TotalRequests: 1,
	}

	// Note: Cost calculation would require model-specific pricing
	// This can be implemented based on the model pricing in SupportedModels

	return &llm.ImageResponse{
		Output: imageData,
		Usage:  usage,
		Cost:   nil,
	}, nil
}

// downloadURL downloads content from a URL and returns it as bytes
func downloadURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}
