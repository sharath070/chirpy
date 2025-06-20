package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sharath070/Chirpy/internal/database"
)

type apiConfig struct {
	fileSeverHits atomic.Int32
	dbQueries     *database.Queries
	jwtSecret     string
}

const (
	port         = "8080"
	fileRootPath = "."
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("DB connection failed")
		return
	}

	cfg := apiConfig{
		fileSeverHits: atomic.Int32{},
		dbQueries:     database.New(db),
		jwtSecret:     os.Getenv("SECRET"),
	}

	mux := http.NewServeMux()
	mux.Handle("GET /app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fileRootPath)))))
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /reset", cfg.handleReset)

	// USERS
	mux.HandleFunc("POST /api/users", cfg.handleCreateUser)
	mux.HandleFunc("POST /api/login", cfg.handleLoginUser)
	mux.HandleFunc("POST /api/refresh", cfg.handleRefresh)

	//  CHIRPS
	mux.HandleFunc("POST /api/chirps", cfg.handleCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handleGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handleGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handleDeleteChirp)

	srv := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving file from %s on port: %s\n", port, fileRootPath)
	log.Fatal(srv.ListenAndServe())
}
