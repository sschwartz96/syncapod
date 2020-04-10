package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RSSFeed is the container to hold all the podcast info
type RSSFeed struct {
	Podcast Podcast `xml:"channel"`
}

// Podcast holds podcast info in the database
type Podcast struct {
	ID            primitive.ObjectID `json:"_id"  bson:"_id"`
	Title         string             `json:"title"  bson:"title"  xml:"title"`
	Author        string             `json:"author"  bson:"author"  xml:"author"`
	Type          string             `json:"type"  bson:"type"  xml:"type"`
	Subtitle      string             `json:"subtitle"  bson:"subtitle"  xml:"subtitle"`
	Summary       string             `json:"summary"  bson:"summary"  xml:"summary"`
	Link          string             `json:"link"  bson:"link"  xml:"Default link"`
	Image         Image              `json:"image"  bson:"image"  xml:"image"`
	Explicit      string             `json:"explicit"  bson:"explicit"  xml:"explicit"`
	Language      string             `json:"locale"  bson:"locale"  xml:"language"`
	Keywords      string             `json:"keywords"  bson:"keywords"  xml:"keywords"`
	Category      []Category         `json:"category"  bson:"category"  xml:"category"`
	PubDate       string             `json:"pubdate"  bson:"pubdate"  xml:"pubDate"`
	LastBuildDate string             `json:"last_build_date"  bson:"last_build_date"  xml:"lastBuildDate"`
	Episodes      []Episode          `json:"episodes"  bson:"episodes"  xml:"item"`
}

// Episode holds information about a single episode of a podcast
type Episode struct {
	ID          primitive.ObjectID `json:"_id"  bson:"_id"  xml:"id"`
	Title       string             `json:"title"  bson:"title"  xml:"title"`
	Subtitle    string             `json:"subtitle" bson:"subtitle" xml:"subtitle"`
	Author      string             `json:"author"  bson:"author"  xml:"author"`
	Type        string             `json:"type"  bson:"type"  xml:"type"`
	Image       Image              `json:"image"  bson:"image"  xml:"image"`
	PubDate     string             `json:"pub_date"  bson:"pub_date"  xml:"pubDate"`
	Description string             `json:"description"  bson:"description"  xml:"description"`
	Summary     string             `json:"summary"  bson:"summary"  xml:"summary"`
	Season      int                `json:"season"  bson:"season"  xml:"season"`
	Episode     int                `json:"episode"  bson:"episode"  xml:"episode"`
	Category    []Category         `json:"category"  bson:"category"  xml:"category"`
	Explicit    string             `json:"explicit"  bson:"explicit"  xml:"explicit"`
	Enclosure   Enclosure          `json:"enclosure" bson:"enclosure" xml:"enclosure"`
}

// Enclosure represents enclosure xml object that contains mp3 data
type Enclosure struct {
	MP3 string `json:"mp3" bson:"mp3" xml:"url,attr"`
}

// Category contains the main category and secondary categories
type Category struct {
	Text     string     `xml:"text,attr"`
	Category []Category `xml:"category"`
}

// Image is the RSS image container
type Image struct {
	Title string `json:"title"  bson:"title"  xml:"title"`
	URL   string `json:"url"  bson:"url"  xml:"url"`
}