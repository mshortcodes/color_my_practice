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
	count     int
	mu        *sync.Mutex
	db        *database.Queries
	platform  string
	jwtSecret string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("couldn't open database: %s", err)
	}

	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		count:     0,
		mu:        &sync.Mutex{},
		db:        dbQueries,
		platform:  platform,
		jwtSecret: jwtSecret,
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

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /count", apiCfg.handlerCountInc)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	mux.HandleFunc("POST /api/logs", apiCfg.handlerLogsCreate)
	mux.HandleFunc("GET /api/logs", apiCfg.handlerLogsGet)
	mux.HandleFunc("GET /api/logs/{logID}", apiCfg.handlerLogsGetByID)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerCountInc(w http.ResponseWriter, r *http.Request) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.count++

	countMsg := fmt.Sprintf("The count is %d", cfg.count)
	w.Write([]byte(countMsg))
}
