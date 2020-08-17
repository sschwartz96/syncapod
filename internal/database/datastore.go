package database

import "github.com/sschwartz96/syncapod/internal/protos"

// Datastore is the interface that shows the behavior of our data storage of the app
type Datastore interface {
	// Datastore specific
	Open() error
	Close() error

	// Auth
	FindSession(key string) (*protos.Session, error)
	InsertSession(session *protos.Session) error
	UpdateSession(session *protos.Session) error
	FindUser(username string) (*protos.User, error) // must accept username OR email

	// Podcast
	FindEpisodes(podcastID *protos.ObjectID, start, end int) ([]*protos.Episode, error)

	// UserEpisode
	FindUserEpisode(userID, episodeID *protos.ObjectID) (*protos.UserEpisode, error)
	UpsertUserEpisode(userEpisode *protos.UserEpisode) error

	// Subscriptions
	FindSubscriptions(userID *protos.ObjectID) ([]*protos.Podcast, error)
	UpsertSubscription(subscription *protos.Subscription) error
}
