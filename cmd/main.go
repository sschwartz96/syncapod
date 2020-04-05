package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/handler"
)

func main() {
	// read config
	config, err := config.ReadConfig("config.json")
	if err != nil {
		log.Fatal("error reading config: ", err)
	}
	fmt.Println("Running syncapd version: ", config.Version)

	// connect to db
	dbClient, err := database.Connect(config.DbUser, config.DbPass, config.DbURI)
	if err != nil {
		log.Fatal("couldn't connect to db: ", err)
	}

	// setup handler
	handler, err := handler.CreateHandler(dbClient)
	if err != nil {
		log.Fatal("could not setup handlers: ", err)
	}

	// start server
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("couldn't not start server: ", err)
	}

}
