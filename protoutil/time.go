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

// TimePtrFromProto converts a proto timestamp to a *time.Time, preserving nil
// for an unset value (rather than the 1970 epoch AsTime gives).
func TimePtrFromProto(value *timestamppb.Timestamp) *time.Time {
	if value == nil {
		return nil
	}
	out := value.AsTime()
	return &out
}

// TimePtrToProto converts an optional time to a proto timestamp, preserving nil.
func TimePtrToProto(value *time.Time) *timestamppb.Timestamp {
	if value == nil {
		return nil
	}
	return timestamppb.New(*value)
}
