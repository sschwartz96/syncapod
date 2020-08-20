package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sschwartz96/syncapod/internal/models"
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
	colPodcast      = "podcast"
	colEpisode      = "episode"
	colSession      = "session"
	colUser         = "user"
	colUserEpisode  = "user_episode"
	colSubscription = "subscription"
	colAuthCode     = "oauth_auth_code"
	colAccessToken  = "oauth_access_token"
)

var (
	collections = []string{
		colPodcast,
		colEpisode,
		colSession,
		colUser,
		colUserEpisode,
		colSubscription,
		colAuthCode,
		colAccessToken,
	}
)

// MongoClient holds the connection to the database
type MongoClient struct {
	*mongo.Client
	collectionMap map[string]*mongo.Collection
}

// ConnectMongo makes a connection with the mongo client
func ConnectMongo(user, pass, URI string) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set up client options
	opts := options.Client().ApplyURI(URI)
	if user != "" {
		opts.Auth.Username = user
		opts.Auth.Password = pass
	}
	opts.SetRegistry(createRegistry())

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

	return &MongoClient{
		Client:        client,
		collectionMap: createCollectionMap(client.Database(DBsyncapod)),
	}, nil
}

func createCollectionMap(db *mongo.Database) map[string]*mongo.Collection {
	collectionMap := make(map[string]*mongo.Collection, len(collections))
	for _, collection := range collections {
		collectionMap[collection] = collectionMap[collection]
	}
	return collectionMap
}

