package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}

// handlerUsersCreate creates a User.
func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type response struct {
		User
	}

	var params parameters

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
