package huma

import (
	"net/http"
	"testing"

	basehuma "github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
)

func TestWithTag(t *testing.T) {
	op := &basehuma.Operation{}

	WithTag("Users")(op)

	assert.Equal(t, []string{"Users"}, op.Tags)
}

func TestStatusCreated(t *testing.T) {
	op := &basehuma.Operation{}

	StatusCreated(op)

	assert.Equal(t, http.StatusCreated, op.DefaultStatus)
}

func TestAffectedRows(t *testing.T) {
	out := AffectedRows(7)

	assert.Equal(t, 7, out.Body.AffectedRows)
}
