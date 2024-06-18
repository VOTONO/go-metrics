package server

import (
	"net/http"

	"github.com/VOTONO/go-metrics/internal/server/handlers"
	"github.com/VOTONO/go-metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type HTTPServer struct {
	storage storage.MetricStorage
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", handlers.AllValueHandler(s.storage))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateHandler(s.storage))
	router.Get("/value/{metricType}/{metricName}", handlers.ValueHandler(s.storage))

	// Serve HTTP using chi router
	router.ServeHTTP(w, r)
}

func New(storage storage.MetricStorage) *HTTPServer {
	return &HTTPServer{
		storage: storage,
	}
}
