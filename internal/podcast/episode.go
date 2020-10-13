package podcast

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
)

// FindEpisodes returns a list of episodes based on podcast id
// returns in chronological order, sectioned by start & end(exclusive)
func FindEpisodesByRange(dbClient db.Database, podcastID *protos.ObjectID, start int64, end int64) ([]*protos.Episode, error) {
	var episodes []*protos.Episode
	filter := &db.Filter{"podcastid": podcastID}
	opts := db.CreateOptions().SetLimit(end-start).SetSkip(start).SetSort("pubdate", -1)
	err := dbClient.FindAll(database.ColEpisode, &episodes, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("FindEpisodesByRange() error finding episodes by range %d - %d: %v", start, end, err)
	}
	return episodes, nil
}

func FindLatestEpisode(dbClient db.Database, podcastID *protos.ObjectID) (*protos.Episode, error) {
	var episode protos.Episode
	filter := &db.Filter{"podcastid": podcastID}
	opts := db.CreateOptions().SetSort("pubdate", -1)
	err := dbClient.FindOne(database.ColEpisode, &episode, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding latest episode: %v", err)
	}
	return &episode, nil
}

func FindEpisodeByID(dbClient db.Database, id *protos.ObjectID) (*protos.Episode, error) {
	var episode protos.Episode
	err := dbClient.FindOne(database.ColEpisode, &episode, &db.Filter{"_id": id}, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding episode by id: %v", err)
	}
	return &episode, nil
}

// FindEpisodeBySeason takes a season episode number returns error if not found
func FindEpisodeBySeason(dbClient db.Database, podID *protos.ObjectID, seasonNum int, episodeNum int) (*protos.Episode, error) {
	var episode protos.Episode

	filter := &db.Filter{
		"podcast_id": podID,
		"season":     seasonNum,
		"episode":    episodeNum,
	}
	err := dbClient.FindOne(database.ColEpisode, &episode, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding episode by season/episode #: %v", err)
	}

	return &episode, nil
}

func UpsertEpisode(dbClient db.Database, episode *protos.Episode) error {
	err := dbClient.Upsert(database.ColEpisode, episode, &db.Filter{"_id": episode.Id})
	if err != nil {
		return fmt.Errorf("error upserting episode: %v", err)
	}
	return nil
}

// helpers
func DoesEpisodeExist(dbClient db.Database, title string, pubDate *timestamp.Timestamp) (bool, error) {
	filter := &db.Filter{
		"title":   title,
		"pubdate": pubDate,
	}
	var episode protos.Episode
	err := dbClient.FindOne(database.ColEpisode, &episode, filter, nil)
	if err != nil && !strings.Contains(err.Error(), "no documents") {
		return false, fmt.Errorf("DoesEpisodeExist() error: %v", err)
	}
	if episode.Id == nil {
		return false, nil
	}
	return true, nil
}
