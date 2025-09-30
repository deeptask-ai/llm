# LLMClient

A unified Go client library for interacting with multiple Large Language Model (LLM) providers through a consistent interface. LLMClient abstracts the complexities of different LLM APIs, providing a seamless experience for developers working with AI models.

[![Go Version](https://img.shields.io/badge/Go-1.24.4+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/easymvp/llmclient)](https://goreportcard.com/report/github.com/easymvp/llmclient)

## Features

### 🌐 **Multi-Provider Support**
- **OpenAI** - GPT models, embeddings, and DALL-E image generation
- **Claude** - Anthropic's Claude models with advanced reasoning
- **Gemini** - Google's Gemini AI models
- **DeepSeek** - DeepSeek's reasoning and coding models
- **Azure OpenAI** - Enterprise-grade OpenAI models via Azure
- **OpenRouter** - Access to multiple models through OpenRouter's API

### 🔄 **Unified Interface**
- Consistent API across all providers
- Seamless provider switching without code changes
- Standardized request/response formats

### 🚀 **Advanced Capabilities**
- **Streaming Support** - Real-time response streaming for better user experience
- **Function Calling** - Tool use and function calling capabilities
- **Embeddings** - Text embedding generation (where supported)
- **Image Generation** - AI image creation (where supported)
- **Cost Calculation** - Built-in pricing and usage cost tracking
- **JSON Schema** - Automatic schema generation for structured outputs

### 🛡️ **Production Ready**
- Input validation and error handling
- Comprehensive test coverage
- Template-based prompt management
- Cache optimization for performance
- Concurrent request support

## Installation

```bash
go get github.com/easymvp/llmclient
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/easymvp/llmclient"
)

func main() {
    // Initialize OpenAI client
    client, err := llmclient.NewOpenAIModel(llmclient.OpenAIModelConfig{
        APIKey: "your-openai-api-key",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a request
    req := &llmclient.ModelRequest{
        Model: "gpt-4",
        Messages: []*llmclient.Message{
            {
                Role:    llmclient.MessageRoleUser,
                Content: "What is the capital of France?",
            },
        },
    }

    // Generate response
    resp, err := client.GenerateContent(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", resp.Content)
    fmt.Printf("Cost: $%.6f\n", *resp.Cost)
}
```

### Streaming Response

```go
func streamExample() {
    client, _ := llmclient.NewOpenAIModel(llmclient.OpenAIModelConfig{
        APIKey: "your-api-key",
    })

    req := &llmclient.ModelRequest{
        Model: "gpt-4",
        Messages: []*llmclient.Message{
            {
                Role:    llmclient.MessageRoleUser,
                Content: "Write a short story about AI",
            },
        },
    }

    stream, err := client.GenerateContentStream(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    for chunk := range stream {
        switch chunk.Type() {
        case llmclient.StreamChunkTypeText:
            fmt.Print(chunk.(*llmclient.StreamTextChunk).Text)
        case llmclient.StreamChunkTypeUsage:
            usage := chunk.(*llmclient.StreamUsageChunk).Usage
            fmt.Printf("\nTokens used: %d\n", usage.TotalTokens)
        }
    }
}
```

### Multi-Provider Example

```go
func multiProviderExample() {
    var models []llmclient.Model

    // Add multiple providers
    openai, _ := llmclient.NewOpenAIModel(llmclient.OpenAIModelConfig{
        APIKey: "openai-key",
    })
    models = append(models, openai)

    claude, _ := llmclient.NewClaudeModel(llmclient.ClaudeModelConfig{
        APIKey: "claude-key",
    })
    models = append(models, claude)

    gemini, _ := llmclient.NewGeminiModel(llmclient.GeminiModelConfig{
        APIKey: "gemini-key",
    })
    models = append(models, gemini)

    req := &llmclient.ModelRequest{
        Model: "gpt-4", // or "claude-3-sonnet", "gemini-pro", etc.
        Messages: []*llmclient.Message{
            {
                Role:    llmclient.MessageRoleUser,
                Content: "Explain quantum computing",
            },
        },
    }

    // Use any provider with the same interface
    for _, model := range models {
        fmt.Printf("Provider: %s\n", model.Name())
        resp, _ := model.GenerateContent(context.Background(), req)
        fmt.Printf("Response: %s\n\n", resp.Content)
    }
}
```

### Function Calling

```go
type WeatherTool struct{}

func (w WeatherTool) Name() string { return "get_weather" }
func (w WeatherTool) Description() string { return "Get current weather for a location" }
func (w WeatherTool) Parameters() interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "location": map[string]interface{}{
                "type": "string",
                "description": "City name",
            },
        },
        "required": []string{"location"},
    }
}

func functionCallingExample() {
    client, _ := llmclient.NewOpenAIModel(llmclient.OpenAIModelConfig{
        APIKey: "your-api-key",
    })

    req := &llmclient.ModelRequest{
        Model: "gpt-4",
        Messages: []*llmclient.Message{
            {
                Role:    llmclient.MessageRoleUser,
                Content: "What's the weather like in Tokyo?",
            },
        },
        Tools: []llmclient.Tool{WeatherTool{}},
    }

    resp, _ := client.GenerateContent(context.Background(), req)
    
    // Handle tool calls in response
    for _, toolCall := range resp.ToolCalls {
        fmt.Printf("Tool: %s, Args: %s\n", toolCall.Name, toolCall.Arguments)
    }
}
```

### Embeddings

```go
func embeddingsExample() {
    client, _ := llmclient.NewOpenAIModel(llmclient.OpenAIModelConfig{
        APIKey: "your-api-key",
    })

    req := &llmclient.EmbeddingRequest{
        Model: "text-embedding-3-small",
        Input: []string{
            "The quick brown fox jumps over the lazy dog",
            "Machine learning is a subset of artificial intelligence",
        },
    }

    resp, err := client.GenerateEmbeddings(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    for i, embedding := range resp.Embeddings {
        fmt.Printf("Text %d embedding dimensions: %d\n", i+1, len(embedding.Values))
    }
}
```

### Image Generation

```go
func imageGenerationExample() {
    client, _ := llmclient.NewOpenAIModel(llmclient.OpenAIModelConfig{
        APIKey: "your-api-key",
    })

    req := &llmclient.ImageRequest{
        Model:  "dall-e-3",
        Prompt: "A futuristic city with flying cars and neon lights",
        Config: &llmclient.ImageModelConfig{
            Size:    "1024x1024",
            Quality: "hd",
            N:       1,
        },
    }

    resp, err := client.GenerateImage(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Generated %d images\n", len(resp.Images))
    for i, img := range resp.Images {
        fmt.Printf("Image %d URL: %s\n", i+1, img.URL)
    }
}
```

## Supported Models

### OpenAI
- **Chat Models**: GPT-4o, GPT-4, GPT-3.5-turbo, o1-preview, o1-mini
- **Embedding Models**: text-embedding-3-large, text-embedding-3-small, text-embedding-ada-002
- **Image Models**: DALL-E 3, DALL-E 2

### Claude (Anthropic)
- **Chat Models**: Claude 3.5 Sonnet, Claude 3 Opus, Claude 3 Sonnet, Claude 3 Haiku

### Gemini (Google)
- **Chat Models**: Gemini Pro, Gemini Pro Vision, Gemini 1.5 Pro, Gemini 1.5 Flash

### DeepSeek
- **Chat Models**: DeepSeek V3, DeepSeek Coder, DeepSeek Chat

### Azure OpenAI
- All OpenAI models available through Azure's enterprise platform

### OpenRouter
- Access to 200+ models from various providers through a single API

## Configuration

### Environment Variables
```bash
export OPENAI_API_KEY="your-openai-key"
export CLAUDE_API_KEY="your-claude-key"
export GEMINI_API_KEY="your-gemini-key"
export DEEPSEEK_API_KEY="your-deepseek-key"
```

### Provider-Specific Configuration

```go
// OpenAI with custom base URL
openai, _ := llmclient.NewOpenAIModel(llmclient.OpenAIModelConfig{
    APIKey:  "key",
    BaseURL: "https://api.openai.com/v1", // Optional
})

// Azure OpenAI
azure, _ := llmclient.NewAzureOpenAIModel(llmclient.AzureOpenAIModelConfig{
    APIKey:     "key",
    BaseURL:    "https://your-resource.openai.azure.com",
    APIVersion: "2024-02-15-preview",
})

// Claude with custom settings
claude, _ := llmclient.NewClaudeModel(llmclient.ClaudeModelConfig{
    APIKey:  "key",
    BaseURL: "https://api.anthropic.com", // Optional
})
```

## Error Handling

```go
resp, err := client.GenerateContent(ctx, req)
if err != nil {
    // Handle different error types
    switch {
    case strings.Contains(err.Error(), "rate limit"):
        // Handle rate limiting
        time.Sleep(time.Minute)
        return retry()
    case strings.Contains(err.Error(), "insufficient_quota"):
        // Handle quota exceeded
        return handleQuotaError()
    default:
        // Handle other errors
        log.Printf("API Error: %v", err)
    }
}
```

## Cost Tracking

```go
resp, _ := client.GenerateContent(ctx, req)

// Access cost information
if resp.Cost != nil {
    fmt.Printf("Request cost: $%.6f\n", *resp.Cost)
}

// Access detailed usage
usage := resp.Usage
fmt.Printf("Input tokens: %d\n", usage.InputTokens)
fmt.Printf("Output tokens: %d\n", usage.OutputTokens)
fmt.Printf("Total tokens: %d\n", usage.TotalTokens)

// Calculate cost manually
modelInfo := client.SupportedModels()[0] // Get model info
cost := llmclient.CalculateCost(modelInfo, usage)
```

## JSON Schema Generation

```go
type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email"`
}

// Generate schema for structured output
schema := llmclient.GenerateSchema[Person]()

req := &llmclient.ModelRequest{
    Model: "gpt-4",
    Messages: []*llmclient.Message{
        {
            Role:    llmclient.MessageRoleUser,
            Content: "Extract person information from: John Doe, 30 years old, john@example.com",
        },
    },
    Config: &llmclient.ModelConfig{
        ResponseFormat: llmclient.ResponseFormatJSON,
        JSONSchema:     schema,
    },
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

1. Clone the repository:
```bash
git clone https://github.com/easymvp/llmclient.git
cd llmclient
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test ./...
```

4. Create your feature branch:
```bash
git checkout -b feature/amazing-feature
```

5. Commit your changes:
```bash
git commit -m 'Add amazing feature'
```

6. Push to the branch:
```bash
git push origin feature/amazing-feature
```

7. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](https://github.com/easymvp/llmclient/wiki)
- 🐛 [Issue Tracker](https://github.com/easymvp/llmclient/issues)
- 💬 [Discussions](https://github.com/easymvp/llmclient/discussions)

## Roadmap

- [ ] Additional provider support (Cohere, Hugging Face, etc.)
- [ ] Advanced retry mechanisms with exponential backoff
- [ ] Built-in prompt templates and management
- [ ] Metrics and monitoring integration
- [ ] Async/batch processing capabilities
- [ ] Model performance benchmarking tools

---

Made with ❤️ by the [EasyMVP](https://github.com/easymvp) team
