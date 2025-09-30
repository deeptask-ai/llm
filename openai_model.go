package llmclient

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"strconv"
)

type OpenAIModel struct {
	client openai.Client
	apiKey string
}

type OpenAIModelConfig struct {
	APIKey string
}

func NewOpenAIModel(config OpenAIModelConfig) (*OpenAIModel, error) {
	if config.APIKey == "" {
		return nil, errors.New("API key cannot be empty")
	}

	// Create the client with API key
	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
	)

	provider := &OpenAIModel{
		client: client,
		apiKey: config.APIKey,
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

func (p *OpenAIModel) GenerateStream(ctx context.Context, req *ModelRequest) (StreamModelResponse, error) {
	params := ToChatCompletionParams(req.Model, req.Instructions, req.Messages, req.Config)
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
	params := ToChatCompletionParams(req.Model, req.Instructions, req.Messages, req.Config)
	resp, err := p.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to complete chat: %w", err)
	}

	// Check if we have any choices in the response
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no completion choices returned")
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
	// For now, return a not implemented error
	// This will be properly implemented once the OpenAI SDK interface is confirmed
	return nil, fmt.Errorf("GenerateEmbeddings not yet implemented for OpenAI model")
}

func (p *OpenAIModel) GenerateImage(ctx context.Context, req *ImageRequest) (*ImageResponse, error) {
	if req.Instructions == "" {
		return nil, fmt.Errorf("no instructions provided for image generation")
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
		return nil, fmt.Errorf("no image data returned")
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

func ToChatCompletionParams(model string, instructions string, messages []*Message, config *ModelConfig) openai.ChatCompletionNewParams {
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

// CalculateCost calculates the cost based on token usage and model pricing information
// This function can be shared across all model implementations
func CalculateCost(modelInfo *ModelInfo, usage *TokenUsage) *float64 {
	if modelInfo == nil {
		return nil
	}

	totalCost := 0.0

	// Calculate input token costs
	cacheReadPrice, err := strconv.ParseFloat(modelInfo.Pricing.InputCacheRead, 64)
	if err != nil {
		cacheReadPrice = 0.0
	}
	promptPrice, err := strconv.ParseFloat(modelInfo.Pricing.Prompt, 64)
	if err != nil {
		return nil
	}

	if cacheReadPrice > 0.0 {
		totalInputTokens := usage.TotalInputTokens - usage.TotalCacheReadTokens
		totalCost += (float64(totalInputTokens) / 1000000.0) * promptPrice
		totalCost += (float64(usage.TotalCacheReadTokens) / 1000000.0) * cacheReadPrice
	} else {
		totalCost += (float64(usage.TotalInputTokens) / 1000000.0) * promptPrice
	}

	// Calculate internal reasoning token costs
	internalReasoningPrice, err := strconv.ParseFloat(modelInfo.Pricing.InternalReasoning, 64)
	if err != nil {
		internalReasoningPrice = 0.0
	}
	if internalReasoningPrice > 0.0 {
		totalCost += (float64(usage.TotalReasoningTokens) / 1000000.0) * internalReasoningPrice
	}

	// Calculate completion token costs
	completionPrice, err := strconv.ParseFloat(modelInfo.Pricing.Completion, 64)
	if err != nil {
		return nil
	}
	totalCost += (float64(usage.TotalOutputTokens) / 1000000.0) * completionPrice

	return &totalCost
}

// getModelInfo returns the ModelInfo for a given model from the supported models
func (p *OpenAIModel) getModelInfo(modelID string) *ModelInfo {
	models := p.SupportedModels()
	for _, model := range models {
		if model.ID == modelID {
			return model
		}
	}
	return nil
}
