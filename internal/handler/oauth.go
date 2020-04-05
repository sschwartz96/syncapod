package handler

import "net/http"

// OauthHandler handles authorization and authentication to oauth clients
type OauthHandler struct{}

// CreateOauthHandler just intantiates an OauthHandler
func CreateOauthHandler() (*OauthHandler, error) {
	return &OauthHandler{}, nil
}

func (h *OauthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
}
