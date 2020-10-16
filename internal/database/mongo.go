package database

import (
	"github.com/sschwartz96/stockpile/mongodb"
	"github.com/sschwartz96/syncapod/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// these constants are for database actions
const (
	// Databases
	DBsyncapod = "syncapod"

	// Collections
	ColPodcast      = "podcast"
	ColEpisode      = "episode"
	ColSession      = "session"
	ColUser         = "user"
	ColUserEpisode  = "user_episode"
	ColSubscription = "subscription"
	ColAuthCode     = "oauth_auth_code"
	ColAccessToken  = "oauth_access_token"
)

var (
	collections = []string{
		ColPodcast,
		ColEpisode,
		ColSession,
		ColUser,
		ColUserEpisode,
		ColSubscription,
		ColAuthCode,
		ColAccessToken,
	}
)

// CreateMongoClient makes a connection with the mongo client
func NewMongoClient(cfg *config.Config) (*mongodb.MongoClient, error) {
	opts := options.Client().ApplyURI(cfg.DbURI)
	opts.SetRegistry(createRegistry())
	client, err := mongodb.NewMongoClient(
		DBsyncapod,
		collections,
		opts,
		map[string](map[string]bool){
			ColPodcast: map[string]bool{
				"author": false, "title": false, "keywords": false, "subtitle": false,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// createCollectionMap creates a map of mongo collections so the program doesn't
// reallocate space for a collection every time a request is called
func createCollectionMap(db *mongo.Database) map[string]*mongo.Collection {
	collectionMap := make(map[string]*mongo.Collection, len(collections))
	for _, collection := range collections {
		collectionMap[collection] = db.Collection(collection)
	}
	return collectionMap
}
