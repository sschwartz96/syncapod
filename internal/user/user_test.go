package user

import (
	"log"
	"reflect"
	"testing"

	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/minimongo/mock"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
)

func createMockDBSession(t *testing.T) (*protos.Session, *mock.DB) {
	mockDB := mock.CreateDB()

	initial := &protos.Session{
		Id:         protos.ObjectIDFromHex("id_1"),
		SessionKey: "key_1",
	}
	err := mockDB.Insert(
		database.ColSession,
		initial,
	)

	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	return initial, mockDB
}

func TestFindSession(t *testing.T) {
	initial, mockDB := createMockDBSession(t)

	type args struct {
		dbClient db.Database
		key      string
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.Session
		wantErr bool
	}{
		{
			name: "FindSession_0",
			args: args{
				dbClient: mockDB,
				key:      "key_1",
			},
			want:    initial,
			wantErr: false,
		}, {
			name: "FindSession_1",
			args: args{
				dbClient: mockDB,
				key:      "not_found",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindSession(tt.args.dbClient, tt.args.key)
			if (err != nil) != tt.wantErr {
				log.Println("db:", tt.args.dbClient)
				t.Errorf("FindSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpsertSession(t *testing.T) {
	_, mockDB := createMockDBSession(t)

	type args struct {
		dbClient db.Database
		session  *protos.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update key to: keychanged",
			args: args{
				mockDB,
				&protos.Session{
					Id:         protos.ObjectIDFromHex("id_1"),
					SessionKey: "keychanged",
				},
			},
			wantErr: false,
		}, {
			name: "add new session",
			args: args{
				mockDB,
				&protos.Session{
					Id:         protos.ObjectIDFromHex("new_id"),
					SessionKey: "new_key",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpsertSession(tt.args.dbClient, tt.args.session); (err != nil) != tt.wantErr {
				t.Errorf("UpsertSession() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				var session protos.Session
				err := tt.args.dbClient.FindOne(database.ColSession, &session, &db.Filter{"sessionkey": tt.args.session.SessionKey}, db.CreateOptions())
				if err != nil {
					t.Errorf("failed to upsert session not found: %v", err)
				}
			}
		})
	}
}

func TestDeleteSession(t *testing.T) {
	initial, mockDB := createMockDBSession(t)

	type args struct {
		dbClient db.Database
		id       *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "not_found",
			args: args{
				dbClient: mockDB,
				id:       protos.ObjectIDFromHex("not_valid_id"),
			},
			wantErr: true,
		},
		{
			name: "delete",
			args: args{
				dbClient: mockDB,
				id:       initial.Id,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Println("dbClient:", tt.args.dbClient)
			if err := DeleteSession(tt.args.dbClient, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteSessionByKey(t *testing.T) {
	type args struct {
		dbClient db.Database
		key      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteSessionByKey(tt.args.dbClient, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSessionByKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindUserByID(t *testing.T) {
	type args struct {
		dbClient db.Database
		id       *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindUserByID(tt.args.dbClient, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindUser(t *testing.T) {
	type args struct {
		dbClient db.Database
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindUser(tt.args.dbClient, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type args struct {
		dbClient db.Database
		id       *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteUser(tt.args.dbClient, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindUserEpisode(t *testing.T) {
	type args struct {
		dbClient  db.Database
		userID    *protos.ObjectID
		episodeID *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.UserEpisode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindUserEpisode(tt.args.dbClient, tt.args.userID, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserEpisode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindLatestUserEpisode(t *testing.T) {
	type args struct {
		dbClient db.Database
		userID   *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.UserEpisode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindLatestUserEpisode(tt.args.dbClient, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindLatestUserEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindLatestUserEpisode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpsertUserEpisode(t *testing.T) {
	type args struct {
		dbClient    db.Database
		userEpisode *protos.UserEpisode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpsertUserEpisode(tt.args.dbClient, tt.args.userEpisode); (err != nil) != tt.wantErr {
				t.Errorf("UpsertUserEpisode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindSubscriptions(t *testing.T) {
	type args struct {
		dbClient db.Database
		userID   *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		want    []*protos.Subscription
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindSubscriptions(tt.args.dbClient, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindSubscriptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpsertSubscription(t *testing.T) {
	type args struct {
		dbClient     db.Database
		subscription *protos.Subscription
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpsertSubscription(tt.args.dbClient, tt.args.subscription); (err != nil) != tt.wantErr {
				t.Errorf("UpsertSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindUserLastPlayed(t *testing.T) {
	type args struct {
		dbClient db.Database
		userID   *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.Podcast
		want1   *protos.Episode
		want2   *protos.UserEpisode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := FindUserLastPlayed(tt.args.dbClient, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserLastPlayed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserLastPlayed() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FindUserLastPlayed() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("FindUserLastPlayed() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestFindOffset(t *testing.T) {
	type args struct {
		dbClient db.Database
		userID   *protos.ObjectID
		epiID    *protos.ObjectID
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindOffset(tt.args.dbClient, tt.args.userID, tt.args.epiID); got != tt.want {
				t.Errorf("FindOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateOffset(t *testing.T) {
	type args struct {
		dbClient db.Database
		uID      *protos.ObjectID
		pID      *protos.ObjectID
		eID      *protos.ObjectID
		offset   int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateOffset(tt.args.dbClient, tt.args.uID, tt.args.pID, tt.args.eID, tt.args.offset); (err != nil) != tt.wantErr {
				t.Errorf("UpdateOffset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateUserEpiPlayed(t *testing.T) {
	type args struct {
		dbClient db.Database
		uID      *protos.ObjectID
		pID      *protos.ObjectID
		eID      *protos.ObjectID
		played   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateUserEpiPlayed(tt.args.dbClient, tt.args.uID, tt.args.pID, tt.args.eID, tt.args.played); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserEpiPlayed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
