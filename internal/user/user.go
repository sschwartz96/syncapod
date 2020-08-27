package user

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/podcast"
	"github.com/sschwartz96/syncapod/internal/protos"
)

// * Auth *
func FindSession(db database.Database, key string) (*protos.Session, error) {
	var session *protos.Session
	err := db.FindOne(database.ColSession, session, &database.Filter{"sessionkey": &key}, nil)

	if err != nil {
		return nil, fmt.Errorf("error finding session: %v", err)
	}
	return session, nil
}

func UpsertSession(db database.Database, session *protos.Session) error {
	if err := db.Upsert(database.ColSession, session, &database.Filter{"_id": session.Id}); err != nil {
		return fmt.Errorf("error upserting session: %v", err)
	}
	return nil
}

func DeleteSession(db database.Database, id *protos.ObjectID) error {
	err := db.Delete(database.ColSession, &database.Filter{"_id": id})
	return fmt.Errorf("error deleting session: %v", err)
}

func DeleteSessionByKey(db database.Database, key string) error {
	err := db.Delete(database.ColSession, &database.Filter{"sessionkey": key})
	return fmt.Errorf("error deleting session by key: %v", err)
}

func FindUserByID(db database.Database, id *protos.ObjectID) (*protos.User, error) {
	var user *protos.User
	err := db.FindOne(database.ColUser, user, &database.Filter{"_id": id}, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding user by id: %v", err)
	}
	return user, nil
}

// FindUser attempts to find user by username/email returns pointer to user or error if not found
func FindUser(db database.Database, username string) (*protos.User, error) {
	var key string
	if strings.Contains(username, "@") {
		key = "email"
		username = strings.ToLower(username)
	} else {
		key = "username"
	}

	var user *protos.User
	err := db.FindOne(database.ColUser, user, &database.Filter{key: username}, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding user by %s: %v", key, err)
	}

	return user, nil
}

func DeleteUser(db database.Database, id *protos.ObjectID) error {
	if err := db.Delete(database.ColUser, &database.Filter{"_id": id}); err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

// UserEpisode

func FindUserEpisode(db database.Database, userID *protos.ObjectID, episodeID *protos.ObjectID) (*protos.UserEpisode, error) {
	var userEpisode protos.UserEpisode
	filter := &database.Filter{
		"userid":    userID,
		"episodeid": episodeID,
	}
	err := db.FindOne(database.ColUserEpisode, &userEpisode, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding user episode details: %v", err)
	}
	return &userEpisode, nil
}

func FindLatestUserEpisode(db database.Database, userID *protos.ObjectID) (*protos.UserEpisode, error) {
	var userEpi *protos.UserEpisode
	filter := &database.Filter{"userid": userID}
	opts := database.CreateOptions().SetSort("lastseen", -1)
	err := db.FindOne(database.ColUserEpisode, userEpi, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding latest user episode: %v", err)
	}
	return userEpi, nil
}

func UpsertUserEpisode(db database.Database, userEpisode *protos.UserEpisode) error {
	err := db.Upsert(database.ColUserEpisode, userEpisode, &database.Filter{"_id": userEpisode.Id})
	if err != nil {
		return fmt.Errorf("error upserting user episode: %v", err)
	}
	return nil
}

// Subscriptions

func FindSubscriptions(db database.Database, userID *protos.ObjectID) ([]*protos.Subscription, error) {
	var subs []*protos.Subscription
	err := db.FindAll(database.ColSubscription, &subs, &database.Filter{"userid": userID}, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding subscriptions: %v", err)
	}
	return subs, nil
}

func UpsertSubscription(db database.Database, subscription *protos.Subscription) error {
	err := db.Upsert(database.ColSubscription, subscription, &database.Filter{"_id": subscription.Id})
	if err != nil {
		return fmt.Errorf("error upserting subscription: %v", err)
	}
	return nil
}

// helpers

// FindUserLastPlayed takes dbClient, userID, returns the latest played episode and offset
func FindUserLastPlayed(db database.Database, userID *protos.ObjectID) (*protos.Podcast, *protos.Episode, *protos.UserEpisode, error) {
	var pod *protos.Podcast
	var epi *protos.Episode

	// find the latest played user_episode
	userEpi, err := FindLatestUserEpisode(db, userID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error finding last user played: %v", err)
	}

	// concurrently
	var poderr, epierr chan error

	// find podcast
	go func() {
		pod, err = podcast.FindPodcastByID(db, userEpi.PodcastID)
		if err != nil {
			poderr <- err
		}
		poderr <- nil
	}()

	// find episode
	go func() {
		epi, err = podcast.FindEpisodeByID(db, userEpi.EpisodeID)
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
func FindOffset(db database.Database, userID, epiID *protos.ObjectID) int64 {
	userEpi, err := FindUserEpisode(db, userID, epiID)
	if err != nil {
		fmt.Println("error finding offset: ", err)
		return 0
	}
	return userEpi.Offset
}

// UpdateOffset takes userID epiID and offset and performs upsert to the UserEpisode collection
func UpdateOffset(db database.Database, uID, pID, eID *protos.ObjectID, offset int64) error {
	userEpi := &protos.UserEpisode{
		UserID:    uID,
		PodcastID: pID,
		EpisodeID: eID,
		Offset:    offset,
		Played:    false,
		LastSeen:  ptypes.TimestampNow(),
	}
	return UpsertUserEpisode(db, userEpi)
}

func UpdateUserEpiPlayed(db database.Database, uID, pID, eID *protos.ObjectID, played bool) error {
	userEpi := &protos.UserEpisode{
		UserID:    uID,
		PodcastID: pID,
		EpisodeID: eID,
		Offset:    0,
		Played:    played,
		LastSeen:  ptypes.TimestampNow(),
	}
	return UpsertUserEpisode(db, userEpi)
}
