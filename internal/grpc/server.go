package grpc

import "github.com/sschwartz96/syncapod/internal/database"

// Server serves gRPC content
type Server struct {
	dbClient *database.Client
}

// CreateServer creates a gRPC Server
func CreateServer(dbClient *database.Client) (*Server, error) {
	s := &Server{
		dbClient: dbClient,
	}
	return s, nil
}
