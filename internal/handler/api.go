package handler

import (
	"fmt"
	"net/http"

	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/database"
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

	var handler func(http.ResponseWriter, *http.Request)

	switch head {
	// if endpoint is alexa then we need to just return cause that is handled with oauth
	case "alexa":
		h.Alexa(res, req)
		return

	// auth handles authentication
	case "auth":
		h.Auth(res, req)

	// the rest need to be authorized first

	default:
		fmt.Fprint(res, "This endpoint is not supported")
		return
	}

	if h.checkAuth(req) {
		handler(res, req)
	}
}

func (h *APIHandler) checkAuth(req *http.Request) bool {
	token := req.URL.Query().Get("access_token")

	if token != "" {
		_, err := auth.ValidateSession(h.dbClient, token)
		if err != nil {
			return false
		}
		return true
	}

	return false
}
