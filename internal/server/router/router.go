package router

import (
	"github.com/VOTONO/go-metrics/internal/compressor"
	"github.com/VOTONO/go-metrics/internal/logger"
	"github.com/VOTONO/go-metrics/internal/server/handlers"
	"github.com/VOTONO/go-metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Router(s storage.Storage, zap zap.SugaredLogger) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(compressor.Compressor)
	router.Use(compressor.Decompressor)

	router.Get("/", logger.WithLogger(handlers.AllValueHandler(s), zap))
	router.Post("/update/", logger.WithLogger(handlers.UpdateHandlerJSON(s), zap))
	router.Post("/value/", logger.WithLogger(handlers.ValueHandlerJSON(s), zap))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", logger.WithLogger(handlers.UpdateHandler(s), zap))
	router.Get("/value/{metricType}/{metricName}", handlers.ValueHandler(s))
	return router
}
