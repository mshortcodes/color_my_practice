package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mshortcodes/color_my_practice/internal/database"
)

type apiConfig struct {
	hits      int
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
		hits:      0,
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

	mux.HandleFunc("GET /app/", apiCfg.handlerStatus)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)

	mux.HandleFunc("POST /api/logs", apiCfg.handlerLogsCreate)
	mux.HandleFunc("GET /api/logs", apiCfg.handlerLogsGet)
	mux.HandleFunc("GET /api/logs/{logID}", apiCfg.handlerLogsGetByID)
	mux.HandleFunc("DELETE /api/logs/{logID}", apiCfg.handlerLogsDelete)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	log.Fatal(srv.ListenAndServe())
}
