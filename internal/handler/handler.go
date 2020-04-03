package handler

import (
	"fmt"
	"net/http"

	"github.com/sschwartz96/syncapod/internal/database"
)

// Handler is the main handler for syncapod
// all routes go through it
type Handler struct {
	dbClient *database.Client
}

// CreateHandler sets up the main handler
func CreateHandler(dbClient *database.Client) (*Handler, error) {
	return &Handler{dbClient: dbClient}, nil
}

// ServeHTTP handles all requests
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello")
}
