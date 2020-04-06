package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Podcast holds podcast info in the database
type Podcast struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	RSSUrl   string             `json:"rss_url" bson:"rss_url"`
	IMGUrl   string             `json:"img_url" bson:"img_url"`
	Episodes []Episode          `json:"episodes" bson:"episodes"`
}

// Episode holds information about a single episode of a podcast
type Episode struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Date        time.Time          `json:"date" bson:"date"`
	Number      int                `json:"number" bson:"number"`
	Description string             `json:"description" bson:"description"`
}
