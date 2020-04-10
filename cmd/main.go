package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/handler"
	"github.com/sschwartz96/syncapod/internal/podcast"
)

func main() {

	// read config
	config, err := config.ReadConfig("config.json")
	if err != nil {
		log.Fatal("error reading config: ", err)
	}
	fmt.Println("Running syncapod version: ", config.Version)

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

	// test podcast
	podcast, err := podcast.ParseRSS("https://feeds.twit.tv/twit.xml")
	if err != nil {
		log.Fatal("error parsing rss feed: ", err)
	}

	dbClient.Insert(database.ColPodcast, &podcast)
	
	// start server
	fmt.Println("starting server")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("couldn't not start server: ", err)
	}

}
