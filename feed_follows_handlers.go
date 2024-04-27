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

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FeedId    uuid.UUID `json:"feed_id"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getFeedFollows(w http.ResponseWriter, _ *http.Request, user database.User) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	retrievedFeedFollows, err := cfg.DB.GetFeedFollows(ctx, user.ID)
	cancel()
	if err != nil {
		log.Printf("coudn't get the feed follows with error %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't get feed follows")
		return
	}

	feedFollows := []FeedFollow{}
	for _, feedFollow := range retrievedFeedFollows {
		feedFollows = append(feedFollows, databaseFeedFollowToJson(feedFollow))
	}
	respondWithJSON(w, http.StatusOK, feedFollows)
}

func (cfg *apiConfig) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type paramaters struct {
		FeedId uuid.UUID `json:"feed_id"`
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
	createdFeedFollow, err := cfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: params.FeedId,
	})
	cancel()
	if err != nil {
		log.Printf("coudn't create a new feed follow with error %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToJson(createdFeedFollow))
}

func (cfg *apiConfig) deleteFeedFollows(w http.ResponseWriter, r *http.Request, _ database.User) {
	id, err := uuid.Parse(r.PathValue("feedFollowID"))
	if err != nil {
		log.Printf("coudn't parse feed follow id with error %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't get the feed follow id")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	err = cfg.DB.DeleteFeedFollow(ctx, id)
	cancel()
	if err != nil {
		log.Printf("coudn't delete the feed follow with error %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't delete feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func databaseFeedFollowToJson(feedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserId:    feedFollow.UserID,
		FeedId:    feedFollow.FeedID,
	}
}
