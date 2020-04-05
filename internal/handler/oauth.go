package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

// OauthHandler handles authorization and authentication to oauth clients
type OauthHandler struct {
	authTemplate *template.Template
}

// CreateOauthHandler just intantiates an OauthHandler
func CreateOauthHandler() (*OauthHandler, error) {
	authTemplate, err := template.ParseFiles("templates/oauth.gohtml")
	if err != nil {
		return nil, err
	}

	return &OauthHandler{
		authTemplate: authTemplate,
	}, nil
}

func (h *OauthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := h.authTemplate.Execute(res, nil)
	if err != nil {
		fmt.Println("error executing auth template: ", err)
	}
}
