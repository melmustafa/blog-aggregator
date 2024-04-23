package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getFeeds(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	retrievedFeeds, err := cfg.DB.GetFeeds(ctx)
	cancel()
	if err != nil {
		log.Printf("coudn't get the feeds with error %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't get feeds")
		return
	}

	type ResponsePayload struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserId    uuid.UUID `json:"user_id"`
	}
	feeds := []ResponsePayload{}
	for _, feed := range retrievedFeeds {
		feeds = append(feeds, ResponsePayload{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserId:    feed.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, feeds)
}
