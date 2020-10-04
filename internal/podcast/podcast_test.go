package podcast

import (
	"reflect"
	"testing"

	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/minimongo/mock"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/protos"
)

func insertOrFail(t *testing.T, mockDB db.Database, collection string, object interface{}) {
	err := mockDB.Insert(collection, object)
	if err != nil {
		t.Fatalf("insertOrFail() error inserting: %v", err)
	}
}

func TestDoesPodcastExist(t *testing.T) {
	mockDB := mock.CreateDB()
	pod := &protos.Podcast{Author: "Sam Schwartz", Rss: "https://somevalidpodcast.com/url.rss"}
	insertOrFail(t, mockDB, database.ColPodcast, pod)

	type args struct {
		dbClient db.Database
		rssURL   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid",
			args: args{
				dbClient: mockDB,
				rssURL:   "https://somevalidpodcast.com/url.rss",
			},
			want: true,
		},
		{
			name: "Invalid",
			args: args{
				dbClient: mockDB,
				rssURL:   "https://someINvalidpodcast.com/url.rss",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := DoesPodcastExist(tt.args.dbClient, tt.args.rssURL)
			if got != tt.want {
				t.Errorf("DoesPodcastExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindPodcastsByRange(t *testing.T) {
	type args struct {
		dbClient db.Database
		start    int
		end      int
	}
	tests := []struct {
		name    string
		args    args
		want    []*protos.Podcast
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FindPodcastsByRange(tt.args.dbClient, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPodcastsByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPodcastsByRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindPodcastByID(t *testing.T) {
	type args struct {
		dbClient db.Database
		id       *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.Podcast
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FindPodcastByID(tt.args.dbClient, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPodcastByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPodcastByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindUserEpisode(t *testing.T) {
	type args struct {
		dbClient db.Database
		userID   *protos.ObjectID
		epiID    *protos.ObjectID
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FindUserEpisode(tt.args.dbClient, tt.args.userID, tt.args.epiID)
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

func TestSearchPodcasts(t *testing.T) {
	type args struct {
		dbClient db.Database
		search   string
	}
	tests := []struct {
		name    string
		args    args
		want    []*protos.Podcast
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := SearchPodcasts(tt.args.dbClient, tt.args.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchPodcasts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchPodcasts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchTitle(t *testing.T) {
	type args struct {
		search   string
		podcasts []protos.Podcast
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			MatchTitle(tt.args.search, tt.args.podcasts)
		})
	}
}

func TestFindLength(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := FindLength(tt.args.url); got != tt.want {
				t.Errorf("FindLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
