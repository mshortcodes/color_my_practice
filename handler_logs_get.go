package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

// handlerLogsGet returns all practice logs.
func (cfg *apiConfig) handlerLogsGet(w http.ResponseWriter, r *http.Request) {
	dbLogs, err := cfg.db.GetLogs(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving logs", err)
		return
	}

	var logs []PracticeLog

	for _, dbLog := range dbLogs {
		log := PracticeLog{
			Id:         dbLog.ID,
			Date:       dbLog.Date.Format(time.DateOnly),
			ColorDepth: dbLog.ColorDepth,
			Confirmed:  dbLog.Confirmed,
			UserID:     dbLog.UserID,
		}

		logs = append(logs, log)
	}

	respondWithJSON(w, http.StatusOK, logs)
}

// handlerLogsGetByID returns a single practice log by ID.
func (cfg *apiConfig) handlerLogsGetByID(w http.ResponseWriter, r *http.Request) {
	type response struct {
		PracticeLog
	}

	logIDStr := r.PathValue("logID")

	logID, err := uuid.Parse(logIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid log ID", err)
		return
	}

	dbLog, err := cfg.db.GetLog(r.Context(), logID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "log not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		PracticeLog: PracticeLog{
			Id:         dbLog.ID,
			Date:       dbLog.Date.Format(time.DateOnly),
			ColorDepth: dbLog.ColorDepth,
			Confirmed:  dbLog.Confirmed,
			UserID:     dbLog.UserID,
		},
	})
}
