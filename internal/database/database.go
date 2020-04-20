package database

import (
	"context"
	"errors"
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
	ColEpisode     = "episode"
	ColSession     = "session"
	ColUser        = "user"
	ColUserEpisode = "user_episode"
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

	if res.InsertedID != nil {
		return nil
	}
	return errors.New("failed to insert object into: " + collection)
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

	return c.FindWithBSON(collection, filter, options.FindOne(), object)
}

// FindAll finds all objects in the collection and inserts them into provided slice
// returns error if the operation fails
func (c *Client) FindAll(collection string, slice interface{}) error {
	col := c.Database(DBsyncapod).Collection(collection)
	cur, err := col.Find(context.Background(), bson.D{{}})
	if err != nil {
		return err
	}
	err = cur.All(context.Background(), slice)
	return err
}

// Upsert updates or inserts object within collection with premade filter
func (c *Client) Upsert(collection string, filter interface{}, object interface{}) error {
	col := c.Database(DBsyncapod).Collection(collection)
	update := bson.M{"$set": object}

	upsert := true
	opts := &options.UpdateOptions{
		Upsert: &upsert,
	}

	res, err := col.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}

	fmt.Println("result: ", res)
	return nil
}

// FindWithBSON takes in object and already made bson filter
func (c *Client) FindWithBSON(collection string, filter interface{}, opts *options.FindOneOptions, object interface{}) error {
	var err error

	// get collection
	col := c.Database(DBsyncapod).Collection(collection)

	// find operation
	if opts == nil {
		opts = options.FindOne()
	}
	result := col.FindOne(context.Background(), filter, opts)
	err = result.Err()
	if err != nil {
		return err
	}
	// decode one
	err = result.Decode(object)

	return err
}

func (c *Client) FindAllWithBSON(collection string, filter interface{}, opts *options.FindOptions, slice interface{}) error {
	// get collection
	col := c.Database(DBsyncapod).Collection(collection)

	// find operation
	cur, err := col.Find(context.Background(), filter, opts)
	if err != nil {
		return err
	}
	// decode all
	err = cur.All(context.Background(), slice)
	return err

}

// UpdateWithBSON takes in collection string & bson filter and update objects
func (c *Client) UpdateWithBSON(collection string, filter, update interface{}) error {
	col := c.Database(DBsyncapod).Collection(collection)
	r, err := col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if r.ModifiedCount != 1 {
		if r.MatchedCount == 1 {
			return errors.New("matched but not updated")
		}
		return errors.New("object failed to update")
	}
	return nil
}

// ExistsByID attempts to find a document in the collection based on its ID
func (c *Client) ExistsByID(collection string, id primitive.ObjectID) (bool, error) {
	return c.Exists(collection, bson.M{"_id": id})
}

// Exists checks if the document exists within the collection based on the filter
func (c *Client) Exists(collection string, filter interface{}) (bool, error) {
	col := c.Database(DBsyncapod).Collection(collection)

	// setup limit in FindOptions
	limit := int64(1)
	opts := options.FindOptions{Limit: &limit}

	cur, err := col.Find(context.Background(), filter, &opts)
	if err != nil {
		return false, err
	}

	return cur.TryNext(context.Background()), nil
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

func (c *Client) Search(collection, search string, object interface{}) error {
	col := c.Database(DBsyncapod).Collection(collection)
	// TODO: maybe dont drop if the index exists?
	col.Indexes().DropAll(context.Background())

	// create index
	indexModel := mongo.IndexModel{Keys: bson.D{{"title", "text"}, {"keywords", "text"}, {"subtitle", "text"}}}
	index, err := col.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		fmt.Println("couldn't create index model: ", err)
	}
	fmt.Println("our index name: ", index)

	// create search filter
	filter := bson.M{"$text": bson.M{"$search": search}}

	// run search
	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		return err
	}

	return cur.All(context.Background(), object)
}
