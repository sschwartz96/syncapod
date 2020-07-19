package grpc

import (
	"context"
	"time"

	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
)

type AuthService struct {
	dbClient *database.Client
}

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
		auth.CreateSession(a.dbClient, user.Id, time.Hour, req.UserAgent)
		res.Success = true
	} else {
		res.Success = false
	}
	return res, nil
}

func (a *AuthService) Authorize(ctx context.Context, req *protos.AuthReq) (*protos.AuthRes, error) {

}
