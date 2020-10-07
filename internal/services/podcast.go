package services

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/syncapod/internal/podcast"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/user"
)

const (
	ctxUserIDVal = "user_id"
)

// PodcastService is the gRPC service for podcast
type PodcastService struct {
	dbClient db.Database
}

// NewPodcastService creates a new *PodcastService
func NewPodcastService(dbClient db.Database) *PodcastService {
	return &PodcastService{dbClient: dbClient}
}

// GetEpisodes returns a list of episodes via podcast id
func (p *PodcastService) GetEpisodes(ctx context.Context, req *protos.Request) (*protos.Episodes, error) {
	var episodes []*protos.Episode
	var err error
	// get the id and validate
	if req.PodcastID != nil || len(req.PodcastID.Hex) > 0 {
		episodes, err = podcast.FindEpisodesByRange(p.dbClient, req.PodcastID, req.Start, req.End)
		if err != nil {
			fmt.Println("error grpc GetEpisodes:", err)
			return &protos.Episodes{Episodes: []*protos.Episode{}}, nil
		}
	} else {
		return &protos.Episodes{Episodes: []*protos.Episode{}}, fmt.Errorf("no podcast id supplied")
	}
	return &protos.Episodes{Episodes: episodes}, nil
}

// GetUserEpisode returns the user playback metadata via episode id & user id
func (p *PodcastService) GetUserEpisode(ctx context.Context, req *protos.Request) (*protos.UserEpisode, error) {
	userID, ok := ctx.Value(ctxUserIDVal).(*protos.ObjectID)
	if !ok {
		fmt.Println("error: empty userID from context")
	}

	userEpi, err := user.FindUserEpisode(p.dbClient, userID, req.EpisodeID)
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
	err := user.UpsertUserEpisode(p.dbClient, userEpi)
	if err != nil {
		fmt.Println("error updating user episode", err)
		return &protos.Response{Success: false, Message: err.Error()}, nil
	}
	return &protos.Response{Success: true, Message: ""}, nil
}

// GetSubscriptions returns a list of podcasts via user id
func (p *PodcastService) GetSubscriptions(ctx context.Context, req *protos.Request) (*protos.Subscriptions, error) {
	userID, ok := ctx.Value(ctxUserIDVal).(*protos.ObjectID)
	if !ok {
		fmt.Println("error: empty userID from context")
	}

	subs, err := user.FindSubscriptions(p.dbClient, userID)
	if err != nil {
		fmt.Println("error getting subs:", err)
		return &protos.Subscriptions{}, nil
	}

	return &protos.Subscriptions{Subscriptions: subs}, nil
}

// GetUserLastPlayed returns the last episode the user was playing & metadata
func (p *PodcastService) GetUserLastPlayed(ctx context.Context, req *protos.Request) (*protos.LastPlayedRes, error) {
	userID, ok := ctx.Value(ctxUserIDVal).(*protos.ObjectID)
	if !ok {
		fmt.Println("error: empty userID from context")
	}

	pod, epi, userEpi, err := user.FindUserLastPlayed(p.dbClient, userID)
	if err != nil {
		fmt.Println("error getting last play:", err)
		return nil, err
	}

	return &protos.LastPlayedRes{
		Podcast: pod,
		Episode: epi,
		Millis:  userEpi.Offset,
	}, nil
}
