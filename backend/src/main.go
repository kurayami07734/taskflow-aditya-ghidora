package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/routers"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
	_ "github.com/lib/pq"

	_ "github.com/kurayami07734/taskflow-aditya-ghidora/docs"
)

// @title Taskflow API
// @version 1.0
// @description API for managing projects and tasks
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cfg := utils.LoadConfig()

	logger := utils.InitLogger("logs/app.log")
	ctx := context.Background()

	logger.Info("Connecting to database", "db", cfg.Db.Name, "host", cfg.Db.Host, "port", cfg.Db.Port)
	db, err := sqlx.Connect("postgres", cfg.GetDbUrl())

	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	router := routers.CreateBaseRouter(db, cfg)

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("Failed to start server", "port", cfg.Port, "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("Server started", "port", cfg.Port)

	<-stop
	logger.Warn("Received SIGTERM, shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}
