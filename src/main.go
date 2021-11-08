package main

import (
	"fmt"
	"os"
	"log"
	"strconv"

	mongo "github.com/lopdan/url-shortener/src/repository/mongodb"
	redis "github.com/lopdan/url-shortener/src/repository/redis"
	"github.com/lopdan/url-shortener/src/shortener"
)

func main() {

}

/** Port in 8000*/
func HttpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

/** Check type of database */
func ChooseRepo() shortener.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		// Look for enviroment variable 
		redisURL := os.Getenv("REDIS_URL")
		repo, err := redis.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		// Look for enviroment variables 
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongo.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}