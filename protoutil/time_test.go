package protoutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimeFromProto(t *testing.T) {
	now := time.Date(2026, 5, 4, 10, 30, 0, 0, time.UTC)

	assert.Equal(t, now, TimeFromProto(timestamppb.New(now)))
	assert.True(t, TimeFromProto(nil).IsZero())
}

func TestTimePtrFromProto(t *testing.T) {
	now := time.Date(2026, 5, 4, 10, 30, 0, 0, time.UTC)
	got := TimePtrFromProto(timestamppb.New(now))
	assert.Equal(t, now, *got)
	assert.Nil(t, TimePtrFromProto(nil))
}

func TestTimePtrToProto(t *testing.T) {
	now := time.Date(2026, 5, 4, 10, 30, 0, 0, time.UTC)
	assert.Equal(t, now.Unix(), TimePtrToProto(&now).AsTime().Unix())
	assert.Nil(t, TimePtrToProto(nil))
}
