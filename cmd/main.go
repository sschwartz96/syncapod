package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sschwartz96/stockpile/db"
	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
	sGRPC "github.com/sschwartz96/syncapod/internal/grpc"
	"github.com/sschwartz96/syncapod/internal/handler"
	"github.com/sschwartz96/syncapod/internal/podcast"
	"github.com/sschwartz96/syncapod/internal/services"
)

func main() {
	// read config
	cfg, err := readConfig("config.json")
	if err != nil {
		log.Fatal("Main() error, could not read config: ", err)
	}

	log.Println("Running syncapod version: ", cfg.Version)

	// connect to db
	log.Println("connecting to db")
	dbClient, err := database.NewMongoClient(cfg)
	if err != nil {
		log.Fatal("couldn't connect to db: ", err)
	}

	// setup & start gRPC server
	grpcServer := sGRPC.NewServer(cfg, dbClient,
		services.NewAuthService(dbClient),
		services.NewPodcastService(dbClient),
	)
	go func() {
		// setup listener
		grpcListener, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.GRPCPort))
		if err != nil {
			log.Fatalf("could not listen on port %d, err: %v", cfg.GRPCPort, err)
		}
		// start server
		err = grpcServer.Start(grpcListener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// start updating podcasts
	go updatePodcasts(dbClient)

	log.Println("setting up handlers")
	// setup handler
	handler, err := handler.CreateHandler(dbClient, cfg)
	if err != nil {
		log.Fatal("could not setup handlers: ", err)
	}

	// start server
	log.Println("starting server")
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

func updatePodcasts(dbClient db.Database) {
	for {
		err := podcast.UpdatePodcasts(dbClient)
		if err != nil {
			log.Println("main/updatePodcasts() error:", err)
		}
		time.Sleep(time.Minute * 15)
	}
}

func redirect(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "https://syncapod.com"+req.RequestURI, http.StatusMovedPermanently)
}

func readConfig(path string) (*config.Config, error) {
	cfgFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("readConfig() error opening file: %v", err)
	}
	return config.ReadConfig(cfgFile)
}
