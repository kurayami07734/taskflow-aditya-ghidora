package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/handlers"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/models"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
)

func CreateAuthRouter(db *sqlx.DB, cfg utils.Config) *chi.Mux {
	r := chi.NewRouter()

	userStore := &models.UserStore{DB: db}
	registerHandler := handlers.RegisterHandler{Store: userStore, Config: cfg}
	loginHandler := handlers.LoginHandler{Store: userStore, Config: cfg}

	r.Use(middleware.Logger)
	r.Post("/register", registerHandler.RegisterHandler)
	r.Post("/login", loginHandler.LoginHandler)

	return r
}
