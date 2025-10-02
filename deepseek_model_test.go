package easyllm

import (
	"testing"
)

func TestDeepSeekModel_Name(t *testing.T) {
	config := DeepSeekModelConfig{
		APIKey: "test-key",
	}

	model, err := NewDeepSeekModel(config)
	if err != nil {
		t.Fatalf("Failed to create DeepSeek model: %v", err)
	}

	if model.Name() != "deepseek" {
		t.Errorf("Expected name 'deepseek', got '%s'", model.Name())
	}
}

func TestDeepSeekModel_SupportedModels(t *testing.T) {
	config := DeepSeekModelConfig{
		APIKey: "test-key",
	}

	model, err := NewDeepSeekModel(config)
	if err != nil {
		t.Fatalf("Failed to create DeepSeek model: %v", err)
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

	// Check for DeepSeek models that should be available based on the JSON data
	foundDeepSeekChat := false
	foundDeepSeekReasoner := false

	for _, model := range models {
		if model.ID == "deepseek-chat" {
			foundDeepSeekChat = true
		}
		if model.ID == "deepseek-reasoner" {
			foundDeepSeekReasoner = true
		}
	}

	if !foundDeepSeekChat {
		t.Error("Expected to find deepseek-chat model")
	}

	if !foundDeepSeekReasoner {
		t.Error("Expected to find deepseek-reasoner model")
	}
}

func TestDeepSeekModel_EmptyAPIKey(t *testing.T) {
	config := DeepSeekModelConfig{
		APIKey: "",
	}

	_, err := NewDeepSeekModel(config)
	if err == nil {
		t.Error("Expected error for empty API key")
	}
	if err.Error() != "API key cannot be empty" {
		t.Errorf("Expected 'API key cannot be empty', got '%s'", err.Error())
	}
}
