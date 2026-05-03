package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToJSONStr(t *testing.T) {
	input := map[string]any{"name": "alice", "age": 30}

	compact, err := ToJSONStr(input, false)
	require.NoError(t, err)
	assert.Contains(t, compact, `"name":"alice"`)

	indented, err := ToJSONStr(input, true)
	require.NoError(t, err)
	assert.Contains(t, indented, "\n  ")
}

func TestToJSONStrError(t *testing.T) {
	_, err := ToJSONStr(func() {}, false)
	require.Error(t, err)
	assert.ErrorContains(t, err, "formatter: to json")
}
