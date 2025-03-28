package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mshortcodes/color_my_practice/internal/auth"
	"github.com/mshortcodes/color_my_practice/internal/database"
)

type PracticeLog struct {
	Id         uuid.UUID `json:"id"`
	Date       string    `json:"date"`
	ColorDepth int32     `json:"color_depth"`
	Confirmed  bool      `json:"confirmed"`
	UserID     uuid.UUID `json:"user_id"`
}

// handlerLogsCreate creates a practice log for a single day.
func (cfg *apiConfig) handlerLogsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Date       string `json:"date"`
		ColorDepth int32  `json:"color_depth"`
	}

	type response struct {
		PracticeLog
	}

	var params parameters

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	parsedDate, err := time.Parse(time.DateOnly, params.Date)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid date format", err)
		return
	}

	accessToken, err := r.Cookie("jwt")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken.Value, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid JWT", err)
		return
	}

	practiceLog, err := cfg.db.CreateLog(r.Context(), database.CreateLogParams{
		Date:       parsedDate,
		ColorDepth: params.ColorDepth,
		UserID:     userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create practice log", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		PracticeLog: PracticeLog{
			Id:         practiceLog.ID,
			Date:       practiceLog.Date.Format(time.DateOnly),
			ColorDepth: practiceLog.ColorDepth,
			Confirmed:  practiceLog.Confirmed,
			UserID:     practiceLog.UserID,
		},
	})
}
