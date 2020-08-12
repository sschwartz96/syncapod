package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// // UserDoc is the model that gets saved into the collection as a document
// type UserDoc struct {
// 	ID       primitive.ObjectID   `json:"_id" bson:"_id"`
// 	Email    string               `json:"email" bson:"email"`
// 	Username string               `json:"username" bson:"username"`
// 	Password string               `json:"password" bson:"password"`
// 	DOB      time.Time            `json:"dob" bson:"dob"`
// 	Subs     []primitive.ObjectID `json:"subs" bson:"subs"`
// }

// // User is the aggregated user type, where Subs is type []Podcast
// type User struct {
// 	ID       primitive.ObjectID `json:"_id" bson:"_id"`
// 	Email    string             `json:"email" bson:"email"`
// 	Username string             `json:"username" bson:"username"`
// 	Password string             `json:"password" bson:"password"`
// 	DOB      time.Time          `json:"dob" bson:"dob"`
// 	Subs     []Podcast          `json:"subs" bson:"subs"`
// }

// // UserEpisode contains the information of the episode specific to the user
// type UserEpisode struct {
// 	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
// 	PodcastID primitive.ObjectID `json:"podcast_id" bson:"podcast_id"`
// 	EpisodeID primitive.ObjectID `json:"episode_id" bson:"episode_id"`
// 	Offset    int64              `json:"offset" bson:"offset"` // milliseconds
// 	Played    bool               `json:"played" bson:"played"`
// 	LastSeen  time.Time          `json:"last_seen" bson:"last_seen"`
// }