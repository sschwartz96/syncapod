package protos

import "go.mongodb.org/mongo-driver/bson/primitive"

// ToBSONID takes a protobuf ObjectID and converts to primitive.ObjectID
func (o *ObjectID) ToBSONID() (*primitive.ObjectID, error) {
	p, err := primitive.ObjectIDFromHex(o.Hex)
	return &p, err
}

// FromBSONID converts a
func FromBSONID(o *primitive.ObjectID) *ObjectID {
	return &ObjectID{Hex: o.Hex()}
}

// NewObjectID generates a new ObjectID based on the mongodb primitive package
func NewObjectID() *ObjectID {
	return &ObjectID{Hex: primitive.NewObjectID().Hex()}
}

// ObjectIDFromHex takes the hex and wraps it in a ObjectID struct
func ObjectIDFromHex(hex string) *ObjectID {
	return &ObjectID{Hex: hex}
}
