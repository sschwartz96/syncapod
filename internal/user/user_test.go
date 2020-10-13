package user

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/stockpile/mock"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createMockDBSession(t *testing.T) (*protos.Session, *mock.DB) {
	mockDB := mock.CreateDB()

	initial := &protos.Session{
		Id:         protos.ObjectIDFromHex("id_1"),
		SessionKey: "key_1",
	}
	insertOrFail(t, mockDB, database.ColSession, initial)

	insertOrFail(t, mockDB, database.ColSession, &protos.Session{
		Id:         protos.ObjectIDFromHex("id_2"),
		SessionKey: "key_2",
	})

	return initial, mockDB
}

func insertOrFail(t *testing.T, db *mock.DB, col string, obj interface{}) {
	err := db.Insert(col, obj)
	if err != nil {
		t.Fatalf("could not insert object into mockDB: %v", err)
	}
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
			if err := DeleteSession(tt.args.dbClient, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteSessionByKey(t *testing.T) {
	initial, mockDB := createMockDBSession(t)

	type args struct {
		dbClient db.Database
		key      string
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
				key:      "not_valid_key",
			},
			wantErr: true,
		},
		{
			name: "delete",
			args: args{
				dbClient: mockDB,
				key:      initial.SessionKey,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteSessionByKey(tt.args.dbClient, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSessionByKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createMockDBUser(t *testing.T) (*protos.User, *mock.DB) {
	mockDB := mock.CreateDB()

	initial := &protos.User{
		Id:       protos.ObjectIDFromHex("id_1"),
		DOB:      timestamppb.Now(),
		Email:    "test@test.org",
		Password: "@somehashedpass",
		Username: "tester",
	}
	insertOrFail(t, mockDB, database.ColUser, initial)

	insertOrFail(t, mockDB, database.ColUser, &protos.User{
		Id:       protos.ObjectIDFromHex("id_2"),
		DOB:      timestamppb.Now(),
		Email:    "test2@test.org",
		Password: "@somehashedpass2",
		Username: "tester2",
	})

	return initial, mockDB
}
func TestFindUserByID(t *testing.T) {
	initial, mockDB := createMockDBUser(t)
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
		{
			name: "not found",
			args: args{
				dbClient: mockDB,
				id:       protos.ObjectIDFromHex("not_found"),
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "find",
			args: args{
				dbClient: mockDB,
				id:       initial.Id,
			},
			want:    initial,
			wantErr: false,
		},
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
	initial, mockDB := createMockDBUser(t)
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
		{
			name: "not found",
			args: args{
				dbClient: mockDB,
				username: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "find user by username(email)",
			args: args{
				dbClient: mockDB,
				username: initial.Email,
			},
			want:    initial,
			wantErr: false,
		},
		{
			name: "find user by username",
			args: args{
				dbClient: mockDB,
				username: initial.Username,
			},
			want:    initial,
			wantErr: false,
		},
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
	initial, mockDB := createMockDBUser(t)
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
			name: "not found",
			args: args{
				dbClient: mockDB,
				id:       protos.ObjectIDFromHex(""),
			},
			wantErr: true,
		},
		{
			name: "find user by username(email)",
			args: args{
				dbClient: mockDB,
				id:       initial.Id,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteUser(tt.args.dbClient, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createMockDBUserEpi(t *testing.T) (*protos.UserEpisode, *mock.DB) {
	mockDB := mock.CreateDB()

	initial := &protos.UserEpisode{
		Id:        protos.ObjectIDFromHex("id1"),
		EpisodeID: protos.ObjectIDFromHex("epi_id1"),
		PodcastID: protos.ObjectIDFromHex("pod_id1"),
		UserID:    protos.ObjectIDFromHex("user_id1"),
		LastSeen:  timestamppb.Now(),
		Offset:    123456,
		Played:    false,
	}
	insertOrFail(t, mockDB, database.ColUserEpisode, initial)

	insertOrFail(t, mockDB, database.ColUserEpisode, &protos.UserEpisode{

		Id:        protos.ObjectIDFromHex("id2"),
		EpisodeID: protos.ObjectIDFromHex("epi_id2"),
		PodcastID: protos.ObjectIDFromHex("pod_id2"),
		UserID:    protos.ObjectIDFromHex("user_id2"),
		LastSeen:  timestamppb.New(time.Now().Add(-time.Minute)),
		Offset:    0,
		Played:    true,
	})

	return initial, mockDB
}

func TestFindUserEpisode(t *testing.T) {
	initial, mockDB := createMockDBUserEpi(t)
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
		{
			name: "not_found",
			args: args{
				dbClient:  mockDB,
				episodeID: protos.ObjectIDFromHex(""),
				userID:    protos.ObjectIDFromHex(""),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "find",
			args: args{
				dbClient:  mockDB,
				episodeID: initial.EpisodeID,
				userID:    initial.UserID,
			},
			want:    initial,
			wantErr: false,
		},
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
	initial, mockDB := createMockDBUserEpi(t)
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
		{
			args: args{
				dbClient: mockDB,
				userID:   protos.ObjectIDFromHex("not_found"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "found",
			args: args{
				dbClient: mockDB,
				userID:   initial.UserID,
			},
			want:    initial,
			wantErr: false,
		},
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
	initial, mockDB := createMockDBUserEpi(t)
	type args struct {
		dbClient    db.Database
		userEpisode *protos.UserEpisode
	}
	want1 := &protos.UserEpisode{Id: initial.Id, EpisodeID: initial.EpisodeID, PodcastID: initial.PodcastID, UserID: initial.UserID, LastSeen: initial.LastSeen, Played: initial.Played, Offset: 8998}
	want2 := &protos.UserEpisode{Id: protos.ObjectIDFromHex("new_id"), EpisodeID: protos.ObjectIDFromHex("new_episode_id"), PodcastID: initial.PodcastID, UserID: initial.UserID, LastSeen: initial.LastSeen, Played: true, Offset: 0}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *protos.UserEpisode
	}{
		{
			name: "update",
			args: args{
				dbClient:    mockDB,
				userEpisode: want1,
			},
			wantErr: false,
			want:    want1,
		},
		{
			name: "insert",
			args: args{
				dbClient:    mockDB,
				userEpisode: want2,
			},
			wantErr: false,
			want:    want2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpsertUserEpisode(tt.args.dbClient, tt.args.userEpisode); (err != nil) != tt.wantErr {
				t.Errorf("UpsertUserEpisode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				var wanted protos.UserEpisode
				err := tt.args.dbClient.FindOne(database.ColUserEpisode, &wanted, &db.Filter{"offset": tt.want.Offset}, db.CreateOptions())
				if err != nil {
					t.Errorf("UpsertUserEpisode() error = upserted episode not found")
				}
			}
		})
	}
}

func createMockDBSub(t *testing.T) (*protos.Subscription, *mock.DB) {
	mockDB := mock.CreateDB()

	initial := &protos.Subscription{
		Id:        protos.ObjectIDFromHex("id_1"),
		UserID:    protos.ObjectIDFromHex("user_id_1"),
		PodcastID: protos.ObjectIDFromHex("pod_id_1"),
	}
	insertOrFail(t, mockDB, database.ColSubscription, initial)

	insertOrFail(t, mockDB, database.ColSubscription, &protos.Subscription{
		Id:        protos.ObjectIDFromHex("id_2"),
		UserID:    protos.ObjectIDFromHex("user_id_2"),
		PodcastID: protos.ObjectIDFromHex("pod_id_2"),
	})

	return initial, mockDB
}

func TestFindSubscriptions(t *testing.T) {
	initial, mockDB := createMockDBSub(t)
	var empty []*protos.Subscription
	found := []*protos.Subscription{initial}
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
		{
			name: "empty_sub",
			args: args{
				dbClient: mockDB,
				userID:   protos.ObjectIDFromHex("empty_sub"),
			},
			want:    empty,
			wantErr: false,
		},
		{
			name: "found",
			args: args{
				dbClient: mockDB,
				userID:   initial.UserID,
			},
			want:    found,
			wantErr: false,
		},
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
	_, mockDB := createMockDBSub(t)
	update := &protos.Subscription{Id: protos.ObjectIDFromHex("id_1"), UserID: protos.ObjectIDFromHex("user_id_1"), PodcastID: protos.ObjectIDFromHex("pod_id_1"), CompletedIDs: []*protos.ObjectID{protos.NewObjectID()}}
	insert := &protos.Subscription{Id: protos.ObjectIDFromHex("id_5"), UserID: protos.ObjectIDFromHex("user_id_1"), PodcastID: protos.ObjectIDFromHex("pod_id_9"), CompletedIDs: []*protos.ObjectID{protos.NewObjectID()}}

	type args struct {
		dbClient     db.Database
		subscription *protos.Subscription
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update",
			args: args{
				dbClient:     mockDB,
				subscription: update,
			},
			wantErr: false,
		},
		{
			name: "insert",
			args: args{
				dbClient:     mockDB,
				subscription: insert,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpsertSubscription(tt.args.dbClient, tt.args.subscription); (err != nil) != tt.wantErr {
				t.Errorf("UpsertSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
			var upserted protos.Subscription
			err := mockDB.FindOne(database.ColSubscription, &upserted, &db.Filter{"_id": tt.args.subscription.Id}, db.CreateOptions())
			if err != nil {
				t.Errorf("UpsertSubscription() error = %v", err)
			}
			if !reflect.DeepEqual(tt.args.subscription, &upserted) {
				t.Error("UpsertSubscription() upserted subscription is not equal error")
			}
		})
	}
}

func createLastPlayedDB(t *testing.T) (*protos.User, *protos.Podcast, *protos.Episode, *protos.UserEpisode, *mock.DB) {
	mockDB := mock.CreateDB()

	user := &protos.User{Id: protos.NewObjectID(), Username: "testUser"}
	pod := &protos.Podcast{Id: protos.NewObjectID(), Author: "Podcast Author"}
	epi := &protos.Episode{Id: protos.NewObjectID(), Author: "Podcast Author"}
	userEpi := &protos.UserEpisode{Id: protos.NewObjectID(), UserID: user.Id, EpisodeID: epi.Id, PodcastID: pod.Id, Offset: 123456, Played: false}

	insertOrFail(t, mockDB, database.ColUser, user)
	insertOrFail(t, mockDB, database.ColPodcast, pod)
	insertOrFail(t, mockDB, database.ColEpisode, epi)
	insertOrFail(t, mockDB, database.ColUserEpisode, userEpi)
	return user, pod, epi, userEpi, mockDB
}

func TestFindUserLastPlayed(t *testing.T) {
	user, pod, epi, userEpi, mockDB := createLastPlayedDB(t)
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
		{
			name: "test1",
			args: args{
				dbClient: mockDB,
				userID:   user.Id,
			},
			want:    pod,
			want1:   epi,
			want2:   userEpi,
			wantErr: false,
		},
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
	user, _, epi, userEpi, mockDB := createLastPlayedDB(t)
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
		{
			name: "one",
			args: args{
				dbClient: mockDB,
				epiID:    epi.Id,
				userID:   user.Id,
			},
			want: userEpi.Offset,
		},
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
	user, pod, epi, _, mockDB := createLastPlayedDB(t)
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
		{
			name: "update",
			args: args{
				dbClient: mockDB,
				eID:      epi.Id,
				offset:   654321,
				pID:      pod.Id,
				uID:      user.Id,
			},
			wantErr: false,
		},
		{
			name: "insert",
			args: args{
				dbClient: mockDB,
				eID:      epi.Id,
				offset:   789123,
				pID:      pod.Id,
				uID:      protos.ObjectIDFromHex("differentUserID"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateOffset(tt.args.dbClient, tt.args.uID, tt.args.pID, tt.args.eID, tt.args.offset); (err != nil) != tt.wantErr {
				t.Errorf("UpdateOffset() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				found := &protos.UserEpisode{}
				err := mockDB.FindOne(database.ColUserEpisode, found, &db.Filter{"offset": tt.args.offset}, db.CreateOptions())
				if err != nil {
					t.Error("UpdateOffset() error = could not updated user episode")
				}
			}
		})
	}
}

func TestUpdateUserEpiPlayed(t *testing.T) {
	user, pod, epi, _, mockDB := createLastPlayedDB(t)
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
		{
			name: "update",
			args: args{
				dbClient: mockDB,
				eID:      epi.Id,
				pID:      pod.Id,
				played:   true,
				uID:      user.Id,
			},
			wantErr: false,
		},
		{
			name: "insert",
			args: args{
				dbClient: mockDB,
				eID:      epi.Id,
				pID:      pod.Id,
				played:   false,
				uID:      protos.ObjectIDFromHex("objID2"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateUserEpiPlayed(tt.args.dbClient, tt.args.uID, tt.args.pID, tt.args.eID, tt.args.played); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserEpiPlayed() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				found := &protos.UserEpisode{}
				err := mockDB.FindOne(database.ColUserEpisode, found, &db.Filter{"played": tt.args.played}, db.CreateOptions())
				if err != nil {
					t.Error("UpdateUserEpiPlayed() error = could not find updated user episode")
				}
			}
		})
	}
}
