package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type apiConfig struct {
	count int
	mu    sync.Mutex
}

func main() {
	mux := http.NewServeMux()
	port := "8080"

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	var apiCfg apiConfig

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	mux.HandleFunc("GET /count", apiCfg.handlerCountInc)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerCountInc(w http.ResponseWriter, r *http.Request) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.count++

	countMsg := fmt.Sprintf("The count is %d", cfg.count)
	w.Write([]byte(countMsg))
}
