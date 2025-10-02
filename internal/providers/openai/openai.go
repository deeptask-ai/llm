package openai

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/easymvp/easyllm/types/completion"
	"github.com/easymvp/easyllm/types/conversation"
	"github.com/easymvp/easyllm/types/embedding"
	"github.com/easymvp/easyllm/types/image"
	"github.com/openai/openai-go/v3/shared"
	"strconv"
	"sync"

	"github.com/easymvp/easyllm/internal/common"
	"github.com/easymvp/easyllm/types"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"

	"github.com/openai/openai-go/v3/option"
)

//go:embed openai.json
var openaiModels []byte

// OpenAIBaseModel provides base functionality for OpenAI models
type OpenAIBaseModel struct {
	client     openai.Client
	apiKey     string
	modelCache map[string]*types.ModelInfo
	cacheMutex sync.RWMutex
}

func NewOpenAIBaseModel(apiKey string, opts ...option.RequestOption) (*OpenAIBaseModel, error) {
	if apiKey == "" {
		return nil, types.ErrAPIKeyEmpty
	}

	// Prepend API key option to any additional options
	allOpts := append([]option.RequestOption{option.WithAPIKey(apiKey)}, opts...)
	client := openai.NewClient(allOpts...)

	return &OpenAIBaseModel{
		client:     client,
		apiKey:     apiKey,
		modelCache: make(map[string]*types.ModelInfo),
		cacheMutex: sync.RWMutex{},
	}, nil
}

func (b *OpenAIBaseModel) Name() string {
	return "openai"
}

func (b *OpenAIBaseModel) SupportedModels() []*types.ModelInfo {
	var models []*types.ModelInfo
	if err := json.Unmarshal(openaiModels, &models); err != nil {
		return nil
	}
	return models
}

