package podcast

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/schollz/closestmatch"
	"github.com/sschwartz96/syncapod/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParseRSS takes in URL path and unmarshals the data
func ParseRSS(path string) (*models.Podcast, error) {
	response, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	var rss models.RSSFeed
	decoder := xml.NewDecoder(response.Body)
	decoder.DefaultSpace = "Default"

	err = decoder.Decode(&rss)
	if err != nil {
		return nil, err
	}

	AddIDs(&rss.Podcast)

	return &rss.Podcast, nil
}

// AddIDs adds missing IDs to the podcast object and episode objects
func AddIDs(podcast *models.Podcast) {
	podcast.ID = primitive.NewObjectID()

	for i, _ := range podcast.Episodes {
		podcast.Episodes[i].ID = primitive.NewObjectID()
	}
}

func MatchTitle(search string, podcasts []models.Podcast) {
	var titles []string
	for _, podcast := range podcasts {
		titles = append(titles, podcast.Title)
	}

	bagSizes := []int{2, 3, 4}

	cm := closestmatch.New(titles, bagSizes)
	fmt.Println(cm)

	return
}
