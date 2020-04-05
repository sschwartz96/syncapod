package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model contructs all user data
type User struct {
	ID       primitive.ObjectID
	Email    string
	Username string
	Password string
	DOB      time.Time
}

// Subscription contains the user specific data on a podcast
type Subscription struct {
	PodcastID   primitive.ObjectID
	CurEpisode  int
	EpisodeTime EpisodeTime
}

// EpisodeTime represents where the user last let off on a specific episode
type EpisodeTime struct {
	Hour   int
	Second int
}
