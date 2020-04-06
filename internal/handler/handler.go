package handler

import (
	"net/http"
	"path"
	"strings"

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

	handler.oauthHandler, err = CreateOauthHandler(dbClient)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

// ServeHTTP handles all requests
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "oauth":
		h.oauthHandler.ServeHTTP(res, req)
	}
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
