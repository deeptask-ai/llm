// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package llm

import (
	"errors"
	"fmt"
)

// Standard error types for better error handling
var (
	// ErrAPIKeyEmpty is returned when API key is empty
	ErrAPIKeyEmpty = errors.New("API key cannot be empty")

	// ErrBaseURLEmpty is returned when base URL is empty
	ErrBaseURLEmpty = errors.New("base URL cannot be empty")

	// ErrAPIVersionEmpty is returned when API version is empty
	ErrAPIVersionEmpty = errors.New("API version cannot be empty")

	// ErrInvalidRequest is returned when request validation fails
	ErrInvalidRequest = errors.New("invalid request")

	// ErrEmptyContent is returned when content is empty
	ErrEmptyContent = errors.New("content cannot be empty")

	ErrEmptyInstructions = errors.New("instructions cannot be empty")

	// ErrInvalidModel is returned when model is invalid
	ErrInvalidModel = errors.New("invalid model specified")

	// ErrContextCanceled is returned when context is canceled
	ErrContextCanceled = errors.New("context canceled")

	// ErrTimeout is returned when operation times out
	ErrTimeout = errors.New("operation timeout")
)

// ValidationError represents a validation error with field details
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e *ValidationError) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("validation failed for field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
	}
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string, value interface{}) error {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// RequestError represents an error from API request
type RequestError struct {
	Provider   string
	StatusCode int
	Message    string
	Err        error
}

func (e *RequestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s API error (status %d): %s: %v", e.Provider, e.StatusCode, e.Message, e.Err)
	}
	return fmt.Sprintf("%s API error (status %d): %s", e.Provider, e.StatusCode, e.Message)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

// NewRequestError creates a new request error
func NewRequestError(provider string, statusCode int, message string, err error) error {
	return &RequestError{
		Provider:   provider,
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

// ResponseError represents an error from API response
type ResponseError struct {
	Provider string
	Message  string
	Err      error
}

func (e *ResponseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s response error: %s: %v", e.Provider, e.Message, e.Err)
	}
	return fmt.Sprintf("%s response error: %s", e.Provider, e.Message)
}

func (e *ResponseError) Unwrap() error {
	return e.Err
}

// NewResponseError creates a new response error
func NewResponseError(provider, message string, err error) error {
	return &ResponseError{
		Provider: provider,
		Message:  message,
		Err:      err,
	}
}

// UnsupportedCapabilityError represents an unsupported feature error
type UnsupportedCapabilityError struct {
	Provider   string
	Capability string
}

func (e *UnsupportedCapabilityError) Error() string {
	if e.Capability == "image generation" {
		return fmt.Sprintf("%s is not supported by %s models", e.Capability, e.Provider)
	}
	return fmt.Sprintf("%s are not supported by %s models", e.Capability, e.Provider)
}

// NewUnsupportedCapabilityError creates a new unsupported capability error
func NewUnsupportedCapabilityError(provider, capability string) error {
	return &UnsupportedCapabilityError{
		Provider:   provider,
		Capability: capability,
	}
}

// StreamError represents an error in streaming
type StreamError struct {
	Provider string
	Message  string
	Err      error
}

func (e *StreamError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s stream error: %s: %v", e.Provider, e.Message, e.Err)
	}
	return fmt.Sprintf("%s stream error: %s", e.Provider, e.Message)
}

func (e *StreamError) Unwrap() error {
	return e.Err
}

// NewStreamError creates a new stream error
func NewStreamError(provider, message string, err error) error {
	return &StreamError{
		Provider: provider,
		Message:  message,
		Err:      err,
	}
}
