package podcast

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/stockpile/mock"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/util"
)

func epiSliceEqual(sliceI, sliceJ interface{}) bool {
	valI := derefencedValue(reflect.ValueOf(sliceI))
	valJ := derefencedValue(reflect.ValueOf(sliceJ))
	if valI.Kind() != reflect.Slice || valJ.Kind() != reflect.Slice {
		return false
	}
	if valI.Len() != valJ.Len() {
		return false
	}
	for i := 0; i < valI.Len(); i++ {
		vI := derefencedValue(valI.Index(i))
		vJ := derefencedValue(valJ.Index(i))
		if strings.EqualFold(vI.Interface().(protos.Episode).Id.Hex, vJ.Interface().(protos.Episode).Id.Hex) {
			return false
		}
	}
	return true
}

func derefencedValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func TestFindEpisodesByRange(t *testing.T) {
	mockDB := mock.CreateDB()
	podId := protos.ObjectIDFromHex("podID")
	epi1 := protos.Episode{Id: protos.ObjectIDFromHex("obj1"), PodcastID: podId}
	epi2 := protos.Episode{Id: protos.ObjectIDFromHex("obj2"), PodcastID: podId}
	epi3 := protos.Episode{Id: protos.ObjectIDFromHex("obj3"), PodcastID: podId}
	epi4 := protos.Episode{Id: protos.ObjectIDFromHex("obj4"), PodcastID: podId}
	insertOrFail(t, mockDB, database.ColEpisode, epi1)
	insertOrFail(t, mockDB, database.ColEpisode, epi2)
	insertOrFail(t, mockDB, database.ColEpisode, epi3)
	insertOrFail(t, mockDB, database.ColEpisode, epi4)

	type args struct {
		dbClient  db.Database
		podcastID *protos.ObjectID
		start     int64
		end       int64
	}
	tests := []struct {
		name    string
		args    args
		want    []*protos.Episode
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				dbClient:  mockDB,
				start:     0,
				end:       3,
				podcastID: podId,
			},
			want:    []*protos.Episode{&epi1, &epi2, &epi3},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				dbClient:  mockDB,
				start:     1,
				end:       5,
				podcastID: podId,
			},
			want:    []*protos.Episode{&epi2, &epi3, &epi4},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindEpisodesByRange(tt.args.dbClient, tt.args.podcastID, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindEpisodesByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if epiSliceEqual(got, tt.want) {
				t.Errorf("FindEpisodesByRange() = \n\t%v, want \n\t%v", got, tt.want)
			}
		})
	}
}

