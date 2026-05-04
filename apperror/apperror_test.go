package apperror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	cause := errors.New("cause")

	err := New(CodeInvalidInput, "invalid input", cause)

	assert.Equal(t, "invalid input: cause", err.Error())
	assert.Equal(t, CodeInvalidInput, err.Code())
	assert.Equal(t, "invalid input", err.Message())
	assert.ErrorIs(t, err, cause)
}

func TestErrorWithoutCause(t *testing.T) {
	err := New(CodeNotFound, "not found", nil)

	assert.Equal(t, "not found", err.Error())
	assert.Nil(t, err.Unwrap())
}

func TestConstructors(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		code Code
	}{
		{name: "internal", err: Internal("message", nil), code: CodeInternal},
		{name: "not found", err: NotFound("message", nil), code: CodeNotFound},
		{name: "already exist", err: AlreadyExist("message", nil), code: CodeAlreadyExist},
		{name: "invalid input", err: InvalidInput("message", nil), code: CodeInvalidInput},
		{name: "unauthorized", err: Unauthorized("message", nil), code: CodeUnauthorized},
		{name: "forbidden", err: Forbidden("message", nil), code: CodeForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.err.Code())
			assert.Equal(t, "message", tt.err.Message())
		})
	}
}
