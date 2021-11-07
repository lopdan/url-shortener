package redis

import (
	"fmt"
	"strconv"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/lopdan/url-shortener/src/shortener"
)

type redisRepository struct {
	client *redis.Client
}

/** Create new Redis Client */
func NewRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	// Return error if no response
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

/** Creates a Repository for Redis */
func NewRedisRepository(redisURL string) (shortener.RedirectRepository, error) {
	repo := &redisRepository{}
	client, err := NewRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	repo.client = client
	return repo, nil
}

/** Create a key for the data stored in database */
func (r *redisRepository) GenerateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

/** Search through the database if the code is correct. */
func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.GenerateKey(code)
	// Return all keys from database
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	// Add data to the struct
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

/** Stores a redirect in database. */
func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	// Generate the key for the data
	key := r.GenerateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}