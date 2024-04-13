package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("couldn't load .env file")
	}

	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/readiness", ready)
	mux.HandleFunc("GET /v1/err", errCheck)

	corsMux := corsMiddleware(mux)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
