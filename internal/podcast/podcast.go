package podcast

import (
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/schollz/closestmatch"
	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/tcolgate/mp3"
)

func DoesPodcastExist(dbClient db.Database, rssURL string) bool {
	var podcast protos.Podcast
	filter := &db.Filter{"rss": rssURL}
	err := dbClient.FindOne(database.ColPodcast, &podcast, filter, nil)
	if err != nil {
		return false
	}
	return true
}

func FindPodcastsByRange(dbClient db.Database, start, end int) ([]*protos.Podcast, error) {
	var podcasts []*protos.Podcast
	opts := db.CreateOptions().SetLimit(int64(end-start)).SetSkip(int64(start)).SetSort("pubdate", -1)

	err := dbClient.FindAll(database.ColPodcast, &podcasts, nil, opts)
	if err != nil {
		return podcasts, fmt.Errorf("error finding podcasts within range %d - %d: %v", start, end, err)
	}
	return podcasts, nil
}

func FindPodcastByID(dbClient db.Database, id *protos.ObjectID) (*protos.Podcast, error) {
	var podcast protos.Podcast
	if err := dbClient.FindOne(database.ColPodcast, &podcast, &db.Filter{"_id": id}, nil); err != nil {
		return nil, fmt.Errorf("error finding podcast by id: %v", err)
	}
	return &podcast, nil
}

// FindUserEpisode takes pointer to database client, userID, epiID
// returns *protos.UserEpisode
func FindUserEpisode(dbClient db.Database, userID, epiID *protos.ObjectID) (*protos.UserEpisode, error) {
	var userEpi protos.UserEpisode
	filter := &db.Filter{
		"userid":    userID,
		"episodeid": epiID,
	}
	err := dbClient.FindOne(database.ColUserEpisode, &userEpi, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding user episodes details, %v", err)
	}
	return &userEpi, nil
}

// SearchPodcasts searches for a podcast given db and text string
func SearchPodcasts(dbClient db.Database, search string) ([]*protos.Podcast, error) {
	var results []*protos.Podcast
	fields := []string{"title", "keywords", "subtitle"}
	err := dbClient.Search(database.ColPodcast, search, fields, &results)
	if err != nil {
		return nil, fmt.Errorf("error SearchPodcasts: %v", err)
	}
	return results, nil
}

// MatchTitle is a helper function to match search with a list of podcasts titles
func MatchTitle(search string, podcasts []protos.Podcast) {
	var titles []string
	for i := range podcasts {
		titles = append(titles, podcasts[i].Title)
	}

	bagSizes := []int{2, 3, 4}

	cm := closestmatch.New(titles, bagSizes)
	fmt.Println(cm)

	return
}

// FindLength attempts to download only the first few frames of the MP3 to figure out its length
func FindLength(url string) int64 {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	clen := resp.ContentLength

	d := mp3.NewDecoder(resp.Body)
	defer resp.Body.Close()

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
				fmt.Println("I think we are VBR: ", dif)
				vbrFlag = true
				maxFrames = 0
			}
		}

		counter++
		if counter == maxFrames {
			break
		}
	}

	fmt.Println("median: ", (bRateHigh - bRateLow))

	fmt.Println("total frames counter: ", counter)

	bitRate := bRateTTL / int64(counter)
	fmt.Println("bitrate: ", bitRate)
	// Just approximate to 128000 if close enough
	if math.Abs(float64(bitRate)-128000) < 1920 {
		bitRate = 128000
		//fmt.Println("bitrate: ", bitRate)
	}
	//fmt.Println("clen - skip: ", clen-skipTTL)
	guess := ((clen - skipTTL) * 8) / bitRate

	if maxFrames == 0 {
		guess += 20
	}
	return guess * 1000
}
