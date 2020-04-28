package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/handler"
	"github.com/sschwartz96/syncapod/internal/podcast"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// tests
	id, _ := primitive.ObjectIDFromHex("5e895b2433b810425c9d1611")
	subs := dbClient.FindUserSubs(id)
	fmt.Println(subs[0].Podcast.Title)
	fmt.Println(subs[0].CurEpi.Title)
	fmt.Println(subs[0].CurEpiDetails.Offset)

	//err = podcast.AddNewPodcast(dbClient, "http://joeroganexp.joerogan.libsynpro.com/rss")
	//if err != nil {
	//	fmt.Println("error adding new podcast: ", err)
	//}

	//err = podcast.AddNewPodcast(dbClient, "http://feeds.99percentinvisible.org/99percentinvisible")
	//if err != nil {
	//	fmt.Println("error adding new podcast: ", err)
	//}

	//err = podcast.AddNewPodcast(dbClient, "http://feeds.twit.tv/twit.xml")
	//if err != nil {
	//	fmt.Println("error adding new podcast: ", err)
	//}

	// start updating podcasts
	go podcast.UpdatePodcasts(dbClient)

	fmt.Println("setting up handlers")
	// setup handler
	handler, err := handler.CreateHandler(dbClient, config)
	if err != nil {
		log.Fatal("could not setup handlers: ", err)
	}

	// start server
	fmt.Println("starting server")
	port := strings.TrimSpace(strconv.Itoa(config.Port))
	if config.Port == 443 {
		// setup redirect server
		go func() {
			if err = http.ListenAndServe(":80", http.HandlerFunc(redirect)); err != nil {
				log.Fatalf("redirect server failed %v\n", err)
			}
		}()

		err = http.ListenAndServeTLS(":"+port, config.CertFile, config.KeyFile, handler)
	} else {
		err = http.ListenAndServe(":"+port, handler)
	}

	if err != nil {
		log.Fatal("couldn't not start server: ", err)
	}

}

func redirect(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "https://syncapod.com"+req.RequestURI, http.StatusMovedPermanently)
}
