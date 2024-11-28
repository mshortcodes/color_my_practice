package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mshortcodes/color_my_practice/internal/database"
)

type apiConfig struct {
	count int
	mu    *sync.Mutex
	db    *database.Queries
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("couldn't open database: %s", err)
	}

	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		count: 0,
		mu:    &sync.Mutex{},
		db:    dbQueries,
	}

	mux := http.NewServeMux()
	port := "8080"

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	mux.HandleFunc("GET /count", apiCfg.handlerCountInc)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerCountInc(w http.ResponseWriter, r *http.Request) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.count++

	countMsg := fmt.Sprintf("The count is %d", cfg.count)
	w.Write([]byte(countMsg))
}
