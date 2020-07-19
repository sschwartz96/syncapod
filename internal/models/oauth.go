package models

import (
	"time"

	"github.com/sschwartz96/syncapod/internal/protos"
)

// Scopes of oauth2.0
var (
	SubScope = Scope{"subscription"}
)

// AuthCode is the authorization code of oauth2.0
type AuthCode struct {
	Code     string           `json:"code" bson:"code"`
	ClientID string           `json:"client_id" bson:"client_id"`
	UserID   *protos.ObjectID `json:"user_id" bson:"user_id"`
	Scope    Scope            `json:"scope" bson:"scope"`
}

// AccessToken contains the information to provide user access
type AccessToken struct {
	AuthCode     string           `json:"auth_code" bson:"auth_code"`
	Token        string           `json:"token" bson:"token"`
	RefreshToken string           `json:"refresh_token" bson:"refresh_token"`
	UserID       *protos.ObjectID `json:"user_id" bson:"user_id"`
	Created      time.Time        `json:"created" bson:"created"`
	Expires      int              `json:"expires" bson:"expires"`
}

// Scope is just a wrapper to the string
type Scope struct{ string }
