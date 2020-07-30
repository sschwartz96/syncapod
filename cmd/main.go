package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/sschwartz96/syncapod/internal/config"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/handler"
	"github.com/sschwartz96/syncapod/internal/podcast"
	"github.com/sschwartz96/syncapod/internal/protos"
	"github.com/sschwartz96/syncapod/internal/services"
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
	dbClient, err := database.Connect(cfg.DbUser, cfg.DbPass, cfg.DbURI)
	if err != nil {
		log.Fatal("couldn't connect to db: ", err)
	}

	// setup gRPC server
	go startGRPC(cfg, dbClient)

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
		log.Fatal("couldn't not start server: ", err)
	}

}

func redirect(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "https://syncapod.com"+req.RequestURI, http.StatusMovedPermanently)
}

func startGRPC(config *config.Config, dbClient *database.Client) {
	// whether or not we are running on the server
	if config.CertFile != "" {

	}

	// setup tls for grpc
	creds, err := credentials.NewClientTLSFromFile(config.CertFile, "syncapod")
	grpcServer := grpc.NewServer()

	// start listener
	grpcListener, err := net.Listen("tcp", ":"+strconv.Itoa(config.GRPCPort))
	if err != nil {
		log.Fatalf("could not listen on port %d, err: %v", config.GRPCPort, err)
	}

	// register services
	reflection.Register(grpcServer)
	protos.RegisterAuthServer(grpcServer, services.NewAuthService(dbClient))
	protos.RegisterPodcastServiceServer(grpcServer, services.NewPodcastService(dbClient))

	// serve
	err = grpcServer.Serve(grpcListener)
	if err != nil {
		log.Fatal("could not serve services:", err)
	}
}
