package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Podcast holds podcast info in the database
type Podcast struct {
	ID       primitive.ObjectID
	Name     string
	RSSUrl   string
	IMGUrl   string
	Episodes []Episode
}

// Episode holds information about a single episode of a podcast
type Episode struct {
	ID          primitive.ObjectID
	Name        string
	Date        time.Time
	Number      int
	Description string
}
