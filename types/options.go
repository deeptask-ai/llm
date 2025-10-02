// Copyright 2025 The Go A2A Authors
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"github.com/openai/openai-go/option"
)

// ModelOption is a functional option for configuring models
type ModelOption func(*ModelOptions)

// ModelOptions contains configuration options for all models
type ModelOptions struct {
	APIKey     string
	BaseURL    string
	APIVersion string // For Azure OpenAI
	Options    []option.RequestOption
}

// WithAPIKey sets the API key
func WithAPIKey(apiKey string) ModelOption {
	return func(o *ModelOptions) {
		o.APIKey = apiKey
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) ModelOption {
	return func(o *ModelOptions) {
		o.BaseURL = baseURL
	}
}

// WithAPIVersion sets the API version (for Azure OpenAI)
func WithAPIVersion(version string) ModelOption {
	return func(o *ModelOptions) {
		o.APIVersion = version
	}
}

// WithRequestOption adds a custom request option from the OpenAI SDK
func WithRequestOption(opt option.RequestOption) ModelOption {
	return func(o *ModelOptions) {
		o.Options = append(o.Options, opt)
	}
}

// WithRequestOptions adds multiple custom request options from the OpenAI SDK
func WithRequestOptions(opts ...option.RequestOption) ModelOption {
	return func(o *ModelOptions) {
		o.Options = append(o.Options, opts...)
	}
}

// applyOptions applies all options to create a ModelOptions struct
func ApplyOptions(opts []ModelOption) *ModelOptions {
	options := &ModelOptions{
		Options: make([]option.RequestOption, 0),
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
