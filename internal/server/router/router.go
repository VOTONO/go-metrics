package router

import (
	"database/sql"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/auth"
	"github.com/VOTONO/go-metrics/internal/compressor"
	"github.com/VOTONO/go-metrics/internal/logger"
	"github.com/VOTONO/go-metrics/internal/server/handlers"
	"github.com/VOTONO/go-metrics/internal/server/repo"
)

func Router(s repo.MetricStorer, db *sql.DB, zap *zap.SugaredLogger, secretKey string) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(auth.HashChecker(secretKey))
	router.Use(compressor.Compressor)
	router.Use(compressor.Decompressor)
	router.Use(auth.HashSigner(secretKey))

	// Your existing application routes
	router.Get("/", logger.WithLogger(handlers.AllValueHandler(s, zap), zap))
	router.Get("/ping", logger.WithLogger(handlers.Ping(db), zap))
	router.Post("/update/", logger.WithLogger(handlers.UpdateHandlerJSON(s), zap))
	router.Post("/updates/", logger.WithLogger(handlers.BatchUpdateHandler(s), zap))
	router.Post("/value/", logger.WithLogger(handlers.ValueHandlerJSON(s), zap))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", logger.WithLogger(handlers.UpdateHandler(s), zap))
	router.Get("/value/{metricType}/{metricName}", handlers.ValueHandler(s))

	// Mount pprof routes directly
	router.Mount("/debug/pprof/", http.DefaultServeMux)

	return router
}
