package services

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/syncapod/internal/podcast"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/user"
	"google.golang.org/grpc/metadata"
)

const (
	ctxUserIDVal = "user_id"
)

// PodcastService is the gRPC service for podcast
type PodcastService struct {
	*protos.UnimplementedPodServer
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
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetUserEpisode() error getting userID: %v", err)
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
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetSubscriptions() error getting user id: %v", err)
	}

	subs, err := user.FindSubscriptions(p.dbClient, userID)
	if err != nil {
		log.Println("GetSubscriptions() error getting subs:", err)
		return &protos.Subscriptions{}, nil
	}

	return &protos.Subscriptions{Subscriptions: subs}, nil
}

// GetUserLastPlayed returns the last episode the user was playing & metadata
func (p *PodcastService) GetUserLastPlayed(ctx context.Context, req *protos.Request) (*protos.LastPlayedRes, error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetUserLastPlayed() error getting user id: %v", err)
	}

	pod, epi, userEpi, err := user.FindUserLastPlayed(p.dbClient, userID)
	if err != nil {
		fmt.Println("GetUserLastPlayed() error getting podcast, episode and userEpi: ", err)
		return nil, err
	}

	return &protos.LastPlayedRes{
		Podcast: pod,
		Episode: epi,
		Millis:  userEpi.Offset,
	}, nil
}

func getUserIDFromContext(ctx context.Context) (*protos.ObjectID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("getUserIDFromContext() error: metadata not valid")
	}
	idHex := md.Get("user_id")
	if len(idHex) == 0 {
		return nil, fmt.Errorf("getUserIDFromContext() error: no user id")
	}
	return protos.ObjectIDFromHex(idHex[0]), nil
}
