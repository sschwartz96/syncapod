package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sschwartz96/syncapod/internal/protos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// these constants are for database actions
const (
	// Databases
	DBsyncapod = "syncapod"

	// Collections
	ColPodcast      = "podcast"
	ColEpisode      = "episode"
	ColSession      = "session"
	ColUser         = "user"
	ColUserEpisode  = "user_episode"
	ColSubscription = "subscription"
	ColAuthCode     = "oauth_auth_code"
	ColAccessToken  = "oauth_access_token"
)

var (
	collections = []string{
		ColPodcast,
		ColEpisode,
		ColSession,
		ColUser,
		ColUserEpisode,
		ColSubscription,
		ColAuthCode,
		ColAccessToken,
	}
)

// mongoClient holds the connection to the database
type mongoClient struct {
	*mongo.Client
	collectionMap map[string]*mongo.Collection
}

// CreateMongoClient makes a connection with the mongo client
func CreateMongoClient(user, pass, URI string) (*mongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set up client options
	opts := options.Client().ApplyURI(URI)
	if user != "" {
		opts.Auth.Username = user
		opts.Auth.Password = pass
	}
	opts = opts.SetRegistry(createRegistry())

	// connect to client
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// confirm the connection with a ping
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &mongoClient{
		Client:        client,
		collectionMap: createCollectionMap(client.Database(DBsyncapod)),
	}, nil
}

func (m *mongoClient) Open(ctx context.Context) error {
	// already opened when using CreateMongoClient
	return nil
}

func (m *mongoClient) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

func (m *mongoClient) FindOne(collection string, object interface{}, filter *Filter, opts *Options) error {
	col := m.collectionMap[collection]
	f := convertToMongoFilter(filter)
	o := convertToFindOneOptions(opts)
	res := col.FindOne(context.Background(), f, o)
	return res.Decode(object)
}

// FindAll finds all within the collection, using filter and options if applicable
func (m *mongoClient) FindAll(collection string, object interface{}, filter *Filter, opts *Options) error {
	col := m.collectionMap[collection]
	f := convertToMongoFilter(filter)
	o := convertToFindOptions(opts)
	cur, err := col.Find(context.Background(), f, o)
	if err != nil {
		return err
	}
	err = cur.All(context.Background(), object)
	return err
}

func (m *mongoClient) Update(collection string, object interface{}, filter *Filter) error {
	col := m.collectionMap[collection]
	f := convertToMongoFilter(filter)
	u := bson.M{"$set": object}
	res, err := col.UpdateOne(context.Background(), f, u)
	if err != nil {
		return err
	}
	if res.MatchedCount > 1 {
		if res.ModifiedCount == 0 {
			return fmt.Errorf("error mongo update: matched %v, but didn't modify", res.MatchedCount)
		}
	} else {
		return fmt.Errorf("error mongo update: did not match any documents")
	}
	return nil
}

// Upsert updates or inserts object within collection with premade filter
func (c *mongoClient) Upsert(collection string, object interface{}, filter *Filter) error {
	col := c.collectionMap[collection]
	update := bson.M{"$set": object}
	f := convertToMongoFilter(filter)

	upsert := true
	opts := &options.UpdateOptions{Upsert: &upsert}

	_, err := col.UpdateOne(context.Background(), f, update, opts)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the certain document based on param and value
func (c *mongoClient) Delete(collection string, filter *Filter) error {
	col := c.collectionMap[collection]
	f := convertToMongoFilter(filter)
	res, err := col.DeleteOne(context.Background(), f)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("error mongo delete: deleted count == 0")
	}
	return nil
}

// convertToMongoFilter converts database.Filter to a bson.M document
func convertToMongoFilter(filter *Filter) interface{} {
	if filter == nil {
		return bson.M{"": ""}
	}
	return bson.M(*filter)
}

// convertToFindOptions converts database.Options to options.FindOptions
func convertToFindOptions(opts *Options) *options.FindOptions {
	if opts == nil {
		return options.Find()
	}
	o := options.Find()
	if opts.limit > 0 {
		o.SetLimit(opts.limit)
	}
	if opts.skip > 0 {
		o.SetSkip(opts.skip)
	}
	if opts.sort != nil {
		o.SetSort(bson.M{opts.sort.key: opts.sort.value})
	}
	return o
}

