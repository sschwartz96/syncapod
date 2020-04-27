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
	PodcastID primitive.ObjectID `json:"podcast_id" bson:"podcast_id"`
	EpisodeID primitive.ObjectID `json:"episode_id" bson:"episode_id"`
	Offset    int64              `json:"offset" bson:"offset"` // milliseconds
	Played    bool               `json:"played" bson:"played"`
	LastSeen  time.Time          `json:"last_seen" bson:"last_seen"`
}

// Subscription contains the user specific data on a podcast
type Subscription struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id"`
	UserID    primitive.ObjectID   `json:"user_id" bson:"user_id"`
	PodcastID primitive.ObjectID   `json:"podcast_id" bson:"podcast_id"`
	CurEpiID  primitive.ObjectID   `json:"cur_epi_id" bson:"cur_epi_id"`
	PlayedIDs []primitive.ObjectID `json:"played_ids" bson:"played_ids"`
}

// FullSubscription is based off of Subscription but contains full podcast and
// current episode objects, to be used in database
type FullSubscription struct {
	ID            primitive.ObjectID   `json:"_id" bson:"_id"`
	UserID        primitive.ObjectID   `json:"user_id" bson:"user_id"`
	Podcast       *Podcast             `json:"podcast" bson:"podcast"`
	CurEpi        *Episode             `json:"cur_epi" bson:"cur_epi"`
	CurEpiDetails *UserEpisode         `json:"cur_epi_details" bson:"cur_epi_details"`
	PlayedIDs     []primitive.ObjectID `json:"played_ids" bson:"played_ids"`
}
