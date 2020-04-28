package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
)

// APIHandler handles calls to the syncapod api
type APIHandler struct {
	dbClient *database.Client
}

// CreateAPIHandler instatiates an APIHandler
func CreateAPIHandler(dbClient *database.Client) (*APIHandler, error) {
	return &APIHandler{
		dbClient: dbClient,
	}, nil
}

// ServeHTTP handles all requests throught /api/* endpoint
func (h *APIHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	var handler func(http.ResponseWriter, *http.Request, *models.User)

	switch head {
	// if endpoint is alexa then we need to just return cause that is handled with oauth
	case "alexa":
		h.Alexa(res, req)
		return

	// auth handles authentication
	case "auth":
		h.Auth(res, req)
		return

	// the rest need to be authorized first
	case "subscriptions":
		handler = h.Subscriptions

	default:
		fmt.Fprint(res, "This endpoint is not supported")
		return
	}

	user, ok := h.checkAuth(req)

	if ok {
		handler(res, req, user)
	}
}

// Subscriptions endpoint returns the users subscriptions
func (h *APIHandler) Subscriptions(res http.ResponseWriter, req *http.Request, user *models.User) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "get":
		subs := h.dbClient.FindUserSubs(user.ID)
		response, _ := json.Marshal(&subs)
		res.Header().Add("Content-Type", "application/json")
		res.Write(response)
	}
}

func (h *APIHandler) checkAuth(req *http.Request) (*models.User, bool) {
	token, _, _ := req.BasicAuth()

	if token != "" {
		u, err := auth.ValidateSession(h.dbClient, token)
		if err != nil {
			return nil, false
		}
		return u, true
	}

	return nil, false
}
