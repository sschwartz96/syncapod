package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sschwartz96/syncapod/internal/auth"
	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/handler"
)

func main() {
	// test create key
	length := 32
	auth.CreateKey(length)

	// read config
	config, err := config.ReadConfig("config.json")
	if err != nil {
		log.Fatal("error reading config: ", err)
	}
	fmt.Println("Running syncapd version: ", config.Version)

	// connect to db
	fmt.Println("connecting to db")
	dbClient, err := database.Connect(config.DbUser, config.DbPass, config.DbURI)
	if err != nil {
		log.Fatal("couldn't connect to db: ", err)
	}

	fmt.Println("setting up handlers")
	// setup handler
	handler, err := handler.CreateHandler(dbClient)
	if err != nil {
		log.Fatal("could not setup handlers: ", err)
	}

	// start server
	fmt.Println("starting server")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("couldn't not start server: ", err)
	}

}