// convertToMongoOne converts database.Options to options.FindOneOptions
func convertToFindOneOptions(opts *Options) *options.FindOneOptions {
	if opts == nil {
		return options.FindOne()
	}
	o := options.FindOne()
	if opts.skip > 0 {
		o.SetSkip(opts.skip)
	}
	if opts.sort != nil {
		o.SetSort(bson.M{opts.sort.key: opts.sort.value})
	}
	return o
}

// createCollectionMap creates a map of mongo collections so the program doesn't
// reallocate space for a collection every time a request is called
func createCollectionMap(db *mongo.Database) map[string]*mongo.Collection {
	collectionMap := make(map[string]*mongo.Collection, len(collections))
	for _, collection := range collections {
		collectionMap[collection] = db.Collection(collection)
	}
	return collectionMap
}

// Insert takes a collection name and interface object and inserts into collection
func (c *mongoClient) Insert(collection string, object interface{}) error {
	col := c.collectionMap[collection]

	res, err := col.InsertOne(context.Background(), object)
	if err != nil {
		return err
	}

	if res.InsertedID != nil {
		return nil
	}
	return errors.New("failed to insert object into: " + collection)
}

// FindByID takes collection name and pointer to object
func (c *mongoClient) FindByID(collection string, objID *protos.ObjectID, object interface{}) error {
	return c.Find(collection, "_id", objID, object)
}

// Find takes collection, param & value to build filter, and object pointer to receive data
func (c *mongoClient) Find(collection, param string, value interface{}, object interface{}) error {
	filter := bson.D{{
		Key:   param,
		Value: value,
	}}

	return c.FindWithBSON(collection, filter, options.FindOne(), object)
}

// FindWithBSON takes in object and already made bson filter
func (c *mongoClient) FindWithBSON(collection string, filter interface{}, opts *options.FindOneOptions, object interface{}) error {
	var err error

	// get collection
	col := c.collectionMap[collection]

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

// FindAllWithBSON takes collection string, bson filter, options.FindOptions
// and decodes into pointer to the slice
func (c *mongoClient) FindAllWithBSON(collection string, filter interface{}, opts *options.FindOptions, slice interface{}) error {
	// get collection
	col := c.collectionMap[collection]

	// find operation
	cur, err := col.Find(context.Background(), filter, opts)
	if err != nil {
		return err
	}
	// decode all
	err = cur.All(context.Background(), slice)
	return err

}

// UpdateWithBSON takes in collection string & bson filter and update object
func (c *mongoClient) UpdateWithBSON(collection string, filter, update interface{}) error {
	col := c.collectionMap[collection]
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
func (c *mongoClient) ExistsByID(collection string, id *protos.ObjectID) (bool, error) {
	return c.Exists(collection, bson.M{"_id": id})
}

// Exists checks if the document exists within the collection based on the filter
func (c *mongoClient) Exists(collection string, filter interface{}) (bool, error) {
	col := c.collectionMap[collection]

	// setup limit in FindOptions
	limit := int64(1)
	opts := options.FindOptions{Limit: &limit}

	cur, err := col.Find(context.Background(), filter, &opts)
	if err != nil {
		return false, err
	}

	return cur.TryNext(context.Background()), nil
}

// Search takes a collection and search string then finds the object and decodes into object
func (c *mongoClient) Search(collection, search string, fields []string, object interface{}) error {
	col := c.collectionMap[collection]

	// drop any previous indexes
	col.Indexes().DropAll(context.Background())

	// create index model
	var indexes bson.M
	for _, field := range fields {
		indexes[field] = "text"
	}
	indexModel := mongo.IndexModel{Keys: indexes}

	index, err := col.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return fmt.Errorf("error creating index model (mongoClient Search): %v", err)
	}
	fmt.Println("created index named: ", index)

	// create search filter
	filter := bson.M{
		"$text": bson.M{
			"$search": search,
		},
		"score": bson.M{"$meta": "textScore"},
	}

	// sort by score
	opts := options.Find().SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})

	// run search
	cur, err := col.Find(context.Background(), filter, opts)
	if err != nil {
		return err
	}

	return cur.All(context.Background(), object)
}

// Aggregate takes in a collection string, filter, pipeline, and pointer to object
// returns error if anything is malformed
func (c *mongoClient) Aggregate(collection string, pipeline mongo.Pipeline, object interface{}) error {
	col := c.collectionMap[collection]
	cur, err := col.Aggregate(context.Background(), pipeline)
	if err != nil {
		return err
	}
	return cur.All(context.Background(), object)
}
