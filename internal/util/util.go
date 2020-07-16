package util

import (
	"time"

	"github.com/sschwartz96/syncapod/internal/protos"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MillisToTimestamp takes milliseconds and converts it to protobuf timestamp
func MillisToTimestamp(millis int64) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: millis / int64(1000),
		Nanos:   int32(millis%1000) * 1000000,
	}
}

// ObjIDToBSONID takes a protobuf ObjectID and converts to primitive.ObjectID
func ObjIDToBSONID(i *protos.ObjectID) (*primitive.ObjectID, error) {
	p, err := primitive.ObjectIDFromHex(i.Hex)
	return &p, err
}

// AddToTimestamp adds duration d to timestamp
func AddToTimestamp(t *timestamppb.Timestamp, d time.Duration) *timestamppb.Timestamp {
	millis := d.Milliseconds()
	t.Seconds += millis / int64(1000)
	t.Nanos += int32(millis%1000) * 1000000
	return t
}
