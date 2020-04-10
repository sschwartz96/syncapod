package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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
		fmt.Println("couldn't unmarshal json to object")
		// TODO: proper response here
	}

	// get the person or user accessToken
	token, err := getAccessToken(&aData)
	if err != nil {
		fmt.Println("no accessToken")
		// TODO: proper response here
	}

	fmt.Println("token: ", token)

	switch aData.Request.Intent.Name {
	case PlayPodcast:
	case PlayLatestPodcast:
	case PlayNthFromLatest:

	case FastForward:
	case Rewind:

	case Pause:

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
	Version float64      `json:"version"`
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
	AlexaSlots AlexaSlots `json:"alexa_slot"`
}

// AlexaSlots are the container for the slots
type AlexaSlots struct {
	Nth     AlexaSlot
	Episode AlexaSlot
	Podcast AlexaSlot
}

// AlexaSlot holds information of the slot for the intent
type AlexaSlot struct {
	Name  string
	Value string
}
