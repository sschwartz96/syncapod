package podcast

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
)

// FindEpisodes returns a list of episodes based on podcast id
// returns in chronological order, sectioned by start & end
func FindEpisodesByRange(db database.Database, podcastID *protos.ObjectID, start int64, end int64) ([]*protos.Episode, error) {
	var episodes []*protos.Episode
	filter := &database.Filter{"podcastid": podcastID}
	opts := database.CreateOptions().SetLimit(end-start).SetSkip(start).SetSort("pubdate", -1)
	err := db.FindAll(database.ColEpisode, &episodes, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding episodes by range %d - %d: %v", start, end, err)
	}
	return episodes, nil
}

func FindAllEpisodes(db database.Database, podcastID *protos.ObjectID) ([]*protos.Episode, error) {
	var episodes []*protos.Episode
	filter := &database.Filter{"podcastid": podcastID}
	opts := database.CreateOptions().SetSort("pubdate", -1)
	err := db.FindAll(database.ColEpisode, &episodes, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding all episodes: %v", err)
	}
	return episodes, nil
}

func FindLatestEpisode(db database.Database, podcastID *protos.ObjectID) (*protos.Episode, error) {
	var episode *protos.Episode
	filter := &database.Filter{"podcastid": podcastID}
	opts := database.CreateOptions().SetSort("pubdate", -1)
	err := db.FindOne(database.ColEpisode, &episode, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding latest episode: %v", err)
	}
	return episode, nil
}

func FindEpisodeByID(db database.Database, id *protos.ObjectID) (*protos.Episode, error) {
	var episode *protos.Episode
	err := db.FindOne(database.ColEpisode, &episode, &database.Filter{"_id": id}, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding episode by id: %v", err)
	}
	return episode, nil
}

// FindEpisodeBySeason takes a season episode number returns error if not found
func FindEpisodeBySeason(db database.Database, id *protos.ObjectID, seasonNum int, episodeNum int) (*protos.Episode, error) {
	var episode protos.Episode

	filter := &database.Filter{
		"podcast_id": id,
		"season":     seasonNum,
		"episode":    episodeNum,
	}
	err := db.FindOne(database.ColEpisode, &episode, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding episode by season/episode #: %v", err)
	}

	return &episode, nil
}

func UpsertEpisode(db database.Database, episode *protos.Episode) error {
	err := db.Upsert(database.ColEpisode, &episode, &database.Filter{"_id": episode.Id})
	if err != nil {
		return fmt.Errorf("error upserting episode: %v", err)
	}
	return nil
}

// helpers
func DoesEpisodeExist(db database.Database, title string, pubDate *timestamp.Timestamp) (bool, error) {
	filter := &database.Filter{
		"title":   title,
		"pubdate": pubDate,
	}
	var episode *protos.Episode
	err := db.FindOne(database.ColUserEpisode, &episode, filter, nil)
	if err != nil {
		return false, fmt.Errorf("error does episode exist: %v", err)
	}
	if episode == nil {
		return false, nil
	}
	return true, nil
}

//// UpdateEpisode takes a pointer to db, podcast, and episode.
//// Attempts to update the episode in the db returning error if not
//func UpdateEpisode(dbClient *database.MongoClient, pod *protos.Podcast, epi *protos.Episode) error {
//	col := dbClient.Database(database.DBsyncapod).Collection(database.ColPodcast)
//
//	filter := bson.D{
//		{Key: "_id", Value: pod.Id},
//		{Key: "episodes._id", Value: epi.Id},
//	}
//
//	update := bson.D{
//		{Key: "$set", Value: bson.M{"episodes.$": epi}},
//	}
//
//	res, err := col.UpdateOne(context.Background(), filter, update)
//	if err != nil {
//		return err
//	}
//	fmt.Println("update result: ", res.ModifiedCount)
//	return nil
//}
//
//// FindEpisode takes a *database.Client and episode ID
//func FindEpisode(dbClient *database.MongoClient, epiID *protos.ObjectID) (*protos.Episode, error) {
//	var epi protos.Episode
//	err := dbClient.FindByID(database.ColEpisode, epiID, &epi)
//	return &epi, err
//}
//
//// FindEpisodeByNumber takes a pointer to database.Client, podcast id, episode #
//func FindEpisodeByNumber(dbClient *database.MongoClient, podID *protos.ObjectID, num int) (*protos.Episode, error) {
//	var epi protos.Episode
//	filter := bson.D{
//		{Key: "podcast_id", Value: podID},
//		{Key: "episode", Value: num},
//	}
//	err := dbClient.FindWithBSON(database.ColEpisode, filter, nil, &epi)
//
//	return &epi, err
//}
//
//// FindLatestEpisode takes a pointer to database.Client and podcast id
//func FindLatestEpisode(dbClient *database.MongoClient, podID *protos.ObjectID) (*protos.Episode, error) {
//	var epi protos.Episode
//	opts := options.FindOne().SetSort(bson.M{"pub_date": -1})
//
//	filter := bson.M{"podcast_id": podID}
//	err := dbClient.FindWithBSON(database.ColEpisode, filter, opts, &epi)
//	return &epi, err
//}
//
//// FindAllEpisodesRange finds the lastest episodes within range(epi # 20-30)
//// s = start, e = end
//func FindAllEpisodesRange(dbClient *database.MongoClient, podID *protos.ObjectID, s, e int) []*protos.Episode {
//	var epis []*protos.Episode
//	filter := bson.M{"podcastid": podID}
//	opts := options.Find().SetLimit(int64(e - s)).SetSkip(int64(s)).SetSort(
//		bson.M{"pubdate": -1},
//	)
//	err := dbClient.FindAllWithBSON(database.ColEpisode, filter, opts, &epis)
//	if err != nil {
//		fmt.Println("error finding all episodes: ", err)
//	}
//	return epis
//}
//
//// FindAllEpisodes takesa pointer to database.Client and a podcast id
//func FindAllEpisodes(dbClient *database.MongoClient, podID *protos.ObjectID) []*protos.Episode {
//	var epis []*protos.Episode
//	filter := bson.M{"podcastid": podID}
//	opts := options.Find().SetSort(bson.M{"pubdate": -1})
//	err := dbClient.FindAllWithBSON(database.ColEpisode, filter, opts, &epis)
//	if err != nil {
//		fmt.Println("error finding all episodes: ", err)
//	}
//	return epis
//}
