package llm

// ModelProvider defines the base interface that all model providers must implement
type ModelProvider interface {
	// Name returns the provider name (e.g., "openai", "claude", "gemini")
	Name() string

	// SupportedModels returns a list of all models supported by this provider
	SupportedModels() []*ModelInfo

	// NewCompletionModel creates a new model instance
	NewCompletionModel(model string, opts ...CompletionOption) (CompletionModel, error)

	NewEmbeddingModel(model string) (EmbeddingModel, error)

	NewImageModel(model string) (ImageModel, error)

	NewConversationModel(model string, opts ...ResponseOption) (ConversationModel, error)
}

type DefaultModelProvider struct {
	name        string
	models      []*ModelInfo
	modelByID   map[string]*ModelInfo
	modelByName map[string]*ModelInfo
}

var _ ModelProvider = (*DefaultModelProvider)(nil)

func NewDefaultModelProvider(name string, models []*ModelInfo) *DefaultModelProvider {
	modelByID := make(map[string]*ModelInfo, len(models))
	modelByName := make(map[string]*ModelInfo, len(models))

	for _, model := range models {
		modelByID[model.ID] = model
		modelByName[model.Name] = model
	}

	return &DefaultModelProvider{
		name:        name,
		models:      models,
		modelByID:   modelByID,
		modelByName: modelByName,
	}
}

func (p *DefaultModelProvider) Name() string {
	return p.name
}

func (p *DefaultModelProvider) SupportedModels() []*ModelInfo {
	return p.models
}

func (p *DefaultModelProvider) GetModelInfo(modelID string) *ModelInfo {
	// O(1) lookup by ID
	if model, exists := p.modelByID[modelID]; exists {
		return model
	}
	// O(1) lookup by Name
	if model, exists := p.modelByName[modelID]; exists {
		return model
	}
	return nil
}

func (p *DefaultModelProvider) NewCompletionModel(model string, opts ...CompletionOption) (CompletionModel, error) {
	return nil, ErrInvalidModel
}

func (p *DefaultModelProvider) NewEmbeddingModel(model string) (EmbeddingModel, error) {
	return nil, ErrInvalidModel
}

func (p *DefaultModelProvider) NewImageModel(model string) (ImageModel, error) {
	return nil, ErrInvalidModel
}

func (p *DefaultModelProvider) NewConversationModel(model string, opts ...ResponseOption) (ConversationModel, error) {
	return nil, ErrInvalidModel
}
