package podcast

import (
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/schollz/closestmatch"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/tcolgate/mp3"
)

func InsertPodcast(db database.Database, podcast *protos.Podcast) error {
	err := db.Insert(database.ColPodcast, podcast)
	if err != nil {
		return fmt.Errorf("error inserting podcast: %v", err)
	}
	return nil
}

func FindAllPodcasts(db database.Database) ([]*protos.Podcast, error) {
	// TODO: get rid of?
	var podcasts []*protos.Podcast
	err := db.FindAll(database.ColPodcast, &podcasts, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding all podcasts: %v", err)
	}
	return podcasts, nil
}

func DoesPodcastExist(db database.Database, rssURL string) (bool, error) {
	var podcast *protos.Podcast
	filter := &database.Filter{"rss": rssURL}
	err := db.FindOne(database.ColPodcast, podcast, filter, nil)
	if err != nil {
		return false, fmt.Errorf("error does podcast exist: %v", err)
	}
	if podcast != nil {
		return false, nil
	}
	return true, nil
}

func FindPodcastsByRange(db database.Database, start, end int) ([]*protos.Podcast, error) {
	var podcasts []*protos.Podcast
	opts := database.CreateOptions().SetLimit(int64(end-start)).SetSkip(int64(start)).SetSort("pubdate", -1)

	err := db.FindAll(database.ColEpisode, &podcasts, nil, opts)
	if err != nil {
		return podcasts, fmt.Errorf("error finding podcasts within range %d - %d: %v", start, end, err)
	}
	return podcasts, nil
}

func FindPodcastByID(db database.Database, id *protos.ObjectID) (*protos.Podcast, error) {
	var podcast *protos.Podcast
	if err := db.FindOne(database.ColPodcast, podcast, &database.Filter{"_id": id}, nil); err != nil {
		return nil, fmt.Errorf("error finding podcast by id: %v", err)
	}
	return podcast, nil
}

// FindUserEpisode takes pointer to database client, userID, epiID
// returns *protos.UserEpisode
func FindUserEpisode(db database.Database, userID, epiID *protos.ObjectID) (*protos.UserEpisode, error) {
	var userEpi protos.UserEpisode
	filter := &database.Filter{
		"userid":    userID,
		"episodeid": epiID,
	}
	err := db.FindOne(database.ColUserEpisode, &userEpi, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding user episodes details, %v", err)
	}
	return &userEpi, nil
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
