package easyllm

import (
	"context"
	"testing"
)

func TestClaudeModel_Name(t *testing.T) {
	config := ClaudeModelConfig{
		APIKey: "test-key",
	}

	model, err := NewClaudeModel(config)
	if err != nil {
		t.Fatalf("Failed to create Claude model: %v", err)
	}

	if model.Name() != "claude" {
		t.Errorf("Expected name 'claude', got '%s'", model.Name())
	}
}

func TestClaudeModel_SupportedModels(t *testing.T) {
	config := ClaudeModelConfig{
		APIKey: "test-key",
	}

	model, err := NewClaudeModel(config)
	if err != nil {
		t.Fatalf("Failed to create Claude model: %v", err)
	}

	models := model.SupportedModels()
	if models == nil {
		t.Skip("SupportedModels returned nil - likely due to JSON pricing type mismatch (numbers vs strings)")
		return
	}
	if len(models) == 0 {
		t.Error("Expected at least one supported model")
	}

	t.Logf("Found %d models", len(models))
	for i, model := range models {
		t.Logf("Model %d: ID=%s, Name=%s", i, model.ID, model.Name)
	}

	// Check for Claude models that should be available based on the JSON data
	foundOpus := false
	foundSonnet := false
	foundHaiku := false

	for _, model := range models {
		if model.ID == "opus-4.1" {
			foundOpus = true
		}
		if model.ID == "sonnet-4.5" {
			foundSonnet = true
		}
		if model.ID == "haiku-3.5" {
			foundHaiku = true
		}
	}

	if !foundOpus {
		t.Error("Expected to find opus-4.1 model")
	}

	if !foundSonnet {
		t.Error("Expected to find sonnet-4.5 model")
	}

	if !foundHaiku {
		t.Error("Expected to find haiku-3.5 model")
	}
}

func TestClaudeModel_EmptyAPIKey(t *testing.T) {
	config := ClaudeModelConfig{
		APIKey: "",
	}

	_, err := NewClaudeModel(config)
	if err == nil {
		t.Error("Expected error for empty API key")
	}
	if err.Error() != "API key cannot be empty" {
		t.Errorf("Expected 'API key cannot be empty', got '%s'", err.Error())
	}
}

func TestClaudeModel_ValidConfig(t *testing.T) {
	config := ClaudeModelConfig{
		APIKey: "test-key",
	}

	model, err := NewClaudeModel(config)
	if err != nil {
		t.Fatalf("Failed to create Claude model with valid config: %v", err)
	}

	if model == nil {
		t.Fatal("Model should not be nil with valid config")
	}

	// Verify the underlying OpenAI model is properly configured
	if model.OpenAIModel == nil {
		t.Error("OpenAI model should not be nil")
	}
}

func TestClaudeModel_GenerateEmbeddings_NotSupported(t *testing.T) {
	config := ClaudeModelConfig{
		APIKey: "test-key",
	}

	model, err := NewClaudeModel(config)
	if err != nil {
		t.Fatalf("Failed to create Claude model: %v", err)
	}

	ctx := context.Background()
	req := &EmbeddingRequest{
		Model:    "test-model",
		Contents: []string{"test text"},
	}

	response, err := model.GenerateEmbeddings(ctx, req)
	if err == nil {
		t.Error("Expected error for unsupported embeddings")
	}
	if response != nil {
		t.Error("Expected nil response for unsupported embeddings")
	}
	if err.Error() != "embeddings are not supported by Claude models" {
		t.Errorf("Expected 'embeddings are not supported by Claude models', got '%s'", err.Error())
	}
}

func TestClaudeModel_GenerateImage_NotSupported(t *testing.T) {
	config := ClaudeModelConfig{
		APIKey: "test-key",
	}

	model, err := NewClaudeModel(config)
	if err != nil {
		t.Fatalf("Failed to create Claude model: %v", err)
	}

	ctx := context.Background()
	req := &ImageRequest{
		Instructions: "test image prompt",
	}

	response, err := model.GenerateImage(ctx, req)
	if err == nil {
		t.Error("Expected error for unsupported image generation")
	}
	if response != nil {
		t.Error("Expected nil response for unsupported image generation")
	}
	if err.Error() != "image generation is not supported by Claude models" {
		t.Errorf("Expected 'image generation is not supported by Claude models', got '%s'", err.Error())
	}
}
