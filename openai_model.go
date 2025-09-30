package easyllm

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIModel struct {
	client     openai.Client
	apiKey     string
	modelCache map[string]*ModelInfo
	cacheMutex sync.RWMutex
}

type OpenAIModelConfig struct {
	APIKey string
}

func NewOpenAIModel(config OpenAIModelConfig) (*OpenAIModel, error) {
	if config.APIKey == "" {
		return nil, ErrAPIKeyEmpty
	}

	// Create the client with API key
	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
	)

	provider := &OpenAIModel{
		client:     client,
		apiKey:     config.APIKey,
		modelCache: make(map[string]*ModelInfo),
		cacheMutex: sync.RWMutex{},
	}

	return provider, nil
}

func (p *OpenAIModel) Name() string {
	return "openai"
}

//go:embed data/openai.json
var openaiModels []byte

func (p *OpenAIModel) SupportedModels() []*ModelInfo {
	var models []*ModelInfo
	if err := json.Unmarshal(openaiModels, &models); err != nil {
		return nil
	}
	return models
}

func (p *OpenAIModel) GenerateContentStream(ctx context.Context, req *ModelRequest) (StreamModelResponse, error) {
	params := ToChatCompletionParams(req.Model, req.Instructions, req.Messages, req.Config, req.Tools)
	stream := p.client.Chat.Completions.NewStreaming(ctx, params)
	chunkChan := make(chan StreamChunk, 1)
	go func() {
		defer close(chunkChan)

		// Use an accumulator to track the full content
		acc := openai.ChatCompletionAccumulator{}
		var fullOutput string

		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)

			if len(chunk.Choices) > 0 {
				if chunk.Choices[0].Delta.Content != "" {
					text := chunk.Choices[0].Delta.Content
					fullOutput += text
					chunkChan <- StreamTextChunk{
						Text: text,
					}
				}
			}
		}

		// Check for errors
		if err := stream.Err(); err != nil {
			// Send an error content as a text chunk
			chunkChan <- StreamTextChunk{
				Text: fmt.Sprintf("Error from OpenAI API: %v", err),
			}
			return
		}

		// Create usage information
		usage := &TokenUsage{
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
		if req.Cost {
			modelInfo := p.getModelInfo(req.Model)
			cost = CalculateCost(modelInfo, usage)
		}

		// Send usage information at the end
		chunkChan <- StreamUsageChunk{
			Usage: usage,
			Cost:  cost,
		}
	}()

	return chunkChan, nil
}

func (p *OpenAIModel) GenerateContent(ctx context.Context, req *ModelRequest) (*ModelResponse, error) {
	params := ToChatCompletionParams(req.Model, req.Instructions, req.Messages, req.Config, req.Tools)
	resp, err := p.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to complete chat: %w", err)
	}

	// Check if we have any choices in the response
	if len(resp.Choices) == 0 {
		return nil, ErrNoCompletionChoices
	}
	usage := &TokenUsage{
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
	var cost *float64
	if req.Cost {
		modelInfo := p.getModelInfo(req.Model)
		cost = CalculateCost(modelInfo, usage)
	}

	output := resp.Choices[0].Message.Content
	return &ModelResponse{
		Output: output,
		Usage:  usage,
		Cost:   cost,
	}, nil
}

