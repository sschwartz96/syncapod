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

	opts := options.Client().ApplyURI(URI)
	if user != "" {
		opts.Auth.Username = user
		opts.Auth.Password = pass
	}
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client}, err
}
