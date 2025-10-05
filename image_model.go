package llm

import (
	"context"
)

// ImageModel defines the interface for image generation operations
type ImageModel interface {
	// GenerateImage generates images from text prompts
	GenerateImage(ctx context.Context, req *ImageRequest) (*ImageResponse, error)
}

type ImageRequest struct {
	Model        string            `json:"model"`
	Instructions string            `json:"instructions"`
	Artifacts    []*ModelArtifact  `json:"artifacts"`
	Config       *ImageModelConfig `json:"config,omitempty"`
}

type ImageModelConfig struct {
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	Style          string `json:"style,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
}

type ImageResponse struct {
	Output []byte      `json:"output"`
	Usage  *TokenUsage `json:"usage,omitempty"`
	Cost   *float64    `json:"cost,omitempty"`
}
