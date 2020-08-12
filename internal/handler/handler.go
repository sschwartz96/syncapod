package handler

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"

	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
)

// Handler is the main handler for syncapod, all routes go through it
type Handler struct {
	dbClient     *database.Client
	oauthHandler *OauthHandler
	apiHandler   *APIHandler
}

// CreateHandler sets up the main handler
func CreateHandler(dbClient *database.Client, config *config.Config) (*Handler, error) {
	handler := &Handler{}
	var err error

	handler.oauthHandler, err = CreateOauthHandler(dbClient, config.AlexaClientID, config.AlexaSecret)
	if err != nil {
		return nil, err
	}

	handler.apiHandler, err = CreateAPIHandler(dbClient)
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
	case "api":
		h.apiHandler.ServeHTTP(res, req)
	}
}

// sendMessageJSON is a helper method in which it decodes an object and sends
// via http.ResponseWriter given, returns error if decode fails
func sendObjectJSON(res http.ResponseWriter, object interface{}) error {
	jsonRes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	res.Header().Add("Content-Type", "application/json")
	_, err = res.Write(jsonRes)
	return err
}

func sendMessageJSON(res http.ResponseWriter, message string) error {
	type Response struct {
		Message string `json:"message"`
	}
	response := Response{Message: message}
	// jsonRes, _ := json.Marshal(&response)
	// res.Header().Add("Content-Type", "application/json")
	// res.Write(jsonRes)
	return sendObjectJSON(res, &response)
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
