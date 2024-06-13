package router

import (
	"github.com/VOTONO/go-metrics/internal/server/handlers"
	"github.com/VOTONO/go-metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(s storage.Storage) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", handlers.AllValueHandler(s))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateHandler(s))
	router.Get("/value/{metricType}/{metricName}", handlers.ValueHandler(s))
	return router
}
