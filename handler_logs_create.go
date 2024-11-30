package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mshortcodes/color_my_practice/internal/database"
)

type PracticeLog struct {
	Id         uuid.UUID `json:"-"`
	Date       string    `json:"date"`
	ColorDepth int32     `json:"color_depth"`
	Confirmed  bool      `json:"-"`
	UserID     uuid.UUID `json:"-"`
}

// handlerLogsCreate creates a practice log for a single day.
func (cfg *apiConfig) handlerLogsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Date       string `json:"date"`
		ColorDepth int32  `json:"color_depth"`
		Email      string `json:"email"`
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

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "no user with that email", err)
		return
	}

	practiceLog, err := cfg.db.CreateLog(r.Context(), database.CreateLogParams{
		Date:       parsedDate,
		ColorDepth: params.ColorDepth,
		UserID:     user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create practice log", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		PracticeLog: PracticeLog{
			Date:       params.Date,
			ColorDepth: practiceLog.ColorDepth,
		},
	})
}
