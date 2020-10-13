package podcast

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/sschwartz96/syncapod/internal/protos"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdatePodcasts attempts to go through the list of podcasts update them via RSS feed
func UpdatePodcasts(dbClient db.Database) error {
	var podcasts []*protos.Podcast
	var err error
	start, end := 0, 10
	for podcasts, err = FindPodcastsByRange(dbClient, start, end); err != nil && len(podcasts) > 0; {
		var wg sync.WaitGroup
		for i := range podcasts {
			pod := podcasts[i]
			wg.Add(1)
			go func() {
				err = updatePodcast(&wg, dbClient, pod)
				if err != nil {
					fmt.Println("UpdatePodcasts() error updating podcast %v, error = %v", pod, err)
				}
			}()
		}
		wg.Wait()
		start = end
		end += 10
	}
	if err != nil {
		return fmt.Errorf("UpdatePodcasts() error retrieving from db: %v", err)
	}
	return nil
}

// updatePodcast updates the given podcast via RSS feed
func updatePodcast(wg *sync.WaitGroup, dbClient db.Database, pod *protos.Podcast) error {
	defer wg.Done()
	newPod, err := ParseRSS(pod.Rss)
	if err != nil {
		fmt.Println("updatePodcast() failed to load podcast rss: ", err)
		return fmt.Errorf("updatePodcast() error parsing RSS: %v", err)
	}

	for e := range newPod.RSSEpisodes {
		epi := convertEpisode(pod.Id, &newPod.RSSEpisodes[e])
		// check if the latest episode is in collection
		exists, err := DoesEpisodeExist(dbClient, epi.Title, epi.PubDate)
		if err != nil {
			fmt.Println("couldn't tell if object exists: ", err)
			continue
		}

		// episode does not exist
		if !exists {
			err = UpsertEpisode(dbClient, epi)
			if err != nil {
				fmt.Println("couldn't insert episode: ", err)
				return fmt.Errorf("updatePodcast() error upserting episode: %v", err)
			}
		} else {
			// assume that if the first podcast exists so do the rest, no need to loop through all
			break
		}
	}
	return nil
}

// AddNewPodcast takes RSS url and downloads contents inserts the podcast and its episodes into the db
// returns error if podcast already exists or connection error
func AddNewPodcast(dbClient db.Database, url string) error {
	// check if podcast already contains that rss url
	exists := DoesPodcastExist(dbClient, url)
	if exists {
		return errors.New("podcast already exists")
	}

	// attempt to download & parse the podcast rss
	rssPod, err := ParseRSS(url)
	if err != nil {
		return err
	}
	pod := convertPodcast(url, rssPod)

	// insert podcast first that way we don't add episodes without podcast
	err = dbClient.Insert(database.ColPodcast, pod)
	if err != nil {
		return fmt.Errorf("error adding new podcast: %v", err)
	}

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

		err = UpsertEpisode(dbClient, epi)
		if err != nil {
			fmt.Println("couldn't insert episode: ", err)
		}
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
	pubDate, err := parseRFC2822ToUTC(e.PubDate)
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

	log.Println("build date:", p.LastBuildDate)
	lBuildDate, err := parseRFC2822ToUTC(p.LastBuildDate)
	if err != nil {
		fmt.Println("couldn't parse podcast build date: ", err)
	}
	buildTimestamp, _ := ptypes.TimestampProto(*lBuildDate)

	pubDate, err := parseRFC2822ToUTC(p.PubDate)
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

func timeZoneAbrToOffset(abr string) string {
	abrMap := map[string]string{"PST": ""}
}

// parseRFC2822ToUTC parses the string in RFC2822 date format
// returns pointer to time object and error
func parseRFC2822ToUTC(s string) (*time.Time, error) {
	var rfc2822 string
	if strings.Contains(s, "+") || strings.Contains(s, "-") {
		rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"
	} else {
		rfc2822 = "Mon, 02 Jan 2006 15:04:05 MST"
	}
	t, err := time.Parse(rfc2822, s)
	if err != nil {
		return &t, err
	}
	log.Println("time location:", t.Location())
	log.Println("time:", t.String())
	log.Println("time(UTC):", t.UTC().String())
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
