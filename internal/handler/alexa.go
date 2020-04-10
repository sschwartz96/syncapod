package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
)

// Intents
const (
	PlayPodcast       = "PlayPodcast"
	PlayLatestPodcast = "PlayLatestPodcast"
	PlayNthFromLatest = "PlayNthFromLatest"
	FastForward       = "FastForward"
	Rewind            = "Rewind"
	Pause             = "AMAZON.PauseIntent"
)

// Alexa handles all requests through /api/alexa endpoint
func (h *APIHandler) Alexa(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("couldn't read the body of the request")
		// TODO: proper response here
		return
	}

	var aData AlexaData
	err = json.Unmarshal(body, &aData)
	if err != nil {
		fmt.Println("couldn't unmarshal json to object: ", err)
		// TODO: proper response here
	}

	// get the person or user accessToken
	token, err := getAccessToken(&aData)
	if err != nil {
		fmt.Println("no accessToken: ", err)
		// TODO: proper response here
	}

	// validate the token and return user
	user, err := auth.ValidateAccessToken(h.dbClient, token)
	if err != nil {
		fmt.Println("error validating token: ", err)
	}
	fmt.Println(user)

	name := aData.Request.Intent.AlexaSlots.Podcast.Value
	fmt.Println("request name of podcast: ", name)

	response := createAlexaResponse(user.ID.Hex())

	switch aData.Request.Intent.Name {
	case PlayPodcast:
		var podcasts []models.Podcast
		err = h.dbClient.Search(database.ColPodcast, name, &podcasts)
		if err != nil {
			fmt.Println("error searching for podcast: ", err)
		}
		if len(podcasts) > 0 {
			response.Response.Directives[0].AudioItem.Stream.URL = podcasts[0].Episodes[0].Enclosure.MP3

			jsonRes, err := json.Marshal(response)
			if err != nil {
				fmt.Println("couldn't marshal alexa response: ", err)
			}

			res.Header().Set("Content-Type", "application/json")
			res.Write(jsonRes)
		} else {
			//TODO: no podcast found
			fmt.Println("no podcast found")
		}

	case PlayLatestPodcast:
	case PlayNthFromLatest:

	case FastForward:
	case Rewind:

	case Pause:

	}
}

func createAlexaResponse(userID string) *AlexaResponseData {
	return &AlexaResponseData{
		Version: "1.0",
		Response: AlexaResponse{
			Directives: []AlexaDirective{
				AlexaDirective{
					Type:         "AudioPlayer.Play",
					PlayBehavior: "REPLACE_ALL",
					AudioItem: AlexaAudioItem{
						Stream: AlexaStream{
							URL:                  "",
							Token:                userID,
							OffsetInMilliseconds: 0,
						},
					},
				},
			},
		},
	}
}

func getAccessToken(data *AlexaData) (string, error) {
	if data.Context.System.Person.AccessToken != "" {
		return data.Context.System.Person.AccessToken, nil
	} else if data.Context.System.User.AccessToken != "" {
		return data.Context.System.User.AccessToken, nil
	}
	return "", errors.New("no accessToken")
}

// AlexaData contains all the informatino and data from request sent from alexa
type AlexaData struct {
	Version string       `json:"version"`
	Context AlexaContext `json:"context"`
	Request AlexaRequest `json:"request"`
}

// AlexaContext contains system
type AlexaContext struct {
	System AlexaSystem `json:"system"`
}

// AlexaSystem is the container for person and user
type AlexaSystem struct {
	Person AlexaPerson `json:"person"`
	User   AlexaUser   `json:"user"`
}

// AlexaPerson holds the info about the person who explicitly called the skill
type AlexaPerson struct {
	PersonID    string `json:"personId"`
	AccessToken string `json:"accessToken"`
}

// AlexaUser contains info about the user that holds the skill
type AlexaUser struct {
	UserID      string `json:"userId"`
	AccessToken string `json:"accessToken"`
}

// AlexaRequest holds all the information and data
type AlexaRequest struct {
	Type                 string      `json:"type"`
	RequestID            string      `json:"requestId"`
	Timestamp            time.Time   `json:"timestamp"`
	Token                string      `json:"token"`
	OffsetInMilliseconds int64       `json:"offsetInMilliseconds"`
	Intent               AlexaIntent `json:"intent"`
}

// AlexaIntent holds information and data of intent sent from alexa
type AlexaIntent struct {
	Name       string     `json:"name"`
	AlexaSlots AlexaSlots `json:"slots"`
}

// AlexaSlots are the container for the slots
type AlexaSlots struct {
	Nth     AlexaSlot `json:"nth"`
	Episode AlexaSlot `json:"episode"`
	Podcast AlexaSlot `json:"podcast"`
}

// AlexaSlot holds information of the slot for the intent
type AlexaSlot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AlexaResponseData contains the version and response
type AlexaResponseData struct {
	Version  string        `json:"version"`
	Response AlexaResponse `json:"response"`
}

// AlexaResponse contains the actual response
type AlexaResponse struct {
	Directives []AlexaDirective `json:"directives"`
}

// AlexaDirective tells alexa what to do
type AlexaDirective struct {
	Type         string         `json:"type"`
	PlayBehavior string         `json:"playBehavior"`
	AudioItem    AlexaAudioItem `json:"audioItem"`
}

// AlexaAudioItem holds information of audio track
type AlexaAudioItem struct {
	Stream AlexaStream `json:"stream"`
}

type AlexaStream struct {
	Token                string `json:"token"`
	URL                  string `json:"url"`
	OffsetInMilliseconds int64  `json:"offsetInMilliseconds"`
}