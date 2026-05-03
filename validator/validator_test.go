package validator

import (
	"errors"
	"testing"

	basevalidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sampleConfig struct {
	Host string `mapstructure:"host" validate:"required"`
	Name string `json:"name" validate:"required"`
	Skip string `mapstructure:"-" json:"-" validate:"required"`
	Port int    `validate:"gte=1"`
}

type customConfig struct {
	Code string `validate:"starts_with_x"`
}

func TestValidate(t *testing.T) {
	val := New()
	err := val.RegisterCustom("starts_with_x", func(fl basevalidator.FieldLevel) bool {
		return len(fl.Field().String()) > 0 && fl.Field().String()[0] == 'x'
	})
	require.NoError(t, err)

	err = val.Validate(sampleConfig{
		Host: "localhost",
		Name: "service",
		Skip: "value",
		Port: 1,
	})
	require.NoError(t, err)

	err = val.Validate(customConfig{
		Code: "x-1",
	})
	require.NoError(t, err)
	assert.NotNil(t, val.Engine())
}

func TestValidateWithErrors(t *testing.T) {
	val := New()

	violations, err := val.ValidateWithErrors(sampleConfig{})
	require.Error(t, err)
	require.Len(t, violations, 4)

	wantFields := map[string]bool{
		"host": true,
		"name": true,
		"Skip": true,
		"Port": true,
	}
	for _, violation := range violations {
		assert.True(t, wantFields[violation.Field], "unexpected violation field %q", violation.Field)
	}
}

func TestValidateWithErrorsReturnsNonValidationError(t *testing.T) {
	val := New()

	violations, err := val.ValidateWithErrors(42)
	require.Error(t, err)
	assert.Nil(t, violations)

	var invalid *basevalidator.InvalidValidationError
	assert.True(t, errors.As(err, &invalid))
}

func TestValidateWithErrorsValid(t *testing.T) {
	val := New()
	violations, err := val.ValidateWithErrors(struct{}{})
	require.NoError(t, err)
	assert.Nil(t, violations)
}
