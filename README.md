# EasyLLM

A unified Go client library for interacting with multiple Large Language Model (LLM) providers through a consistent interface. EasyLLM abstracts the complexities of different LLM APIs, providing a seamless experience for developers working with AI models.

[![Go Version](https://img.shields.io/badge/Go-1.24.4+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/easymvp/easyllm)](https://goreportcard.com/report/github.com/easymvp/easyllm)

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
go get github.com/easymvp/easyllm
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/easymvp/easyllm"
)

func main() {
    // Initialize OpenAI client
    client, err := easyllm.NewOpenAIModel(easyllm.OpenAIModelConfig{
        APIKey: "your-openai-api-key",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a request
    req := &easyllm.ModelRequest{
        Model: "gpt-4",
        Messages: []*easyllm.Message{
            {
                Role:    easyllm.MessageRoleUser,
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
    client, _ := easyllm.NewOpenAIModel(easyllm.OpenAIModelConfig{
        APIKey: "your-api-key",
    })

    req := &easyllm.ModelRequest{
        Model: "gpt-4",
        Messages: []*easyllm.Message{
            {
                Role:    easyllm.MessageRoleUser,
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
        case easyllm.StreamChunkTypeText:
            fmt.Print(chunk.(*easyllm.StreamTextChunk).Text)
        case easyllm.StreamChunkTypeUsage:
            usage := chunk.(*easyllm.StreamUsageChunk).Usage
            fmt.Printf("\nTokens used: %d\n", usage.TotalTokens)
        }
    }
}
```

### Multi-Provider Example

```go
func multiProviderExample() {
    var models []easyllm.Model

    // Add multiple providers
    openai, _ := easyllm.NewOpenAIModel(easyllm.OpenAIModelConfig{
        APIKey: "openai-key",
    })
    models = append(models, openai)

    claude, _ := easyllm.NewClaudeModel(easyllm.ClaudeModelConfig{
        APIKey: "claude-key",
    })
    models = append(models, claude)

    gemini, _ := easyllm.NewGeminiModel(easyllm.GeminiModelConfig{
        APIKey: "gemini-key",
    })
    models = append(models, gemini)

    req := &easyllm.ModelRequest{
        Model: "gpt-4", // or "claude-3-sonnet", "gemini-pro", etc.
        Messages: []*easyllm.Message{
            {
                Role:    easyllm.MessageRoleUser,
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
    client, _ := easyllm.NewOpenAIModel(easyllm.OpenAIModelConfig{
        APIKey: "your-api-key",
    })

    req := &easyllm.ModelRequest{
        Model: "gpt-4",
        Messages: []*easyllm.Message{
            {
                Role:    easyllm.MessageRoleUser,
                Content: "What's the weather like in Tokyo?",
            },
        },
        Tools: []easyllm.Tool{WeatherTool{}},
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
    client, _ := easyllm.NewOpenAIModel(easyllm.OpenAIModelConfig{
        APIKey: "your-api-key",
    })

    req := &easyllm.EmbeddingRequest{
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
    client, _ := easyllm.NewOpenAIModel(easyllm.OpenAIModelConfig{
        APIKey: "your-api-key",
    })

    req := &easyllm.ImageRequest{
        Model:  "dall-e-3",
        Prompt: "A futuristic city with flying cars and neon lights",
        Config: &easyllm.ImageModelConfig{
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
openai, _ := easyllm.NewOpenAIModel(easyllm.OpenAIModelConfig{
    APIKey:  "key",
    BaseURL: "https://api.openai.com/v1", // Optional
})

// Azure OpenAI
azure, _ := easyllm.NewAzureOpenAIModel(easyllm.AzureOpenAIModelConfig{
    APIKey:     "key",
    BaseURL:    "https://your-resource.openai.azure.com",
    APIVersion: "2024-02-15-preview",
})

// Claude with custom settings
claude, _ := easyllm.NewClaudeModel(easyllm.ClaudeModelConfig{
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
cost := easyllm.CalculateCost(modelInfo, usage)
```

## JSON Schema Generation

```go
type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email"`
}

// Generate schema for structured output
schema := easyllm.GenerateSchema[Person]()

req := &easyllm.ModelRequest{
    Model: "gpt-4",
    Messages: []*easyllm.Message{
        {
            Role:    easyllm.MessageRoleUser,
            Content: "Extract person information from: John Doe, 30 years old, john@example.com",
        },
    },
    Config: &easyllm.ModelConfig{
        ResponseFormat: easyllm.ResponseFormatJSON,
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
git clone https://github.com/easymvp/easyllm.git
cd easyllm
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

- 📖 [Documentation](https://github.com/easymvp/easyllm/wiki)
- 🐛 [Issue Tracker](https://github.com/easymvp/easyllm/issues)
- 💬 [Discussions](https://github.com/easymvp/easyllm/discussions)

## Roadmap

- [ ] Additional provider support (Cohere, Hugging Face, etc.)
- [ ] Advanced retry mechanisms with exponential backoff
- [ ] Built-in prompt templates and management
- [ ] Metrics and monitoring integration
- [ ] Async/batch processing capabilities
- [ ] Model performance benchmarking tools

---

Made with ❤️ by the [EasyMVP](https://github.com/easymvp) team
