package podcast

import (
	"encoding/xml"
	"net/http"

	"github.com/sschwartz96/syncapod/internal/models"
)

// ParseRSS takes in URL path and unmarshals the data
func ParseRSS(path string) (*models.RSSPodcast, error) {
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

	//AddEpiIDs(&rss.Podcast)

	return &rss.Podcast, nil
}
