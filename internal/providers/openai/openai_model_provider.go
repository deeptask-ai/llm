// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package openai

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/internal/common"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
	"strconv"

	"github.com/openai/openai-go/v3/option"
)

//go:embed openai.json
var openaiModels []byte

// OpenAIModelProvider provides base functionality for OpenAI models
type OpenAIModelProvider struct {
	*llm.DefaultModelProvider
	client openai.Client
}

var _ llm.ModelProvider = (*OpenAIModelProvider)(nil)

func NewOpenAIModelProvider(opts ...llm.ModelOption) (*OpenAIModelProvider, error) {
	var models []*llm.ModelInfo
	if err := json.Unmarshal(openaiModels, &models); err != nil {
		return nil, errors.New("failed to read model info")
	}
	config := llm.ApplyOptions(opts)
	requestOpts := []option.RequestOption{}
	requestOpts = append(requestOpts, option.WithAPIKey(config.APIKey))

	return NewBaseOpenAIModelProvider("openai", models, requestOpts)
}

func NewBaseOpenAIModelProvider(name string, models []*llm.ModelInfo, reqOpts []option.RequestOption) (*OpenAIModelProvider, error) {
	client := openai.NewClient(reqOpts...)

	provider := llm.NewDefaultModelProvider(name, models)

	return &OpenAIModelProvider{
		DefaultModelProvider: provider,
		client:               client,
	}, nil
}

func (p *OpenAIModelProvider) NewCompletionModel(model string, opts ...llm.CompletionOption) (llm.CompletionModel, error) {
	info := p.GetModelInfo(model)
	if info == nil {
		return nil, errors.New("model not found")
	}
	return NewOpenAICompletionModel(model, info, p.client, opts...)
}

func (p *OpenAIModelProvider) NewEmbeddingModel(model string) (llm.EmbeddingModel, error) {
	info := p.GetModelInfo(model)
	if info == nil {
		return nil, errors.New("model not found")
	}
	return NewOpenAIEmbeddingModel(model, info, p.client)
}

func (p *OpenAIModelProvider) NewImageModel(model string) (llm.ImageModel, error) {
	info := p.GetModelInfo(model)
	if info == nil {
		return nil, errors.New("model not found")
	}
	return NewOpenAIImageModel(model, info, p.client)
}

func (p *OpenAIModelProvider) NewConversationModel(model string, opts ...llm.ResponseOption) (llm.ConversationModel, error) {
	info := p.GetModelInfo(model)
	if info == nil {
		return nil, errors.New("model not found")
	}
	return NewOpenAIConversationModel(model, info, p.client, opts...)
}

// OpenAICompletionModel implements CompletionModel interface
type OpenAICompletionModel struct {
	name      string
	modelInfo *llm.ModelInfo
	client    openai.Client
	options   []llm.CompletionOption
}

func NewOpenAICompletionModel(name string, modelInfo *llm.ModelInfo, client openai.Client, opts ...llm.CompletionOption) (*OpenAICompletionModel, error) {
	return &OpenAICompletionModel{
		name:      name,
		modelInfo: modelInfo,
		client:    client,
		options:   opts,
	}, nil
}

