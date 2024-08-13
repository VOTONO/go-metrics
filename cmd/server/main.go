package main

import (
	"context"
	"database/sql"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"github.com/VOTONO/go-metrics/internal/server/router"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	zapLogger := *logger.Sugar()
	config := getConfig()

	storer, db, err := createStorer(&zapLogger, config)
	if err != nil {
		log.Fatalf("can't initialize metric storer: %v", err)
	}
	defer db.Close()
	rout := router.Router(storer, db, &zapLogger)

	zapLogger.Infow(
		"Starting server",
		"address", config.address,
		"DSN", config.DSN,
		"fileStoragePath", config.fileStoragePath,
		"storeInterval", config.storeInterval,
		"restore", config.restore,
	)

	httpServer := &http.Server{
		Addr:    config.address,
		Handler: rout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL)
	defer stop()

	go func() {
		<-ctx.Done()
		zapLogger.Infow("Shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		metrics, err := storer.All(ctx)
		if err != nil {
			zapLogger.Errorw(
				"failed get metrics from storage before writing to file",
				"filePath", config.fileStoragePath,
				"err", err.Error())
			// Ensure the server shutdown is attempted even if there's an error retrieving metrics
		} else {
			repo.RewriteFile(config.fileStoragePath, metrics, &zapLogger)
		}

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			zapLogger.Errorw(
				"Server shutdown failed",
				"error", err,
			)
		} else {
			zapLogger.Infow("Server gracefully stopped")
		}
	}()

	repo.StartWriting(ctx, storer, &zapLogger, config.storeInterval, config.fileStoragePath)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zapLogger.Errorw(
			"Fail start server",
			"error", err,
		)
	}
}

func createStorer(logger *zap.SugaredLogger, config Config) (repo.MetricStorer, *sql.DB, error) {
	if config.DSN != "" {
		db, err := sql.Open("pgx", config.DSN)
		if err != nil {
			logger.Errorw(
				"Fail open db",
				"address", config.DSN,
			)
			return nil, nil, err
		}
		storer, err := repo.NewPostgresMetricStorer(logger, db)

		if err != nil {
			logger.Errorw(
				"Fail create storer",
				"error", err.Error(),
			)
		}
		return storer, db, nil
	}

	if config.storeInterval == 0 {
		return repo.NewFileMetricStorer(config.fileStoragePath, logger), nil, nil
	}

	return repo.NewLocalMetricStorer(config.restore, config.fileStoragePath, logger), nil, nil
}
