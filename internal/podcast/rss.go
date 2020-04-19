package podcast

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdatePodcasts attempts to go through the list of podcasts update them via RSS feed
func UpdatePodcasts(dbClient *database.Client) {
	for {
		var podcasts []models.Podcast
		// TODO: use mongo "skip" and "limit" to access only a few podcasts say 100 at a time
		err := dbClient.FindAll(database.ColPodcast, &podcasts)
		if err != nil {
			fmt.Println("error getting all podcasts: ", err)
		}

		for _, pod := range podcasts {
			go UpdatePodcast(dbClient, &pod)
		}
		time.Sleep(time.Minute * 15)
	}
}

// UpdatePodcast updates the given podcast via RSS feed
func UpdatePodcast(dbClient *database.Client, pod *models.Podcast) {
	newPod, err := ParseRSS(pod.RSS)
	if err != nil {
		fmt.Println("failed to load podcast rss: ", err)
		return
	}

	for e := range newPod.RSSEpisodes {
		epi := &newPod.RSSEpisodes[e]
		// TODO: maybe check if the episode has the same title but different size
		// TODO: hopefully the podcast just uses the same URL if they update it
		// check if the latest episode is in collection
		filter := bson.D{
			{Key: "title", Value: epi.Title},
			{Key: "pub_date", Value: epi.PubDate},
		}
		exists, err := dbClient.Exists(database.ColEpisode, filter)
		if err != nil {
			fmt.Println("couldn't tell if object exists: ", err)
			continue
		}

		// episode does exists
		if !exists {
			epi.PodcastID = pod.ID
			mongoEpi := convertEpisode(epi)
			err = dbClient.Insert(database.ColEpisode, &mongoEpi)
			if err != nil {
				fmt.Println("couldn't insert episode: ", err)
			}
		}
	}
}

// AddNewPodcast takes RSS url and downloads contents inserts the podcast and its episodes into the db
// returns error if podcast already exists or connection error
func AddNewPodcast(dbClient *database.Client, url string) error {
	// check if podcast already contains that rss url
	filter := bson.D{{Key: "rss", Value: url}}
	exists, err := dbClient.Exists(database.ColPodcast, filter)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("podcast already exists")
	}

	// attempt to download & parse the podcast rss
	pod, err := ParseRSS(url)
	if err != nil {
		return err
	}
	pod.ID = primitive.NewObjectID()
	pod.RSS = url

	// loop through episodes and save them
	for i := range pod.Episodes {
		epi := pod.Episodes[i]
		epi.ID = primitive.NewObjectID()
		epi.PodcastID = pod.ID

		err = dbClient.Insert(database.ColEpisode, &epi)
		if err != nil {
			fmt.Println("couldn't insert episode: ", err)
		}
	}

	// Set episodes to nil and save podcast info to collection
	pod.Episodes = nil
	err = dbClient.Insert(database.ColPodcast, pod)
	if err != nil {
		fmt.Println("couldn't insert podcast: ", err)
		return err
	}

	return nil
}

// ParseRSS takes in URL path and unmarshals the data
func ParseRSS(path string) (*models.RSSPodcast, error) {
	// make the connection
	response, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	// set up rss feed object and decoder
	var rss models.RSSFeed
	decoder := xml.NewDecoder(response.Body)
	decoder.DefaultSpace = "Default"

	// decode
	err = decoder.Decode(&rss)
	if err != nil {
		return nil, err
	}

	return &rss.RSSPodcast, nil
}

// convertEpisode takes in id of parent podcast and RSSEpisode
// and returns a pointer to Episode
func convertEpisode(pID primitive.ObjectID, e *models.RSSEpisode) *models.Episode {
	pubDate, err := parseRFC2822(e.PubDate)
	if err != nil {
		fmt.Println("error converting episode: ", err)
	}

	var explicit bool
	if strings.Contains(e.Explicit, "yes") || strings.Contains(e.Explicit, "explicit") {
		explicit = false
	}

	epi := &models.Episode{
		ID:             primitive.NewObjectID(),
		PodcastID:      pID,
		Title:          e.Title,
		Description:    e.Description,
		Subtitle:       e.Subtitle,
		Author:         e.Author,
		Type:           e.Type,
		Image:          e.Image,
		PubDate:        *pubDate,
		Summary:        e.Summary,
		Season:         e.Season,
		Episode:        e.Episode,
		Category:       e.Category,
		Explicit:       explicit,
		URL:            e.Enclosure.MP3,
		DurationMillis: parseDuration(e.Duration),
	}

	return epi
}

// parseRFC2822 parses the string in RFC2822 date format
// returns pointer to time object and error
func parseRFC2822(s string) (*time.Time, error) {
	const rfc2822 = "Mon, 02 Jan 15:04:05 2006 MST"
	t, err := time.Parse(rfc2822, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

//parseDuration takes in the string duration and returns the duration in millis
func parseDuration(d string) int64 {
	// check if they just applied the seconds
	if !strings.Contains(d, ":") {
		sec, err := strconv.Atoi(d)
		if err != nil {
			fmt.Println("error converting duration of episode: ", err)
			return 0
		}
		return int64(sec) * int64(1000)
	}
	var millis int64
	multiplier := int64(1000)

	// format hh:mm:ss || mm:ss
	split := strings.Split(d, ":")

	for i := len(split) - 1; i >= 0; i-- {
		v, _ := strconv.Atoi(split[i])
		millis += int64(v) * multiplier
		multiplier *= int64(60)
	}

	return millis
}
