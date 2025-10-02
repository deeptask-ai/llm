# Folder Structure Refactoring Plan

## Current Issues
1. All provider implementations are in the root directory (cluttered)
2. No clear separation between public API and internal implementation
3. Test files mixed with implementation files
4. No examples demonstrating usage
5. Data files in root-level data/ folder

## Proposed Structure (Go Best Practices)

```
easyllm/
├── go.mod
├── go.sum
├── README.md
├── LICENSE
├── IMPROVEMENTS.md
├── .gitignore
│
├── types.go                    # Public type definitions (Model, ModelInfo, etc.)
├── errors.go                   # Public error definitions
├── options.go                  # Public configuration options
├── client.go                   # Main client factory functions
│
├── pkg/                        # Public reusable packages
│   ├── schema/                 # JSON schema generation
│   │   ├── schema.go
│   │   └── schema_test.go
│   └── pricing/                # Pricing calculations
│       ├── pricing.go
│       └── pricing_test.go
│
├── internal/                   # Private implementation
│   ├── providers/              # Provider implementations
│   │   ├── openai/
│   │   │   ├── openai.go
│   │   │   ├── openai_test.go
│   │   │   ├── completion.go
│   │   │   ├── embedding.go
│   │   │   ├── image.go
│   │   │   └── converter.go
│   │   ├── claude/
│   │   │   ├── claude.go
│   │   │   ├── claude_test.go
│   │   │   └── converter.go
│   │   ├── gemini/
│   │   │   ├── gemini.go
│   │   │   ├── gemini_test.go
│   │   │   └── converter.go
│   │   ├── deepseek/
│   │   │   ├── deepseek.go
│   │   │   ├── deepseek_test.go
│   │   │   └── converter.go
│   │   ├── azure/
│   │   │   ├── azure.go
│   │   │   ├── azure_test.go
│   │   │   └── converter.go
│   │   └── openrouter/
│   │       ├── openrouter.go
│   │       ├── openrouter_test.go
│   │       └── converter.go
│   │
│   ├── models/                 # Model data and metadata
│   │   ├── data/
│   │   │   ├── openai.json
│   │   │   ├── claude.json
│   │   │   ├── gemini.json
│   │   │   └── deepseek.json
│   │   ├── loader.go           # Model data loader
│   │   └── cache.go            # Model info caching
│   │
│   ├── http/                   # HTTP client utilities
│   │   ├── client.go
│   │   └── middleware.go
│   │
│   ├── validation/             # Input validation
│   │   ├── request.go
│   │   ├── embedding.go
│   │   └── image.go
│   │
│   └── common/                 # Common utilities
│       ├── constants.go
│       └── utils.go
│
├── examples/                   # Usage examples
│   ├── basic/
│   │   └── main.go
│   ├── streaming/
│   │   └── main.go
│   ├── embeddings/
│   │   └── main.go
│   ├── images/
│   │   └── main.go
│   ├── function_calling/
│   │   └── main.go
│   └── multi_provider/
│       └── main.go
│
└── docs/                       # Additional documentation
    ├── ARCHITECTURE.md
    └── CONTRIBUTING.md
```

## Key Improvements

### 1. Clear Public API Surface
- Root-level files expose only the public API
- Types, errors, and options clearly defined
- Easy for users to understand what's available

### 2. Internal Package Organization
- `internal/providers/`: Each provider in its own package
- `internal/models/`: Model data and metadata management
- `internal/validation/`: Input validation logic
- `internal/common/`: Shared utilities

### 3. Public Packages (pkg/)
- `pkg/schema/`: JSON schema generation (public utility)
- `pkg/pricing/`: Pricing calculation (public utility)

### 4. Examples
- Real, runnable examples for each major feature
- Helps users understand how to use the library

### 5. Better Testing
- Test files co-located with implementation
- Clear separation of unit tests and integration tests

## Migration Steps

1. Create new directory structure
2. Move provider implementations to internal/providers/
3. Move model data to internal/models/data/
4. Move utilities to appropriate internal packages
5. Create public API files at root (types.go, client.go)
6. Update all import paths
7. Create example programs
8. Update documentation
9. Run tests to verify everything works
10. Update CI/CD if applicable

## Benefits

- **Maintainability**: Clear organization makes code easier to navigate
- **Encapsulation**: Internal packages prevent external dependencies on implementation details
- **Testability**: Co-located tests improve test organization
- **Discoverability**: Examples help users get started quickly
- **Standards Compliance**: Follows Go community best practices
- **Scalability**: Easy to add new providers or features

## Backward Compatibility

Since this is a library, we need to ensure:
1. Public API remains stable (no breaking changes to exported types)
2. Import paths for public packages remain consistent
3. Gradual migration if needed (deprecate old, introduce new)
