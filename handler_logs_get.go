package main

import (
	"net/http"
	"time"
)

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
