// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package providers

import (
	"github.com/easymvp/easyllm"
)

// Provider defines the base interface that all provider implementations must satisfy
type Provider interface {
	easyllm.BaseModel
}

// CompletionProvider extends Provider with completion capabilities
type CompletionProvider interface {
	Provider
	easyllm.CompletionModel
}

// EmbeddingProvider extends Provider with embedding capabilities
type EmbeddingProvider interface {
	Provider
	easyllm.EmbeddingModel
}

// ImageProvider extends Provider with image generation capabilities
type ImageProvider interface {
	Provider
	easyllm.ImageModel
}

// FullProvider combines all provider capabilities
type FullProvider interface {
	CompletionProvider
	EmbeddingProvider
	ImageProvider
}

// Config contains common configuration for all providers
type Config struct {
	APIKey     string
	BaseURL    string
	APIVersion string
	// Provider-specific options can be added here
}

// BaseProvider provides common functionality for all providers
type BaseProvider struct {
	config *Config
}

// NewBaseProvider creates a new base provider with the given configuration
func NewBaseProvider(config *Config) *BaseProvider {
	return &BaseProvider{
		config: config,
	}
}

// GetConfig returns the provider configuration
func (b *BaseProvider) GetConfig() *Config {
	return b.config
}
