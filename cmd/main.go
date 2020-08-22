package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
	sGRPC "github.com/sschwartz96/syncapod/internal/grpc"
	"github.com/sschwartz96/syncapod/internal/handler"
	"github.com/sschwartz96/syncapod/internal/podcast"
)

func main() {
	// read config
	cfg, err := config.ReadConfig("config.json")
	if err != nil {
		log.Fatal("error reading config: ", err)
	}
	fmt.Println("Running syncapod version: ", cfg.Version)

	// connect to db
	fmt.Println("connecting to db")
	dbClient, err := database.CreateMongoClient(cfg.DbUser, cfg.DbPass, cfg.DbURI)
	if err != nil {
		log.Fatal("couldn't connect to db: ", err)
	}

	// setup & start gRPC server
	grpcServer := sGRPC.NewServer(cfg, dbClient)
	go func() {
		err = grpcServer.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// start updating podcasts
	go podcast.UpdatePodcasts(dbClient)

	fmt.Println("setting up handlers")
	// setup handler
	handler, err := handler.CreateHandler(dbClient, cfg)
	if err != nil {
		log.Fatal("could not setup handlers: ", err)
	}

	// start server
	fmt.Println("starting server")
	port := strings.TrimSpace(strconv.Itoa(cfg.Port))
	if cfg.Port == 443 {
		// setup redirect server
		go func() {
			if err = http.ListenAndServe(":80", http.HandlerFunc(redirect)); err != nil {
				log.Fatalf("redirect server failed %v\n", err)
			}
		}()

		err = http.ListenAndServeTLS(":"+port, cfg.CertFile, cfg.KeyFile, handler)
	} else {
		err = http.ListenAndServe(":"+port, handler)
	}

	if err != nil {
		log.Fatal("couldn't not start server:", err)
	}

}

func redirect(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "https://syncapod.com"+req.RequestURI, http.StatusMovedPermanently)
}
