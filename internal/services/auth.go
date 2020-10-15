package services

import (
	"context"
	"fmt"

	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/user"

	"github.com/sschwartz96/syncapod/internal/auth"
)

// AuthService is the gRPC service for authentication and authorization
type AuthService struct {
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
		res.Message = fmt.Sprint("failed on error: ", err)
		return res, nil
	}
	// authenticate
	if auth.Compare(user.Password, req.Password) {
		// create session
		key, err := auth.CreateSession(a.dbClient, user.Id, req.UserAgent, req.StayLoggedIn)
		if err != nil {
			fmt.Println("error creating session:", err)
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
	fmt.Println("received grpc authorize request")
	res := &protos.AuthRes{}

	user, err := auth.ValidateSession(a.dbClient, req.SessionKey)
	if err != nil {
		//res.Message = fmt.Sprint("Authorize() error validating user session:", err)
		fmt.Println("Authorize() error validating user session:", err)
		res.Success = false
	} else {
		res.Success = true
		res.SessionKey = req.SessionKey
		user.Password = ""
		res.User = user
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
