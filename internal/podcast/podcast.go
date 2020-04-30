package podcast

import (
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/schollz/closestmatch"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/tcolgate/mp3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddEpiIDs adds missing IDs to the podcast object and episode objects
// func AddEpiIDs(podcast *models.RSSPodcast) {
// 	podcast.ID = primitive.NewObjectID()

// 	for i := range podcast.Episodes {
// 		podcast.Episodes[i].ID = primitive.NewObjectID()
// 	}
// }

// FindOffset takes database client and pointers to user and episode to lookup episode details and offset
func FindOffset(dbClient *database.Client, userID, epiID primitive.ObjectID) int64 {
	var userEpi models.UserEpisode
	filter := bson.D{{Key: "user_id", Value: userID}, {Key: "episode_id", Value: epiID}}
	err := dbClient.FindWithBSON(database.ColUserEpisode, filter, nil, &userEpi)
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

// FindPodcast takes a *database.Client and podcast ID
func FindPodcast(dbClient *database.Client, podID primitive.ObjectID) (*models.Podcast, error) {
	var pod models.Podcast
	err := dbClient.FindByID(database.ColPodcast, podID, &pod)
	return &pod, err
}

// FindUserLastPlayed takes dbClient, userID, returns the latest played episode and offset
func FindUserLastPlayed(dbClient *database.Client, userID primitive.ObjectID) (*models.Podcast, *models.Episode, int64, error) {
	var userEp models.UserEpisode
	var pod models.Podcast
	var epi models.Episode

	// find the latest played user_episode
	filter := bson.M{"user_id": userID}
	opts := options.FindOne().SetSort(bson.M{"last_seen": -1})

	err := dbClient.FindWithBSON(database.ColUserEpisode, filter, opts, &userEp)
	if err != nil {
		fmt.Println("error finding user_episod: ", err)
		return nil, nil, 0, err
	}

	// find podcast
	err = dbClient.FindByID(database.ColPodcast, userEp.PodcastID, &pod)
	if err != nil {
		fmt.Println("couldn't find podcast: ", err)
		return nil, nil, 0, err
	}

	// find episode
	err = dbClient.FindByID(database.ColEpisode, userEp.EpisodeID, &epi)
	if err != nil {
		fmt.Println("couldn't find podcast: ", err)
		return nil, nil, 0, err
	}

	return &pod, &epi, userEp.Offset, err
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

// UpdateUserEpiOffset changes the offset in the collection
func UpdateUserEpiOffset(dbClient *database.Client, userID, epiID primitive.ObjectID, offset int64) error {
	return UpdateUserEpi(dbClient, userID, epiID, "offset", offset)
}

// UpdateUserEpiPlayed marks the episode as played in db
func UpdateUserEpiPlayed(dbClient *database.Client, userID, epiID primitive.ObjectID, played bool) error {
	return UpdateUserEpi(dbClient, userID, epiID, "played", played)
}

// UpdateUserEpi updates the user's episode data based on param and data
func UpdateUserEpi(dbClient *database.Client, userID, epiID primitive.ObjectID, param string, data interface{}) error {
	filter := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "episode_id", Value: epiID},
	}

	update := bson.M{
		"$set": bson.M{param: data},
	}

	err := dbClient.UpdateWithBSON(database.ColUserEpisode, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// TODO:

// MatchTitle is a helper function to match search with a list of podcasts titles
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
