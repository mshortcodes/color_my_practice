package main

import (
	"net/http"
	"time"

	"github.com/mshortcodes/color_my_practice/internal/auth"
)

// handlerRefresh creates a new JWT for the given user
// after validating the refresh token.
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't find token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get user", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't make new JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

// handlerRevoke revokes a refresh token.
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get token", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
