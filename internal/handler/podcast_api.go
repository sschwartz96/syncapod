package handler

import (
	"fmt"
	"net/http"

	"github.com/sschwartz96/syncapod/internal/models"
)

// Podcast handles all request on /api/podcast/*
func (h *APIHandler) Podcast(res http.ResponseWriter, req *http.Request, user *models.User) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "subscriptions":
		h.Subscription(res, req, user)
	default:
		fmt.Fprint(res, "This endpoint is not supported")
	}
}

// Subscription handles requests on /api/podcast/subscription/*
func (h *APIHandler) Subscription(res http.ResponseWriter, req *http.Request, user *models.User) {
	var err error
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "get":
		subs := h.dbClient.FindUserSubs(user.ID)
		var pods []models.Podcast
		for i := range subs {
			pods = append(pods, *subs[i].Podcast)
		}
		err = sendObjectJSON(res, pods)
	default:
		fmt.Fprint(res, "This endpoint is not supported")
	}

	if err != nil {
		fmt.Println("error sending json object: ", err)
		fmt.Fprint(res, "internal error: ")
	}
}
