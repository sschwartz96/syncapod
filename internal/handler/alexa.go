package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

	switch aData.Request.Type {
	case "AudioPlayer.Play":
		// get the title and episode details for podcast
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

// AlexaData
type AlexaData struct {
	Version float64      `json:"version"`
	Context AlexaContext `json:"context"`
	Request AlexaRequest `json:"request"`
}

type AlexaContext struct {
	System AlexaSystem `json:"system"`
}

type AlexaSystem struct {
	Person AlexaPerson `json:"person"`
	User   AlexaUser   `json:"user"`
}

type AlexaPerson struct {
	PersonID    string `json:"personId"`
	AccessToken string `json:"accessToken"`
}

type AlexaUser struct {
	UserID      string `json:"userId"`
	AccessToken string `json:"accessToken"`
}

type AlexaRequest struct {
	Type                 string    `json:"type"`
	RequestID            string    `json:"requestId"`
	Timestamp            time.Time `json:"timestamp"`
	Token                string    `json:"token"`
	OffsetInMilliseconds int64     `json:"offsetInMilliseconds"`
}
