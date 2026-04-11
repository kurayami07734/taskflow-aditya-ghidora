package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/handlers"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/middleware"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/models"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
)

func CreateProjectRouter(db *sqlx.DB, cfg utils.Config) *chi.Mux {
	r := chi.NewRouter()

	projectStore := &models.ProjectStore{DB: db}
	taskStore := &models.TaskStore{DB: db}
	projectHandler := handlers.ProjectHandler{Store: projectStore}
	taskHandler := handlers.TaskHandler{TaskStore: taskStore, ProjectStore: projectStore}

	r.Use(middleware.AuthMiddleware(cfg.JwtSecret))
	r.Get("/", projectHandler.ListProjects)
	r.Post("/", projectHandler.CreateProject)
	r.Get("/{id}", projectHandler.GetProject)
	r.Patch("/{id}", projectHandler.UpdateProject)
	r.Delete("/{id}", projectHandler.DeleteProject)
	r.Post("/{id}/tasks", taskHandler.CreateTask)

	return r
}
