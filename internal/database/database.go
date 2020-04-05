package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// these constants are for database actions
const (
	// Databases
	DBsyncapod = "syncapod"

	// Collections
	ColSession = "session"
	ColUser    = "user"
)

// Client holds the connection to the database
type Client struct{ *mongo.Client }

// Connect makes a connection with the database client
func Connect(user, pass, URI string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set up client options
	opts := options.Client().ApplyURI(URI)
	if user != "" {
		opts.Auth.Username = user
		opts.Auth.Password = pass
	}

	// connect to client
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	// confirm the connection with a ping
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client}, err
}

// Insert takes a collection name and interface object and inserts into collection
func (c *Client) Insert(collection string, object interface{}) error {
	col := c.Database(DBsyncapod).Collection(collection)

	res, err := col.InsertOne(context.Background(), object)
	if err != nil {
		return err
	}

	fmt.Println("inserted object successfully with ID: ", res.InsertedID)
	return nil
}

// FindByID takes collection name and pointer to object
func (c *Client) FindByID(collection string, objID primitive.ObjectID, object interface{}) error {
	return c.Find(collection, "_id", objID, object)
}

// Find takes collection, param & value to build fitler, and object pointer
func (c *Client) Find(collection string, param string, value interface{}, object interface{}) error {
	filter := bson.D{{
		Key:   param,
		Value: value,
	}}

	col := c.Database(DBsyncapod).Collection(collection)
	result := col.FindOne(context.Background(), filter)
	err := result.Decode(object)
	if err != nil {
		return err
	}

	return nil
}
