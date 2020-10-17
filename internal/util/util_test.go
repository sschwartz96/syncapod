package util

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMillisToTimestamp(t *testing.T) {
	type args struct {
		millis int64
	}
	tests := []struct {
		name string
		args args
		want *timestamppb.Timestamp
	}{
		{
			name: "negative",
			args: args{
				millis: -3600,
			},
			want: &timestamppb.Timestamp{Seconds: 3, Nanos: 600000000},
		},
		{
			name: "positive",
			args: args{
				millis: 3600,
			},
			want: &timestamppb.Timestamp{Seconds: 3, Nanos: 600000000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MillisToTimestamp(tt.args.millis); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MillisToTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddToTimestamp(t *testing.T) {
	type args struct {
		t *timestamppb.Timestamp
		d time.Duration
	}
	hourAdd, _ := ptypes.TimestampProto(time.Now().Add(time.Hour))
	tests := []struct {
		name string
		args args
		want *timestamppb.Timestamp
	}{
		{
			name: "nil",
			args: args{t: nil},
			want: ptypes.TimestampNow(),
		},
		{
			name: "hour",
			args: args{
				t: ptypes.TimestampNow(),
				d: time.Hour,
			},
			want: hourAdd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddToTimestamp(tt.args.t, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddToTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}
