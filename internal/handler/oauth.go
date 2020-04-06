package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

// OauthHandler handles authorization and authentication to oauth clients
type OauthHandler struct {
	loginTemplate *template.Template
	authTemplate  *template.Template
}

// CreateOauthHandler just intantiates an OauthHandler
func CreateOauthHandler() (*OauthHandler, error) {
	loginT, err := template.ParseFiles("templates/oauth/login.gohtml")
	authT, err := template.ParseFiles("templates/oauth/auth.gohtml")
	if err != nil {
		return nil, err
	}

	return &OauthHandler{
		loginTemplate: loginT,
		authTemplate:  authT,
	}, nil
}

func (h *OauthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		var head string
		var err error

		head, req.URL.Path = ShiftPath(req.URL.Path)

		switch head {
		case "login":
			err = h.loginTemplate.Execute(res, nil)
		case "authorize":
			err = h.authTemplate.Execute(res, nil)
		}

		if err != nil {
			fmt.Println("error executing template: ", err)
		}
	}
}
