package podcast

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/schollz/closestmatch"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/tcolgate/mp3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddEpiIDs adds missing IDs to the podcast object and episode objects
// func AddEpiIDs(podcast *protos.RSSPodcast) {
// 	podcast.ID = *protos.NewObjectID()

// 	for i := range podcast.Episodes {
// 		podcast.Episodes[i].ID = *protos.NewObjectID()
// 	}
// }

// FindUserEpisode takes pointer to database client, userID, epiID
// returns *protos.UserEpisode
func FindUserEpisode(dbClient *database.Client, userID, epiID *protos.ObjectID) (*protos.UserEpisode, error) {
	var userEpi protos.UserEpisode
	filter := bson.D{{Key: "userid", Value: userID}, {Key: "episodeid", Value: epiID}}
	err := dbClient.FindWithBSON(database.ColUserEpisode, filter, nil, &userEpi)
	if err != nil {
		return nil, fmt.Errorf("error finding user episodes details, %v", err)
	}
	return &userEpi, nil
}

// FindOffset takes database client and pointers to user and episode to lookup episode details and offset
func FindOffset(dbClient *database.Client, userID, epiID *protos.ObjectID) int64 {
	userEpi, err := FindUserEpisode(dbClient, userID, epiID)
	if err != nil {
		fmt.Println("error finding offset: ", err)
		return 0
	}
	return userEpi.Offset
}

// UpdateOffset takes userID epiID and offset and performs upsert to the UserEpisode collection
func UpdateOffset(dbClient *database.Client, uID, pID, eID *protos.ObjectID, offset int64) error {
	userEpi := &protos.UserEpisode{
		UserID:    uID,
		PodcastID: pID,
		EpisodeID: eID,
		Offset:    offset,
		Played:    false,
		LastSeen:  ptypes.TimestampNow(),
	}

	err := dbClient.Upsert(database.ColUserEpisode, bson.D{
		{Key: "userid", Value: uID},
		{Key: "podcastid", Value: pID},
		{Key: "episodeid", Value: eID}},
		userEpi)
	if err != nil {
		fmt.Println("error upserting offset: ", err)
		return err
	}
	return nil
}

// FindPodcast takes a *database.Client and podcast ID
func FindPodcast(dbClient *database.Client, podID *protos.ObjectID) (*protos.Podcast, error) {
	var pod protos.Podcast
	err := dbClient.FindByID(database.ColPodcast, podID, &pod)
	return &pod, err
}

// FindUserLastPlayed takes dbClient, userID, returns the latest played episode and offset
func FindUserLastPlayed(dbClient *database.Client, userID *protos.ObjectID) (*protos.Podcast, *protos.Episode, int64, error) {
	var userEp protos.UserEpisode
	var pod protos.Podcast
	var epi protos.Episode

	// find the latest played user_episode
	filter := bson.M{"userid": userID}
	opts := options.FindOne().SetSort(bson.M{"lastseen": -1})

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

// GetSubscriptions returns a list of subscriptions via userID
func GetSubscriptions(dbClient *database.Client, userID *protos.ObjectID) ([]*protos.Subscription, error) {
	var subs []*protos.Subscription
	err := dbClient.FindAllWithBSON(database.ColSubscription, bson.M{"userid": userID}, nil, &subs)
	return subs, err
}

// UpdateUserEpiOffset changes the offset in the collection
func UpdateUserEpiOffset(dbClient *database.Client, userID, epiID *protos.ObjectID, offset int64) error {
	return UpdateUserEpiParam(dbClient, userID, epiID, "offset", offset)
}

// UpdateUserEpiPlayed marks the episode as played in db
func UpdateUserEpiPlayed(dbClient *database.Client, userID, epiID *protos.ObjectID, played bool) error {
	return UpdateUserEpiParam(dbClient, userID, epiID, "played", played)
}

// UpdateUserEpiParam updates the user's episode data based on param and data
func UpdateUserEpiParam(dbClient *database.Client, userID, epiID *protos.ObjectID, param string, data interface{}) error {
	filter := bson.D{
		{Key: "userid", Value: userID},
		{Key: "episodeid", Value: epiID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.M{param: data}},
		{Key: "$set", Value: bson.M{"lastseen": time.Now()}},
	}

	return dbClient.Upsert(database.ColUserEpisode, filter, update)
}

// UpdateUserEpi updates an entire UserEpisode
func UpdateUserEpi(dbClient *database.Client, userEpi *protos.UserEpisode) error {
	filter := bson.D{
		{Key: "userid", Value: userEpi.UserID},
		{Key: "episodeid", Value: userEpi.EpisodeID},
	}

	// update := bson.D{
	// 	{Key: "$set", Value: bson.D{
	// 		{Key: "offset", Value: userEpi.Offset},
	// 		{Key: "played", Value: userEpi.Played},
	// 	}},
	// }

	return dbClient.Upsert(database.ColUserEpisode, filter, userEpi)
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
