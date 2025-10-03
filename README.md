# EasyLLM

A unified Go client library for interacting with multiple Large Language Model (LLM) providers through a consistent interface. EasyLLM abstracts the complexities of different LLM APIs, providing a seamless experience for developers working with AI models.

[![Go Version](https://img.shields.io/badge/Go-1.24.4+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/easymvp-ai/llm)](https://goreportcard.com/report/github.com/easymvp-ai/llm)

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
go get github.com/easymvp-ai/llm
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/easymvp-ai/llm"
    "github.com/easymvp-ai/llm/types"
    "github.com/easymvp-ai/llm/types/completion"
)

func main() {
    // Initialize OpenAI client
    model, err := llm.NewOpenAIModel(
        types.WithAPIKey("your-openai-api-key"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create a request
    req := &completion.CompletionRequest{
        Model:        "gpt-4o-mini",
        Instructions: "You are a helpful assistant.",
        Messages: []*types.ModelMessage{
            {
                Role:    types.MessageRoleUser,
                Content: "What is the capital of France?",
            },
        },
        Options: []completion.CompletionOption{
            completion.WithCost(true),
        },
    }

    // Generate response
    resp, err := model.Complete(context.Background(), req, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Response:", resp.Output)
    if resp.Cost != nil {
        fmt.Printf("Cost: $%.6f\n", *resp.Cost)
    }
}
```

### Streaming Response

```go
func streamExample() {
    model, _ := llm.NewOpenAIModel(
        types.WithAPIKey("your-api-key"),
    )

    req := &completion.CompletionRequest{
        Model:        "gpt-4o-mini",
        Instructions: "You are a helpful assistant.",
        Messages: []*types.ModelMessage{
            {
                Role:    types.MessageRoleUser,
                Content: "Write a short story about AI",
            },
        },
    }

    stream, err := model.StreamComplete(context.Background(), req, nil)
    if err != nil {
        log.Fatal(err)
    }

    for chunk := range stream {
        switch c := chunk.(type) {
        case types.StreamTextChunk:
            fmt.Print(c.Text)
        case types.StreamUsageChunk:
            fmt.Printf("\nTokens used: %d\n", c.Usage.TotalInputTokens+c.Usage.TotalOutputTokens)
        }
    }
}
```

### Multi-Provider Example

```go
func multiProviderExample() {
    // Initialize multiple providers
    openai, _ := llm.NewOpenAIModel(
        types.WithAPIKey("openai-key"),
    )

    deepseek, _ := llm.NewDeepSeekModel(
        types.WithAPIKey("deepseek-key"),
    )

    req := &completion.CompletionRequest{
        Model:        "gpt-4o-mini", // or "deepseek-chat", etc.
        Instructions: "You are a helpful assistant.",
        Messages: []*types.ModelMessage{
            {
                Role:    types.MessageRoleUser,
                Content: "Explain quantum computing",
            },
        },
    }

    // Use OpenAI
    resp1, _ := openai.Complete(context.Background(), req, nil)
    fmt.Printf("OpenAI Response: %s\n\n", resp1.Output)

    // Use DeepSeek with same request structure
    req.Model = "deepseek-chat"
    resp2, _ := deepseek.Complete(context.Background(), req, nil)
    fmt.Printf("DeepSeek Response: %s\n\n", resp2.Output)
}
```

### Reasoning Models

```go
func reasoningExample() {
    // Use reasoning model with completion API
    model, _ := llm.NewOpenAIModel(
        types.WithAPIKey("your-api-key"),
    )

    req := &completion.CompletionRequest{
        Model:        "o4-mini",
        Instructions: "You are a helpful assistant.",
        Messages: []*types.ModelMessage{
            {
                Role:    types.MessageRoleUser,
                Content: "Solve this logic puzzle: If all A are B, and all B are C, what can we conclude?",
            },
        },
        Options: []completion.CompletionOption{
            completion.WithReasoningEffort(completion.ReasoningEffortLow),
        },
    }

    resp, _ := model.Complete(context.Background(), req, nil)
    fmt.Printf("Response: %s\n", resp.Output)
    
    // Access reasoning tokens if available
    if resp.Usage != nil {
        fmt.Printf("Reasoning tokens: %d\n", resp.Usage.TotalReasoningTokens)
    }
}
```

### Conversation API (Reasoning Models)

```go
func conversationExample() {
    // Use conversation API for advanced reasoning
    model, _ := llm.NewOpenAIConversationModel(
        types.WithAPIKey("your-api-key"),
    )

    req := &conversation.ConversationRequest{
        Model: "o4-mini",
        Input: "Explain the theory of relativity in simple terms.",
        Options: []conversation.ResponseOption{
            conversation.WithReasoningEffort(conversation.ReasoningEffortMedium),
            conversation.WithReasoningSummary("detailed"),
        },
    }

    stream, err := model.StreamResponse(context.Background(), req, nil)
    if err != nil {
        log.Fatal(err)
    }

    for chunk := range stream {
        switch c := chunk.(type) {
        case types.StreamTextChunk:
            fmt.Print(c.Text)
        }
    }
    fmt.Println()
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
openai, _ := llm.NewOpenAIModel(
    types.WithAPIKey("key"),
    types.WithBaseURL("https://api.openai.com/v1"), // Optional
)

// DeepSeek
deepseek, _ := llm.NewDeepSeekModel(
    types.WithAPIKey("key"),
)

// Using environment variables
model, _ := llm.NewOpenAIModel(
    types.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
)
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
// Enable cost tracking with options
req := &completion.CompletionRequest{
    Model:        "gpt-4o-mini",
    Instructions: "You are a helpful assistant.",
    Messages: []*types.ModelMessage{
        {
            Role:    types.MessageRoleUser,
            Content: "Hello!",
        },
    },
    Options: []completion.CompletionOption{
        completion.WithCost(true),
        completion.WithUsage(true),
    },
}

resp, _ := model.Complete(ctx, req, nil)

// Access cost information
if resp.Cost != nil {
    fmt.Printf("Request cost: $%.6f\n", *resp.Cost)
}

// Access detailed usage
if resp.Usage != nil {
    fmt.Printf("Input tokens: %d\n", resp.Usage.TotalInputTokens)
    fmt.Printf("Output tokens: %d\n", resp.Usage.TotalOutputTokens)
    fmt.Printf("Reasoning tokens: %d\n", resp.Usage.TotalReasoningTokens)
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
git clone https://github.com/easymvp-ai/llm.git
cd llm
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

- 📖 [Documentation](https://github.com/easymvp-ai/llm/wiki)
- 🐛 [Issue Tracker](https://github.com/easymvp-ai/llm/issues)
- 💬 [Discussions](https://github.com/easymvp-ai/llm/discussions)

## Roadmap

- [ ] Additional provider support (Cohere, Hugging Face, etc.)
- [ ] Advanced retry mechanisms with exponential backoff
- [ ] Built-in prompt templates and management
- [ ] Metrics and monitoring integration
- [ ] Async/batch processing capabilities
- [ ] Model performance benchmarking tools

---

Made with ❤️ by the [EasyMVP](https://github.com/easymvp-ai) team
