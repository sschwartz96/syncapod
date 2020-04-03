package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client holds the connection to the database
type Client struct{ *mongo.Client }

// Connect makes a connection with the database client
func Connect(user, pass, URI string) (*Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(URI).
		SetAuth(options.Credential{Username: user, Password: pass}))
	if err != nil {
		log.Fatal(err)
	}

	return &Client{Client: client}, err
}
