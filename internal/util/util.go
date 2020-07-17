package util

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// MillisToTimestamp takes milliseconds and converts it to protobuf timestamp
func MillisToTimestamp(millis int64) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: millis / int64(1000),
		Nanos:   int32(millis%1000) * 1000000,
	}
}

// AddToTimestamp adds duration d to timestamp
func AddToTimestamp(t *timestamppb.Timestamp, d time.Duration) *timestamppb.Timestamp {
	millis := d.Milliseconds()
	t.Seconds += millis / int64(1000)
	t.Nanos += int32(millis%1000) * 1000000
	return t
}
