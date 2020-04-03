package handler

import "github.com/sschwartz96/syncapod/internal/database"

// Handler is the main handler for syncapod
// all routes go through it
type Handler struct {
	database database.Client
}
