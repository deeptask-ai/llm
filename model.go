package llmclient

import (
	"context"
)

type ModelInfo struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Pricing ModelPricing `json:"pricing"`
}

type ModelPricing struct {
	Prompt            string `json:"prompt"`
	Completion        string `json:"completion"`
	Request           string `json:"request"`
	Image             string `json:"image"`
	WebSearch         string `json:"webSearch"`
	InternalReasoning string `json:"internalReasoning"`
	InputCacheRead    string `json:"inputCacheRead"`
	InputCacheWrite   string `json:"inputCacheWrite"`
}

type Model interface {
	Name() string
	SupportedModels() []*ModelInfo
	StreamGenerateContent(ctx context.Context, req *ModelRequest) (StreamModelResponse, error)
	GenerateContent(ctx context.Context, req *ModelRequest) (*ModelResponse, error)
	GenerateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
}
