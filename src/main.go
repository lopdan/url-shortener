package main

import (
	"fmt"
	"os"
	"os/signal"
	"log"
	"strconv"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"syscall"

	h "github.com/lopdan/url-shortener/src/api"
	mongo "github.com/lopdan/url-shortener/src/repository/mongodb"
	redis "github.com/lopdan/url-shortener/src/repository/redis"
	"github.com/lopdan/url-shortener/src/shortener"
)

func main() {
	// Set up the service
	repo := ChooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)
	// Check services and requests
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)
	// Start the server
	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(HttpPort(), r)

	}()
	// Check exit command
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Printf("Terminated %s", <-errs)
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