// getModelInfo returns the ModelInfo for a given model with caching for better performance
func (b *OpenAIBaseModel) getModelInfo(modelID string) *types.ModelInfo {
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

// ClearModelCache clears the model info cache
func (b *OpenAIBaseModel) ClearModelCache() {
	b.cacheMutex.Lock()
	b.modelCache = make(map[string]*types.ModelInfo)
	b.cacheMutex.Unlock()
}

// OpenAICompletionModel implements CompletionModel interface
type OpenAICompletionModel struct {
	*OpenAIBaseModel
}

func NewOpenAICompletionModel(apiKey string, opts ...option.RequestOption) (*OpenAICompletionModel, error) {
	base, err := NewOpenAIBaseModel(apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &OpenAICompletionModel{OpenAIBaseModel: base}, nil
}

func (p *OpenAICompletionModel) StreamComplete(ctx context.Context, req *completion.CompletionRequest, tools []types.ModelTool) (completion.StreamCompletionResponse, error) {
	// Parse options
	opts := completion.ApplyCompletionOptions(req.Options)

	params, err := ToChatCompletionParams(req.Model, req.Instructions, req.Messages, opts, tools)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion params: %w", err)
	}

	stream := p.client.Chat.Completions.NewStreaming(ctx, params)
	chunkChan := make(chan types.StreamChunk, 10) // Increased buffer to reduce blocking

	go func() {
		defer close(chunkChan)

		// Use an accumulator to track the full content
		acc := openai.ChatCompletionAccumulator{}

		for stream.Next() {
			// Check for context cancellation
			select {
			case <-ctx.Done():
				// Context was canceled, send error and return
				select {
				case chunkChan <- types.StreamTextChunk{
					Text: fmt.Sprintf("Stream canceled: %v", ctx.Err()),
				}:
				default:
					// Channel full or closed, just return
				}
				return
			default:
				// Continue processing
			}

			chunk := stream.Current()
			acc.AddChunk(chunk)

			if len(chunk.Choices) > 0 {
				if chunk.Choices[0].Delta.Content != "" {
					text := chunk.Choices[0].Delta.Content
					select {
					case chunkChan <- types.StreamTextChunk{
						Text: text,
					}:
					case <-ctx.Done():
						// Context canceled while sending
						return
					}
				} else if f, ok := chunk.Choices[0].Delta.JSON.ExtraFields["reasoning_content"]; ok {
					reasoning := f.Raw()
					reasoning = reasoning[1 : len(reasoning)-1]
					select {
					case chunkChan <- types.StreamReasoningChunk{
						Reasoning: reasoning,
					}:
					case <-ctx.Done():
						// Context canceled while sending
						return
					}
				}
			}
		}

		// Check for errors from the stream
		if err := stream.Err(); err != nil {
			// Check if error is due to context cancellation
			if ctx.Err() != nil {
				return
			}

			select {
			case chunkChan <- types.StreamTextChunk{
				Text: fmt.Sprintf("Error from OpenAI API: %v", err),
			}:
			case <-ctx.Done():
				return
			}
			return
		}

		// Check if usage information should be included
		if opts.WithUsage != nil && *opts.WithUsage {
			// Create usage information
			usage := &types.TokenUsage{
				TotalInputTokens:      acc.ChatCompletion.Usage.PromptTokens,
				TotalOutputTokens:     acc.ChatCompletion.Usage.CompletionTokens,
				TotalReasoningTokens:  acc.ChatCompletion.Usage.CompletionTokensDetails.ReasoningTokens,
				TotalImages:           0,
				TotalWebSearches:      0,
				TotalRequests:         1,
				TotalCacheReadTokens:  acc.ChatCompletion.Usage.PromptTokensDetails.CachedTokens,
				TotalCacheWriteTokens: 0,
			}

			// Calculate cost if requested
			var cost *float64
			if opts.WithCost != nil && *opts.WithCost {
				modelInfo := p.getModelInfo(req.Model)
				cost = common.CalculateCost(modelInfo, usage)
			}

			// Send usage information at the end
			select {
			case chunkChan <- types.StreamUsageChunk{
				Usage: usage,
				Cost:  cost,
			}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return chunkChan, nil
}

func (p *OpenAICompletionModel) Complete(ctx context.Context, req *completion.CompletionRequest, tools []types.ModelTool) (*completion.CompletionResponse, error) {
	// Parse options
	opts := completion.ApplyCompletionOptions(req.Options)

	params, err := ToChatCompletionParams(req.Model, req.Instructions, req.Messages, opts, tools)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion params: %w", err)
	}
	resp, err := p.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to complete chat: %w", err)
	}

	// Check if we have any choices in the response
	if len(resp.Choices) == 0 {
		return nil, types.ErrEmptyContent
	}

	var usage *types.TokenUsage
	var cost *float64

	// Include usage information if requested
	if opts.WithUsage != nil && *opts.WithUsage {
		usage = &types.TokenUsage{
			TotalInputTokens:      resp.Usage.PromptTokens,
			TotalOutputTokens:     resp.Usage.CompletionTokens,
			TotalReasoningTokens:  resp.Usage.CompletionTokensDetails.ReasoningTokens,
			TotalImages:           0,
			TotalWebSearches:      0,
			TotalRequests:         1,
			TotalCacheReadTokens:  resp.Usage.PromptTokensDetails.CachedTokens,
			TotalCacheWriteTokens: 0,
		}

		// Calculate cost if requested
		if opts.WithCost != nil && *opts.WithCost {
			modelInfo := p.getModelInfo(req.Model)
			cost = common.CalculateCost(modelInfo, usage)
		}
	}

	output := resp.Choices[0].Message.Content
	return &completion.CompletionResponse{
		Output: output,
		Usage:  usage,
		Cost:   cost,
	}, nil
}

// OpenAIEmbeddingModel implements EmbeddingModel interface
type OpenAIEmbeddingModel struct {
	*OpenAIBaseModel
}

func NewOpenAIEmbeddingModel(apiKey string, opts ...option.RequestOption) (*OpenAIEmbeddingModel, error) {
	base, err := NewOpenAIBaseModel(apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &OpenAIEmbeddingModel{OpenAIBaseModel: base}, nil
}

func (p *OpenAIEmbeddingModel) GenerateEmbeddings(ctx context.Context, req *embedding.EmbeddingRequest) (*embedding.EmbeddingResponse, error) {
	// Validate the request
	if err := common.ValidateEmbeddingRequest(req); err != nil {
		return nil, err
	}

	// Set up parameters for embedding generation
	params := openai.EmbeddingNewParams{
		Model: req.Model,
	}

	// Handle input - use any interface{} for the union type
	var input any
	if len(req.Contents) == 1 {
		input = req.Contents[0]
	} else {
		input = req.Contents
	}
	params.Input = input.(openai.EmbeddingNewParamsInputUnion)

	// Apply config if provided
	if req.Config != nil {
		if req.Config.Dimensions > 0 {
			params.Dimensions = openai.Int(int64(req.Config.Dimensions))
		}
		if req.Config.EncodingFormat != "" {
			switch req.Config.EncodingFormat {
			case types.EmbeddingEncodingFormatFloat:
				params.EncodingFormat = openai.EmbeddingNewParamsEncodingFormatFloat
			case types.EmbeddingEncodingFormatBase64:
				params.EncodingFormat = openai.EmbeddingNewParamsEncodingFormatBase64
			}
		}
	}

	// Generate embeddings
	resp, err := p.client.Embeddings.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Convert OpenAI embeddings to our format
	embeddings := make([]embedding.Embedding, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = embedding.Embedding{
			Index:     int(data.Index),
			Embedding: data.Embedding,
			Object:    string(data.Object),
		}
	}

	// Create usage information
	usage := &types.TokenUsage{
		TotalInputTokens:  resp.Usage.PromptTokens,
		TotalOutputTokens: 0, // Embeddings don't have output tokens
		TotalRequests:     1,
	}

	// Calculate cost if model info is available
	var cost *float64
	modelInfo := p.getModelInfo(req.Model)
	if modelInfo != nil {
		cost = common.CalculateCost(modelInfo, usage)
	}

	return &embedding.EmbeddingResponse{
		Embeddings: embeddings,
		Usage:      usage,
		Cost:       cost,
	}, nil
}

// OpenAIImageModel implements ImageModel interface
type OpenAIImageModel struct {
	*OpenAIBaseModel
}

func NewOpenAIImageModel(apiKey string, opts ...option.RequestOption) (*OpenAIImageModel, error) {
	base, err := NewOpenAIBaseModel(apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &OpenAIImageModel{OpenAIBaseModel: base}, nil
}

func (p *OpenAIImageModel) GenerateImage(ctx context.Context, req *image.ImageRequest) (*image.ImageResponse, error) {
	if req.Instructions == "" {
		return nil, types.ErrEmptyInstructions
	}

	// Set up parameters for image generation using instructions as prompt
	params := openai.ImageGenerateParams{
		Prompt:         req.Instructions,
		Model:          openai.ImageModelDallE3,                         // Default to DALL-E 3
		ResponseFormat: openai.ImageGenerateParamsResponseFormatB64JSON, // Return base64 to get []byte
		N:              openai.Int(1),
	}

	// Apply config if provided
	if req.Config != nil {
		if req.Config.Size != "" {
			// Map size strings to OpenAI size constants
			switch req.Config.Size {
			case "1024x1024":
				params.Size = openai.ImageGenerateParamsSize1024x1024
			case "1792x1024":
				params.Size = openai.ImageGenerateParamsSize1792x1024
			case "1024x1792":
				params.Size = openai.ImageGenerateParamsSize1024x1792
			}
		}
		if req.Config.Quality != "" {
			switch req.Config.Quality {
			case "standard":
				params.Quality = openai.ImageGenerateParamsQualityStandard
			case "hd":
				params.Quality = openai.ImageGenerateParamsQualityHD
			}
		}
		if req.Config.Style != "" {
			switch req.Config.Style {
			case "vivid":
				params.Style = openai.ImageGenerateParamsStyleVivid
			case "natural":
				params.Style = openai.ImageGenerateParamsStyleNatural
			}
		}
	}

	// Generate the image
	resp, err := p.client.Images.Generate(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, types.ErrEmptyContent
	}

	// Decode base64 image data
	imageBytes, err := base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image data: %w", err)
	}

	// Create usage information
	usage := &types.TokenUsage{
		TotalImages:   1,
		TotalRequests: 1,
	}

	// Calculate cost if requested - image generation has fixed pricing
	var cost *float64
	if req.Config != nil {
		modelInfo := p.getModelInfo("dall-e-3") // Use DALL-E 3 pricing
		if modelInfo != nil {
			imagePrice, err := strconv.ParseFloat(modelInfo.Pricing.Image, 64)
			if err == nil {
				totalCost := imagePrice
				cost = &totalCost
			}
		}
	}

	return &image.ImageResponse{
		Output: imageBytes,
		Usage:  usage,
		Cost:   cost,
	}, nil
}

// OpenAIConversationModel implements ConversationModel interface
type OpenAIConversationModel struct {
	*OpenAIBaseModel
}

func NewOpenAIConversationModel(apiKey string, opts ...option.RequestOption) (*OpenAIConversationModel, error) {
	base, err := NewOpenAIBaseModel(apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &OpenAIConversationModel{OpenAIBaseModel: base}, nil
}

func (p *OpenAIConversationModel) StreamResponse(ctx context.Context, req *conversation.ConversationRequest, tools []types.ModelTool) (conversation.StreamConversationResponse, error) {
	// Parse options
	opts := conversation.ApplyResponseOptions(req.Options)

	params, err := ToResponseNewParams(req.Model, req.Input, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create response params: %w", err)
	}

	stream := p.client.Responses.NewStreaming(ctx, params)
	chunkChan := make(chan types.StreamChunk, 10)

	go func() {
		defer close(chunkChan)

		for stream.Next() {
			// Check for context cancellation
			select {
			case <-ctx.Done():
				select {
				case chunkChan <- types.StreamTextChunk{
					Text: fmt.Sprintf("Stream canceled: %v", ctx.Err()),
				}:
				default:
				}
				return
			default:
			}

			data := stream.Current()

			// Send delta content
			if data.Delta != "" {
				select {
				case chunkChan <- types.StreamTextChunk{
					Text: data.Delta,
				}:
				case <-ctx.Done():
					return
				}
			}

			// Check if we have complete text (end of stream)
			if data.JSON.Text.Valid() {
				break
			}
		}

		// Check for errors from the stream
		if err := stream.Err(); err != nil {
			if ctx.Err() != nil {
				return
			}

			select {
			case chunkChan <- types.StreamTextChunk{
				Text: fmt.Sprintf("Error from OpenAI API: %v", err),
			}:
			case <-ctx.Done():
				return
			}
			return
		}

		// Note: Usage information is not easily accessible from the streaming response
		// TODO: Implement usage tracking once the proper SDK field structure is clarified
	}()

	return chunkChan, nil
}

func (p *OpenAIConversationModel) Response(ctx context.Context, req *conversation.ConversationRequest, tools []types.ModelTool) (*conversation.ConversationResponse, error) {
	// Parse options
	opts := conversation.ApplyResponseOptions(req.Options)

	params, err := ToResponseNewParams(req.Model, req.Input, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create response params: %w", err)
	}

	resp, err := p.client.Responses.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	output := resp.OutputText()
	if output == "" {
		return nil, types.ErrEmptyContent
	}

	var usage *types.TokenUsage
	var cost *float64

	if opts.WithUsage != nil && *opts.WithUsage {
		usage = &types.TokenUsage{
			TotalInputTokens:  resp.Usage.InputTokens,
			TotalOutputTokens: resp.Usage.OutputTokens,
			TotalRequests:     1,
		}

		if opts.WithCost != nil && *opts.WithCost {
			modelInfo := p.getModelInfo(req.Model)
			cost = common.CalculateCost(modelInfo, usage)
		}
	}

	return &conversation.ConversationResponse{
		Output: output,
		Usage:  usage,
		Cost:   cost,
	}, nil
}

// OpenAIModel is a composite model that implements all OpenAI capabilities
type OpenAIModel struct {
	*OpenAICompletionModel
	*OpenAIConversationModel
	*OpenAIEmbeddingModel
	*OpenAIImageModel
}

func NewOpenAIModel(opts ...types.ModelOption) (*OpenAIModel, error) {
	config := types.ApplyOptions(opts)

	if config.APIKey == "" {
		return nil, types.ErrAPIKeyEmpty
	}

	// Build request options list
	requestOpts := config.Options

	// Set base URL if provided
	if config.BaseURL != "" {
		requestOpts = append([]option.RequestOption{option.WithBaseURL(config.BaseURL)}, requestOpts...)
	}

	// Create base model
	base, err := NewOpenAIBaseModel(config.APIKey, requestOpts...)
	if err != nil {
		return nil, err
	}

	return &OpenAIModel{
		OpenAICompletionModel:   &OpenAICompletionModel{OpenAIBaseModel: base},
		OpenAIConversationModel: &OpenAIConversationModel{OpenAIBaseModel: base},
		OpenAIEmbeddingModel:    &OpenAIEmbeddingModel{OpenAIBaseModel: base},
		OpenAIImageModel:        &OpenAIImageModel{OpenAIBaseModel: base},
	}, nil
}

// Override base methods to avoid ambiguity
func (m *OpenAIModel) Name() string {
	return m.OpenAICompletionModel.Name()
}

func (m *OpenAIModel) SupportedModels() []*types.ModelInfo {
	return m.OpenAICompletionModel.SupportedModels()
}

func (m *OpenAIModel) ClearModelCache() {
	m.OpenAICompletionModel.ClearModelCache()
}

// Helper functions

func ToChatCompletionParams(model string, instructions string, messages []*types.ModelMessage, opts *completion.CompletionOptions, tools []types.ModelTool) (openai.ChatCompletionNewParams, error) {
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages)+1)

	// Add system content if provided
	if instructions != "" {
		openaiMessages = append(openaiMessages, openai.SystemMessage(instructions))
	}

	// Add the rest of the messages
	for _, msg := range messages {
		openaiMsg, err := ToChatCompletionMessage(msg)
		if err != nil {
			return openai.ChatCompletionNewParams{}, fmt.Errorf("failed to convert message: %w", err)
		}
		openaiMessages = append(openaiMessages, openaiMsg)
	}

	params := openai.ChatCompletionNewParams{
		Messages: openaiMessages,
		Model:    model,
	}

	// Add tools if provided
	if len(tools) > 0 {

	}

	if opts != nil {
		if opts.Temperature != nil && *opts.Temperature != 0 {
			params.Temperature = openai.Float(*opts.Temperature)
		}
		if opts.TopP != nil && *opts.TopP != 0 {
			params.TopP = openai.Float(*opts.TopP)
		}
		if opts.MaxTokens != nil && *opts.MaxTokens != 0 {
			params.MaxTokens = openai.Int(int64(*opts.MaxTokens))
		}
		if opts.PresencePenalty != nil && *opts.PresencePenalty != 0 {
			params.PresencePenalty = openai.Float(*opts.PresencePenalty)
		}
		if opts.FrequencyPenalty != nil && *opts.FrequencyPenalty != 0 {
			params.FrequencyPenalty = openai.Float(*opts.FrequencyPenalty)
		}
		if opts.Seed != nil && *opts.Seed != 0 {
			params.Seed = openai.Int(*opts.Seed)
		}
		if opts.ReasoningEffort != nil {
			switch *opts.ReasoningEffort {
			case completion.ReasoningEffortLow:
				params.ReasoningEffort = openai.ReasoningEffortLow
			case completion.ReasoningEffortMedium:
				params.ReasoningEffort = openai.ReasoningEffortMedium
			case completion.ReasoningEffortHigh:
				params.ReasoningEffort = openai.ReasoningEffortHigh
			default:
				params.ReasoningEffort = openai.ReasoningEffortLow
			}
		}
		if len(opts.Stop) > 0 {
			params.Stop = openai.ChatCompletionNewParamsStopUnion{
				OfStringArray: opts.Stop,
			}
		}
		if opts.ResponseFormat != nil {
			if *opts.ResponseFormat == completion.ResponseFormatJson {
				params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
					OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
				}
			} else {
				params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
					OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
						JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
							Name:   "response_schema",
							Schema: opts.JSONSchema,
						},
					},
				}
			}
		}
	}

	return params, nil
}

