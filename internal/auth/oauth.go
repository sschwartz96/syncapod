package auth

import (
	"errors"
	"time"

	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateAuthorizationCode creates and saves an authorization code with the client & user id
func CreateAuthorizationCode(dbClient *database.Client, userID primitive.ObjectID, clientID string) string {
	code := models.AuthCode{
		Code:     CreateKey(64),
		ClientID: clientID,
		UserID:   userID,
		Scope:    models.SubScope,
	}

	dbClient.Insert(database.ColAuthCode, &code)
	return code.Code
}

// CreateAccessToken creates and saves an access token with a year of validity
func CreateAccessToken(dbClient *database.Client, authCode *models.AuthCode) *models.AccessToken {
	token := models.AccessToken{
		AuthCode:     authCode.Code,
		Token:        CreateKey(32),
		RefreshToken: CreateKey(32),
		UserID:       authCode.UserID,
		Created:      time.Now(),
		Expires:      time.Now().Add(time.Minute * 6),
	}

	dbClient.Insert(database.ColAccessToken, &token)
	return &token
}

// ValidateAuthCode takes pointer to db client and code string, finds the code and returns it
func ValidateAuthCode(dbClient *database.Client, code string) (*models.AuthCode, error) {
	var authCode models.AuthCode

	err := dbClient.Find(database.ColAuthCode, "code", code, &authCode)
	if err != nil {
		return nil, err
	}

	return &authCode, nil
}

// ValidateAccessToken takes pointer to dbclient and access_token and checks its validity
func ValidateAccessToken(dbClient *database.Client, token string) (*models.User, error) {
	var tokenObj models.AccessToken
	err := dbClient.Find(database.ColAccessToken, "token", token, &tokenObj)
	if err != nil {
		return nil, err
	}

	// if expired
	if tokenObj.Created.Add(time.Second * time.Duration(tokenObj.Expires).Before(time.Now()) {
		return nil, errors.New("expired token")
	}

	var user models.User
	err = dbClient.FindByID(database.ColUser, tokenObj.UserID, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
