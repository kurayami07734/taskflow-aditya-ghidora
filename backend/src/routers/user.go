package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/handlers"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/models"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
)

func CreateUserRouter(db *sqlx.DB, cfg utils.Config) *chi.Mux {
	r := chi.NewRouter()

	userStore := &models.UserStore{DB: db}
	userHandler := handlers.UserHandler{Store: userStore}

	r.Get("/search", userHandler.SearchUsers)
	r.Get("/{id}", userHandler.GetUser)

	return r
}
