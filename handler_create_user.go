package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/melmustafa/blog-aggregator/internal/database"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Name string `json:"name"`
	}
	params := paramaters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("coudn't unmarshal the request with error %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't unmarshal the request")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	createdUser, err := cfg.DB.CreateUser(ctx, database.CreateUserParams{
		ID:   uuid.New(),
		Name: params.Name,
	})
	cancel()
	if err != nil {
		log.Printf("coudn't create a new user with error %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	type responsePayload struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}
	respondWithJSON(w, http.StatusCreated, responsePayload{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Name:      createdUser.Name,
		ApiKey:    createdUser.ApiKey,
	})
}
