package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dnieln7/go-examples/rss-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT not found")
	}

	apiConfig := setupDatabase()
	router := setupRouter()

	router.Get("/timestamp", getTimestamp)
	router.Post("/users", apiConfig.postUser)
	router.Get("/users", apiConfig.middlewareAuth(apiConfig.getUser))
	router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.postFeed))
	router.Get("/feeds", apiConfig.getFeeds)
	router.Post("/feeds/follows", apiConfig.middlewareAuth(apiConfig.postFeedFollow))
	router.Get("/feeds/follows", apiConfig.middlewareAuth(apiConfig.getFeedFollows))
	router.Delete("/feeds/follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.deleteFeedFollow))
	router.Get("/posts", apiConfig.middlewareAuth(apiConfig.getPostsByUserID))

	setupWorkers(apiConfig.DB)

	log.Println("Starting server on port: ", port)
	
	http.ListenAndServe(":"+port, router)
}

func setupDatabase() *ApiConfig {
	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("DB_URL not found")
	}

	connection, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("Could not connect to database")
	}

	queries := database.New(connection)

	return &ApiConfig{
		DB: queries,
	}
}

func setupRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	return router
}

type ApiConfig struct {
	DB *database.Queries
}

func setupWorkers(db *database.Queries)  {
	go startScraping(db, 2, time.Minute)
}
