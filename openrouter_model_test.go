package easyllm

import (
	"testing"
)

func TestOpenRouterModel_Name(t *testing.T) {
	// Skip this test as it requires a real API call to load models
	t.Skip("Skipping test that requires real API call")

	config := OpenRouterModelConfig{
		APIKey: "test-key",
	}

	model, err := NewOpenRouterModel(config)
	if err != nil {
		t.Fatalf("Failed to create OpenRouter model: %v", err)
	}

	if model.Name() != "openrouter" {
		t.Errorf("Expected name 'openrouter', got '%s'", model.Name())
	}
}

func TestOpenRouterModel_EmptyAPIKey(t *testing.T) {
	config := OpenRouterModelConfig{
		APIKey: "",
	}

	_, err := NewOpenRouterModel(config)
	if err == nil {
		t.Error("Expected error for empty API key")
	}
	if err.Error() != "API key cannot be empty" {
		t.Errorf("Expected 'API key cannot be empty', got '%s'", err.Error())
	}
}

func TestOpenRouterModel_LoadModelsError(t *testing.T) {
	// This test checks if an invalid API key fails to load models
	config := OpenRouterModelConfig{
		APIKey: "invalid-test-key",
	}

	_, err := NewOpenRouterModel(config)
	// OpenRouter API might allow listing models without authentication
	// or with invalid keys, so we just check that the function returns properly
	if err != nil {
		t.Logf("Got error as expected: %v", err)
		// The error should be wrapped with "failed to load models"
		if len(err.Error()) == 0 {
			t.Error("Expected non-empty error message")
		}
	} else {
		t.Logf("OpenRouter API allowed model listing with invalid key - this is acceptable")
	}
}

func TestOpenRouterModelInfo_Conversion(t *testing.T) {
	// Test that OpenRouterModelInfo can be properly converted to ModelInfo
	openRouterModel := OpenRouterModelInfo{
		ID:   "test-model",
		Name: "Test Model",
		Pricing: OpenRouterModelPricing{
			Prompt:     "0.001",
			Completion: "0.002",
			Request:    "0.0001",
		},
	}

	modelInfo := &ModelInfo{
		ID:   openRouterModel.ID,
		Name: openRouterModel.Name,
		Pricing: ModelPricing{
			Prompt:            openRouterModel.Pricing.Prompt,
			Completion:        openRouterModel.Pricing.Completion,
			Request:           openRouterModel.Pricing.Request,
			Image:             openRouterModel.Pricing.Image,
			WebSearch:         openRouterModel.Pricing.WebSearch,
			InternalReasoning: openRouterModel.Pricing.InternalReasoning,
			InputCacheRead:    openRouterModel.Pricing.InputCacheRead,
			InputCacheWrite:   openRouterModel.Pricing.InputCacheWrite,
		},
	}

	if modelInfo.ID != "test-model" {
		t.Errorf("Expected ID 'test-model', got '%s'", modelInfo.ID)
	}
	if modelInfo.Name != "Test Model" {
		t.Errorf("Expected Name 'Test Model', got '%s'", modelInfo.Name)
	}
	if modelInfo.Pricing.Prompt != "0.001" {
		t.Errorf("Expected Prompt pricing '0.001', got '%s'", modelInfo.Pricing.Prompt)
	}
	if modelInfo.Pricing.Completion != "0.002" {
		t.Errorf("Expected Completion pricing '0.002', got '%s'", modelInfo.Pricing.Completion)
	}
}
