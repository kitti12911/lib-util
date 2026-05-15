package huma

import (
	"errors"
	"testing"

	basehuma "github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCErrorMapsStatusCodes(t *testing.T) {
	tests := []struct {
		name string
		code codes.Code
		want int
	}{
		{name: "invalid argument", code: codes.InvalidArgument, want: 400},
		{name: "out of range", code: codes.OutOfRange, want: 400},
		{name: "unauthenticated", code: codes.Unauthenticated, want: 401},
		{name: "permission denied", code: codes.PermissionDenied, want: 403},
		{name: "not found", code: codes.NotFound, want: 404},
		{name: "already exists", code: codes.AlreadyExists, want: 409},
		{name: "aborted", code: codes.Aborted, want: 409},
		{name: "failed precondition", code: codes.FailedPrecondition, want: 412},
		{name: "resource exhausted", code: codes.ResourceExhausted, want: 429},
		{name: "unimplemented", code: codes.Unimplemented, want: 501},
		{name: "unavailable", code: codes.Unavailable, want: 503},
		{name: "deadline exceeded", code: codes.DeadlineExceeded, want: 504},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GRPCError(status.Error(tt.code, "boom"))

			var statusErr *basehuma.ErrorModel
			require.ErrorAs(t, err, &statusErr)
			assert.Equal(t, tt.want, statusErr.Status)
			assert.Equal(t, "boom", statusErr.Detail)
		})
	}
}

func TestGRPCErrorReturnsOriginalError(t *testing.T) {
	err := errors.New("plain")
	internalErr := status.Error(codes.Internal, "boom")

	assert.Same(t, err, GRPCError(err))
	assert.Same(t, internalErr, GRPCError(internalErr))
}
