package podcast

import (
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/tcolgate/mp3"
)

func DoesPodcastExist(dbClient db.Database, rssURL string) bool {
	var podcast protos.Podcast
	filter := &db.Filter{"rss": rssURL}
	err := dbClient.FindOne(database.ColPodcast, &podcast, filter, nil)
	return err == nil
}

func FindPodcastsByRange(dbClient db.Database, start, end int) ([]*protos.Podcast, error) {
	var podcasts []*protos.Podcast
	fmt.Println("skiP:", start)
	opts := db.CreateOptions().SetLimit(int64(end-start)).SetSkip(int64(start)).SetSort("pubdate", -1)

	err := dbClient.FindAll(database.ColPodcast, &podcasts, nil, opts)
	if err != nil {
		return podcasts, fmt.Errorf("error finding podcasts within range %d - %d: %v", start, end, err)
	}
	return podcasts, nil
}

func FindPodcastByID(dbClient db.Database, id *protos.ObjectID) (*protos.Podcast, error) {
	podcast := &protos.Podcast{}
	if err := dbClient.FindOne(database.ColPodcast, podcast, &db.Filter{"_id": id}, nil); err != nil {
		return nil, fmt.Errorf("FindPodcastByID() error: %v", err)
	}
	return podcast, nil
}

func FindPodcastsByIDs(dbClient db.Database, ids []*protos.ObjectID) ([]*protos.Podcast, error) {
	podcasts := []*protos.Podcast{}
	filter := &db.Filter{"_id": db.Filter{"$in": ids}}
	err := dbClient.FindAll(database.ColPodcast, &podcasts, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("FindPodcastsByIDs() error: %v", err)
	}
	return podcasts, nil
}

// SearchPodcasts searches for a podcast given db and text string
func SearchPodcasts(dbClient db.Database, search string) ([]*protos.Podcast, error) {
	var results []*protos.Podcast
	fields := []string{"author", "title", "keywords", "subtitle"}
	err := dbClient.Search(database.ColPodcast, search, fields, &results)
	if err != nil {
		return nil, fmt.Errorf("error SearchPodcasts: %v", err)
	}
	return results, nil
}

// MatchTitle is a helper function to match search with a list of podcasts titles
// func MatchTitle(search string, podcasts []protos.Podcast) {
// 	var titles []string
// 	for i := range podcasts {
// 		titles = append(titles, podcasts[i].Title)
// 	}
// 	bagSizes := []int{2, 3, 4}
// 	cm := closestmatch.New(titles, bagSizes)
// 	fmt.Println(cm)
// 	return
// }

func GetPodcastResp(url string) (io.ReadCloser, int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("GetPodcastResp() error: %v", err)
	}
	return resp.Body, resp.ContentLength, nil
}

// FindLength attempts to download only the first few frames of the MP3 to figure out its length
func FindLength(r io.Reader, fileLength int64) int64 {
	d := mp3.NewDecoder(r)

	var f mp3.Frame
	var skipTTL int64
	skipped := 0
	vbrFlag := false

	bRateTTL := int64(0)
	bRateLow := 0
	bRateHigh := 0

	counter := 0

	maxFrames := 512

	for {
		if skipped > 0 {
			skipTTL += int64(skipped)
		}
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return 0
		}

		bRateTTL += int64(f.Header().BitRate())

		if !vbrFlag {
			if int(f.Header().BitRate()) > bRateHigh && counter > 0 {
				bRateHigh = int(f.Header().BitRate())
			}
			if int(f.Header().BitRate()) < bRateLow || bRateLow == 0 && counter > 0 {
				bRateLow = int(f.Header().BitRate())
			}
			dif := float32(bRateHigh-bRateLow) / float32(bRateHigh)
			if dif < .9 && dif != 0 {
				vbrFlag = true
				maxFrames = 0
			}
		}

		counter++
		if counter == maxFrames {
			break
		}
	}

	bitRate := bRateTTL / int64(counter)
	// Just approximate to 128000 if close enough
	if math.Abs(float64(bitRate)-128000) < 1920 {
		bitRate = 128000
	}
	guess := ((fileLength - skipTTL) * 8) / bitRate

	if maxFrames == 0 {
		guess += 20
	}
	return guess * 1000
}
