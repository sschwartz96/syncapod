package handler

import (
	"net/http"

	"github.com/sschwartz96/syncapod/internal/database"
)

// APIHandler handles calls to the syncapod api
type APIHandler struct {
	dbClient *database.Client
}

// ServeHTTP handles all requests throught /api/* endpoint
func (h *APIHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	accessToken := req.URL.Query().Get("access_token")

}
