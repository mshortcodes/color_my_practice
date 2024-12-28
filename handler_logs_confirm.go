package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mshortcodes/color_my_practice/internal/auth"
	"github.com/mshortcodes/color_my_practice/internal/database"
)

// handlerLogsConfirm sets the confirmed field to true for all logs given in the array.
func (cfg *apiConfig) handlerLogsConfirm(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		LogIDs   []string `json:"log_ids"`
		Password string   `json:"password"`
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

	user, err := cfg.db.GetUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get user", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var params parameters

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect password", err)
		return
	}

	logIDs := make([]uuid.UUID, 0, len(params.LogIDs))

	for _, logIDStr := range params.LogIDs {
		logID, err := uuid.Parse(logIDStr)

		if err != nil {
			continue
		}

		logIDs = append(logIDs, logID)
	}

	confirmedDbLogs, err := cfg.db.ConfirmLogs(r.Context(), database.ConfirmLogsParams{
		UserID:  userID,
		Column2: logIDs,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error confirming logs", err)
		return
	}

	var confirmedLogs []PracticeLog

	for _, dbLog := range confirmedDbLogs {
		confirmedLog := PracticeLog{
			Id:         dbLog.ID,
			Date:       dbLog.Date.Format(time.DateOnly),
			ColorDepth: dbLog.ColorDepth,
			Confirmed:  dbLog.Confirmed,
			UserID:     dbLog.UserID,
		}

		confirmedLogs = append(confirmedLogs, confirmedLog)
	}

	respondWithJSON(w, http.StatusOK, confirmedLogs)
}
