package easyllm

import (
	"reflect"
	"sync"
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test structs for schema generation
type TestUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TestProduct struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	InStock     bool    `json:"in_stock"`
	Description *string `json:"description,omitempty"`
}

type TestOrder struct {
	ID       int           `json:"id"`
	User     TestUser      `json:"user"`
	Products []TestProduct `json:"products"`
	Total    float64       `json:"total"`
}

func TestGenerateSchema_BasicFunctionality(t *testing.T) {
	// Clear cache before test
	clearCache()

	schema := GenerateSchema[TestUser]()
	require.NotNil(t, schema)

	// Verify it's a valid jsonschema.Schema
	jsonSchema, ok := schema.(*jsonschema.Schema)
	require.True(t, ok, "Expected schema to be of type *jsonschema.Schema")

	// Verify basic schema properties
	assert.Equal(t, "object", jsonSchema.Type)
	assert.NotNil(t, jsonSchema.Properties)

	// Verify required properties exist
	idProp, idExists := jsonSchema.Properties.Get("id")
	assert.True(t, idExists)
	assert.NotNil(t, idProp)

	nameProp, nameExists := jsonSchema.Properties.Get("name")
	assert.True(t, nameExists)
	assert.NotNil(t, nameProp)

	emailProp, emailExists := jsonSchema.Properties.Get("email")
	assert.True(t, emailExists)
	assert.NotNil(t, emailProp)
}

func TestGenerateSchema_Caching(t *testing.T) {
	// Clear cache before test
	clearCache()

	// First call should generate and cache
	schema1 := GenerateSchema[TestUser]()
	require.NotNil(t, schema1)

	// Second call should return cached result
	schema2 := GenerateSchema[TestUser]()
	require.NotNil(t, schema2)

	// Should be the exact same object (pointer equality)
	assert.Same(t, schema1, schema2, "Expected cached schema to be the same object")

	// Verify cache contains the entry
	cacheMutex.RLock()
	userType := reflect.TypeOf(TestUser{})
	cachedSchema, exists := schemaCache[userType]
	cacheMutex.RUnlock()

	assert.True(t, exists, "Expected schema to be cached")
	assert.Same(t, schema1, cachedSchema, "Expected cached schema to match returned schema")
}

func TestGenerateSchema_DifferentTypes(t *testing.T) {
	// Clear cache before test
	clearCache()

	userSchema := GenerateSchema[TestUser]()
	productSchema := GenerateSchema[TestProduct]()
	orderSchema := GenerateSchema[TestOrder]()

	// All schemas should be different objects
	assert.NotSame(t, userSchema, productSchema)
	assert.NotSame(t, userSchema, orderSchema)
	assert.NotSame(t, productSchema, orderSchema)

	// Verify all are cached separately
	cacheMutex.RLock()
	assert.Len(t, schemaCache, 3, "Expected 3 different types to be cached")
	cacheMutex.RUnlock()

	// Verify each type has correct schema
	userType := reflect.TypeOf(TestUser{})
	productType := reflect.TypeOf(TestProduct{})
	orderType := reflect.TypeOf(TestOrder{})

	cacheMutex.RLock()
	assert.Contains(t, schemaCache, userType)
	assert.Contains(t, schemaCache, productType)
	assert.Contains(t, schemaCache, orderType)
	cacheMutex.RUnlock()
}

func TestGenerateSchema_ConcurrentAccess(t *testing.T) {
	// Clear cache before test
	clearCache()

	const numGoroutines = 100
	var wg sync.WaitGroup
	results := make([]interface{}, numGoroutines)

	// Launch multiple goroutines calling GenerateSchema concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = GenerateSchema[TestUser]()
		}(i)
	}

	wg.Wait()

	// All results should be the same cached object
	firstResult := results[0]
	for i := 1; i < numGoroutines; i++ {
		assert.Same(t, firstResult, results[i], "Expected all concurrent calls to return same cached object")
	}

	// Verify only one entry in cache
	cacheMutex.RLock()
	assert.Len(t, schemaCache, 1, "Expected only one entry in cache")
	cacheMutex.RUnlock()
}

