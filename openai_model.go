package llmclient

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"text/template"
)

type OpenAIModel struct {
	client openai.Client
	apiKey string
}

type OpenAIModelConfig struct {
	APIKey string
}

var _ Model = (*OpenAIModel)(nil)

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

func (p *OpenAIModel) StreamGenerateContent(ctx context.Context, req *ModelRequest) (StreamModelResponse, error) {
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

		// Send usage information at the end (no cost calculation for OpenAI)
		chunkChan <- StreamUsageChunk{
			Usage: usage,
			Cost:  nil,
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
	output := resp.Choices[0].Message.Content
	return &ModelResponse{
		Output: output,
		Usage:  usage,
		Cost:   nil,
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

func GetPrompts(prompt string, params map[string]interface{}) (string, error) {
	tmpl, err := template.New("prompt").Parse(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, params)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
