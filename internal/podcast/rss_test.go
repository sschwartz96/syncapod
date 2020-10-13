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
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdatePodcasts(tt.args.dbClient); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePodcasts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updatePodcast(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			updatePodcast(tt.args.wg, tt.args.dbClient, tt.args.pod)
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
		t.Run(tt.name, func(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
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
		want    time.Time
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				s: "Thu, 08 Oct 2020 15:30:00 +0000",
			},
			want:    time.Unix(1602171000, 0),
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				s: "Tue, 06 Oct 2020 20:00:00 PDT",
			},
			want:    time.Unix(1602039600, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRFC2822ToUTC(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRFC2822() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.UTC(), tt.want.UTC()) {
				t.Errorf("parseRFC2822() = %v, want %v", got.UTC(), tt.want.UTC())
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
		{
			name: "seconds",
			args: args{
				d: "5400",
			},
			want: 5400000,
		},
		{
			name: "hh:mm:ss",
			args: args{
				d: "01:30:30",
			},
			want: 5430000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		{
			name: "one",
			args: args{
				cats: []models.Category{
					{Text: "CatOne", Category: []models.Category{{Text: "SubCatOne", Category: []models.Category{}}}},
					{Text: "CatTwo", Category: []models.Category{}},
				},
			},
			want: []*protos.Category{
				{Text: "CatOne", Category: []*protos.Category{{Text: "SubCatOne", Category: []*protos.Category{}}}},
				{Text: "CatTwo", Category: []*protos.Category{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		{
			name: "one",
			args: args{
				cat: models.Category{Text: "CatOne", Category: []models.Category{{Text: "SubCatOne", Category: []models.Category{}}}},
			},
			want: &protos.Category{Text: "CatOne", Category: []*protos.Category{{Text: "SubCatOne", Category: []*protos.Category{}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertCategory(tt.args.cat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