func ToChatCompletionMessage(msg *types.ModelMessage) (openai.ChatCompletionMessageParamUnion, error) {
	if msg == nil {
		return openai.UserMessage(""), types.NewValidationError("message", "cannot be nil", nil)
	}

	switch msg.Role {
	case types.MessageRoleUser:
		return openai.UserMessage(msg.Content), nil

	case types.MessageRoleAssistant:
		if msg.ToolCall == nil {
			return openai.AssistantMessage(msg.Content), nil
		}
		jsonBytes, err := json.Marshal(msg.ToolCall)
		if err != nil {
			return openai.AssistantMessage(""), fmt.Errorf("failed to marshal tool call: %w", err)
		}
		return openai.AssistantMessage("call tool: ```" + string(jsonBytes) + "```"), nil

	case types.MessageRoleTool:
		jsonBytes, err := json.Marshal(msg.ToolCall)
		if err != nil {
			return openai.UserMessage(""), fmt.Errorf("failed to marshal tool call results: %w", err)
		}
		return openai.UserMessage("call tool results: ```" + string(jsonBytes) + "```"), nil

	default:
		return openai.UserMessage(""), types.NewValidationError("role", "unknown role", string(msg.Role))
	}
}

func ToResponseNewParams(model string, input string, opts *conversation.ResponseOptions) (responses.ResponseNewParams, error) {
	params := responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(input)},
		Model: model,
	}

	if opts != nil {
		if opts.Temperature != nil && *opts.Temperature != 0 {
			params.Temperature = openai.Float(*opts.Temperature)
		}
		if opts.TopP != nil && *opts.TopP != 0 {
			params.TopP = openai.Float(*opts.TopP)
		}
		if opts.MaxOutputTokens != nil && *opts.MaxOutputTokens != 0 {
			params.MaxOutputTokens = openai.Int(int64(*opts.MaxOutputTokens))
		}

		if opts.ReasoningSummary != nil {
			params.Reasoning.Summary = shared.ReasoningSummary(*opts.ReasoningSummary)
		}
		if opts.ReasoningEffort != nil {
			params.Reasoning.Effort = shared.ReasoningEffort(*opts.ReasoningEffort)
		}
		if opts.ParallelToolCalls != nil {
			params.ParallelToolCalls = openai.Bool(*opts.ParallelToolCalls)
		}
		if opts.Store != nil {
			params.Store = openai.Bool(*opts.Store)
		}
		if opts.TopLogprobs != nil && *opts.TopLogprobs != 0 {
			params.TopLogprobs = openai.Int(int64(*opts.TopLogprobs))
		}
	}

	return params, nil
}
