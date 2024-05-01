package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/melmustafa/blog-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("couldn't load .env file with error %s\n", err)
	}

	dbUrl := os.Getenv("DATABASE_URI")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("couldn't connect to the database with error %s\n", err)
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		DB: dbQueries,
	}

	go apiCfg.FetchService(10, time.Second*10)
	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/readiness", ready)
	mux.HandleFunc("GET /v1/err", errCheck)
	mux.HandleFunc("POST /v1/users", apiCfg.createUser)
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.getUser))
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.createFeed))
	mux.HandleFunc("GET /v1/feeds", apiCfg.getFeeds)
	mux.HandleFunc("POST /v1/feed-follows", apiCfg.middlewareAuth(apiCfg.createFeedFollow))
	mux.HandleFunc("GET /v1/feed-follows", apiCfg.middlewareAuth(apiCfg.getFeedFollows))
	mux.HandleFunc(
		"DELETE /v1/feed-follows/{feedFollowID}",
		apiCfg.middlewareAuth(apiCfg.deleteFeedFollows),
	)
	corsMux := corsMiddleware(mux)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
