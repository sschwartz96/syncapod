package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sschwartz96/syncapod/internal/models"
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
	ColPodcast     = "podcast"
	ColSession     = "session"
	ColUser        = "user"
	ColAuthCode    = "auth_code"
	ColAccessToken = "access_token"
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

// Delete deletes the certain document based on param and value
func (c *Client) Delete(collection, param string, value interface{}) error {
	filter := bson.D{{
		Key:   param,
		Value: value,
	}}

	res, err := c.Database(DBsyncapod).Collection(collection).DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	fmt.Printf("successfully deleted: %v documents\n", res.DeletedCount)
	return nil
}

// FindByID takes collection name and pointer to object
func (c *Client) FindByID(collection string, objID primitive.ObjectID, object interface{}) error {
	return c.Find(collection, "_id", objID, object)
}

// Find takes collection, param & value to build fitler, and object pointer
func (c *Client) Find(collection, param string, value interface{}, object interface{}) error {
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

// FindUser attempts to find user by username/email returns pointer to user or error if not found
func (c *Client) FindUser(username string) (*models.User, error) {
	var param string
	if strings.Contains(username, "@") {
		param = "email"
		username = strings.ToLower(username)
	} else {
		param = "username"
	}

	var user models.User
	err := c.Find(ColUser, param, username, &user)

	return &user, err
}
