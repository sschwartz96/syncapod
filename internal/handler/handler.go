package handler

import (
	"fmt"
	"net/http"

	"github.com/sschwartz96/syncapod/internal/database"
)

// Handler is the main handler for syncapod, all routes go through it
type Handler struct {
	dbClient     *database.Client
	oauthHandler *OauthHandler
}

// CreateHandler sets up the main handler
func CreateHandler(dbClient *database.Client) (*Handler, error) {
	handler := &Handler{}
	var err error

	handler.oauthHandler, err = CreateOauthHandler()
	if err != nil {
		return nil, err
	}

	return handler, nil
}

// ServeHTTP handles all requests
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	keys := req.URL.Query()
	user := keys.Get("user")

	fmt.Println("user trying to access: ", user)

	fmt.Fprintln(res, "user: ", user)
}
