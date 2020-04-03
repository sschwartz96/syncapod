package main

import (
	"fmt"
	"log"

	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
)

func main() {
	// read config
	config, err := config.ReadConfig("config.json")
	if err != nil {
		log.Fatal("error reading config: ", err)
	}
	fmt.Println("Running syncapd version: ", config.Version)

	// connect to db
	dbClient, err := database.Connect(config.DbUser, config.DbPass, config.URI)

	fmt.Println("connected to db: ", dbClient)

}
