package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/melmustafa/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		user, err := cfg.DB.GetUserByApiKey(ctx, apiKey)
		cancel()
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "unauthenticated")
			return
		}
		handler(w, r, user)
	})
}
