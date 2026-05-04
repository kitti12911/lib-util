package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterOpFromString(t *testing.T) {
	tests := map[string]int32{
		"":                  FilterOpExact,
		"exact":             FilterOpExact,
		"like":              FilterOpLike,
		"gt":                FilterOpGT,
		"lt":                FilterOpLT,
		"gte":               FilterOpGTE,
		"lte":               FilterOpLTE,
		"null":              FilterOpNull,
		"not_null":          FilterOpNotNull,
		"like_ci":           FilterOpLikeCI,
		"in":                FilterOpIn,
		"between":           FilterOpBetween,
		"between_exclusive": FilterOpBetweenExclusive,
		"LIKE":              FilterOpLike,
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, want, FilterOpFromString[int32](input))
		})
	}
}

func TestOrderDirectionFromString(t *testing.T) {
	assert.Equal(t, OrderDirectionAsc, OrderDirectionFromString[int32](""))
	assert.Equal(t, OrderDirectionAsc, OrderDirectionFromString[int32]("asc"))
	assert.Equal(t, OrderDirectionDesc, OrderDirectionFromString[int32]("desc"))
	assert.Equal(t, OrderDirectionDesc, OrderDirectionFromString[int32]("DESC"))
}
