package database

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/util"
)

var (
	timestampKind = reflect.ValueOf(timestamppb.Timestamp{}).Kind()
	objectIDKind  = reflect.ValueOf(protos.ObjectID{}).Kind()

	timestampType = reflect.ValueOf(timestamppb.Timestamp{}).Type()
	objectIDType  = reflect.ValueOf(protos.ObjectID{}).Type()
)

// ObjectIDCodec is the codec responsible for encoding and decoding the custom
// protobuf objectID to the mongodb primitive.Object
type ObjectIDCodec struct{}

// EncodeValue encodes protos.ObjectID to primitive.ObjectID
func (o *ObjectIDCodec) EncodeValue(en bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if val.Kind() != objectIDKind {
		return bsoncodec.ValueEncoderError{Name: "ObjectIDEncodeValue", Kinds: []reflect.Kind{objectIDKind}, Received: val}
	}

	obj := &protos.ObjectID{Hex: val.Interface().(protos.ObjectID).Hex}
	pObj, err := primitive.ObjectIDFromHex(obj.Hex)
	if err != nil {
		return err
	}

	return vw.WriteObjectID(pObj)
}

// DecodeValue decodes primitive.ObjectID to protos.ObjectID
func (o *ObjectIDCodec) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Kind() != objectIDKind {
		return bsoncodec.ValueDecoderError{Name: "ObjectIDDecodeValue", Kinds: []reflect.Kind{objectIDKind}, Received: val}
	}

	if vr.Type() != bsontype.ObjectID {
		return fmt.Errorf("cannot decode %v into protobuf objectID type", vr.Type())
	}

	pObj, err := vr.ReadObjectID()
	if err != nil {
		return err
	}

	val.Set(reflect.ValueOf(protos.ObjectID{
		Hex: pObj.Hex(),
	}))
	return nil
}

// TimestampCodec is the codec responsible for encoding and decoding timestamp values
// implements ValueCodec which in turn implements ValueEncoder & ValueDecoder. All of
// the bsoncodec package of mongo driver
type TimestampCodec struct{}

// EncodeValue encodes a timestamp to bson
func (t *TimestampCodec) EncodeValue(en bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if val.Kind() != timestampKind {
		return bsoncodec.ValueEncoderError{Name: "TimestampEncodeValue", Kinds: []reflect.Kind{timestampKind}, Received: val}
	}

	seconds := val.Interface().(timestamppb.Timestamp).Seconds
	nanos := val.Interface().(timestamppb.Timestamp).Nanos

	return vw.WriteDateTime(seconds*1000 + int64(nanos)/1000000)
}

// DecodeValue decodes the bson to timestamp
func (t *TimestampCodec) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Kind() != timestampKind {
		return bsoncodec.ValueDecoderError{Name: "TimestampDecodeValue", Kinds: []reflect.Kind{timestampKind}, Received: val}
	}

	if vr.Type() != bsontype.DateTime {
		return fmt.Errorf("cannot decode %v into timestamp type", vr.Type())
	}

	millis, err := vr.ReadDateTime()
	if err != nil {
		return err
	}

	val.Set(reflect.ValueOf(
		*util.MillisToTimestamp(millis),
	))

	return nil
}

func createRegistry() *bsoncodec.Registry {
	// register our custom codec
	rb := bsoncodec.NewRegistryBuilder().
		RegisterCodec(timestampType, &TimestampCodec{}).
		RegisterCodec(objectIDType, &ObjectIDCodec{})

	// register the default encoders & decoders, and primitive codecs
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bson.PrimitiveCodecs{}.RegisterPrimitiveCodecs(rb)

	return rb.Build()
}
