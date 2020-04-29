package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Hash takes pwd string and returns hash type string
func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v", err)
		return "", err
	}
	return string(hash), nil
}

// Compare takes a password and hash compares and returns true for match
func Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Comparing passwords failed: %v", err)
		return false
	}
	return true
}

// CreateSession creates a session and stores into database
func CreateSession(dbClient *database.Client, userID primitive.ObjectID,
	expires time.Duration, userAgent string) (string, error) {
	// Create key
	key := CreateKey(32)

	if userAgent == "" {
		userAgent = "unknown"
	}

	// Create Session object
	session := models.Session{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		SessionKey:   key,
		LoginTime:    time.Now(),
		LastSeenTime: time.Now(),
		Expires:      time.Now().Add(expires),
		UserAgent:    userAgent,
	}

	// Store session in database
	err := dbClient.Insert(database.ColSession, &session)
	if err != nil {
		return "", err
	}

	return key, nil
}

// CreateKey takes in a key length and returns base64 encoding
func CreateKey(l int) string {
	key := make([]byte, l)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Printf("Could not make key with err: %v\n", err)
	}
	return base64.URLEncoding.EncodeToString(key)[:l]
}

// ValidateSession looks up session key, check if its valid and returns a pointer to the user
// returns error if the key doesn't exist, or has expired
func ValidateSession(dbClient *database.Client, key string) (*models.User, error) {
	// Find the key
	var sesh models.Session
	err := dbClient.Find(database.ColSession, "session_key", key, &sesh)
	if err != nil {
		fmt.Println("validate sesion, couldn't find session key")
		return nil, err
	}

	// Check if expired
	if sesh.Expires.Before(time.Now()) {
		err := dbClient.Delete(database.ColSession, "_id", sesh.ID)
		if err != nil {
			fmt.Println("couldn't delete session: ", err)
		}
		return nil, errors.New("session expired")
	}

	sesh.LastSeenTime = time.Now()
	sesh.Expires.Add(time.Hour * 1)
	go dbClient.Upsert(database.ColSession, bson.M{"_id": sesh.ID}, sesh)

	// Find the user
	//var user models.User
	//err = dbClient.FindByID(database.ColUser, sesh.UserID, &user)
	user, err := FindUser(dbClient, sesh.UserID)
	if err != nil {
		fmt.Println("validate sesion, couldn't find user")
		return nil, err
	}

	return user, nil
}

// FindUser takes a pointer to database.Client and userID and returns user if
// found, error otherwise
func FindUser(dbClient *database.Client, userID primitive.ObjectID) (*models.User, error) {
	match := bson.D{{Key: "$match", Value: bson.M{"_id": userID}}}
	lookup := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: database.ColPodcast},
		{Key: "localField", Value: "subs"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "subs"},
	}}}
	pipeline := mongo.Pipeline{match, lookup}

	var user []models.User
	err := dbClient.Aggregate(database.ColUser, pipeline, &user)
	if err != nil {
		return nil, err
	}

	if len(user) == 1 {
		return &user[0], nil
	}

	return nil, errors.New("user not found")
}
