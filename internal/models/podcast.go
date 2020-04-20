package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Podcast holds info of a podcast
type Podcast struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Title         string             `json:"title" bson:"title"`
	Author        string             `json:"author" bson:"author"`
	Type          string             `json:"type" bson:"type"`
	Subtitle      string             `json:"subtitle" bson:"subtitle"`
	Link          string             `json:"link" bson:"link"`
	Image         Image              `json:"image" bson:"image"`
	Explicit      string             `json:"explicit" bson:"explicit"`
	Language      string             `json:"locale" bson:"locale"`
	Keywords      []string           `json:"keywords" bson:"keywords"`
	Category      []Category         `json:"category" bson:"category"`
	PubDate       time.Time          `json:"pub_date" bson:"pub_date"`
	LastBuildDate time.Time          `json:"last_build_date" bson:"last_build_date"`
	RSS           string             `json:"rss" bson:"rss"`
}

// Episode holds info of the episode
type Episode struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	PodcastID      primitive.ObjectID `json:"podcast_id" bson:"podcast_id"`
	Title          string             `json:"title" bson:"title"`
	Subtitle       string             `json:"subtitle" bson:"subtitle"`
	Author         string             `json:"author" bson:"author"`
	Type           string             `json:"type" bson:"type"`
	Image          Image              `json:"image" bson:"image"`
	PubDate        time.Time          `json:"pub_date" bson:"pub_date"`
	Description    string             `json:"description" bson:"description"`
	Summary        string             `json:"summary" bson:"summary"`
	Season         int                `json:"season" bson:"season"`
	Episode        int                `json:"episode" bson:"episode"`
	Category       []Category         `json:"category" bson:"category"`
	Explicit       string             `json:"explicit" bson:"explicit"`
	URL            string             `json:"url" bson:"url"`
	DurationMillis int64              `json:"duration_millis" bson:"duration_millis"`
}
