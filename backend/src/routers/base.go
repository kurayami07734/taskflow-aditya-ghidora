package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	slogm "github.com/kurayami07734/taskflow-aditya-ghidora/src/middleware"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/kurayami07734/taskflow-aditya-ghidora/docs"
)

func CreateBaseRouter(db *sqlx.DB, cfg utils.Config) *chi.Mux {
	logger := utils.InitLogger("logs/app.log")

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:80", "http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Recoverer)
	r.Use(slogm.RecoverWithLog(logger))
	r.Use(middleware.RequestLogger(slogm.NewSlogLogger(logger)))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusFound)
	})
	r.Get("/docs/*", httpSwagger.Handler())

	authR := CreateAuthRouter(db, cfg)
	r.Mount("/auth", authR)

	projectR := CreateProjectRouter(db, cfg)
	r.Mount("/projects", projectR)

	taskR := CreateTaskRouter(db, cfg)
	r.Mount("/tasks", taskR)

	return r
}
