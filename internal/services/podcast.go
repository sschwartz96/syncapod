package services

import "github.com/sschwartz96/syncapod/internal/database"

// PodcastService is the gRPC service for podcast
type PodcastService struct {
	dbClient *database.Client
}

// NewPodcastService creates a new *PodcastService
func NewPodcastService(dbClient *database.Client) *PodcastService {
	return &PodcastService{dbClient: dbClient}
}
