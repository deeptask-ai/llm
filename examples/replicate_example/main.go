// Copyright 2025 The DeepTask Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"github.com/easyagent-dev/llm"
	"github.com/easyagent-dev/llm/providers"
	"log"
	"os"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("REPLICATE_API_KEY")
	if apiKey == "" {
		log.Fatal("REPLICATE_API_KEY environment variable is required")
	}

	// Create Replicate image model client
	provider, err := providers.NewReplicateModelProvider(llm.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create Replicate model: %v", err)
	}

	ctx := context.Background()

	// Generate a cat image using FLUX 1.1 Pro
	fmt.Println("=== Generating Cat Image with FLUX 1.1 Pro ===")
	req := &llm.ImageRequest{
		Model:        "black-forest-labs/flux-1.1-pro",
		Instructions: "A cute cat sitting on a windowsill, looking outside at a beautiful sunset. Photorealistic, high quality, detailed fur.",
		Config: &llm.ImageModelConfig{
			Size: "1024x1024",
		},
	}
	model, err := provider.NewImageModel("black-forest-labs/flux-1.1-pro")
	if err != nil {
		log.Fatalf("Failed to create image model: %v", err)
	}
	fmt.Println("Sending request to Replicate...")
	resp, err := model.GenerateImage(ctx, req)
	if err != nil {
		log.Fatalf("Image generation failed: %v", err)
	}

	// Save the generated image
	outputPath := "generated_cat.png"
	err = os.WriteFile(outputPath, resp.Output, 0644)
	if err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	fmt.Printf("✓ Image generated successfully!\n")
	fmt.Printf("✓ Saved to: %s\n", outputPath)

	if resp.Usage != nil {
		fmt.Printf("✓ Usage: %d images, %d requests\n", resp.Usage.TotalImages, resp.Usage.TotalRequests)
	}

	if resp.Cost != nil {
		fmt.Printf("✓ Cost: $%.4f\n", *resp.Cost)
	}
}
