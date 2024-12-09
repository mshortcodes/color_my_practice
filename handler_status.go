package main

import (
	"fmt"
	"net/http"
)

// handlerStatus serves a simple status summary of the page hits, users, and logs.
func (cfg *apiConfig) handlerStatus(w http.ResponseWriter, r *http.Request) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.hits++

	dbUsers, err := cfg.db.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving users", err)
		return
	}

	dbLogs, err := cfg.db.GetLogs(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving logs", err)
		return
	}

	statusMsg := fmt.Sprintf(`Page hits: %d
Users: %d
Logs: %d
`, cfg.hits, len(dbUsers), len(dbLogs))

	w.Write([]byte(statusMsg))
}
