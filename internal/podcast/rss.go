package podcast

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/sschwartz96/syncapod/internal/protos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdatePodcasts attempts to go through the list of podcasts update them via RSS feed
func UpdatePodcasts(dbClient *database.MongoClient) {
	for {
		var podcasts []protos.Podcast
		// TODO: use mongo "skip" and "limit" to access only a few podcasts say 100 at a time
		err := dbClient.FindAll(database.ColPodcast, &podcasts)
		if err != nil {
			fmt.Println("error getting all podcasts: ", err)
		}

		var wg sync.WaitGroup

		for i := range podcasts {
			pod := &podcasts[i]
			wg.Add(1)
			go UpdatePodcast(&wg, dbClient, pod)
		}

		wg.Wait()
		fmt.Println("finished updating podcast: waiting...")
		time.Sleep(time.Minute * 15)
	}
}

// UpdatePodcast updates the given podcast via RSS feed
func UpdatePodcast(wg *sync.WaitGroup, dbClient *database.MongoClient, pod *protos.Podcast) {
	defer wg.Done()
	newPod, err := ParseRSS(pod.Rss)
	if err != nil {
		fmt.Println("failed to load podcast rss: ", err)
		return
	}

	for e := range newPod.RSSEpisodes {
		epi := convertEpisode(pod.Id, &newPod.RSSEpisodes[e])
		// TODO: maybe check if the episode has the same title but different size
		// TODO: hopefully the podcast just uses the same URL if they update it
		// check if the latest episode is in collection
		filter := bson.D{
			{Key: "title", Value: epi.Title},
			{Key: "pubdate", Value: epi.PubDate},
		}
		exists, err := dbClient.Exists(database.ColEpisode, filter)
		if err != nil {
			fmt.Println("couldn't tell if object exists: ", err)
			continue
		}

		// episode does not exist
		if !exists {
			err = dbClient.Insert(database.ColEpisode, &epi)
			if err != nil {
				fmt.Println("couldn't insert episode: ", err)
			}
		} else {
			// assume that if the first podcast exists so do the rest, no need to loop through all
			break
		}
	}
}

// AddNewPodcast takes RSS url and downloads contents inserts the podcast and its episodes into the db
// returns error if podcast already exists or connection error
func AddNewPodcast(dbClient *database.MongoClient, url string) error {
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
	rssPod, err := ParseRSS(url)
	if err != nil {
		return err
	}
	pod := convertPodcast(url, rssPod)

	rssEpisodes := rssPod.RSSEpisodes

	// loop through episodes and save them
	for i := range rssEpisodes {
		rssEpi := rssEpisodes[i]
		rssEpi.ID = primitive.NewObjectID()
		rssEpi.PodcastID = rssPod.ID
		if rssEpi.Author == "" {
			rssEpi.Author = pod.Author
		}

		epi := convertEpisode(pod.Id, &rssEpi)

		err = dbClient.Insert(database.ColEpisode, epi)
		if err != nil {
			fmt.Println("couldn't insert episode: ", err)
		}
	}

	// Set episodes to nil and save podcast info to collection
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
func convertEpisode(pID *protos.ObjectID, e *models.RSSEpisode) *protos.Episode {
	pubDate, err := parseRFC2822(e.PubDate)
	if err != nil {
		fmt.Println("error converting episode: ", err)
	}
	// no error since we are checking for one above
	pubTimestamp, _ := ptypes.TimestampProto(*pubDate)

	image := &protos.Image{Title: "", Url: e.Image.HREF}

	return &protos.Episode{
		Id:             protos.NewObjectID(),
		PodcastID:      pID,
		Title:          e.Title,
		Description:    e.Description,
		Subtitle:       e.Subtitle,
		Author:         e.Author,
		Type:           e.Type,
		Image:          image,
		PubDate:        pubTimestamp,
		Summary:        e.Summary,
		Season:         int32(e.Season),
		Episode:        int32(e.Episode),
		Category:       convertCategories(e.Category),
		Explicit:       e.Explicit,
		MP3URL:         e.Enclosure.MP3,
		DurationMillis: parseDuration(e.Duration),
	}
}

// convertPodcast
func convertPodcast(url string, p *models.RSSPodcast) *protos.Podcast {

	keywords := strings.Split(p.Keywords, ",")
	for w := range keywords {
		keywords[w] = strings.TrimSpace(keywords[w])
	}

	lBuildDate, err := parseRFC2822(p.LastBuildDate)
	if err != nil {
		fmt.Println("couldn't parse podcast build date: ", err)
	}
	buildTimestamp, _ := ptypes.TimestampProto(*lBuildDate)

	pubDate, err := parseRFC2822(p.PubDate)
	if err != nil {
		fmt.Println("couldn't parse podcast pubdate: ", err)
	}
	pubTimestamp, _ := ptypes.TimestampProto(*pubDate)

	return &protos.Podcast{
		Id:            protos.NewObjectID(),
		Author:        p.Author,
		Category:      convertCategories(p.Category),
		Explicit:      p.Explicit,
		Image:         &protos.Image{Title: p.Image.Title, Url: p.Image.URL},
		Keywords:      keywords,
		Language:      p.Language,
		LastBuildDate: buildTimestamp,
		PubDate:       pubTimestamp,
		Link:          p.Link,
		Rss:           url,
		Subtitle:      p.Subtitle,
		Title:         p.Title,
		Type:          p.Type,
	}
}

// parseRFC2822 parses the string in RFC2822 date format
// returns pointer to time object and error
func parseRFC2822(s string) (*time.Time, error) {
	var rfc2822 string
	if strings.Contains(s, "+") || strings.Contains(s, "-") {
		rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"
	} else {
		rfc2822 = "Mon, 02 Jan 2006 15:04:05 MST"
	}
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

func convertCategories(cats []models.Category) []*protos.Category {
	protoCats := make([]*protos.Category, len(cats))
	for i := range cats {
		protoCats[i] = convertCategory(cats[i])
	}
	return protoCats
}

func convertCategory(cat models.Category) *protos.Category {
	newCat := &protos.Category{
		Text:     cat.Text,
		Category: convertCategories(cat.Category),
	}
	return newCat
}
