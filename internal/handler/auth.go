package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
)

// AuthHandler handles all authentication into the api
type AuthHandler struct {
	dbClient *database.Client
}

// CreateAuthHandler creates an instance of auth handler
func CreateAuthHandler(dbClient *database.Client) *AuthHandler {
	return &AuthHandler{dbClient: dbClient}
}

// Auth is the api endpoint that handles all authentication and authorization
func (h *APIHandler) Auth(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var head string
		head, req.URL.Path = ShiftPath(req.URL.Path)
		switch head {
		// api/auth/*
		case "authenticate":
			h.Authenticate(res, req)
			return
		case "authorize":
			h.Authorize(res, req)
		}
	} else {
		fmt.Fprint(res, "Protocol not supported")
	}
}

// Authorize takes the request and validates the access token
func (h *APIHandler) Authorize(res http.ResponseWriter, req *http.Request) {
	type AuthRequest struct {
		Token string `json:"token"`
	}
	var authReq AuthRequest
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &authReq)
	if err != nil {
		fmt.Println("error unmarshalling json: ", err)
		fmt.Fprint(res, "invalid json")
		return
	}

	type AuthRes struct {
		Valid bool         `json:"valid"`
		User  *models.User `json:"user"`
	}
	authRes := AuthRes{Valid: false}
	user, err := auth.ValidateSession(h.dbClient, authReq.Token)
	if err == nil {
		authRes.Valid = true
	}
	user.Password = ""
	authRes.User = user

	// marshal json
	response, _ := json.Marshal(&authRes)
	res.Write(response)
}

// Authenticate takes the request and attempts to authenticate the user
func (h *APIHandler) Authenticate(res http.ResponseWriter, req *http.Request) {
	type UserCred struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expires  bool   `json:"expires"`
	}

	var userCred UserCred
	reqBody, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(reqBody, &userCred)
	if err != nil {
		fmt.Println("error unmarshalling json: ", err)
		fmt.Fprint(res, "invalid json")
		return
	}

	// find in db
	user, err := h.dbClient.FindUser(userCred.Email)
	if err != nil {
		fmt.Println("couldn't find user in db: ", err)
		fmt.Fprint(res, "user not found")
		return
	}

	// set up response
	type Response struct {
		Authenticated bool         `json:"authenticated"`
		AccessToken   string       `json:"access_token"`
		User          *models.User `json:"user"`
	}
	response := Response{Authenticated: false, AccessToken: ""}
	// match the passwords
	if auth.Compare(user.Password, userCred.Password) {
		// create the proper session token
		if userCred.Expires {
			response.AccessToken, err = auth.CreateSession(h.dbClient,
				user.ID, time.Hour*24, req.UserAgent())
		} else {
			response.AccessToken, err = auth.CreateSession(h.dbClient,
				user.ID, time.Hour*24*365*10, req.UserAgent())
		}
		// check for error
		if err != nil {
			fmt.Println("error creating user session")
			fmt.Fprint(res, "error creating user session")
			return
		}
		// authenticated true
		response.Authenticated = true
		user.Password = ""
		response.User = user

		// marshal data and send off
		marshalRes, _ := json.Marshal(&response)
		res.Write(marshalRes)
	} else {
		fmt.Fprint(res, "Wrong password")
	}
}