func TestFindLatestEpisode(t *testing.T) {
	mockDB := mock.CreateDB()
	podId := protos.NewObjectID()
	epi1 := &protos.Episode{Id: protos.ObjectIDFromHex("obj1"), PodcastID: podId, PubDate: util.AddToTimestamp(ptypes.TimestampNow(), time.Minute*-1)}
	epi2 := &protos.Episode{Id: protos.ObjectIDFromHex("obj2"), PodcastID: podId, PubDate: ptypes.TimestampNow()}
	insertOrFail(t, mockDB, database.ColEpisode, epi1)
	insertOrFail(t, mockDB, database.ColEpisode, epi2)

	type args struct {
		dbClient  db.Database
		podcastID *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.Episode
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				dbClient:  mockDB,
				podcastID: podId,
			},
			want:    epi2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindLatestEpisode(tt.args.dbClient, tt.args.podcastID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindLatestEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindLatestEpisode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindEpisodeByID(t *testing.T) {
	mockDB := mock.CreateDB()
	podId := protos.NewObjectID()
	epi1 := &protos.Episode{Id: protos.ObjectIDFromHex("obj1"), PodcastID: podId}
	epi2 := &protos.Episode{Id: protos.ObjectIDFromHex("obj2"), PodcastID: podId}
	insertOrFail(t, mockDB, database.ColEpisode, epi1)
	insertOrFail(t, mockDB, database.ColEpisode, epi2)

	type args struct {
		dbClient db.Database
		id       *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.Episode
		wantErr bool
	}{
		{
			name: "invalid",
			args: args{
				dbClient: mockDB,
				id:       protos.ObjectIDFromHex("invalid"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid",
			args: args{
				dbClient: mockDB,
				id:       epi2.Id,
			},
			want:    epi2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindEpisodeByID(tt.args.dbClient, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindEpisodeByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindEpisodeByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindEpisodeBySeason(t *testing.T) {
	mockDB := mock.CreateDB()
	podId := protos.NewObjectID()
	epi1 := &protos.Episode{Id: protos.ObjectIDFromHex("obj1"), PodcastID: podId, Season: 1, Episode: 1}
	epi2 := &protos.Episode{Id: protos.ObjectIDFromHex("obj2"), PodcastID: podId, Season: 1, Episode: 2}
	insertOrFail(t, mockDB, database.ColEpisode, epi1)
	insertOrFail(t, mockDB, database.ColEpisode, epi2)

	type args struct {
		dbClient   db.Database
		id         *protos.ObjectID
		seasonNum  int
		episodeNum int
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.Episode
		wantErr bool
	}{
		{
			name: "1.invalid episode",
			args: args{
				dbClient:   mockDB,
				seasonNum:  1,
				episodeNum: 3,
				id:         podId,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "2.invalid podcast",
			args: args{
				dbClient:   mockDB,
				seasonNum:  1,
				episodeNum: 1,
				id:         protos.ObjectIDFromHex("invalid podcast id"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "3.valid episode 1",
			args: args{
				dbClient:   mockDB,
				seasonNum:  1,
				episodeNum: 1,
				id:         podId,
			},
			want:    epi1,
			wantErr: false,
		},
		{
			name: "3.valid episode 2",
			args: args{
				dbClient:   mockDB,
				seasonNum:  1,
				episodeNum: 2,
				id:         podId,
			},
			want:    epi2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindEpisodeBySeason(tt.args.dbClient, tt.args.id, tt.args.seasonNum, tt.args.episodeNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindEpisodeBySeason() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindEpisodeBySeason() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpsertEpisode(t *testing.T) {
	mockDB := mock.CreateDB()
	podId := protos.NewObjectID()
	epi1 := &protos.Episode{Id: protos.ObjectIDFromHex("obj1"), PodcastID: podId, Author: "Sam Schwartz"}
	insertOrFail(t, mockDB, database.ColEpisode, epi1)
	epi1Copy := *epi1
	epi1Copy.Author = "Simon Schwartz"
	type args struct {
		dbClient db.Database
		episode  *protos.Episode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1_insert",
			args: args{
				dbClient: mockDB,
				episode:  &protos.Episode{Author: "Test Author"},
			},
			wantErr: false,
		},
		{
			name: "1_update",
			args: args{
				dbClient: mockDB,
				episode:  &epi1Copy,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpsertEpisode(tt.args.dbClient, tt.args.episode); (err != nil) != tt.wantErr {
				t.Errorf("UpsertEpisode() error = %v, wantErr %v", err, tt.wantErr)
			}
			var found protos.Episode
			err := tt.args.dbClient.FindOne(database.ColEpisode, &found, &db.Filter{"author": tt.args.episode.Author}, nil)
			if err != nil {
				t.Errorf("UpsertEpisode() error = %v, could not find upserted episode: %v", err, tt.args.episode)
			}
		})
	}
}

func TestDoesEpisodeExist(t *testing.T) {
	mockDB := mock.CreateDB()
	podId := protos.NewObjectID()
	pubDate := ptypes.TimestampNow()
	epi1 := &protos.Episode{Id: protos.ObjectIDFromHex("obj1"), PodcastID: podId, Title: "Cool Episode Title", PubDate: pubDate}
	insertOrFail(t, mockDB, database.ColEpisode, epi1)

	type args struct {
		dbClient db.Database
		title    string
		pubDate  *timestamp.Timestamp
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "1_not_exists",
			args: args{
				dbClient: mockDB,
				pubDate:  ptypes.TimestampNow(),
				title:    "Cool",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "2_exists",
			args: args{
				dbClient: mockDB,
				pubDate:  pubDate,
				title:    "Cool Episode Title",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoesEpisodeExist(tt.args.dbClient, tt.args.title, tt.args.pubDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoesEpisodeExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DoesEpisodeExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
