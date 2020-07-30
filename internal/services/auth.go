package services

import (
	"context"
	"fmt"

	"github.com/sschwartz96/syncapod/internal/protos"

	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/database"
)

// AuthService is the gRPC service for authentication and authorization
type AuthService struct {
	dbClient *database.Client
}

// NewAuthService creates a new *AuthService
func NewAuthService(dbClient *database.Client) *AuthService {
	return &AuthService{dbClient: dbClient}
}

// Authenticate handles the authentication to syncapod and returns response
func (a *AuthService) Authenticate(ctx context.Context, req *protos.AuthReq) (*protos.AuthRes, error) {
	res := &protos.AuthRes{}

	// find user from database
	user, err := a.dbClient.FindUser(req.Username)
	if err != nil {
		res.Success = false
	}

	// authenticate
	if auth.Compare(user.Password, req.Password) {
		// create session
		key, err := auth.CreateSession(a.dbClient, user.Id, req.UserAgent, req.StayLoggedIn)
		if err != nil {
			fmt.Println("error creating session:", err)
			res.Success = false
		} else {
			res.Success = true
			res.User = user
			res.SessionKey = key
			res.User.Password = ""
		}
	} else {
		res.Success = false
	}
	return res, nil
}

// Authorize authorizes user based on a session key
func (a *AuthService) Authorize(ctx context.Context, req *protos.AuthReq) (*protos.AuthRes, error) {
	fmt.Println("received grpc authorize request")
	res := &protos.AuthRes{}

	user, err := auth.ValidateSession(a.dbClient, req.SessionKey)
	if err != nil {
		fmt.Println("error validating user session:", err)
		res.Success = false
	} else {
		res.Success = true
		res.SessionKey = req.SessionKey
		res.User = user
	}

	return res, nil
}

// Logout removes the given session key
func (a *AuthService) Logout(ctx context.Context, req *protos.AuthReq) (*protos.AuthRes, error) {
	success := true
	err := a.dbClient.Delete(database.ColSession, "sessionkey", req.SessionKey)
	if err != nil {
		fmt.Println("error logging out:", err)
		success = false
	}
	return &protos.AuthRes{Success: success}, nil
}
