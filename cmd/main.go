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
)

func main() {
	// test the find length
	//	dur := podcast.FindLength("http://traffic.libsyn.com/joeroganexp/p1458.mp3")
	//	fmt.Printf("duration: %v\n\n", dur)
	//	dur = podcast.FindLength("https://dts.podtrac.com/redirect.mp3/media.blubrry.com/99percentinvisible/dovetail.prxu.org/96/905fcf0b-2c8b-4c2a-9832-a3e0b10d0472/01_398_Unsheltered_in_Place_pt01.mp3")
	//	fmt.Printf("duration: %v\n\n", dur)

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
