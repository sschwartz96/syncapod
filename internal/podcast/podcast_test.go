package podcast

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/stockpile/mock"
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
	mockDB := mock.CreateDB()
	pod1 := &protos.Podcast{Id: protos.NewObjectID(), Author: "Sam Schwartz", Rss: "https://valid.com/rss"}
	pod2 := &protos.Podcast{Id: protos.NewObjectID(), Author: "Simon Schwartz", Rss: "https://cool.com/rss"}
	pod3 := &protos.Podcast{Id: protos.NewObjectID(), Author: "Joe Rogan", Rss: "https://joerogan.com/rss"}
	pod4 := &protos.Podcast{Id: protos.NewObjectID(), Author: "Conan O'Brien", Rss: "https://conan.com/rss"}
	insertOrFail(t, mockDB, database.ColPodcast, pod1)
	insertOrFail(t, mockDB, database.ColPodcast, pod2)
	insertOrFail(t, mockDB, database.ColPodcast, pod3)
	insertOrFail(t, mockDB, database.ColPodcast, pod4)

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
		{
			name: "Valid",
			args: args{
				dbClient: mockDB,
				start:    0,
				end:      4,
			},
			want:    []*protos.Podcast{pod1, pod2, pod3, pod4},
			wantErr: false,
		},
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
	mockDB := mock.CreateDB()
	pod1 := &protos.Podcast{Id: protos.NewObjectID()}
	pod2 := &protos.Podcast{Id: protos.NewObjectID()}
	insertOrFail(t, mockDB, database.ColPodcast, pod1)
	insertOrFail(t, mockDB, database.ColPodcast, pod2)

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
		{
			name: "valid",
			args: args{
				dbClient: mockDB,
				id:       pod1.Id,
			},
			want:    pod1,
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				dbClient: mockDB,
				id:       protos.NewObjectID(),
			},
			want:    nil,
			wantErr: true,
		},
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

func TestSearchPodcasts(t *testing.T) {
	mockDB := mock.CreateDB()
	pod1 := &protos.Podcast{Id: protos.NewObjectID(), Author: "Sam Schwartz"}
	pod2 := &protos.Podcast{Id: protos.NewObjectID(), Title: "The Tech Podcast"}
	pod3 := &protos.Podcast{Id: protos.NewObjectID(), Keywords: []string{"food", "taste"}}
	insertOrFail(t, mockDB, database.ColPodcast, pod1)
	insertOrFail(t, mockDB, database.ColPodcast, pod2)
	insertOrFail(t, mockDB, database.ColPodcast, pod3)

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
		{
			name: "valid1",
			args: args{
				dbClient: mockDB,
				search:   "Sam Schwartz",
			},
			want:    []*protos.Podcast{pod1},
			wantErr: false,
		},
		{
			name: "valid2",
			args: args{
				dbClient: mockDB,
				search:   "Tech",
			},
			want:    []*protos.Podcast{pod2},
			wantErr: false,
		},
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

func TestFindLength(t *testing.T) {
	mp3File, err := os.Open("./test/sample.mp3")
	if err != nil {
		t.Fatalf("TestFindLength() unable to open file: %v", err)
	}
	mp3FileInfo, err := mp3File.Stat()
	if err != nil {
		t.Fatalf("TestFindLength() unable to get file stat: %v", err)
	}

	type args struct {
		r          io.Reader
		fileLength int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "sample",
			args: args{
				r:          mp3File,
				fileLength: mp3FileInfo.Size(),
			},
			want: 10000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindLength(tt.args.r, tt.args.fileLength); got != tt.want {
				t.Errorf("FindLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindPodcastsByIDs(t *testing.T) {
	mockDB := mock.CreateDB()
	insertOrFail(t, mockDB, database.ColPodcast, &protos.Podcast{Id: protos.ObjectIDFromHex("pod_id1")})
	insertOrFail(t, mockDB, database.ColPodcast, &protos.Podcast{Id: protos.ObjectIDFromHex("pod_id2")})
	insertOrFail(t, mockDB, database.ColPodcast, &protos.Podcast{Id: protos.ObjectIDFromHex("pod_id3")})

	type args struct {
		dbClient db.Database
		ids      []*protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		want    []*protos.Podcast
		wantErr bool
	}{
		{
			name: "2&3",
			args: args{
				dbClient: mockDB,
				ids:      []*protos.ObjectID{protos.ObjectIDFromHex("pod_id2"), protos.ObjectIDFromHex("pod_id3")},
			},
			want:    []*protos.Podcast{{Id: protos.ObjectIDFromHex("pod_id2")}, {Id: protos.ObjectIDFromHex("pod_id3")}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindPodcastsByIDs(tt.args.dbClient, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPodcastsByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPodcastsByIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
