package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	SubScope = Scope{"subscription"}
)

type AuthCode struct {
	Code     string             `json:"code" bson:"code"`
	ClientID string             `json:"client_id" bson:"client_id"`
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id"`
	Scope    Scope              `json:"scope" bson:"scope"`
}

type AccessToken struct {
	AuthCode     string             `json:"auth_code" bson:"auth_code"`
	Token        string             `json:"token" bson:"token"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	Created      time.Time          `json:"created" bson:"created"`
	Expires      time.Time          `json:"expires" bson:"expires"`
}

type Scope struct{ string }
