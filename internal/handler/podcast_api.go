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
	case "user_episode":
		h.UserEpisode(res, req, user)
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
		id, err := primitive.ObjectIDFromHex(jReq.PodID)
		if err != nil {
			sendMessageJSON(res, "invalid object id")
			return
		}

		epis := podcast.FindAllEpisodesRange(h.dbClient, id, jReq.Start, jReq.End)
		err = sendObjectJSON(res, epis)
	case "latest":
		// retrieve the last episode
		pod, epi, offset, err := podcast.FindUserLastPlayed(h.dbClient, user.ID)
		if err != nil {
			fmt.Println("error getting the last user played: ", err)
			sendMessageJSON(res, "User has no last played")
		}

		// setup information
		type LatestEpisode struct {
			Podcast *models.Podcast `json:"podcast"`
			Episode *models.Episode `json:"episode"`
			Offset  int64           `json:"offset"`
		}
		latestEpisode := &LatestEpisode{Podcast: pod, Episode: epi, Offset: offset}

		// marshal and send
		jsonRes, _ := json.Marshal(latestEpisode)
		res.Write(jsonRes)

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

// UserEpisode handles all requests at /api/podcasts/user_episode/*
func (h *APIHandler) UserEpisode(res http.ResponseWriter, req *http.Request, user *models.User) {
	var err error
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	info, err := getJSONObj(req)
	if err != nil {
		fmt.Println("couldn't parse the json body of request")
		sendMessageJSON(res, "couldn't parse the json body of the request")
		return
	}
	podID, _ := primitive.ObjectIDFromHex(info.PodID)
	epiID, _ := primitive.ObjectIDFromHex(info.EpiID)

	switch head {
	case "get":
		userEpi, err := podcast.FindUserEpisode(h.dbClient, user.ID, epiID)
		if err != nil {
			fmt.Println("error trying to get userEpi: ", err)
			sendMessageJSON(res, fmt.Sprint("error trying to get userEpi"))
		}
		jsonRes, err := json.Marshal(&userEpi)
		if err != nil {
			sendMessageJSON(res, fmt.Sprint("error marshalling userEpi"))
		}
		res.Write(jsonRes)

	case "update":
		err = podcast.UpdateOffset(h.dbClient, user.ID, podID, epiID, info.Offset)
		if err != nil {
			sendMessageJSON(res, "error updating offest")
			return
		}
		sendMessageJSON(res, "success")
	default:
		fmt.Fprint(res, "Method not supported")
	}
}

// JSONReq is what we could receive in a json request
type JSONReq struct {
	PodID  string `json:"pod_id,omitempty"` // podcast id, could be empty
	EpiID  string `json:"epi_id,omitempty"` // episode id, could be empty
	Start  int    `json:"start,omitempty"`  // used for getting a list of episodes
	End    int    `json:"end,omitempty"`    // used for getting list of episodes
	Offset int64  `json:"offset,omitempty"` // position is the time in milliseconds
	Played bool   `json:"played,omitempty"` // if the user episode should be marked as played
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
