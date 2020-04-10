package handler

import (
	"net/http"

	"github.com/sschwartz96/syncapod/internal/database"
)

// APIHandler handles calls to the syncapod api
type APIHandler struct {
	dbClient *database.Client
}

func CreateAPIHandler(dbClient *database.Client) (*APIHandler, error) {
	return &APIHandler{
		dbClient: dbClient,
	}, nil
}

// ServeHTTP handles all requests throught /api/* endpoint
func (h *APIHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	switch head {
	case "alexa":
		h.Alexa(res, req)
	}
}
