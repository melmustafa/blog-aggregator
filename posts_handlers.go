package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/melmustafa/blog-aggregator/internal/database"
)

type JsonPost struct {
	Id          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	FeedID      *uuid.UUID `json:"feed_id"`
}

func (cfg *apiConfig) getPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		log.Printf("something went wrong when parsing the query params. error: %v\n", err)
		log.Println(r.URL.Query())
	}
	if limit == 0 {
		limit = 10
	}
	posts, err := cfg.DB.GetPostsByUser(context.TODO(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		log.Printf("couldn't get the posts from the database with error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	var resPayload []JsonPost
	for _, post := range posts {
		resPayload = append(resPayload, DatabaseToJsonPost(post))
	}
	respondWithJSON(w, http.StatusOK, resPayload)
}

func DatabaseToJsonPost(dbPost database.Post) JsonPost {
	var description *string
	var feedId *uuid.UUID
	var publishedAt *time.Time
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}
	if dbPost.PublishedAt.Valid {
		publishedAt = &dbPost.PublishedAt.Time
	}
	if dbPost.FeedID.Valid {
		feedId = &dbPost.FeedID.UUID
	}
	return JsonPost{
		Id:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		Description: description,
		PublishedAt: publishedAt,
		FeedID:      feedId,
	}
}
