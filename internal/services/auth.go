package services

import (
	"context"
	"fmt"
	"log"

	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/user"

	"github.com/sschwartz96/syncapod/internal/auth"
)

// AuthService is the gRPC service for authentication and authorization
type AuthService struct {
	*protos.UnimplementedAuthServer
	dbClient db.Database
}

// NewAuthService creates a new *AuthService
func NewAuthService(dbClient db.Database) *AuthService {
	return &AuthService{dbClient: dbClient}
}

// Authenticate handles the authentication to syncapod and returns response
func (a *AuthService) Authenticate(ctx context.Context, req *protos.AuthReq) (*protos.AuthRes, error) {
	res := &protos.AuthRes{Success: false}
	// find user from database
	user, err := user.FindUser(a.dbClient, req.Username)
	if err != nil {
		return nil, fmt.Errorf("Authenticate(), error finding user: %v", err)
	}
	// authenticate
	if auth.Compare(user.Password, req.Password) {
		// create session
		key, err := auth.CreateSession(a.dbClient, user.Id, req.UserAgent, req.StayLoggedIn)
		if err != nil {
			log.Println("error creating session:", err)
			return nil, fmt.Errorf("Authenticate(), error creating session: %v", err)
		} else {
			res.Success = true
			res.User = user
			res.SessionKey = key
			res.User.Password = ""
		}
	}
	return res, nil
}

// Authorize authorizes user based on a session key
func (a *AuthService) Authorize(ctx context.Context, req *protos.AuthReq) (*protos.AuthRes, error) {
	user, err := auth.ValidateSession(a.dbClient, req.SessionKey)
	if err != nil {
		log.Println("error in authroize")
		return nil, fmt.Errorf("Authorize() error validating user session: %v", err)
	}
	user.Password = ""
	res := &protos.AuthRes{
		Success:    true,
		SessionKey: req.SessionKey,
		User:       user,
	}
	return res, nil
}

// Logout removes the given session key
func (a *AuthService) Logout(ctx context.Context, req *protos.AuthReq) (*protos.AuthRes, error) {
	success := true
	message := ""
	err := user.DeleteSessionByKey(a.dbClient, req.SessionKey)
	if err != nil {
		message = fmt.Sprintf("Logout() error deleting session: %v", err)
		fmt.Println("Logout() error deleting session:", err)
		success = false
	}
	return &protos.AuthRes{Success: success, Message: message}, nil
}
