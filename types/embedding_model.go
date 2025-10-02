package types

import "context"

// EmbeddingModel defines the interface for embedding generation operations
type EmbeddingModel interface {
	BaseModel
	// GenerateEmbeddings generates embeddings from input text
	GenerateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
}

type EmbeddingRequest struct {
	Model    string                `json:"model"`
	Contents []string              `json:"contents"`
	Config   *EmbeddingModelConfig `json:"config,omitempty"`
}

type EmbeddingModelConfig struct {
	EncodingFormat EmbeddingEncodingFormat `json:"encoding_format,omitempty"`
	Dimensions     int                     `json:"dimensions,omitempty"`
	User           string                  `json:"user,omitempty"`
}

type Embedding struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
	Object    string    `json:"object"`
}

type EmbeddingResponse struct {
	Embeddings []Embedding `json:"embeddings"`
	Usage      *TokenUsage `json:"usage,omitempty"`
	Cost       *float64    `json:"cost,omitempty"`
}
