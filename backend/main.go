package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreateBaseRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}

func main() {
	cfg := LoadConfig()
	r := CreateBaseRouter()

	log.Printf("Connecting to %s database at %s:%d...", cfg.DbName, cfg.DbHost, cfg.DbPort)
	db, err := sqlx.Connect("postgres", cfg.getDbUrl())

	if err != nil {
		log.Fatalf("Failed to connect to database :%v", err)
	}

	defer db.Close()

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on %d...", cfg.Port)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server on %d: %v", cfg.Port, err)
	}
}
