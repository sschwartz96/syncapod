package user

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/podcast"
	"github.com/sschwartz96/syncapod/internal/protos"
)

// * Auth *
func FindSession(dbClient db.Database, key string) (*protos.Session, error) {
	session := &protos.Session{}
	err := dbClient.FindOne(database.ColSession, session, &db.Filter{"sessionkey": key}, nil)

	if err != nil {
		return nil, fmt.Errorf("error finding session: %v", err)
	}
	return session, nil
}

func UpsertSession(dbClient db.Database, session *protos.Session) error {
	if err := dbClient.Upsert(database.ColSession, session, &db.Filter{"_id": session.Id}); err != nil {
		return fmt.Errorf("error upserting session: %v", err)
	}
	return nil
}

func DeleteSession(dbClient db.Database, id *protos.ObjectID) error {
	err := dbClient.Delete(database.ColSession, &db.Filter{"_id": id})
	if err != nil {
		return fmt.Errorf("error deleting session: %v", err)
	}
	return nil
}

func DeleteSessionByKey(dbClient db.Database, key string) error {
	err := dbClient.Delete(database.ColSession, &db.Filter{"sessionkey": key})
	return fmt.Errorf("error deleting session by key: %v", err)
}

func FindUserByID(dbClient db.Database, id *protos.ObjectID) (*protos.User, error) {
	user := &protos.User{}
	err := dbClient.FindOne(database.ColUser, user, &db.Filter{"_id": id}, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding user by id: %v", err)
	}
	return user, nil
}

// FindUser attempts to find user by username/email returns pointer to user or error if not found
func FindUser(dbClient db.Database, username string) (*protos.User, error) {
	var key string
	if strings.Contains(username, "@") {
		key = "email"
		username = strings.ToLower(username)
	} else {
		key = "username"
	}

	user := &protos.User{}
	err := dbClient.FindOne(database.ColUser, user, &db.Filter{key: username}, nil)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return nil, fmt.Errorf("error finding user by %s: user does not exist", key)
		}
		return nil, fmt.Errorf("error finding user by %s: %v", key, err)
	}

	return user, nil
}

func DeleteUser(dbClient db.Database, id *protos.ObjectID) error {
	if err := dbClient.Delete(database.ColUser, &db.Filter{"_id": id}); err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

// UserEpisode

func FindUserEpisode(dbClient db.Database, userID *protos.ObjectID, episodeID *protos.ObjectID) (*protos.UserEpisode, error) {
	userEpisode := &protos.UserEpisode{}
	filter := &db.Filter{
		"userid":    userID,
		"episodeid": episodeID,
	}
	err := dbClient.FindOne(database.ColUserEpisode, userEpisode, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding user episode details: %v", err)
	}
	return userEpisode, nil
}

func FindLatestUserEpisode(dbClient db.Database, userID *protos.ObjectID) (*protos.UserEpisode, error) {
	userEpi := &protos.UserEpisode{}
	filter := &db.Filter{"userid": userID}
	opts := db.CreateOptions().SetSort("lastseen", -1)
	err := dbClient.FindOne(database.ColUserEpisode, userEpi, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding latest user episode: %v", err)
	}
	return userEpi, nil
}

func UpsertUserEpisode(dbClient db.Database, userEpisode *protos.UserEpisode) error {
	err := dbClient.Upsert(database.ColUserEpisode, userEpisode, &db.Filter{"_id": userEpisode.Id})
	if err != nil {
		return fmt.Errorf("error upserting user episode: %v", err)
	}
	return nil
}

// Subscriptions

func FindSubscriptions(dbClient db.Database, userID *protos.ObjectID) ([]*protos.Subscription, error) {
	var subs []*protos.Subscription
	err := dbClient.FindAll(database.ColSubscription, &subs, &db.Filter{"userid": userID}, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding subscriptions: %v", err)
	}
	return subs, nil
}

func UpsertSubscription(dbClient db.Database, subscription *protos.Subscription) error {
	err := dbClient.Upsert(database.ColSubscription, subscription, &db.Filter{"_id": subscription.Id})
	if err != nil {
		return fmt.Errorf("error upserting subscription: %v", err)
	}
	return nil
}

// helpers

// FindUserLastPlayed takes dbClient, userID, returns the latest played episode and offset
func FindUserLastPlayed(dbClient db.Database, userID *protos.ObjectID) (*protos.Podcast, *protos.Episode, *protos.UserEpisode, error) {
	pod := &protos.Podcast{}
	epi := &protos.Episode{}

	// find the latest played user_episode
	userEpi, err := FindLatestUserEpisode(dbClient, userID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error finding last user played: %v", err)
	}

	// concurrently
	poderr := make(chan error)
	epierr := make(chan error)

	// find podcast
	go func() {
		pod, err = podcast.FindPodcastByID(dbClient, userEpi.PodcastID)
		if err != nil {
			poderr <- err
		}
		poderr <- nil
	}()

	// find episode
	go func() {
		epi, err = podcast.FindEpisodeByID(dbClient, userEpi.EpisodeID)
		if err != nil {
			epierr <- err
		}
		epierr <- nil
	}()

	err = <-poderr
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error finding laster user episode: %v", err)
	}
	err = <-epierr
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error finding laster user episode: %v", err)
	}

	return pod, epi, userEpi, nil
}

// FindOffset takes database client and pointers to user and episode to lookup episode details and offset
func FindOffset(dbClient db.Database, userID, epiID *protos.ObjectID) int64 {
	userEpi, err := FindUserEpisode(dbClient, userID, epiID)
	if err != nil {
		fmt.Println("error finding offset: ", err)
		return 0
	}
	return userEpi.Offset
}

// UpdateOffset takes userID epiID and offset and performs upsert to the UserEpisode collection
func UpdateOffset(dbClient db.Database, uID, pID, eID *protos.ObjectID, offset int64) error {
	userEpi := &protos.UserEpisode{
		UserID:    uID,
		PodcastID: pID,
		EpisodeID: eID,
		Offset:    offset,
		Played:    false,
		LastSeen:  ptypes.TimestampNow(),
	}
	return UpsertUserEpisode(dbClient, userEpi)
}

func UpdateUserEpiPlayed(dbClient db.Database, uID, pID, eID *protos.ObjectID, played bool) error {
	userEpi := &protos.UserEpisode{
		UserID:    uID,
		PodcastID: pID,
		EpisodeID: eID,
		Offset:    0,
		Played:    played,
		LastSeen:  ptypes.TimestampNow(),
	}
	return UpsertUserEpisode(dbClient, userEpi)
}
