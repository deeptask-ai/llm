package llm

import "sync"

// ModelProvider defines the base interface that all model providers must implement
type ModelProvider interface {
	// Name returns the provider name (e.g., "openai", "claude", "gemini")
	Name() string

	// SupportedModels returns a list of all models supported by this provider
	SupportedModels() []*ModelInfo

	// NewCompletionModel creates a new model instance
	NewCompletionModel(model string, opts ...CompletionOption) (CompletionModel, error)

	NewEmbeddingModelModel(model string) (EmbeddingModel, error)

	NewImageModelModel(model string) (ImageModel, error)

	NewConversationModel(model string, opts ...ResponseOption) (ConversationModel, error)
}

type DefaultModelProvider struct {
	name       string
	modelCache map[string]*ModelInfo
	cacheMutex sync.RWMutex
}

var _ ModelProvider = (*DefaultModelProvider)(nil)

func NewDefaultModelProvider(name string, models []*ModelInfo) *DefaultModelProvider {
	modelCache := make(map[string]*ModelInfo)
	for _, model := range models {
		modelCache[model.Name] = model
	}
	return &DefaultModelProvider{
		name:       name,
		modelCache: modelCache,
	}
}

func (p *DefaultModelProvider) Name() string {
	return p.name
}

func (p *DefaultModelProvider) SupportedModels() []*ModelInfo {
	models := make([]*ModelInfo, 0, len(p.modelCache))
	for _, model := range p.modelCache {
		models = append(models, model)
	}
	return models
}

func (b *DefaultModelProvider) GetModelInfo(modelID string) *ModelInfo {
	// Try to get from cache first (read lock)
	b.cacheMutex.RLock()
	if modelInfo, exists := b.modelCache[modelID]; exists {
		b.cacheMutex.RUnlock()
		return modelInfo
	}
	b.cacheMutex.RUnlock()

	// Not in cache, search through supported models
	models := b.SupportedModels()
	for _, model := range models {
		if model.ID == modelID {
			// Cache the result (write lock)
			b.cacheMutex.Lock()
			b.modelCache[modelID] = model
			b.cacheMutex.Unlock()
			return model
		}
	}
	return nil
}

func (p *DefaultModelProvider) NewCompletionModel(model string, opts ...CompletionOption) (CompletionModel, error) {
	return nil, ErrInvalidModel
}

func (p *DefaultModelProvider) NewEmbeddingModelModel(model string) (EmbeddingModel, error) {
	return nil, ErrInvalidModel
}

func (p *DefaultModelProvider) NewImageModelModel(model string) (ImageModel, error) {
	return nil, ErrInvalidModel
}

func (p *DefaultModelProvider) NewConversationModel(model string, opts ...ResponseOption) (ConversationModel, error) {
	return nil, ErrInvalidModel
}
