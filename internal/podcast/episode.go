package podcast

import (
	"context"
	"fmt"

	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

// FindEpisode takes a *database.Client and episode ID
func FindEpisode(dbClient *database.Client, epiID primitive.ObjectID) (*models.Episode, error) {
	var epi models.Episode
	err := dbClient.FindByID(database.ColEpisode, epiID, &epi)
	return &epi, err
}

// FindEpisodeByNumber takes a pointer to database.Client, podcast id, episode #
func FindEpisodeByNumber(dbClient *database.Client, podID primitive.ObjectID, num int) (*models.Episode, error) {
	var epi models.Episode
	filter := bson.D{
		{Key: "podcast_id", Value: podID},
		{Key: "episode", Value: num},
	}
	err := dbClient.FindWithBSON(database.ColEpisode, filter, nil, &epi)

	return &epi, err
}

// FindLatestEpisode takes a pointer to database.Client and podcast id
func FindLatestEpisode(dbClient *database.Client, podID primitive.ObjectID) (*models.Episode, error) {
	var epi models.Episode
	opts := options.FindOne().SetSort(bson.M{"pub_date": -1})

	filter := bson.M{"podcast_id": podID}
	err := dbClient.FindWithBSON(database.ColEpisode, filter, opts, &epi)
	return &epi, err
}

// FindAllEpisodesRange finds the lastest episodes within range(epi # 20-30)
// s = start, e = end
func FindAllEpisodesRange(dbClient *database.Client, podID primitive.ObjectID, s, e int) []models.Episode {
	var epis []models.Episode
	filter := bson.M{"podcast_id": podID}
	opts := options.Find().SetLimit(int64(e - s)).SetSkip(int64(s)).SetSort(
		bson.M{"pub_date": -1},
	)
	err := dbClient.FindAllWithBSON(database.ColEpisode, filter, opts, &epis)
	if err != nil {
		fmt.Println("error finding all episodes: ", err)
	}
	return epis
}

// FindAllEpisodes takesa pointer to database.Client and a podcast id
func FindAllEpisodes(dbClient *database.Client, podID primitive.ObjectID) []models.Episode {
	var epis []models.Episode
	filter := bson.M{"podcast_id": podID}
	err := dbClient.FindAllWithBSON(database.ColEpisode, filter, nil, &epis)
	if err != nil {
		fmt.Println("error finding all episodes: ", err)
	}
	return epis
}
