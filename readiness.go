package main

import "net/http"

func ready(w http.ResponseWriter, r *http.Request) {
	type status struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, status{
		Status: "ok",
	})
}
