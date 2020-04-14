package podcast

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/schollz/closestmatch"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"go.mongodb.org/mongo-driver/bson"
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

// TODO
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

// FindOffset takes database client and pointers to user and episode to lookup episode details and offset
func FindOffset(dbClient *database.Client, user *models.User, episode *models.Episode) int64 {
	var userEpi models.UserEpisode
	filter := bson.D{{Key: "user_id", Value: user.ID}, {Key: "episode_id", Value: episode.ID}}
	err := dbClient.FindWithBSON(database.ColUserEpisode, filter, &userEpi)
	if err != nil {
		fmt.Println("error finding user episode details: ", err)
		return 0
	}
	return userEpi.Offset
}

// UpdateOffset takes userID epiID and offset and performs upsert to the UserEpisode collection
func UpdateOffset(dbClient *database.Client, userID, epiID string, offset int64) {
	uID, _ := primitive.ObjectIDFromHex(userID)
	eID, _ := primitive.ObjectIDFromHex(epiID)

	userEpi := &models.UserEpisode{UserID: uID, EpisodeID: eID, Offset: offset, Played: false}

	err := dbClient.Upsert(database.ColUserEpisode, bson.D{
		{Key: "user_id", Value: uID},
		{Key: "episode_id", Value: eID}},
		userEpi)
	if err != nil {
		fmt.Println("error upserting offset: ", err)
	}
}

// FindPodcastEpisode takes a *database.Client, podcast and episode ID
func FindPodcastEpisode(dbClient *database.Client, podID, epiID string) (*models.Podcast, *models.Episode, error) {
	pID, _ := primitive.ObjectIDFromHex(podID)
	eID, _ := primitive.ObjectIDFromHex(epiID)

	var pod models.Podcast
	err := dbClient.FindByID(database.ColPodcast, pID, &pod)
	if err != nil {
		return nil, nil, err
	}

	for i, _ := range pod.Episodes {
		if pod.Episodes[i].ID == eID {
			return &pod, &pod.Episodes[i], nil
		}
	}

	return nil, nil, errors.New("podcast episode not found")
}

// FindLength
func FindLength(epi *models.Episode) *time.Duration {

}
