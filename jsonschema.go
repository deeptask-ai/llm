package llmclient

import (
	"encoding/json"
	"reflect"
	"sync"

	"github.com/invopop/jsonschema"
)

var (
	schemaCache = make(map[reflect.Type]interface{})
	cacheMutex  sync.RWMutex
)

func GenerateSchemaText[T any]() string {
	schema := GenerateSchema[T]()
	jsonBytes, err := json.Marshal(schema)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

type NoneSchema struct{}

func EmptySchema() interface{} {
	return GenerateSchema[NoneSchema]()
}

func GenerateSchema[T any]() interface{} {
	var v T
	typeKey := reflect.TypeOf(v)

	// Check cache first
	cacheMutex.RLock()
	if cachedSchema, exists := schemaCache[typeKey]; exists {
		cacheMutex.RUnlock()
		return cachedSchema
	}
	cacheMutex.RUnlock()

	// Generate schema if not cached
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	schema := reflector.Reflect(v)

	// Cache the result
	cacheMutex.Lock()
	schemaCache[typeKey] = schema
	cacheMutex.Unlock()

	return schema
}
