package llmclient

import (
	"testing"
)

func TestGeminiModel_Name(t *testing.T) {
	// Test creating a Gemini model
	config := GeminiModelConfig{
		APIKey: "test-key",
	}

	model, err := NewGeminiModel(config)
	if err != nil {
		t.Fatalf("Failed to create Gemini model: %v", err)
	}

	// Test name
	if model.Name() != "gemini" {
		t.Errorf("Expected name 'gemini', got '%s'", model.Name())
	}
}

func TestGeminiModel_SupportedModels(t *testing.T) {
	config := GeminiModelConfig{
		APIKey: "test-key",
	}

	model, err := NewGeminiModel(config)
	if err != nil {
		t.Fatalf("Failed to create Gemini model: %v", err)
	}

	// Test supported models
	models := model.SupportedModels()
	if models == nil {
		t.Fatal("SupportedModels returned nil - JSON unmarshaling failed")
	}
	if len(models) == 0 {
		t.Error("Expected at least one supported model")
	}

	t.Logf("Found %d models", len(models))
	for i, model := range models {
		t.Logf("Model %d: ID=%s, Name=%s", i, model.ID, model.Name)
	}

	// Check for specific models
	foundGemini25 := false
	foundGemini15Pro := false

	for _, model := range models {
		if model.ID == "gemini-2.5-flash" {
			foundGemini25 = true
		}
		if model.ID == "gemini-1.5-pro" {
			foundGemini15Pro = true
		}
	}

	if !foundGemini25 {
		t.Error("Expected to find gemini-2.5-flash model")
	}

	if !foundGemini15Pro {
		t.Error("Expected to find gemini-1.5-pro model")
	}
}

func TestGeminiModel_EmptyAPIKey(t *testing.T) {
	config := GeminiModelConfig{
		APIKey: "",
	}

	_, err := NewGeminiModel(config)
	if err == nil {
		t.Error("Expected error for empty API key")
	}
	if err.Error() != "API key cannot be empty" {
		t.Errorf("Expected 'API key cannot be empty', got '%s'", err.Error())
	}
}

func TestGeminiModel_ValidConfig(t *testing.T) {
	config := GeminiModelConfig{
		APIKey: "test-key",
	}

	model, err := NewGeminiModel(config)
	if err != nil {
		t.Fatalf("Failed to create Gemini model with valid config: %v", err)
	}

	if model == nil {
		t.Fatal("Model should not be nil with valid config")
	}
}
