package main

import (
	"fmt"
	"net/http"
)

// handlerReset resets the database and sets page hits to 0.
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("reset is only allowed in dev environment"))
		return
	}

	err := cfg.db.Reset(r.Context())
	if err != nil {
		msg := fmt.Sprintf("error resetting database: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.hits = 0

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("database reset and page hits set to 0"))
}
