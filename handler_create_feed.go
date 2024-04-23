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

func (cfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type paramaters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
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
	createdFeed, err := cfg.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   params.Name,
		Url:    params.Url,
		UserID: user.ID,
	})
	cancel()
	if err != nil {
		log.Printf("coudn't create a new feed with error %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create feed")
		return
	}

	type responsePayload struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserId    uuid.UUID `json:"user_id"`
	}
	respondWithJSON(w, http.StatusCreated, responsePayload{
		ID:        createdFeed.ID,
		CreatedAt: createdFeed.CreatedAt,
		UpdatedAt: createdFeed.UpdatedAt,
		Name:      createdFeed.Name,
		Url:       createdFeed.Url,
		UserId:    createdFeed.UserID,
	})
}
