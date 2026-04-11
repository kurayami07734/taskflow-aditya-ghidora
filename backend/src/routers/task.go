package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/handlers"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/middleware"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/models"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
)

func CreateTaskRouter(db *sqlx.DB, cfg utils.Config) *chi.Mux {
	r := chi.NewRouter()

	taskStore := &models.TaskStore{DB: db}
	projectStore := &models.ProjectStore{DB: db}
	taskHandler := handlers.TaskHandler{TaskStore: taskStore, ProjectStore: projectStore}

	r.Use(middleware.AuthMiddleware(cfg.JwtSecret))
	r.Patch("/{id}", taskHandler.UpdateTask)

	return r
}
