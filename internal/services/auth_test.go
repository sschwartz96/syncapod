package services

import (
	"context"
	"log"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/minimongo/mock"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	// cfgFile, err := os.Open("config.json")
	// if err != nil {
	// log.Fatalf("init() error opening config file: %v", err)
	// }
	// cfg, err := config.ReadConfig(cfgFile)
	// if err != nil {
	// log.Fatalf("init() error reading config file: %v", err)
	// }
	//

	mockDB := mock.CreateDB()
	err := initDB(mockDB)

	lis = bufconn.Listen(bufSize)
	s := gogrpc.NewServer()
	protos.RegisterAuthService(s, protos.NewAuthService(NewAuthService(mockDB)))

	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func initDB(dbClient db.Database) error {
	err := dbClient.Insert(database.ColUser, &protos.User{
		Id:       protos.NewObjectID(),
		Username: "user",
		Password: "$2a$04$Rxbh4f5cUjABPp2RE8o8PuvOafWNeYRsvYI/2t1lSL/DD/IYmWsfe",
		DOB:      ptypes.TimestampNow(),
		Email:    "user@example.com",
	})
	if err != nil {
		return err
	}
	err := dbClient.Insert(database.ColSession, &protos.Session{Id: protos.NewObjectID(), Expires: time.Now().Add(time.Hour), SessionKey: "secret"})
	return err
}

func TestAuthenticate_OLD(t *testing.T) {
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

func setupMockAuthClient() (authClient *protos.AuthClient, cleanup func() error) {
	ctx := context.Background()
	conn, err := gogrpc.DialContext(ctx, "bufnet",
		gogrpc.WithContextDialer(bufDialer),
		gogrpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	client := protos.NewAuthClient(conn)
	return &client, conn.Close
}

func TestAuthService_Authenticate(t *testing.T) {
	authClient, cleanupFunc := setupMockAuthClient()
	defer cleanupFunc()
	type args struct {
		ctx context.Context
		req *protos.AuthReq
	}
	tests := []struct {
		name    string
		client  *protos.AuthClient
		args    args
		want    *protos.AuthRes
		wantErr bool
	}{
		{
			name:    "invalid",
			args:    args{ctx: context.Background(), req: &protos.AuthReq{Username: "user", Password: "wrong"}},
			client:  authClient,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Authenticate(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Authorize(t *testing.T) {
	type args struct {
		ctx context.Context
		req *protos.AuthReq
	}
	tests := []struct {
		name    string
		a       *AuthService
		args    args
		want    *protos.AuthRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Authorize(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	type args struct {
		ctx context.Context
		req *protos.AuthReq
	}
	tests := []struct {
		name    string
		a       *AuthService
		args    args
		want    *protos.AuthRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Logout(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Logout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.Logout() = %v, want %v", got, tt.want)
			}
		})
	}
}
