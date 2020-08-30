package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/user"
)

// OauthHandler handles authorization and authentication to oauth clients
type OauthHandler struct {
	dbClient      db.Database
	loginTemplate *template.Template
	authTemplate  *template.Template
	// only used for alexa, need these in database if suppport more than one client
	clientID     string
	clientSecret string
}

// CreateOauthHandler just intantiates an OauthHandler
func CreateOauthHandler(dbClient db.Database, clientID, clientSecret string) (*OauthHandler, error) {
	loginT, err := template.ParseFiles("templates/oauth/login.gohtml")
	authT, err := template.ParseFiles("templates/oauth/auth.gohtml")
	if err != nil {
		return nil, err
	}

	return &OauthHandler{
		dbClient:      dbClient,
		loginTemplate: loginT,
		authTemplate:  authT,
		clientID:      clientID,
		clientSecret:  clientSecret,
	}, nil
}

func (h *OauthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// path: /oauth/*
	if req.Method == http.MethodGet {
		h.Get(res, req)
	} else if req.Method == http.MethodPost {
		h.Post(res, req)
	}
}

// Get is the handler for /api/oauth
func (h *OauthHandler) Get(res http.ResponseWriter, req *http.Request) {
	var head string
	var err error

	head, req.URL.Path = ShiftPath(req.URL.Path)

	// path: /oauth/*
	switch head {
	case "login":
		err = h.loginTemplate.Execute(res, nil)
	case "authorize":
		key := strings.TrimSpace(req.URL.Query().Get("sesh_key"))
		_, err := auth.ValidateSession(h.dbClient, key)
		if err != nil {
			fmt.Println("couldn't not validate, redirecting to login page: ", err)
			http.Redirect(res, req, "/oauth/login", http.StatusSeeOther)
			return
		}
		err = h.authTemplate.Execute(res, nil)
	}

	if err != nil {
		fmt.Println("error executing template: ", err)
	}
}

// Post hanldes all post request at the oauth endpoint
func (h *OauthHandler) Post(res http.ResponseWriter, req *http.Request) {
	var head string

	head, req.URL.Path = ShiftPath(req.URL.Path)

	// path: /oauth/*
	switch head {
	case "login":
		h.Login(res, req)
	case "authorize":
		h.Authorize(res, req)
	case "token":
		h.Token(res, req)
	}
}

// Login handles the post and get request of a login page
func (h *OauthHandler) Login(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Println("couldn't parse post values: ", err)
		h.loginTemplate.Execute(res, true)
		return
	}

	username := req.FormValue("uname")
	password := req.FormValue("pass")

	userObj, err := user.FindUser(h.dbClient, username)
	if err != nil {
		h.loginTemplate.Execute(res, true)
		return
	}

	if auth.Compare(userObj.Password, password) {
		key, err := auth.CreateSession(h.dbClient, userObj.Id, req.UserAgent(), false)
		if err != nil {
			h.loginTemplate.Execute(res, true)
			return
		}
		req.Method = http.MethodGet

		values := url.Values{}

		values.Add("sesh_key", key)
		values.Add("client_id", req.URL.Query().Get("client_id"))
		values.Add("redirect_uri", req.URL.Query().Get("redirect_uri"))
		values.Add("state", req.URL.Query().Get("state"))

		http.Redirect(res, req, "/oauth/authorize"+"?"+values.Encode(), http.StatusSeeOther)
		return
	}

	h.loginTemplate.Execute(res, true)
}

// Authorize takes a session(access) token and validates it and sents back user info
func (h *OauthHandler) Authorize(res http.ResponseWriter, req *http.Request) {
	// get session key, validate and get user info
	seshKey := strings.TrimSpace(req.URL.Query().Get("sesh_key"))
	userObj, err := auth.ValidateSession(h.dbClient, seshKey)
	if err != nil {
		fmt.Println("couldn't not validate, redirecting to login page: ", err)
		http.Redirect(res, req, "/oauth/login", http.StatusSeeOther)
		return
	}

	// create auth code
	clientID := strings.TrimSpace(req.URL.Query().Get("client_id"))
	authCode, err := auth.CreateAuthorizationCode(h.dbClient, userObj.Id, clientID)
	if err != nil {
		//TODO: handle this error properly
		fmt.Printf("error creating oauth authorization code: %v\n", err)
	}

	// setup redirect url
	redirectURI := strings.TrimSpace(req.URL.Query().Get("redirect_uri"))

	// add query params
	values := url.Values{}
	values.Add("state", req.URL.Query().Get("state"))
	values.Add("code", authCode)

	// redirect
	fmt.Println("auth: redirecting to: ", redirectURI+"?"+values.Encode())
	http.Redirect(res, req, redirectURI+"?"+values.Encode(), http.StatusSeeOther)
}

// Token handles authenticating the oauth client with the given token
func (h *OauthHandler) Token(res http.ResponseWriter, req *http.Request) {
	// authenticate client
	id, sec, ok := req.BasicAuth()
	if !ok {
		fmt.Println("not using basic authentication?")
		return
	}
	fmt.Printf("id: %v & secret: %v\n", id, sec)
	if id != h.clientID || sec != h.clientSecret {
		fmt.Println("incorrect credentials")
		return
	}

	// ^^^^^^^^^^ client is authenticated after above ^^^^^^^^^^
	var queryCode string

	// find grant type: refresh_token or authorization_code
	grantType := req.FormValue("grant_type")

	if strings.ToLower(grantType) == "refresh_token" {
		refreshToken := req.FormValue("refresh_token")
		accessToken, err := auth.FindOauthAccessToken(h.dbClient, refreshToken)
		if err != nil {
			fmt.Println("couldn't find token based on refresh: ", err)
			http.Redirect(res, req, "/oauth/login", http.StatusSeeOther)
			//TODO: fail gracefully??
			return
		}
		queryCode = accessToken.AuthCode

		// delete the token
		go func() {
			err := auth.DeleteOauthAccessToken(h.dbClient, accessToken.Token)
			if err != nil {
				fmt.Println("error oauth handler(Token):", err)
			}
		}()
	} else {
		queryCode = req.FormValue("code")
	}

	// validate auth code
	authCode, err := auth.ValidateAuthCode(h.dbClient, queryCode)
	if err != nil {
		fmt.Println("couldn't find auth code: ", err)
		http.Redirect(res, req, "/oauth/login", http.StatusSeeOther)
		// TODO: send more appropriate error response
		return
	}

	// create access token
	token, err := auth.CreateAccessToken(h.dbClient, authCode)
	if err != nil {
		fmt.Println("error oauth handler(Token), could not create access token:", err)
		// TODO: send error message back
		http.Redirect(res, req, "/oauth/login", http.StatusSeeOther)
		return
	}

	// setup json
	type tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	tRes := &tokenResponse{
		AccessToken:  token.Token,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    3600,
	}

	// marshal data and send off
	json, _ := json.Marshal(&tRes)
	res.Header().Set("Content-Type", "application/json")
	res.Write(json)
}
