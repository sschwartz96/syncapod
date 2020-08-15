package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	config   *config.Config
	dbClient *database.Client
}

func NewServer(config *config.Config, dbClient *database.Client) *Server {
	return &Server{config: config, dbClient: dbClient}
}

func (s *Server) Start() error {
	var grpcServer *grpc.Server

	// setup server
	gOptCreds := getTransportCreds(s.config)
	gOptInter := grpc.UnaryInterceptor(s.Intercept())
	grpcServer = grpc.NewServer(gOptCreds, gOptInter)

	// setup listener
	grpcListener, err := net.Listen("tcp", ":"+strconv.Itoa(s.config.GRPCPort))
	if err != nil {
		return errors.New(fmt.Sprintf("could not listen on port %d, err: %v", s.config.GRPCPort, err))
	}

	// register services
	reflection.Register(grpcServer)
	protos.RegisterAuthServer(grpcServer, services.NewAuthService(s.dbClient))
	protos.RegisterPodcastServiceServer(grpcServer, services.NewPodcastService(s.dbClient))

	// serve
	return grpcServer.Serve(grpcListener)
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
		token, ok := ctx.Value("token").(string)
		if !ok {
			return nil, errors.New("invalid access token, not string format")
		}
		if token == "" {
			return nil, errors.New("invalid access token, empty")
		}

		user, err := auth.ValidateAccessToken(s.dbClient, token)
		if err != nil {
			return nil, errors.New("invalid access token")
		}

		userCtx := context.WithValue(ctx, "user_id", user.Id)

		return handler(userCtx, req)
	}
}