func (p *OpenAIModel) GenerateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	// Validate the request
	if err := ValidateEmbeddingRequest(req); err != nil {
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
			params.Dimensions = openai.Int(req.Config.Dimensions)
		}
		if req.Config.EncodingFormat != "" {
			switch req.Config.EncodingFormat {
			case EmbeddingEncodingFormatFloat:
				params.EncodingFormat = openai.EmbeddingNewParamsEncodingFormatFloat
			case EmbeddingEncodingFormatBase64:
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
	embeddings := make([]Embedding, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = Embedding{
			Index:     int(data.Index),
			Embedding: data.Embedding,
			Object:    string(data.Object),
		}
	}

	// Create usage information
	usage := &TokenUsage{
		TotalInputTokens:  resp.Usage.PromptTokens,
		TotalOutputTokens: 0, // Embeddings don't have output tokens
		TotalRequests:     1,
	}

	// Calculate cost if model info is available
	var cost *float64
	modelInfo := p.getModelInfo(req.Model)
	if modelInfo != nil {
		cost = CalculateCost(modelInfo, usage)
	}

	return &EmbeddingResponse{
		Embeddings: embeddings,
		Usage:      usage,
		Cost:       cost,
	}, nil
}

func (p *OpenAIModel) GenerateImage(ctx context.Context, req *ImageRequest) (*ImageResponse, error) {
	if req.Instructions == "" {
		return nil, ErrNoInstructions
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
	image, err := p.client.Images.Generate(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	if len(image.Data) == 0 {
		return nil, ErrNoImageData
	}

	// Decode base64 image data
	imageBytes, err := base64.StdEncoding.DecodeString(image.Data[0].B64JSON)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image data: %w", err)
	}

	// Create usage information
	usage := &TokenUsage{
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

	return &ImageResponse{
		Output: imageBytes,
		Usage:  usage,
		Cost:   cost,
	}, nil
}

func ToChatCompletionParams(model string, instructions string, messages []*Message, config *ModelConfig, tools []Tool) openai.ChatCompletionNewParams {
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages)+1)

	// Add system content if provided
	if instructions != "" {
		openaiMessages = append(openaiMessages, openai.SystemMessage(instructions))
	}

	// Add the rest of the messages
	for _, msg := range messages {
		openaiMessages = append(openaiMessages, ToChatCompletionMessage(msg))
	}

	params := openai.ChatCompletionNewParams{
		Messages: openaiMessages,
		Model:    model,
	}

	// Add tools if provided
	if len(tools) > 0 {
		openaiTools := make([]openai.ChatCompletionToolParam, 0, len(tools))
		for _, tool := range tools {
			openaiTool := openai.ChatCompletionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        tool.Name(),
					Description: openai.String(tool.Description()),
				},
			}
			openaiTools = append(openaiTools, openaiTool)
		}
		params.Tools = openaiTools
	}

	if config != nil {
		if config.Temperature != 0 {
			params.Temperature = openai.Float(config.Temperature)
		}
		if config.TopP != 0 {
			params.TopP = openai.Float(config.TopP)
		}
		if config.MaxTokens != 0 {
			params.MaxTokens = openai.Int(int64(config.MaxTokens))
		}
		if config.PresencePenalty != 0 {
			params.PresencePenalty = openai.Float(config.PresencePenalty)
		}
		if config.FrequencyPenalty != 0 {
			params.FrequencyPenalty = openai.Float(config.FrequencyPenalty)
		}
		if config.Seed != 0 {
			params.Seed = openai.Int(config.Seed)
		}
		if config.ReasoningEffort != "" {
			switch config.ReasoningEffort {
			case "low":
				params.ReasoningEffort = openai.ReasoningEffortLow
			case "medium":
				params.ReasoningEffort = openai.ReasoningEffortMedium
			case "high":
				params.ReasoningEffort = openai.ReasoningEffortHigh
			default:
				params.ReasoningEffort = openai.ReasoningEffortLow
			}
		} else {
			params.ReasoningEffort = openai.ReasoningEffortLow
		}
		if len(config.Stop) > 0 {
			params.Stop = openai.ChatCompletionNewParamsStopUnion{
				OfStringArray: config.Stop,
			}
		}
		if config.ResponseFormat != "" {
			if config.ResponseFormat == ResponseFormatJson {
				params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
					OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
				}
			} else {
				params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
					OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
						JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
							Name:   "response_schema",
							Schema: config.JSONSchema,
						},
					},
				}
			}
		}
	}

	return params
}

func ToChatCompletionMessage(msg *Message) openai.ChatCompletionMessageParamUnion {
	if string(msg.Role) == string(openai.MessageRoleUser) {
		return openai.UserMessage(msg.Content)
	} else if string(msg.Role) == string(openai.MessageRoleAssistant) {
		if msg.ToolCall == nil {
			return openai.AssistantMessage(msg.Content)
		} else {
			jsonBytes, err := json.Marshal(msg.ToolCall)
			toolCallJSON := "{}"
			if err == nil {
				toolCallJSON = string(jsonBytes)
			}
			return openai.AssistantMessage("call tool: ```" + toolCallJSON + "```")
		}
	} else if string(msg.Role) == "tool" {
		jsonBytes, err := json.Marshal(msg.ToolCall)
		toolCallJSON := "{}"
		if err == nil {
			toolCallJSON = string(jsonBytes)
		}
		return openai.UserMessage("call tool results: ```" + toolCallJSON + "```")
	} else {
		panic("unknown role")
	}
}

// getModelInfo returns the ModelInfo for a given model with caching for better performance
func (p *OpenAIModel) getModelInfo(modelID string) *ModelInfo {
	// Try to get from cache first (read lock)
	p.cacheMutex.RLock()
	if modelInfo, exists := p.modelCache[modelID]; exists {
		p.cacheMutex.RUnlock()
		return modelInfo
	}
	p.cacheMutex.RUnlock()

	// Not in cache, search through supported models
	models := p.SupportedModels()
	for _, model := range models {
		if model.ID == modelID {
			// Cache the result (write lock)
			p.cacheMutex.Lock()
			p.modelCache[modelID] = model
			p.cacheMutex.Unlock()
			return model
		}
	}
	return nil
}

// ClearModelCache clears the model info cache
func (p *OpenAIModel) ClearModelCache() {
	p.cacheMutex.Lock()
	p.modelCache = make(map[string]*ModelInfo)
	p.cacheMutex.Unlock()
}
