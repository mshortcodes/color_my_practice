package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mshortcodes/color_my_practice/internal/auth"
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

	port := os.Getenv("PORT")
	log.Printf("debugging port: %s\n", port)
	if port == "" {
		log.Fatal("PORT not set")
	}

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

	_, err := auth.GetTokenIssuer()
	if err != nil {
		log.Fatal("TOKEN_ISSUER not set")
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

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	mux.HandleFunc("GET /app/", apiCfg.handlerStatus)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)

	mux.HandleFunc("POST /api/logs", apiCfg.handlerLogsCreate)
	mux.HandleFunc("GET /api/logs", apiCfg.handlerLogsGet)
	mux.HandleFunc("GET /api/logs/{logID}", apiCfg.handlerLogsGetByID)
	mux.HandleFunc("PUT /api/logs/confirm", apiCfg.handlerLogsConfirm)
	mux.HandleFunc("DELETE /api/logs/{logID}", apiCfg.handlerLogsDelete)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	log.Fatal(srv.ListenAndServe())
}
