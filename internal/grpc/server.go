package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"strings"

	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// Server is truly needed for its Intercept method which authenticates users before
// accessing services
type Server struct {
	server *grpc.Server
	db     db.Database
}

func NewServer(config *config.Config, dbClient db.Database) *Server {
	var grpcServer *grpc.Server
	s := &Server{db: dbClient}

	// setup server
	gOptCreds := getTransportCreds(config)
	gOptInter := grpc.UnaryInterceptor(s.Intercept())
	grpcServer = grpc.NewServer(gOptCreds, gOptInter)
	s.server = grpcServer

	// register services
	reflection.Register(grpcServer)

	as := services.NewAuthService(dbClient)
	protos.RegisterAuthService(
		grpcServer,
		protos.NewAuthService(as),
	)

	pd := services.NewPodcastService(dbClient)
	protos.RegisterPodService(
		grpcServer,
		protos.NewPodService(pd),
	)

	return s
}

func (s *Server) Start(list net.Listener) error {
	return s.server.Serve(list)
}

func getTransportCreds(config *config.Config) grpc.ServerOption {
	var creds credentials.TransportCredentials

	// whether or not we are running tls
	if config.CertFile != "" {
		creds, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
		if err != nil {
			log.Fatal("error setting up creds for grpc:", creds)
		}
	}

	return grpc.Creds(creds)
}

func (s *Server) Intercept() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// if this is going to the Auth service allow through
		if strings.Contains(info.FullMethod, "protos.Auth") {
			return handler(ctx, req)
		}

		token, ok := ctx.Value("token").(string)
		if !ok {
			return nil, errors.New("invalid access token, not string format")
		}
		if token == "" {
			return nil, errors.New("invalid access token, empty")
		}

		user, err := auth.ValidateAccessToken(s.db, token)
		if err != nil {
			return nil, errors.New("invalid access token")
		}

		userCtx := context.WithValue(ctx, "user_id", user.Id)

		return handler(userCtx, req)
	}
}
