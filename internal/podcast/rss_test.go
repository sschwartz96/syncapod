package podcast

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/sschwartz96/syncapod/internal/protos"
)

func TestUpdatePodcasts(t *testing.T) {
	type args struct {
		dbClient db.Database
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
			UpdatePodcasts(tt.args.dbClient)
		})
	}
}

func TestUpdatePodcast(t *testing.T) {
	type args struct {
		wg       *sync.WaitGroup
		dbClient db.Database
		pod      *protos.Podcast
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
			UpdatePodcast(tt.args.wg, tt.args.dbClient, tt.args.pod)
		})
	}
}

func TestAddNewPodcast(t *testing.T) {
	type args struct {
		dbClient db.Database
		url      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := AddNewPodcast(tt.args.dbClient, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("AddNewPodcast() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseRSS(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.RSSPodcast
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseRSS(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRSS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRSS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertEpisode(t *testing.T) {
	type args struct {
		pID *protos.ObjectID
		e   *models.RSSEpisode
	}
	tests := []struct {
		name string
		args args
		want *protos.Episode
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := convertEpisode(tt.args.pID, tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertEpisode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertPodcast(t *testing.T) {
	type args struct {
		url string
		p   *models.RSSPodcast
	}
	tests := []struct {
		name string
		args args
		want *protos.Podcast
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := convertPodcast(tt.args.url, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertPodcast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseRFC2822(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseRFC2822(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRFC2822() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRFC2822() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDuration(t *testing.T) {
	type args struct {
		d string
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
			if got := parseDuration(tt.args.d); got != tt.want {
				t.Errorf("parseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertCategories(t *testing.T) {
	type args struct {
		cats []models.Category
	}
	tests := []struct {
		name string
		args args
		want []*protos.Category
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := convertCategories(tt.args.cats); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertCategory(t *testing.T) {
	type args struct {
		cat models.Category
	}
	tests := []struct {
		name string
		args args
		want *protos.Category
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := convertCategory(tt.args.cat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
