package podcast

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/syncapod/internal/protos"
)

func TestFindEpisodesByRange(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindEpisodesByRange(tt.args.dbClient, tt.args.podcastID, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindEpisodesByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindEpisodesByRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindAllEpisodes(t *testing.T) {
	type args struct {
		dbClient  db.Database
		podcastID *protos.ObjectID
	}
	tests := []struct {
		name    string
		args    args
		want    []*protos.Episode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindAllEpisodes(tt.args.dbClient, tt.args.podcastID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAllEpisodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindLatestEpisode(t *testing.T) {
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
	type args struct {
		dbClient db.Database
		episode  *protos.Episode
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
			if err := UpsertEpisode(tt.args.dbClient, tt.args.episode); (err != nil) != tt.wantErr {
				t.Errorf("UpsertEpisode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDoesEpisodeExist(t *testing.T) {
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
		// TODO: Add test cases.
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
