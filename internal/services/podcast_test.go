package services

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/stockpile/mock"
	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/grpc"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/util"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

// createAuthServiceMockDB fails on error and returns db.Database and *protos.User
func createPodcastServiceMockDB(t *testing.T) db.Database {
	dbClient := mock.CreateDB()
	podcast := &protos.Podcast{Id: protos.ObjectIDFromHex("pod_id"), Author: "Sam Schwartz", Title: "Mock Podcast"}
	err := dbClient.Insert(database.ColPodcast, podcast)
	if err != nil {
		t.Fatalf("createAuthSerivceMockDB() error inserting mock podcast: %v", err)
	}
	episode := &protos.Episode{PodcastID: podcast.Id, Id: protos.ObjectIDFromHex("epi_id"), Author: "Sam Schwartz", Title: "Mock Episode"}
	err = dbClient.Insert(database.ColEpisode, episode)
	if err != nil {
		t.Fatalf("createAuthSerivceMockDB() error inserting mock episode: %v", err)
	}
	user := &protos.User{
		Id:       protos.ObjectIDFromHex("user_id"),
		Username: "user",
		Password: "$2a$04$Rxbh4f5cUjABPp2RE8o8PuvOafWNeYRsvYI/2t1lSL/DD/IYmWsfe",
		DOB:      ptypes.TimestampNow(),
		Email:    "user@example.com",
	}
	err = dbClient.Insert(database.ColUser, user)
	if err != nil {
		t.Fatalf("createAuthSerivceMockDB() error inserting mock user: %v", err)
	}
	err = dbClient.Insert(database.ColSession, &protos.Session{Id: protos.NewObjectID(), Expires: util.AddToTimestamp(ptypes.TimestampNow(), time.Hour), SessionKey: "secret", UserID: user.Id})
	if err != nil {
		t.Fatalf("createAuthSerivceMockDB() error inserting mock session: %v", err)
	}
	err = dbClient.Insert(database.ColUserEpisode, &protos.UserEpisode{
		Id: protos.ObjectIDFromHex("userepi_id"), EpisodeID: protos.ObjectIDFromHex("epi_id"),
		UserID: protos.ObjectIDFromHex("user_id"), PodcastID: protos.ObjectIDFromHex("pod_id")})
	if err != nil {
		t.Fatalf("createAuthSerivceMockDB() error inserting mock user episode: %v", err)
	}
	err = dbClient.Insert(database.ColSubscription, &protos.Subscription{
		Id:            protos.ObjectIDFromHex("sub_id"),
		UserID:        protos.ObjectIDFromHex("user_id"),
		PodcastID:     protos.ObjectIDFromHex("pod_id"),
		InProgressIDs: []*protos.ObjectID{protos.ObjectIDFromHex("epi_id")},
	})
	if err != nil {
		t.Fatalf("createAuthSerivceMockDB() error inserting mock subscription: %v", err)
	}
	return dbClient
}

