package protoutil

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimeFromProto(value *timestamppb.Timestamp) time.Time {
	if value == nil {
		return time.Time{}
	}
	return value.AsTime()
}
