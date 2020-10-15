package services

import (
	"context"
	"log"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/stockpile/mock"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/util"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(ctx context.Context, s string) (net.Conn, error) {
	log.Println("bufDialer ctx", ctx)
	return lis.Dial()
}

// createAuthServiceMockDB fails on error and returns db.Database and *protos.User
func createAuthServiceMockDB(t *testing.T) db.Database {
	dbClient := mock.CreateDB()
	user := &protos.User{
		Id:       protos.NewObjectID(),
		Username: "user",
		Password: "$2a$04$Rxbh4f5cUjABPp2RE8o8PuvOafWNeYRsvYI/2t1lSL/DD/IYmWsfe",
		DOB:      ptypes.TimestampNow(),
		Email:    "user@example.com",
	}
	err := dbClient.Insert(database.ColUser, user)
	if err != nil {
		t.Fatalf("createAuthSerivceMockDB() error inserting mock user: %v", err)
	}
	err = dbClient.Insert(database.ColSession, &protos.Session{Id: protos.NewObjectID(), Expires: util.AddToTimestamp(ptypes.TimestampNow(), time.Hour), SessionKey: "secret", UserID: user.Id})
	if err != nil {
		t.Fatalf("createAuthSerivceMockDB() error inserting mock session: %v", err)
	}
	return dbClient
}

func createMockAuthClient(t *testing.T) (authClient protos.AuthClient, cleanup func() error) {
	ctx := context.Background()
	conn, err := gogrpc.DialContext(ctx, "bufnet",
		gogrpc.WithContextDialer(bufDialer),
		gogrpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	client := protos.NewAuthClient(conn)
	return client, conn.Close
}

func TestAuthService(t *testing.T) {
	// setup mock database and mock server
	mockDB := createAuthServiceMockDB(t)

	lis = bufconn.Listen(bufSize)
	s := gogrpc.NewServer()
	protos.RegisterAuthServer(s, NewAuthService(mockDB))

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// setup mock client used for gRPC requests
	authClient, cleanupFunc := createMockAuthClient(t)
	defer func() {
		err := cleanupFunc()
		if err != nil {
			t.Fatalf("TestAuthService() error cleanupFunc: %v", err)
		}
	}()

	// go through tests
	testAuthService_Authenticate(t, authClient)
	testAuthService_Authorize(t, authClient)
	testAuthService_Logout(t, authClient)
}

func testAuthService_Authenticate(t *testing.T, authClient protos.AuthClient) {
	type args struct {
		ctx context.Context
		req *protos.AuthReq
	}
	tests := []struct {
		name    string
		client  protos.AuthClient
		args    args
		want    *protos.AuthRes
		wantErr bool
	}{
		{
			name:    "authenticate_invalid",
			args:    args{ctx: context.Background(), req: &protos.AuthReq{Username: "user", Password: "123wrong"}},
			client:  authClient,
			want:    &protos.AuthRes{Success: false},
			wantErr: false,
		},
		{
			name:    "authenticate_valid",
			args:    args{ctx: context.Background(), req: &protos.AuthReq{Username: "user", Password: "password"}},
			client:  authClient,
			want:    &protos.AuthRes{Success: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.Authenticate(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Success != tt.want.Success {
				t.Errorf("AuthService.Authenticate() = %v, want %v", got.String(), tt.want.String())
			}
		})
	}
}

func testAuthService_Authorize(t *testing.T, authClient protos.AuthClient) {
	type args struct {
		ctx context.Context
		req *protos.AuthReq
	}
	tests := []struct {
		name       string
		authClient protos.AuthClient
		args       args
		want       *protos.AuthRes
		wantErr    bool
	}{
		{
			name:       "authorize_invalid",
			authClient: authClient,
			args: args{
				ctx: context.Background(),
				req: &protos.AuthReq{
					SessionKey: "not_correct_secret",
				},
			},
			want:    &protos.AuthRes{Success: false},
			wantErr: false,
		},
		{
			name:       "authorize_valid",
			authClient: authClient,
			args: args{
				ctx: context.Background(),
				req: &protos.AuthReq{
					SessionKey: "secret",
				},
			},
			want:    &protos.AuthRes{Success: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.authClient.Authorize(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Success != tt.want.Success {
				t.Errorf("AuthService.Authorize() = %v, want %v", got.Success, tt.want.Success)
			}
		})
	}
}

func testAuthService_Logout(t *testing.T, authClient protos.AuthClient) {
	type args struct {
		ctx context.Context
		req *protos.AuthReq
	}
	tests := []struct {
		name        string
		authClient  protos.AuthClient
		args        args
		wantSuccess bool
	}{
		{
			name: "logout_invalid",
			args: args{
				ctx: context.Background(),
				req: &protos.AuthReq{SessionKey: "invalid_key"},
			},
			authClient:  authClient,
			wantSuccess: false,
		},
		{
			name: "logout_valid",
			args: args{
				ctx: context.Background(),
				req: &protos.AuthReq{SessionKey: "secret"},
			},
			authClient:  authClient,
			wantSuccess: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.authClient.Logout(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(got.Success, tt.wantSuccess) {
				t.Errorf("AuthService.Logout() = %v, want %v", got, tt.wantSuccess)
			}
		})
	}
}

func testAuthenticate_OLD(t *testing.T) {
	ctx := context.Background()
	conn, err := gogrpc.DialContext(ctx, "bufnet",
		gogrpc.WithContextDialer(bufDialer),
		gogrpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := protos.NewAuthClient(conn)
	resp, err := client.Authenticate(ctx,
		&protos.AuthReq{Username: "user", Password: "password"},
	)
	if err != nil {
		t.Fatalf("TestAuthenticate() response error: %v", err)
	}
	log.Println("response:", resp)
}
