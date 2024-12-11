package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mshortcodes/color_my_practice/internal/auth"
)

func (cfg *apiConfig) handlerLogsDelete(w http.ResponseWriter, r *http.Request) {
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

	logIDStr := r.PathValue("logID")

	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid log ID", err)
		return
	}

	log, err := cfg.db.GetLog(r.Context(), logID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "log not found", err)
		return
	}

	if log.UserID != userID {
		respondWithError(w, http.StatusForbidden, "you can't delete this log", err)
		return
	}

	err = cfg.db.DeleteLog(r.Context(), log.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete log", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
