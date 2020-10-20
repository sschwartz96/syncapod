package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

// Server is truly needed for its Intercept method which authenticates users before accessing services,
// but also useful to have all the grpc server boilerplate contained within NewServer function
type Server struct {
	server *grpc.Server
	db     db.Database
}

func NewServer(cfg *config.Config, dbClient db.Database, aS protos.AuthServer, pS protos.PodServer) *Server {
	var grpcServer *grpc.Server
	s := &Server{db: dbClient}
	// setup server
	gOptCreds := getTransportCreds(cfg)
	gOptInter := grpc.UnaryInterceptor(s.Intercept())
	grpcServer = grpc.NewServer(gOptCreds, gOptInter)
	s.server = grpcServer
	// register services
	reflection.Register(grpcServer)
	protos.RegisterAuthServer(s.server, aS)
	protos.RegisterPodServer(s.server, pS)
	return s
}

func (s *Server) Start(lis net.Listener) error {
	return s.server.Serve(lis)
}

func getTransportCreds(config *config.Config) grpc.ServerOption {
	var creds credentials.TransportCredentials
	var err error
	// whether or not we are running tls
	if config.CertFile != "" {
		creds, err = credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
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

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("invalid metadata")
		}
		token := md.Get("token")
		if len(token) == 0 {
			return nil, errors.New("no access token sent")
		}

		user, err := auth.ValidateSession(s.db, token[0])
		if err != nil {
			return nil, fmt.Errorf("invalid access token: %v", err)
		}

		//md.Set("user_id", user.Id.Hex) // causes errors
		newMD := md.Copy()
		newMD.Set("user_id", user.Id.Hex)
		newCtx := metadata.NewIncomingContext(ctx, newMD)

		return handler(newCtx, req)
	}
}
