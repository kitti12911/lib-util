package huma

import (
	"encoding/json"
	"testing"

	basehuma "github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatchUnmarshalValue(t *testing.T) {
	var patch Patch[string]

	require.NoError(t, json.Unmarshal([]byte(`"hello"`), &patch))

	assert.True(t, patch.Set)
	assert.False(t, patch.Null)
	assert.Equal(t, "hello", patch.Value)
}

func TestPatchUnmarshalNull(t *testing.T) {
	patch := Patch[string]{Value: "old"}

	require.NoError(t, json.Unmarshal([]byte(`null`), &patch))

	assert.True(t, patch.Set)
	assert.True(t, patch.Null)
	assert.Empty(t, patch.Value)
}

func TestPatchUnmarshalEmptyData(t *testing.T) {
	patch := Patch[string]{Value: "old"}

	require.NoError(t, patch.UnmarshalJSON(nil))

	assert.False(t, patch.Set)
	assert.False(t, patch.Null)
	assert.Equal(t, "old", patch.Value)
}

func TestPatchUnmarshalInvalidValue(t *testing.T) {
	var patch Patch[string]

	err := patch.UnmarshalJSON([]byte(`{`))

	require.Error(t, err)
	assert.True(t, patch.Set)
	assert.False(t, patch.Null)
}

func TestPatchMarshal(t *testing.T) {
	unset, err := json.Marshal(Patch[string]{})
	require.NoError(t, err)
	assert.Equal(t, `null`, string(unset))

	nullValue, err := json.Marshal(Patch[string]{Set: true, Null: true, Value: "ignored"})
	require.NoError(t, err)
	assert.Equal(t, `null`, string(nullValue))

	value, err := json.Marshal(Patch[string]{Set: true, Value: "hello"})
	require.NoError(t, err)
	assert.Equal(t, `"hello"`, string(value))
}

func TestPatchSchemaIsNullable(t *testing.T) {
	registry := basehuma.NewMapRegistry("#/components/schemas/", basehuma.DefaultSchemaNamer)

	schema := Patch[string]{}.Schema(registry)

	require.NotNil(t, schema)
	assert.True(t, schema.Nullable)
}

func TestPatchSchemaFallbackForNilType(t *testing.T) {
	registry := basehuma.NewMapRegistry("#/components/schemas/", basehuma.DefaultSchemaNamer)

	schema := Patch[any]{}.Schema(registry)

	require.NotNil(t, schema)
	assert.True(t, schema.Nullable)
}
