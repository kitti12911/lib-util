package huma

import (
	basehuma "github.com/danielgtaylor/huma/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return basehuma.Error400BadRequest(st.Message(), err)
	case codes.Unauthenticated:
		return basehuma.Error401Unauthorized(st.Message(), err)
	case codes.PermissionDenied:
		return basehuma.Error403Forbidden(st.Message(), err)
	case codes.NotFound:
		return basehuma.Error404NotFound(st.Message(), err)
	case codes.AlreadyExists, codes.Aborted:
		return basehuma.Error409Conflict(st.Message(), err)
	case codes.Unavailable:
		return basehuma.Error503ServiceUnavailable(st.Message(), err)
	default:
		return err
	}
}
