package services

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/podcast"
	"github.com/sschwartz96/syncapod/internal/protos"
)

// PodcastService is the gRPC service for podcast
type PodcastService struct {
	dbClient *database.Client
}

// NewPodcastService creates a new *PodcastService
func NewPodcastService(dbClient *database.Client) *PodcastService {
	return &PodcastService{dbClient: dbClient}
}

// GetEpisodes returns a list of episodes via episode id
func (p *PodcastService) GetEpisodes(ctx context.Context, req *protos.Request) (*protos.Episodes, error) {
	var episodes []*protos.Episode

	// get the id and validate
	if req.PodcastID != nil || len(req.PodcastID.Hex) > 0 {
		episodes = podcast.FindAllEpisodesRange(p.dbClient, req.PodcastID, req.start, req.end)
	}
	return &protos.Episodes{Episodes: episodes}, nil
}

// GetUserEpisode returns the user playback metadata via episode id & user id
func (p *PodcastService) GetUserEpisode(ctx context.Context, req *protos.Request) (*protos.UserEpisode, error) {
	userEpi, err := podcast.FindUserEpisode(p.dbClient, req.UserID, req.EpisodeID)
	if err != nil {
		fmt.Println("error finding userEpi:", err)
	}
	return userEpi, nil
}

// UpdateUserEpisode updates the user playback metadata via episode id & user id
func (p *PodcastService) UpdateUserEpisode(ctx context.Context, req *protos.UserEpisodeReq) (*protos.Response, error) {
	if req.LastSeen == nil {
		req.LastSeen = ptypes.TimestampNow()
	}
	userEpi := &protos.UserEpisode{

		EpisodeID: req.EpisodeID,
		PodcastID: req.PodcastID,
		Played:    req.Played,
		Offset:    req.Offset,
	}
	err := podcast.UpdateUserEpi(p.dbClient, req)
	if err != nil {
		fmt.Println("error updating user episode", err)
		return &protos.Response{Success: false, Message: err.Error()}, nil
	}
	return &protos.Response{Success: true, Message: ""}, nil
}

// GetSubscriptions returns a list of podcasts via user id
func (p *PodcastService) GetSubscriptions(ctx context.Context, req *protos.Request) (*protos.Subscriptions, error) {
	subs, err := podcast.GetSubscriptions(p.dbClient, req.UserID)
	if err != nil {
		fmt.Println("error getting subs:", err)
		return &protos.Subscriptions{}, nil
	}

	return &protos.Subscriptions{Subscriptions: subs}, nil
}

// GetUserLastPlayed returns the last episode the user was playing & metadata
func (p *PodcastService) GetUserLastPlayed(ctx context.Context, req *protos.Request) (*protos.LastPlayedRes, error) {
	pod, epi, millis, err := podcast.FindUserLastPlayed(p.dbClient, req.UserID)
	if err != nil {
		fmt.Println("error getting last play:", err)
		return nil, err
	}
	return &protos.LastPlayedRes{
		Podcast: pod,
		Episode: epi,
		Millis:  millis,
	}, nil
}
