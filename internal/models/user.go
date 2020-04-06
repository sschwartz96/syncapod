package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model contructs all user data
type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	DOB      time.Time          `json:"dob" bson:"dob"`
}

// Subscription contains the user specific data on a podcast
type Subscription struct {
	PodcastID   primitive.ObjectID `json:"podcast_id" bson:"podcast_id"`
	CurEpisode  int                `json:"cur_episode" bson:"cur_episode"`
	EpisodeTime EpisodeTime        `json:"episode_time" bson:"episode_time"`
}

// EpisodeTime represents where the user last let off on a specific episode
type EpisodeTime struct {
	Hour   int `json:"hour" bson:"hour"`
	Second int `json:"second" bson:"second"`
}
