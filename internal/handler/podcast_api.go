package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/sschwartz96/syncapod/internal/podcast"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Podcast handles all request on /api/podcast/*
func (h *APIHandler) Podcast(res http.ResponseWriter, req *http.Request, user *models.User) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "subscriptions":
		h.Subscription(res, req, user)
	case "episodes":
		h.Episodes(res, req, user)
	default:
		fmt.Fprint(res, "This endpoint is not supported")
	}
}

// Episodes handles requests on /api/podcast/episodes/*
func (h *APIHandler) Episodes(res http.ResponseWriter, req *http.Request, user *models.User) {
	var err error
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "get":
		jReq, err := getJSONObj(req)
		if err != nil {
			sendMessageJSON(res, fmt.Sprint("Error sending episodes: ", err))
			return
		}
		id, err := primitive.ObjectIDFromHex(jReq.ID)
		if err != nil {
			sendMessageJSON(res, "invalid object id")
			return
		}

		epis := podcast.FindAllEpisodesRange(h.dbClient, id, jReq.Start, jReq.End)
		err = sendObjectJSON(res, epis)
	default:
		sendMessageJSON(res, "This endpoint is not supported")
	}

	if err != nil {
		fmt.Println("error sending json object: ", err)
		sendMessageJSON(res, "internal error: ")
	}
}

// Subscription handles requests on /api/podcast/subscription/*
func (h *APIHandler) Subscription(res http.ResponseWriter, req *http.Request, user *models.User) {
	var err error
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "get":
		sendMessageJSON(res, "subscriptions are sent with user object")
	default:
		sendMessageJSON(res, "This endpoint is not supported")
	}

	if err != nil {
		fmt.Println("error sending json object: ", err)
		sendMessageJSON(res, "internal error: ")
	}
}

// JSONReq is what we could receive in a json request
type JSONReq struct {
	ID    string `json:"id,omitempty"`
	Start int    `json:"start,omitempty"`
	End   int    `json:"end,omitempty"`
}

func getJSONObj(req *http.Request) (*JSONReq, error) {
	// setup the request body for reading
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	// unmarshal json
	var request JSONReq
	err = json.Unmarshal(body, &request)
	if err != nil {
		return nil, err
	}
	return &request, err
}
