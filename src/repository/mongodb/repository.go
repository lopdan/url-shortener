package mongo

import (
	"context"
	"time"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/lopdan/url-shortener/src/shortener"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

/** Create new Mongo Client */
func NewMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	// Timeout if there is no response
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	// Read on database connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

/** Creates a Repository for MongoDB */
func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (shortener.RedirectRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	// Timeout if there is no response
	client, err := NewMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepository")
	}
	repo.client = client
	return repo, nil
}

/** Search through the database if the code is correct. */
func (r *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	// Timeout if there is no response
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	redirect := &shortener.Redirect{}
	// Return redirects from database and search the code
	collection := r.client.Database(r.database).Collection("redirects")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	return redirect, nil
}

/** Stores a redirect in database. */
func (r *mongoRepository) Store(redirect *shortener.Redirect) error {
	// Timeout if there is no response
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	// Return redirects from database and insert the redirect to the collection
	collection := r.client.Database(r.database).Collection("redirects")
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":       redirect.Code,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}