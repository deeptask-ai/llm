package replicate

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/replicate"
)

// NewReplicateImageModel creates a new Replicate image model
func NewReplicateImageModel(apiKey string, opts ...llm.ModelOption) (llm.ImageModel, error) {
	return replicate.NewReplicateImageModel(apiKey, opts...)
}
