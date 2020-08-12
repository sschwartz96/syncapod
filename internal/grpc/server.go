package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

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

func (s *Server) Start() {
	var creds credentials.TransportCredentials
	// whether or not we are running tls
	if s.config.CertFile != "" {
		creds, err := credentials.NewServerTLSFromFile(s.config.CertFile, s.config.KeyFile)
		if err != nil {
			log.Fatal("error setting up creds for grpc:", creds)
		}
	}

	// setup server
	var gOptCred grpc.ServerOption
	gOptInter := grpc.UnaryInterceptor(s.Intercept())
	if creds != nil {
		gOptCred = grpc.Creds(creds)
	}
	grpcServer := grpc.NewServer(gOptCred, gOptInter)

	// start listener
	grpcListener, err := net.Listen("tcp", ":"+strconv.Itoa(s.config.GRPCPort))
	if err != nil {
		log.Fatalf("could not listen on port %d, err: %v", s.config.GRPCPort, err)
	}

	// register services
	reflection.Register(grpcServer)
	protos.RegisterAuthServer(grpcServer, services.NewAuthService(s.dbClient))
	protos.RegisterPodcastServiceServer(grpcServer, services.NewPodcastService(s.dbClient))

	// serve
	err = grpcServer.Serve(grpcListener)
	if err != nil {
		log.Fatal("could not serve services:", err)
	}
}

func (s *Server) Intercept() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("server gRPC interceptor function")
		return nil, nil
	}
}