func createMockPodcastClient(t *testing.T) (authClient protos.PodClient, cleanup func() error) {
	ctx := context.Background()
	conn, err := gogrpc.DialContext(ctx, "bufnet",
		gogrpc.WithContextDialer(bufDialer),
		gogrpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	client := protos.NewPodClient(conn)
	return client, conn.Close
}

func TestPodcastService(t *testing.T) {
	// setup mock database and mock server
	mockDB := createPodcastServiceMockDB(t)

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer(&config.Config{}, mockDB, NewAuthService(mockDB), NewPodcastService(mockDB))

	go func() {
		if err := s.Start(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// setup mock client used for gRPC requests
	podcastClient, cleanupFunc := createMockPodcastClient(t)
	defer func() {
		err := cleanupFunc()
		if err != nil {
			t.Fatalf("TestAuthService() error cleanupFunc: %v", err)
		}
	}()

	// go through tests
	testPodcastService_GetEpisodes(t, podcastClient)
	testPodcastService_GetUserEpisode(t, podcastClient)
	testPodcastService_UpdateUserEpisode(t, podcastClient)
	testPodcastService_GetSubscriptions(t, podcastClient)
	testPodcastService_GetUserLastPlayed(t, podcastClient)
}

func testPodcastService_GetEpisodes(t *testing.T, podClient protos.PodClient) {
	type args struct {
		ctx context.Context
		req *protos.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.Episodes
		wantErr bool
	}{
		{
			name: "GetEpisodes_invalid",
			args: args{
				ctx: metadata.AppendToOutgoingContext(context.Background(), "token", "invalid"),
				req: &protos.Request{PodcastID: protos.ObjectIDFromHex("pod_id"), Start: 0, End: 10},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "GetEpisodes_valid",
			args: args{
				ctx: metadata.AppendToOutgoingContext(context.Background(), "token", "secret"),
				req: &protos.Request{PodcastID: protos.ObjectIDFromHex("pod_id"), Start: 0, End: 10},
			},
			want:    &protos.Episodes{Episodes: []*protos.Episode{{Id: protos.ObjectIDFromHex("epi_id"), PodcastID: protos.ObjectIDFromHex("pod_id"), Title: "Mock Episode", Author: "Sam Schwartz"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := podClient.GetEpisodes(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodcastService.GetEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("PodcastService.GetEpisodes() = \n\t%v, want \n\t%v", got.String(), tt.want.String())
			}
		})
	}
}

func testPodcastService_GetUserEpisode(t *testing.T, podClient protos.PodClient) {
	type args struct {
		ctx context.Context
		req *protos.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.UserEpisode
		wantErr bool
	}{
		{
			name: "GetUserEpisode_valid",
			args: args{
				ctx: metadata.AppendToOutgoingContext(context.Background(), "token", "secret"),
				req: &protos.Request{
					EpisodeID: protos.ObjectIDFromHex("epi_id"),
					PodcastID: protos.ObjectIDFromHex("pod_id"),
				},
			},
			want: &protos.UserEpisode{
				Id: protos.ObjectIDFromHex("userepi_id"), EpisodeID: protos.ObjectIDFromHex("epi_id"),
				UserID: protos.ObjectIDFromHex("user_id"), PodcastID: protos.ObjectIDFromHex("pod_id"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := podClient.GetUserEpisode(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodcastService.GetUserEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("PodcastService.GetUserEpisode() = %v, want %v", got.String(), tt.want.String())
			}
		})
	}
}

func testPodcastService_UpdateUserEpisode(t *testing.T, podClient protos.PodClient) {
	type args struct {
		ctx context.Context
		req *protos.UserEpisodeReq
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.Response
		wantErr bool
	}{
		{
			name: "UpdateUserEpisode_valid",
			args: args{
				ctx: metadata.AppendToOutgoingContext(context.Background(), "token", "secret"),
				req: &protos.UserEpisodeReq{
					EpisodeID: protos.ObjectIDFromHex("epi_id"),
					Offset:    11111,
					Played:    true,
				},
			},
			want:    &protos.Response{Success: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := podClient.UpdateUserEpisode(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodcastService.UpdateUserEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("PodcastService.UpdateUserEpisode() = %v, want %v", got.String(), tt.want.String())
			}
		})
	}
}

func testPodcastService_GetSubscriptions(t *testing.T, podClient protos.PodClient) {
	type args struct {
		ctx context.Context
		req *protos.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.Subscriptions
		wantErr bool
	}{
		{
			name: "GetSubscriptions_valid",
			args: args{
				ctx: metadata.AppendToOutgoingContext(context.Background(), "token", "secret"),
				req: &protos.Request{
					PodcastID: protos.ObjectIDFromHex("pod_id"),
				},
			},
			want: &protos.Subscriptions{
				Subscriptions: []*protos.Subscription{{
					Id:            protos.ObjectIDFromHex("sub_id"),
					UserID:        protos.ObjectIDFromHex("user_id"),
					PodcastID:     protos.ObjectIDFromHex("pod_id"),
					InProgressIDs: []*protos.ObjectID{protos.ObjectIDFromHex("epi_id")},
				}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := podClient.GetSubscriptions(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodcastService.GetSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("PodcastService.GetSubscriptions() = %v, want %v", got.String(), tt.want.String())
			}
		})
	}
}

func testPodcastService_GetUserLastPlayed(t *testing.T, podClient protos.PodClient) {
	type args struct {
		ctx context.Context
		req *protos.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.LastPlayedRes
		wantErr bool
	}{
		{
			name: "GetUserLastPlayed_valid",
			args: args{
				ctx: metadata.AppendToOutgoingContext(context.Background(), "token", "secret"),
				req: &protos.Request{},
			},
			want:    &protos.LastPlayedRes{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := podClient.GetUserLastPlayed(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodcastService.GetUserLastPlayed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodcastService.GetUserLastPlayed() = %v, want %v", got, tt.want)
			}
		})
	}
}
