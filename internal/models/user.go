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

// UserEpisode contains the information of the episode specific to the user
type UserEpisode struct {
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	EpisodeID primitive.ObjectID `json:"episode_id" bson:"episode_id"`
	Offset    int64              `json:"offset" bson:"offset"` // milliseconds
	Played    bool               `json:"played" bson:"played"`
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
