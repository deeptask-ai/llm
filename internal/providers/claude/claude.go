package claude

import (
	_ "embed"
	"encoding/json"
	"github.com/easymvp/easyllm/internal/providers/openai"
	"github.com/easymvp/easyllm/types"
	"github.com/openai/openai-go/option"
)

type ClaudeModel struct {
	*openai.OpenAICompletionModel
}

func NewClaudeModel(opts ...types.ModelOption) (*ClaudeModel, error) {
	config := types.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, types.ErrAPIKeyEmpty
	}

	// Build request options list with defaults
	requestOpts := []option.RequestOption{
		option.WithHeader("anthropic-version", "2023-06-01"),
	}

	// Set base URL (use default if not provided)
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1/"
	}
	requestOpts = append(requestOpts, option.WithBaseURL(baseURL))

	// Append any custom options
	requestOpts = append(requestOpts, config.Options...)

	// Create the completion model with Claude's API endpoint and required headers
	completionModel, err := openai.NewOpenAICompletionModel(config.APIKey, requestOpts...)
	if err != nil {
		return nil, err
	}

	return &ClaudeModel{
		OpenAICompletionModel: completionModel,
	}, nil
}

//go:embed claude.json
var claudeModels []byte

func (p *ClaudeModel) SupportedModels() []*types.ModelInfo {
	var models []*types.ModelInfo
	if err := json.Unmarshal(claudeModels, &models); err != nil {
		return nil
	}
	return models
}

func (p *ClaudeModel) Name() string {
	return "claude"
}
