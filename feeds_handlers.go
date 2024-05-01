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

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserId        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func (cfg *apiConfig) getFeeds(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	retrievedFeeds, err := cfg.DB.GetNextFeedsToFetch(ctx, 20)
	cancel()
	if err != nil {
		log.Printf("coudn't get the feeds with error %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't get feeds")
		return
	}

	feeds := []Feed{}
	for _, feed := range retrievedFeeds {
		feeds = append(feeds, databaseFeedToJson(feed))
	}
	respondWithJSON(w, http.StatusOK, feeds)
}

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
	createdFeedFollow, _ := cfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: createdFeed.ID,
	})
	cancel()
	if err != nil {
		log.Printf("coudn't create a new feed with error %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create feed")
		return
	}

	type responsePayload struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}
	respondWithJSON(w, http.StatusCreated, responsePayload{
		Feed:       databaseFeedToJson(createdFeed),
		FeedFollow: databaseFeedFollowToJson(createdFeedFollow),
	})
}

func databaseFeedToJson(feed database.Feed) Feed {
	var timeValue *time.Time
	if feed.LastFetchedAt.Valid {
		timeValue = &feed.LastFetchedAt.Time
	}
	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserId:        feed.UserID,
		LastFetchedAt: timeValue,
	}
}
