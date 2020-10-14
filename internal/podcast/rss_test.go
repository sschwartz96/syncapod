package podcast

import (
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/stockpile/mock"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/sschwartz96/syncapod/internal/protos"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdatePodcasts(t *testing.T) {
	mockDB := mock.CreateDB()
	insertOrFail(t, mockDB, database.ColPodcast, &protos.Podcast{Id: protos.NewObjectID(), Title: "Go Time", Author: "Changelog Media", Type: "", Subtitle: "", Link: "https://changelog.com/gotime", Image: &protos.Image{Title: "", Url: ""}, Explicit: "no", Language: "en-us", Keywords: []string{"go", "golang", "open source", "software", "development", "devops", "architecture", "docker", "kubernetes"}, Rss: "https://changelog.com/gotime/feed"})
	insertOrFail(t, mockDB, database.ColEpisode, &protos.Episode{})

	type args struct {
		dbClient db.Database
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				dbClient: mockDB,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := UpdatePodcasts(tt.args.dbClient); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePodcasts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updatePodcast(t *testing.T) {
	mockDB := mock.CreateDB()
	type args struct {
		dbClient db.Database
		pod      *protos.Podcast
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				dbClient: mockDB,
				pod:      &protos.Podcast{Id: protos.NewObjectID(), Title: "Go Time", Author: "Changelog Media", Type: "", Subtitle: "", Link: "https://changelog.com/gotime", Image: &protos.Image{Title: "", Url: ""}, Explicit: "no", Language: "en-us", Keywords: []string{"go", "golang", "open source", "software", "development", "devops", "architecture", "docker", "kubernetes"}, Rss: "https://changelog.com/gotime/feed"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := updatePodcast(tt.args.dbClient, tt.args.pod); (err != nil) != tt.wantErr {
				t.Errorf("updatePodcast() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddNewPodcast(t *testing.T) {
	mockDB := mock.CreateDB()
	type args struct {
		dbClient db.Database
		url      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				dbClient: mockDB,
				url:      "https://changelog.com/gotime/feed",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := AddNewPodcast(tt.args.dbClient, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("AddNewPodcast() error = %v, wantErr %v", err, tt.wantErr)
			}
			// validate in the database
			newPod := &protos.Podcast{}
			err := tt.args.dbClient.FindOne(database.ColPodcast, newPod, &db.Filter{}, nil)
			if err != nil {
				t.Errorf("AddNewPodcast() error finding podcast: %v", err)
			}
			newEpis := []*protos.Episode{}
			err = tt.args.dbClient.FindAll(database.ColEpisode, &newEpis, &db.Filter{"podcastid": newPod.Id}, nil)
			if err != nil {
				t.Errorf("AddNewPodcast() error finding episodes: %v", err)
			}
			if len(newEpis) == 0 {
				t.Errorf("AddNewPodcast() error no valid episodes added")
			}
		})
	}
}

func Test_parseRSS(t *testing.T) {
	rssFile, err := os.Open("./test/feed.xml")
	if err != nil {
		t.Fatalf("Test_parseRSS() error opening test file: %v", err)
	}
	defer rssFile.Close()

	type args struct {
		r io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    *models.RSSPodcast
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				r: rssFile,
			},
			want:    &models.RSSPodcast{ID: primitive.ObjectID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, Title: "Go Time", Author: "Changelog Media", Type: "", Subtitle: "", Summary: "Your source for diverse discussions from around the Go community  Panelists include Mat Ryer, Ashley McNamara, Johnny Boursiquot, Carmen Andoh, Jaana B. Dogan (JBD), Mark Bates, and Jon Calhoun.\n\n\t\tThis show records LIVE every Tuesday at 3pm US Eastern. Join the Golang community and chat with us during the show in the #gotimefm channel of Gophers slack.\n\n\t\tWe discuss cloud infrastructure, distributed systems, microservices, Kubernetes, Docker... oh and also Go!\n\n\t\tSome people search for GoTime or GoTimeFM and can't find the show, so now the strings GoTime and GoTimeFM are in our description too.", Link: "https://changelog.com/gotime", Image: models.Image{Title: "", URL: ""}, Explicit: "no", Language: "en-us", Keywords: "go, golang, open source, software, development, devops, architecture, docker, kubernetes", Category: []models.Category{models.Category{Text: "Technology", Category: []models.Category{models.Category{Text: "Software How-To", Category: []models.Category(nil)}, models.Category{Text: "Tech News", Category: []models.Category(nil)}}}}, PubDate: "", LastBuildDate: "", RSSEpisodes: []models.RSSEpisode{models.RSSEpisode{ID: primitive.ObjectID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, PodcastID: primitive.ObjectID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, Title: "There's a lot to learn about teaching Go", Subtitle: " Mat, Jon, Johnny, & Mark", Author: "Mat Ryer, Jon Calhoun, Johnny Boursiquot, and Mark Bates", Type: "", Image: models.EpiImage{HREF: "https://cdn.changelog.com/uploads/covers/go-time-original.png?v=63725770357"}, Thumbnail: models.EpiThumbnail{URL: ""}, PubDate: "Thu, 01 Oct 2020 15:00:00 +0000", Description: "In this episode we dive into teaching Go, asking questions like, “What techniques work well for teaching programming?”, “What role does community play in education?”, and “What are the best ways to improve at Go as a beginner/intermediate/senior dev?” ", Summary: "In this episode we dive into teaching Go, asking questions like, “What techniques work well for teaching programming?”, “What role does community play in education?”, and “What are the best ways to improve at Go as a beginner/intermediate/senior dev?” ", Season: 0, Episode: 149, Category: []models.Category(nil), Explicit: "no", Enclosure: models.Enclosure{MP3: "https://cdn.changelog.com/uploads/gotime/149/go-time-149.mp3"}, Duration: "1:16:18"}}, NewFeedURL: "", RSS: ""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRSS(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRSS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRSS() = \n\t%#v,\nwant\n\t%v", got, tt.want)
			}
		})
	}
}

func Test_convertEpisode(t *testing.T) {
	podID := protos.NewObjectID()
	timeData := time.Unix(1602171000, 0)
	tStampData, _ := ptypes.TimestampProto(timeData)
	type args struct {
		pID *protos.ObjectID
		e   *models.RSSEpisode
	}
	tests := []struct {
		name string
		args args
		want *protos.Episode
	}{
		{
			name: "valid",
			args: args{
				pID: podID,
				e:   &models.RSSEpisode{Author: "Sam Schwartz", Title: "Title to Episode", PubDate: "Thu, 08 Oct 2020 15:30:00 +0000"},
			},
			want: &protos.Episode{PodcastID: podID, Author: "Sam Schwartz", Title: "Title to Episode", PubDate: tStampData, Image: &protos.Image{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := convertEpisode(tt.args.pID, tt.args.e)
			got.Id = nil
			if !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("convertEpisode() = \n\t%v,\nwant \n\t%v", got.String(), tt.want.String())
			}
		})
	}
}

func Test_convertPodcast(t *testing.T) {
	timeData := time.Unix(1602171000, 0)
	tStampData, _ := ptypes.TimestampProto(timeData)
	type args struct {
		url string
		p   *models.RSSPodcast
	}
	tests := []struct {
		name string
		args args
		want *protos.Podcast
	}{
		{
			name: "valid",
			args: args{
				url: "https://example.com",
				p:   &models.RSSPodcast{Author: "Sam Schwartz", PubDate: "Thu, 08 Oct 2020 15:30:00 +0000", LastBuildDate: "Thu, 08 Oct 2020 15:30:00 +0000"},
			},
			want: &protos.Podcast{Rss: "https://example.com", Author: "Sam Schwartz", LastBuildDate: tStampData, PubDate: tStampData, Image: &protos.Image{}, Keywords: []string{""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := convertPodcast(tt.args.url, tt.args.p)
			got.Id = nil
			if !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("convertPodcast() = \n\t%v,\n want \n\t%v", got.String(), tt.want.String())
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
			t.Parallel()
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
			t.Parallel()
			if got := convertCategory(tt.args.cat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDuration(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
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
			t.Parallel()
			got, err := parseDuration(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
