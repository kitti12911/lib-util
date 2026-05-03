package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrom(t *testing.T) {
	value := "hello"
	assert.Equal(t, "hello", From(&value))
	assert.Empty(t, From[string](nil))
}

func TestValueOr(t *testing.T) {
	value := 42
	assert.Equal(t, 42, ValueOr(&value, 7))
	assert.Equal(t, 7, ValueOr(nil, 7))
}
