package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/user"
)

// CreateAuthorizationCode creates and saves an authorization code with the client & user id
func CreateAuthorizationCode(dbClient db.Database, userID *protos.ObjectID, clientID string) (string, error) {
	code := models.AuthCode{
		Code:     CreateKey(64),
		ClientID: clientID,
		UserID:   userID,
		Scope:    models.SubScope,
	}

	err := insertAuthCode(dbClient, &code)
	if err != nil {
		return "", fmt.Errorf("error creating auth code: %v", err)
	}

	return code.Code, nil
}

// CreateAccessToken creates and saves an access token with a year of validity
func CreateAccessToken(dbClient db.Database, authCode *models.AuthCode) (*models.AccessToken, error) {
	token := models.AccessToken{
		AuthCode:     authCode.Code,
		Token:        CreateKey(32),
		RefreshToken: CreateKey(32),
		UserID:       authCode.UserID,
		Created:      time.Now(),
		Expires:      3600,
	}

	err := insertAccessToken(dbClient, &token)
	if err != nil {
		return nil, fmt.Errorf("error creating access token: %v", err)
	}

	return &token, nil
}

// ValidateAuthCode takes pointer to db client and code string, finds the code and returns it
func ValidateAuthCode(dbClient db.Database, code string) (*models.AuthCode, error) {
	authCode, err := findAuthCode(dbClient, code)
	if err != nil {
		return nil, fmt.Errorf("error validating auth code: %v", err)
	}

	return authCode, nil
}

// ValidateAccessToken takes pointer to dbclient and access_token and checks its validity
func ValidateAccessToken(dbClient db.Database, token string) (*protos.User, error) {
	tokenObj, err := FindOauthAccessToken(dbClient, token)
	if err != nil {
		return nil, err
	}

	// if expired
	if tokenObj.Created.Add(time.Second * time.Duration(tokenObj.Expires)).Before(time.Now()) {
		return nil, errors.New("expired access token")
	}

	u, err := user.FindUserByID(dbClient, tokenObj.UserID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func insertAuthCode(dbClient db.Database, code *models.AuthCode) error {
	if err := dbClient.Insert(database.ColAuthCode, code); err != nil {
		return fmt.Errorf("error inserting auth code: %v", err)
	}
	return nil
}

func findAuthCode(dbClient db.Database, code string) (*models.AuthCode, error) {
	var authCode *models.AuthCode
	err := dbClient.FindOne(database.ColAuthCode, authCode, &db.Filter{"auth_code": code}, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding auth code: %v", err)
	}
	return authCode, nil
}

func insertAccessToken(dbClient db.Database, token *models.AccessToken) error {
	if err := dbClient.Insert(database.ColAccessToken, token); err != nil {
		return fmt.Errorf("error inserting access token: %v", err)
	}
	return nil
}

func FindOauthAccessToken(dbClient db.Database, token string) (*models.AccessToken, error) {
	var accessToken *models.AccessToken
	err := dbClient.FindOne(database.ColAccessToken, accessToken, &db.Filter{"token": token}, nil)
	if err != nil {
		return nil, fmt.Errorf("error finding access token: %v", err)
	}
	return accessToken, nil
}

func DeleteOauthAccessToken(dbClient db.Database, token string) error {
	err := dbClient.Delete(database.ColAccessToken, &db.Filter{"token": token})
	if err != nil {
		return fmt.Errorf("error deleting oauth access token: %v", err)
	}
	return nil
}
