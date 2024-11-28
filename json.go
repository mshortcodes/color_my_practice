package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithJSON marshals the payload and writes the response.
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(code)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// respondWithError logs an error and calls respondWithJSON.
func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}
