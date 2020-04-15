package podcast

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"sort"
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
	err := dbClient.FindWithBSON(database.ColUserEpisode, filter, &userEpi, false)
	if err != nil {
		fmt.Println("error finding user episode details: ", err)
		return 0
	}
	return userEpi.Offset
}

// UpdateOffset takes userID epiID and offset and performs upsert to the UserEpisode collection
func UpdateOffset(dbClient *database.Client, uID, pID, eID primitive.ObjectID, offset int64) {
	userEpi := &models.UserEpisode{UserID: uID, PodcastID: pID, EpisodeID: eID, Offset: offset, Played: false}

	err := dbClient.Upsert(database.ColUserEpisode, bson.D{
		{Key: "user_id", Value: uID},
		{Key: "podcast_id", Value: pID},
		{Key: "episode_id", Value: eID}},
		userEpi)
	if err != nil {
		fmt.Println("error upserting offset: ", err)
	}
}

// FindPodcastEpisode takes a *database.Client, podcast and episode ID
func FindPodcastEpisode(dbClient *database.Client, podID, epiID primitive.ObjectID) (*models.Podcast, *models.Episode, error) {
	var pod models.Podcast
	err := dbClient.FindByID(database.ColPodcast, podID, &pod)
	if err != nil {
		return nil, nil, err
	}

	for i, _ := range pod.Episodes {
		if pod.Episodes[i].ID == epiID {
			return &pod, &pod.Episodes[i], nil
		}
	}

	return nil, nil, errors.New("podcast episode not found")
}

// FindUserLastPlayed takes dbClient, userID, returns the latest played episode and offset
func FindUserLastPlayed(dbClient *database.Client, userID primitive.ObjectID) (*models.Podcast, *models.Episode, int64, error) {
	var userEps []models.UserEpisode
	// look up all the UserEpisodes with the user id
	err := dbClient.Find(database.ColUserEpisode, "user_id", userID, &userEps, true)
	if err != nil {
		fmt.Println("find user last played error: ", err)
		return nil, nil, 0, err
	}
	if len(userEps) == 0 {
		return nil, nil, 0, errors.New("User has no currently played episodes")
	}

	// sort array
	sort.Slice(userEps, func(i, j int) bool {
		return userEps[i].LastSeen.After(userEps[j].LastSeen)
	})

	// gather ids and perform db lookups
	podID := userEps[0].PodcastID
	epiID := userEps[0].EpisodeID

	// find the episode
	pod, epi, err := FindPodcastEpisode(dbClient, podID, epiID)
	if err != nil {
		return nil, nil, 0, err
	}

	return pod, epi, userEps[0].Offset, err
}

// FindLength
func FindLength(epi *models.Episode) *time.Duration {
	zero := time.Duration(0)
	return &zero
}
