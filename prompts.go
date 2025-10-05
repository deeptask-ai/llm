package llm

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
)

// Template cache for better performance using sync.Map for concurrent read-heavy workloads
var promptCache sync.Map

// GetPrompts executes a template with caching for better performance
func GetPrompts(prompt string, params map[string]interface{}) (string, error) {
	// Try to get cached template
	if cached, ok := promptCache.Load(prompt); ok {
		tmpl := cached.(*template.Template)
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, params); err != nil {
			return "", fmt.Errorf("failed to execute template: %w", err)
		}
		return buf.String(), nil
	}

	// Parse new template
	tmpl, err := template.New("prompt").Parse(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Use LoadOrStore to handle race conditions
	actual, _ := promptCache.LoadOrStore(prompt, tmpl)

	var buf bytes.Buffer
	if err := actual.(*template.Template).Execute(&buf, params); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ClearPromptCache clears the template cache to free memory
func ClearPromptCache() {
	promptCache = sync.Map{}
}
