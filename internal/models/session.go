package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session controls user logins
type Session struct {
	ID           primitive.ObjectID `json:"_id"  bson:"_id"`
	UserID       primitive.ObjectID `json:"user_id"  bson:"user_id"`
	SessionKey   string             `json:"session_key"  bson:"session_key"`
	LoginTime    time.Time          `json:"login_time"  bson:"login_time"`
	LastSeenTime time.Time          `json:"last_seen_time"  bson:"last_seen_time"`
	Expires      time.Time          `json:"expires"  bson:"expires"`
}