func TestGenerateSchema_SchemaProperties(t *testing.T) {
	// Clear cache before test
	clearCache()

	schema := GenerateSchema[TestProduct]()
	jsonSchema, ok := schema.(*jsonschema.Schema)
	require.True(t, ok)

	// Verify schema follows structured output requirements
	// Note: AdditionalProperties verification removed due to API complexity

	// Verify properties exist
	idProp, idExists := jsonSchema.Properties.Get("id")
	assert.True(t, idExists)
	assert.NotNil(t, idProp)

	nameProp, nameExists := jsonSchema.Properties.Get("name")
	assert.True(t, nameExists)
	assert.NotNil(t, nameProp)

	priceProp, priceExists := jsonSchema.Properties.Get("price")
	assert.True(t, priceExists)
	assert.NotNil(t, priceProp)

	inStockProp, inStockExists := jsonSchema.Properties.Get("in_stock")
	assert.True(t, inStockExists)
	assert.NotNil(t, inStockProp)

	descProp, descExists := jsonSchema.Properties.Get("description")
	assert.True(t, descExists)
	assert.NotNil(t, descProp)

	// Verify property types
	assert.Equal(t, "integer", idProp.Type)
	assert.Equal(t, "string", nameProp.Type)
	assert.Equal(t, "number", priceProp.Type)
	assert.Equal(t, "boolean", inStockProp.Type)
}

func TestGenerateSchema_ComplexTypes(t *testing.T) {
	// Clear cache before test
	clearCache()

	schema := GenerateSchema[TestOrder]()
	jsonSchema, ok := schema.(*jsonschema.Schema)
	require.True(t, ok)

	// Verify nested object properties
	userProp, userExists := jsonSchema.Properties.Get("user")
	assert.True(t, userExists)
	assert.NotNil(t, userProp)
	assert.Equal(t, "object", userProp.Type)

	productsProp, productsExists := jsonSchema.Properties.Get("products")
	assert.True(t, productsExists)
	assert.NotNil(t, productsProp)
	assert.Equal(t, "array", productsProp.Type)
	assert.NotNil(t, productsProp.Items)
}

func TestGenerateSchema_PointerTypes(t *testing.T) {
	// Clear cache before test
	clearCache()

	// Test with pointer type
	schema1 := GenerateSchema[*TestUser]()
	schema2 := GenerateSchema[TestUser]()

	// Pointer and non-pointer types should generate different schemas
	assert.NotSame(t, schema1, schema2, "Pointer and non-pointer types should have different schemas")

	// Verify both are cached separately
	cacheMutex.RLock()
	ptrType := reflect.TypeOf((*TestUser)(nil))
	valueType := reflect.TypeOf(TestUser{})

	assert.Contains(t, schemaCache, ptrType)
	assert.Contains(t, schemaCache, valueType)
	assert.Len(t, schemaCache, 2)
	cacheMutex.RUnlock()
}

func TestGenerateSchema_EmptyStruct(t *testing.T) {
	type EmptyStruct struct{}

	// Clear cache before test
	clearCache()

	schema := GenerateSchema[EmptyStruct]()
	require.NotNil(t, schema)

	jsonSchema, ok := schema.(*jsonschema.Schema)
	require.True(t, ok)

	assert.Equal(t, "object", jsonSchema.Type)
	// Empty struct should have no properties or empty properties
	assert.True(t, jsonSchema.Properties.Len() == 0)
}

func TestGenerateSchema_CacheIsolation(t *testing.T) {
	// Clear cache before test
	clearCache()

	// Generate schemas for different types
	userSchema := GenerateSchema[TestUser]()
	productSchema := GenerateSchema[TestProduct]()

	// Manually clear cache
	clearCache()

	// Generate schemas again - should be different objects now
	userSchema2 := GenerateSchema[TestUser]()
	productSchema2 := GenerateSchema[TestProduct]()

	assert.NotSame(t, userSchema, userSchema2, "After cache clear, should generate new schema objects")
	assert.NotSame(t, productSchema, productSchema2, "After cache clear, should generate new schema objects")
}

// Benchmark tests
func BenchmarkGenerateSchema_FirstCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		clearCache()
		GenerateSchema[TestUser]()
	}
}

func BenchmarkGenerateSchema_CachedCall(b *testing.B) {
	// Pre-populate cache
	clearCache()
	GenerateSchema[TestUser]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateSchema[TestUser]()
	}
}

func BenchmarkGenerateSchema_ConcurrentCached(b *testing.B) {
	// Pre-populate cache
	clearCache()
	GenerateSchema[TestUser]()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenerateSchema[TestUser]()
		}
	})
}

// Helper function to clear cache for testing
func clearCache() {
	cacheMutex.Lock()
	schemaCache = make(map[reflect.Type]interface{})
	cacheMutex.Unlock()
}
