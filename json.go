package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Msg string `json:"message"`
	}
	respondWithJSON(w, code, errorResponse{
		Msg: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling to json: %s\n", err)
		log.Printf("payload: %v\n", payload)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(status)
	w.Write(body)
}
