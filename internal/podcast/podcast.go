package podcast

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/closestmatch"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/tcolgate/mp3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdatePodcasts attempts to go through the list of podcasts and add
// episodes to the collection
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

// UpdatePodcast updates the given podcast
func UpdatePodcast(dbClient *database.Client, pod *models.Podcast) {
	newPod, err := ParseRSS(pod.RSS)
	if err != nil {
		fmt.Println("failed to load podcast rss: ", err)
		return
	}

	for e := range newPod.Episodes {
		epi := newPod.Episodes[e]
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

		// episode exists
		if exists {
			fmt.Println("episode already exists")
			break
		} else {
			epi.PodcastID = pod.ID
			err = dbClient.Insert(database.ColEpisode, &epi)
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

func parseDuration(d string) int64 {
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

// AddEpiIDs adds missing IDs to the podcast object and episode objects
// func AddEpiIDs(podcast *models.RSSPodcast) {
// 	podcast.ID = primitive.NewObjectID()

// 	for i := range podcast.Episodes {
// 		podcast.Episodes[i].ID = primitive.NewObjectID()
// 	}
// }

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

// FindPodcast takes a *database.Client and podcast ID
func FindPodcast(dbClient *database.Client, podID primitive.ObjectID) (*models.Podcast, error) {
	var pod models.Podcast
	err := dbClient.FindByID(database.ColPodcast, podID, &pod)
	return &pod, err
}

// FindEpisode takes a *database.Client and episode ID
func FindEpisode(dbClient *database.Client, epiID primitive.ObjectID) (*models.Episode, error) {
	var epi models.Episode
	err := dbClient.FindByID(database.ColEpisode, epiID, &epi)
	return &epi, err
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

	// find the podcast
	pod, err := FindPodcast(dbClient, podID)
	if err != nil {
		return nil, nil, 0, err
	}

	// find the episode
	epi, err := FindEpisode(dbClient, epiID)
	if err != nil {
		return nil, nil, 0, err
	}

	return pod, epi, userEps[0].Offset, err
}

// UpdateEpisode takes a pointer to db, podcast, and episode.
// Attempts to update the episode in the db returning error if not
func UpdateEpisode(dbClient *database.Client, pod *models.Podcast, epi *models.Episode) error {
	col := dbClient.Database(database.DBsyncapod).Collection(database.ColPodcast)

	filter := bson.D{
		{Key: "_id", Value: pod.ID},
		{Key: "episodes._id", Value: epi.ID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.M{"episodes.$": epi}},
	}

	res, err := col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	fmt.Println("update result: ", res.ModifiedCount)
	return nil
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
