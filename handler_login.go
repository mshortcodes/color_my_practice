package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mshortcodes/color_my_practice/internal/auth"
	"github.com/mshortcodes/color_my_practice/internal/database"
)

// handlerLogin authenticates a user by checking their password.
// Once authenticated, a JWT and refresh token are created.
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "no user with that email", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't make refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't save refresh token", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    accessToken,
		Path:     "/api/",
		Expires:  time.Now().UTC().Add(time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/",
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 60),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
