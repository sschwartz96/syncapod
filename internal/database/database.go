package database

import (
	"context"

	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/sschwartz96/syncapod/internal/protos"
)

type Filter map[string]interface{}

// Database defines database functionality
type Database interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error

	Insert(collection string, object interface{}) error
	Find(collection string, object interface{}, filter Filter) error
	Update(collection string, object interface{}, filter Filter) error
	Upsert(collection string, object interface{}, filter Filter) error
	Delete(collection string, filter Filter) error
}

// SyncapodStore is the interface that shows the behavior of our data storage of the app
type SyncapodStore interface {
	Open(context.Context) error
	Close(context.Context) error

	// * oAuth *
	InsertAuthCode(code *models.AuthCode) error
	FindAuthCode(code string) (*models.AuthCode, error)
	InsertAccessToken(token *models.AccessToken) error
	FindAccessToken(token string) (*models.AccessToken, error)

	// * Auth *

	FindSession(key string) (*protos.Session, error)
	UpsertSession(session *protos.Session) error
	DeleteSession(id *protos.ObjectID) error
	// FindUser finds the user based on username OR email
	FindUser(username string) (*protos.User, error)
	FindUserByID(id *protos.ObjectID) (*protos.User, error)
	DeleteUser(id *protos.ObjectID) error

	// * Podcast *

	FindAllPodcasts() ([]*protos.Podcast, error)
	FindPodcastsByRange(start, end int) ([]*protos.Podcast, error)
	FindPodcastByID(id *protos.ObjectID) (*protos.Podcast, error)

	// FindEpisodes returns a list of episodes based on podcast id
	// returns in chronological order, sectioned by start & end

	// * Episode *

	FindEpisodesByRange(podcastID *protos.ObjectID, start, end int) ([]*protos.Episode, error)
	FindAllEpisodes(podcastID *protos.ObjectID) ([]*protos.Episode, error)
	FindLatestEpisode(podcastID *protos.ObjectID) (*protos.Episode, error)
	FindEpisodeByID(id *protos.ObjectID) (*protos.Episode, error)
	// FindEpisodeBySeason takes a season episode number returns error if not found
	FindEpisodeBySeason(id *protos.ObjectID, seasonNum, episodeNum int) (*protos.Episode, error)
	UpsertEpisode(episode *protos.Episode) error

	// * UserEpisode *

	FindUserEpisode(userID, episodeID *protos.ObjectID) (*protos.UserEpisode, error)
	FindLatestUserEpisode(userID *protos.ObjectID) (*protos.UserEpisode, error)
	UpsertUserEpisode(userEpisode *protos.UserEpisode) error

	// Subscriptions

	FindSubscriptions(userID *protos.ObjectID) ([]*protos.Subscription, error)
	UpsertSubscription(subscription *protos.Subscription) error
}