// Insert takes a collection name and interface object and inserts into collection
func (c *MongoClient) Insert(collection string, object interface{}) error {
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

// Delete deletes the certain document based on param and value
func (c *MongoClient) Delete(collection, param string, value interface{}) error {
	filter := bson.D{{
		Key:   param,
		Value: value,
	}}

	res, err := c.collectionMap[collection].DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	fmt.Printf("successfully deleted: %v documents\n", res.DeletedCount)
	return nil
}

// FindByID takes collection name and pointer to object
func (c *MongoClient) FindByID(collection string, objID *protos.ObjectID, object interface{}) error {
	return c.Find(collection, "_id", objID, object)
}

// Find takes collection, param & value to build filter, and object pointer to receive data
func (c *MongoClient) Find(collection, param string, value interface{}, object interface{}) error {
	filter := bson.D{{
		Key:   param,
		Value: value,
	}}

	return c.FindWithBSON(collection, filter, options.FindOne(), object)
}

// FindAll finds all objects in the collection and inserts them into provided slice
// returns error if the operation fails
func (c *MongoClient) FindAll(collection string, slice interface{}) error {
	col := c.collectionMap[collection]
	cur, err := col.Find(context.Background(), bson.D{{}})
	if err != nil {
		return err
	}
	err = cur.All(context.Background(), slice)
	return err
}

// Upsert updates or inserts object within collection with premade filter
func (c *MongoClient) Upsert(collection string, filter interface{}, object interface{}) error {
	col := c.collectionMap[collection]
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
func (c *MongoClient) FindWithBSON(collection string, filter interface{}, opts *options.FindOneOptions, object interface{}) error {
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
func (c *MongoClient) FindAllWithBSON(collection string, filter interface{}, opts *options.FindOptions, slice interface{}) error {
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

// UpdateWithBSON takes in collection string & bson filter and update objects
func (c *MongoClient) UpdateWithBSON(collection string, filter, update interface{}) error {
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
func (c *MongoClient) ExistsByID(collection string, id *protos.ObjectID) (bool, error) {
	return c.Exists(collection, bson.M{"_id": id})
}

// Exists checks if the document exists within the collection based on the filter
func (c *MongoClient) Exists(collection string, filter interface{}) (bool, error) {
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
func (c *MongoClient) Search(collection, search string, object interface{}) error {
	col := c.collectionMap[collection]
	// TODO: maybe dont drop if the index exists?
	col.Indexes().DropAll(context.Background())

	// create index
	indexModel := mongo.IndexModel{Keys: bson.D{
		{Key: "title", Value: "text"},
		{Key: "keywords", Value: "text"},
		{Key: "subtitle", Value: "text"},
	}}
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

// Aggregate takes in a collection string, filter, pipeline, and pointer to object
// returns error if anything is malformed
func (c *MongoClient) Aggregate(collection string, pipeline mongo.Pipeline, object interface{}) error {
	col := c.collectionMap[collection]
	cur, err := col.Aggregate(context.Background(), pipeline)
	if err != nil {
		return err
	}
	return cur.All(context.Background(), object)
}

// * Database *
func (c *MongoClient) Open(ctx context.Context) error {
	fmt.Println("connection should be opened")
	return nil
}

func (c *MongoClient) Close(ctx context.Context) error {
	if err := c.Disconnect(ctx); err != nil {
		return fmt.Errorf("error disconnecting mongo client: %v", err)
	}
	return nil
}

// * oAuth *
func (c *MongoClient) InsertAuthCode(code *models.AuthCode) error {
	if err := c.Insert(colAuthCode, code); err != nil {
		return fmt.Errorf("error inserting auth code: %v", err)
	}
	return nil
}

func (c *MongoClient) FindAuthCode(code string) (*models.AuthCode, error) {
	var authCode *models.AuthCode
	err := c.Find(colAuthCode, "auth_code", &code, authCode)
	if err != nil {
		return nil, fmt.Errorf("error finding auth code: %v", err)
	}
	return authCode, nil
}

func (c *MongoClient) InsertAccessToken(token *models.AccessToken) error {
	if err := c.Insert(colAccessToken, token); err != nil {
		return fmt.Errorf("error inserting access token: %v", err)
	}
	return nil
}

func (c *MongoClient) FindAccessToken(token string) (*models.AccessToken, error) {
	var accessToken *models.AccessToken
	err := c.Find(colAccessToken, "token", &token, accessToken)
	if err != nil {
		return nil, fmt.Errorf("error finding access token: %v", err)
	}
	return accessToken, nil
}

// * Auth *
func (c *MongoClient) FindSession(key string) (*protos.Session, error) {
	var session *protos.Session
	err := c.Find(colSession, "sessionkey", &key, session)

	if err != nil {
		return nil, fmt.Errorf("error finding session: %v", err)
	}
	return session, nil
}

func (c *MongoClient) UpsertSession(session *protos.Session) error {
	if err := c.Upsert(colSession, bson.M{"_id": session.Id}, session); err != nil {
		return fmt.Errorf("error upserting session: %v", err)
	}
	return nil
}

func (c *MongoClient) DeleteSession(id *protos.ObjectID) error {
	err := c.Delete(colSession, "_id", id)
	return fmt.Errorf("error deleting session: %v", err)
}

func (c *MongoClient) FindUserByID(id *protos.ObjectID) (*protos.User, error) {
	var user *protos.User
	err := c.Find(colUser, "_id", id, user)
	if err != nil {
		return nil, fmt.Errorf("error finding user by id: %v", err)
	}
	return user, nil
}

// FindUser attempts to find user by username/email returns pointer to user or error if not found
func (c *MongoClient) FindUser(username string) (*protos.User, error) {
	var param string
	if strings.Contains(username, "@") {
		param = "email"
		username = strings.ToLower(username)
	} else {
		param = "username"
	}

	var user protos.User
	err := c.Find(colUser, param, username, &user)
	if err != nil {
		return nil, fmt.Errorf("error finding user by %s: %v", param, err)
	}

	return &user, nil
}

func (c *MongoClient) DeleteUser(id *protos.ObjectID) error {
	if err := c.Delete(colUser, "_id", id); err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

// * Podcast *
func (c *MongoClient) FindAllPodcasts() ([]*protos.Podcast, error) {
	// TODO: get rid of?
	var podcasts []*protos.Podcast
	err := c.FindAll(colPodcast, &podcasts)
	if err != nil {
		return nil, fmt.Errorf("error finding all podcasts: %v", err)
	}
	return podcasts, nil
}

func (c *MongoClient) FindPodcastsByRange(start, end int) ([]*protos.Podcast, error) {
	var podcasts []*protos.Podcast
	filter := bson.D{{}}
	opts := options.Find().SetLimit(int64(end - start)).SetSkip(int64(start)).SetSort(
		bson.M{"pubdate": -1},
	)
	err := c.FindAllWithBSON(colEpisode, filter, opts, &podcasts)
	if err != nil {
		return podcasts, fmt.Errorf("error finding podcasts within range %d - %d: %v", start, end, err)
	}
	return podcasts, nil
}

func (c *MongoClient) FindPodcastByID(id *protos.ObjectID) (*protos.Podcast, error) {
	var podcast *protos.Podcast
	if err := c.Find(colPodcast, "_id", id, podcast); err != nil {
		return nil, fmt.Errorf("error finding podcast by id: %v", err)
	}
	return podcast, nil
}

// FindEpisodes returns a list of episodes based on podcast id
// returns in chronological order, sectioned by start & end
// * Episode *
func (c *MongoClient) FindEpisodesByRange(podcastID *protos.ObjectID, start int, end int) ([]*protos.Episode, error) {
	var episodes []*protos.Episode
	filter := bson.M{"podcastid": podcastID}
	opts := options.Find().SetLimit(int64(end - start)).SetSkip(int64(start)).SetSort(
		bson.M{"pubdate": -1},
	)
	err := c.FindAllWithBSON(colEpisode, filter, opts, &episodes)
	if err != nil {
		return nil, fmt.Errorf("error finding episodes by range %d - %d: %v", start, end, err)
	}
	return episodes, nil
}

func (c *MongoClient) FindAllEpisodes(podcastID *protos.ObjectID) ([]*protos.Episode, error) {
	var episodes []*protos.Episode
	filter := bson.M{"podcastid": podcastID}
	opts := options.Find().SetSort(bson.M{"pubdate": -1})
	err := c.FindAllWithBSON(colEpisode, filter, opts, &episodes)
	if err != nil {
		return nil, fmt.Errorf("error finding all episodes: %v", err)
	}
	return episodes, nil
}

func (c *MongoClient) FindLatestEpisode(podcastID *protos.ObjectID) (*protos.Episode, error) {
	var episode *protos.Episode
	col := c.collectionMap[colEpisode]
	filter := bson.M{"podcastid": podcastID}
	opts := options.FindOne().SetSort(bson.M{"pubdate": -1})
	res := col.FindOne(context.Background(), filter, opts)
	if err := res.Decode(&episode); err != nil {
		return nil, fmt.Errorf("error finding latest episode: %v", err)
	}
	return episode, nil
}

func (c *MongoClient) FindEpisodeByID(id *protos.ObjectID) (*protos.Episode, error) {
	var episode *protos.Episode
	err := c.Find(colEpisode, "_id", id, &episode)
	if err != nil {
		return nil, fmt.Errorf("error finding episode by id: %v", err)
	}
	return episode, nil
}

// FindEpisodeBySeason takes a season episode number returns error if not found
func (c *MongoClient) FindEpisodeBySeason(id *protos.ObjectID, seasonNum int, episodeNum int) (*protos.Episode, error) {
	var episode protos.Episode
	filter := bson.D{
		{Key: "podcast_id", Value: id},
		{Key: "season", Value: seasonNum},
		{Key: "episode", Value: episodeNum},
	}
	err := c.FindWithBSON(colEpisode, filter, nil, &episode)

	return &episode, err
}

func (c *MongoClient) UpsertEpisode(episode *protos.Episode) error {
	panic("not implemented") // TODO: Implement
}

// * UserEpisode *
func (c *MongoClient) FindUserEpisode(userID *protos.ObjectID, episodeID *protos.ObjectID) (*protos.UserEpisode, error) {
	panic("not implemented") // TODO: Implement
}

func (c *MongoClient) FindLatestUserEpisode(userID *protos.ObjectID) (*protos.UserEpisode, error) {
	panic("not implemented") // TODO: Implement
}

func (c *MongoClient) UpsertUserEpisode(userEpisode *protos.UserEpisode) error {
	panic("not implemented") // TODO: Implement
}

// Subscriptions
func (c *MongoClient) FindSubscriptions(userID *protos.ObjectID) ([]*protos.Podcast, error) {
	panic("not implemented") // TODO: Implement
}

func (c *MongoClient) UpsertSubscription(subscription *protos.Subscription) error {
	panic("not implemented") // TODO: Implement
}
