package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/routers"
	_ "github.com/lib/pq"
)

func main() {
	cfg := LoadConfig()
	router := routers.CreateBaseRouter()

	log.Printf("Connecting to %s database at %s:%d...", cfg.DbName, cfg.DbHost, cfg.DbPort)
	db, err := sqlx.Connect("postgres", cfg.getDbUrl())

	if err != nil {
		log.Fatalf("Failed to connect to database :%v", err)
	}

	defer db.Close()

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start server on %d: %v", cfg.Port, err)
		}
	}()

	log.Printf("Started server on %d...", cfg.Port)

	<-stop
	log.Println("Recieved SIGTERM shutting down..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

}
