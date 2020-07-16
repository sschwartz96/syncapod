package grpc

import (
	"context"

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

	// match user

}

func (a *AuthService) Authorize(ctx context.Context, req *protos.AuthReq) (*protos.AuthRes, error) {

}
