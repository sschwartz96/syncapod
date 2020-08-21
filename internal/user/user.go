package user

import (
	"fmt"
	"strings"

	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
)

// * Auth *
func FindSession(db database.Database, key string) (*protos.Session, error) {
	var session *protos.Session
	err := db.Find(database.ColSession, session, database.Filter{"sessionkey": &key})

	if err != nil {
		return nil, fmt.Errorf("error finding session: %v", err)
	}
	return session, nil
}

func UpsertSession(db database.Database, session *protos.Session) error {
	if err := db.Upsert(database.ColSession, session, database.Filter{"_id": session.Id}); err != nil {
		return fmt.Errorf("error upserting session: %v", err)
	}
	return nil
}

func DeleteSession(db database.Database, id *protos.ObjectID) error {
	err := db.Delete(database.ColSession, database.Filter{"_id": id})
	return fmt.Errorf("error deleting session: %v", err)
}

func FindUserByID(db database.Database, id *protos.ObjectID) (*protos.User, error) {
	var user *protos.User
	err := db.Find(database.ColUser, user, database.Filter{"_id": id})
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
	err := db.Find(database.ColUser, user, database.Filter{key: username})
	if err != nil {
		return nil, fmt.Errorf("error finding user by %s: %v", key, err)
	}

	return user, nil
}

func DeleteUser(db database.Database, id *protos.ObjectID) error {
	if err := db.Delete(database.ColUser, database.Filter{"_id": id}); err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}