func (p *OpenAICompletionModel) StreamComplete(ctx context.Context, req *llm.CompletionRequest) (llm.StreamCompletionResponse, error) {
	// Parse options
	opts := llm.ApplyCompletionOptions(p.options)

	params, err := ToChatCompletionParams(p.name, req.Instructions, req.Messages, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat llm params: %w", err)
	}

	stream := p.client.Chat.Completions.NewStreaming(ctx, params)
	chunkChan := make(chan llm.StreamChunk, 10) // Increased buffer to reduce blocking

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
				case chunkChan <- llm.StreamTextChunk{
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
					case chunkChan <- llm.StreamTextChunk{
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
					case chunkChan <- llm.StreamReasoningChunk{
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
			case chunkChan <- llm.StreamTextChunk{
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
			usage := &llm.TokenUsage{
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
				cost = common.CalculateCost(p.modelInfo, usage)
			}

			// Send usage information at the end
			select {
			case chunkChan <- llm.StreamUsageChunk{
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

func (p *OpenAICompletionModel) Complete(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	// Parse options
	opts := llm.ApplyCompletionOptions(p.options)

	params, err := ToChatCompletionParams(p.name, req.Instructions, req.Messages, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat llm params: %w", err)
	}
	resp, err := p.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to complete chat: %w", err)
	}

	// Check if we have any choices in the response
	if len(resp.Choices) == 0 {
		return nil, llm.ErrEmptyContent
	}

	var usage *llm.TokenUsage
	var cost *float64

	// Include usage information if requested
	if opts.WithUsage != nil && *opts.WithUsage {
		usage = &llm.TokenUsage{
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
			cost = common.CalculateCost(p.modelInfo, usage)
		}
	}

	output := resp.Choices[0].Message.Content
	return &llm.CompletionResponse{
		Output: output,
		Usage:  usage,
		Cost:   cost,
	}, nil
}

// OpenAIEmbeddingModel implements EmbeddingModel interface
type OpenAIEmbeddingModel struct {
	name      string
	modelInfo *llm.ModelInfo
	client    openai.Client
}

func NewOpenAIEmbeddingModel(name string, modelInfo *llm.ModelInfo, client openai.Client) (*OpenAIEmbeddingModel, error) {
	return &OpenAIEmbeddingModel{
		name:      name,
		modelInfo: modelInfo,
		client:    client,
	}, nil
}

func (p *OpenAIEmbeddingModel) GenerateEmbeddings(ctx context.Context, req *llm.EmbeddingRequest) (*llm.EmbeddingResponse, error) {
	// Validate the request
	if err := common.ValidateEmbeddingRequest(req); err != nil {
		return nil, err
	}

	// Set up parameters for llm generation
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
			case llm.EmbeddingEncodingFormatFloat:
				params.EncodingFormat = openai.EmbeddingNewParamsEncodingFormatFloat
			case llm.EmbeddingEncodingFormatBase64:
				params.EncodingFormat = openai.EmbeddingNewParamsEncodingFormatBase64
			}
		}
	}

	// Generate llms
	resp, err := p.client.Embeddings.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate llms: %w", err)
	}

	// Convert OpenAI llms to our format
	llms := make([]llm.Embedding, len(resp.Data))
	for i, data := range resp.Data {
		llms[i] = llm.Embedding{
			Index:     int(data.Index),
			Embedding: data.Embedding,
			Object:    string(data.Object),
		}
	}

	// Create usage information
	usage := &llm.TokenUsage{
		TotalInputTokens:  resp.Usage.PromptTokens,
		TotalOutputTokens: 0, // Embeddings don't have output tokens
		TotalRequests:     1,
	}

	// Calculate cost if model info is available
	var cost *float64
	if p.modelInfo != nil {
		cost = common.CalculateCost(p.modelInfo, usage)
	}

	return &llm.EmbeddingResponse{
		Embeddings: llms,
		Usage:      usage,
		Cost:       cost,
	}, nil
}

// OpenAIImageModel implements ImageModel interface
type OpenAIImageModel struct {
	name      string
	modelInfo *llm.ModelInfo
	client    openai.Client
}

func NewOpenAIImageModel(name string, modelInfo *llm.ModelInfo, client openai.Client) (*OpenAIImageModel, error) {
	return &OpenAIImageModel{
		name:      name,
		modelInfo: modelInfo,
		client:    client,
	}, nil
}

func (p *OpenAIImageModel) GenerateImage(ctx context.Context, req *llm.ImageRequest) (*llm.ImageResponse, error) {
	if req.Instructions == "" {
		return nil, llm.ErrEmptyInstructions
	}

	// Set up parameters for llm generation using instructions as prompt
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

	// Generate the llm
	resp, err := p.client.Images.Generate(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate llm: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, llm.ErrEmptyContent
	}

	// Decode base64 llm data
	llmBytes, err := base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
	if err != nil {
		return nil, fmt.Errorf("failed to decode llm data: %w", err)
	}

	// Create usage information
	usage := &llm.TokenUsage{
		TotalImages:   1,
		TotalRequests: 1,
	}

	// Calculate cost if requested - llm generation has fixed pricing
	var cost *float64
	if req.Config != nil {
		if p.modelInfo != nil {
			llmPrice, err := strconv.ParseFloat(p.modelInfo.Pricing.Image, 64)
			if err == nil {
				totalCost := llmPrice
				cost = &totalCost
			}
		}
	}

	return &llm.ImageResponse{
		Output: llmBytes,
		Usage:  usage,
		Cost:   cost,
	}, nil
}

// OpenAIConversationModel implements ConversationModel interface
type OpenAIConversationModel struct {
	name      string
	modelInfo *llm.ModelInfo
	client    openai.Client
	options   []llm.ResponseOption
}

func NewOpenAIConversationModel(name string, modelInfo *llm.ModelInfo, client openai.Client, opts ...llm.ResponseOption) (*OpenAIConversationModel, error) {
	return &OpenAIConversationModel{
		name:      name,
		modelInfo: modelInfo,
		client:    client,
		options:   opts,
	}, nil
}

func (p *OpenAIConversationModel) StreamResponse(ctx context.Context, req *llm.ConversationRequest) (llm.StreamConversationResponse, error) {
	// Parse options
	opts := llm.ApplyResponseOptions(p.options)

	params, err := ToResponseNewParams(p.name, req.Input, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create response params: %w", err)
	}

	stream := p.client.Responses.NewStreaming(ctx, params)
	chunkChan := make(chan llm.StreamChunk, 10)

	go func() {
		defer close(chunkChan)

		for stream.Next() {
			// Check for context cancellation
			select {
			case <-ctx.Done():
				select {
				case chunkChan <- llm.StreamTextChunk{
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
				case chunkChan <- llm.StreamTextChunk{
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
			case chunkChan <- llm.StreamTextChunk{
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

func (p *OpenAIConversationModel) Response(ctx context.Context, req *llm.ConversationRequest) (*llm.ConversationResponse, error) {
	// Parse options
	opts := llm.ApplyResponseOptions(p.options)

	params, err := ToResponseNewParams(p.name, req.Input, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create response params: %w", err)
	}

	resp, err := p.client.Responses.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	output := resp.OutputText()
	if output == "" {
		return nil, llm.ErrEmptyContent
	}

	var usage *llm.TokenUsage
	var cost *float64

	if opts.CompletionOptions.WithUsage != nil && *opts.CompletionOptions.WithUsage {
		usage = &llm.TokenUsage{
			TotalInputTokens:  resp.Usage.InputTokens,
			TotalOutputTokens: resp.Usage.OutputTokens,
			TotalRequests:     1,
		}

		if opts.CompletionOptions.WithCost != nil && *opts.CompletionOptions.WithCost {
			cost = common.CalculateCost(p.modelInfo, usage)
		}
	}

	return &llm.ConversationResponse{
		Output: output,
		Usage:  usage,
		Cost:   cost,
	}, nil
}

// Helper functions
func ToChatCompletionParams(model string, instructions string, messages []*llm.ModelMessage, opts *llm.CompletionOptions) (openai.ChatCompletionNewParams, error) {
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
			case llm.ReasoningEffortLow:
				params.ReasoningEffort = openai.ReasoningEffortLow
			case llm.ReasoningEffortMedium:
				params.ReasoningEffort = openai.ReasoningEffortMedium
			case llm.ReasoningEffortHigh:
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
			if *opts.ResponseFormat == llm.ResponseFormatJson {
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

func ToChatCompletionMessage(msg *llm.ModelMessage) (openai.ChatCompletionMessageParamUnion, error) {
	if msg == nil {
		return openai.UserMessage(""), llm.NewValidationError("message", "cannot be nil", nil)
	}

	switch msg.Role {
	case llm.RoleUser:
		return openai.UserMessage(msg.Content), nil

	case llm.RoleAssistant:
		if msg.ToolCall == nil {
			return openai.AssistantMessage(msg.Content), nil
		}
		jsonBytes, err := json.Marshal(msg.ToolCall)
		if err != nil {
			return openai.AssistantMessage(""), fmt.Errorf("failed to marshal tool call: %w", err)
		}
		return openai.AssistantMessage("call tool: ```" + string(jsonBytes) + "```"), nil

	case llm.RoleTool:
		jsonBytes, err := json.Marshal(msg.ToolCall)
		if err != nil {
			return openai.UserMessage(""), fmt.Errorf("failed to marshal tool call results: %w", err)
		}
		return openai.UserMessage("call tool results: ```" + string(jsonBytes) + "```"), nil

	default:
		return openai.UserMessage(""), llm.NewValidationError("role", "unknown role", string(msg.Role))
	}
}

func ToResponseNewParams(model string, input string, opts *llm.ResponseOptions) (responses.ResponseNewParams, error) {
	params := responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(input)},
		Model: model,
	}

	if opts != nil {
		completionOptions := opts.CompletionOptions
		if completionOptions.Temperature != nil && *completionOptions.Temperature != 0 {
			params.Temperature = openai.Float(*completionOptions.Temperature)
		}
		if completionOptions.TopP != nil && *completionOptions.TopP != 0 {
			params.TopP = openai.Float(*completionOptions.TopP)
		}
		if completionOptions.MaxOutputTokens != nil && *completionOptions.MaxOutputTokens != 0 {
			params.MaxOutputTokens = openai.Int(int64(*completionOptions.MaxOutputTokens))
		}

		if opts.ReasoningSummary != nil {
			params.Reasoning.Summary = shared.ReasoningSummary(*opts.ReasoningSummary)
		}
		if completionOptions.ReasoningEffort != nil {
			params.Reasoning.Effort = shared.ReasoningEffort(*completionOptions.ReasoningEffort)
		}
		if completionOptions.ParallelToolCalls != nil {
			params.ParallelToolCalls = openai.Bool(*completionOptions.ParallelToolCalls)
		}
		if opts.Store != nil {
			params.Store = openai.Bool(*opts.Store)
		}
		if completionOptions.TopLogprobs != nil && *completionOptions.TopLogprobs != 0 {
			params.TopLogprobs = openai.Int(int64(*completionOptions.TopLogprobs))
		}
	}

	return params, nil
}
