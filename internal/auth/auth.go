package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Hash takes pwd []byte and returns hash type []byte
func Hash(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v", err)
		return nil, err
	}
	return hash, nil
}

// Compare takes a []byte password and hash compares and returns true for match
func Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Comparing passwords failed: %v", err)
		return false
	}
	return true
}

// CreateSession creates a session and stores into database
func CreateSession(dbClient *database.Client, userID primitive.ObjectID) (string, error) {
	// Create key
	key := CreateKey(32)

	// Create Session object
	session := models.Session{
		UserID:       userID,
		SessionKey:   key,
		LoginTime:    time.Now(),
		LastSeenTime: time.Now(),
		Expires:      time.Now().Add(time.Hour * 24 * 30),
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
func ValidateSession(dbClient *database.Client, key string) (*models.User, error) {
	// Find the key
	var sesh models.Session
	err := dbClient.Find(database.ColSession, "session_key", key, &sesh)
	if err != nil {
		return nil, err
	}

	// Check if expired
	if sesh.Expires.After(time.Now()) {
		return nil, errors.New("session expired")
	}

	// Find the user
	var user models.User
	err = dbClient.FindByID(database.ColUser, sesh.UserID, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
