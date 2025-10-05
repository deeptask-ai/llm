package providers

import (
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/providers/replicate"
)

// NewReplicateModelProvider creates a new Replicate image model
func NewReplicateModelProvider(opts ...llm.ModelOption) (llm.ModelProvider, error) {
	return replicate.NewReplicateModelProvider(opts...)
}
