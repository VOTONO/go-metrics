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

	db, err := sql.Open("pgx", config.dbAddress)
	if err != nil {
		zapLogger.Errorw(
			"Fail open db",
			"address", config.dbAddress,
		)
		return
	}
	defer db.Close()

	shouldSyncWriteToFile := config.storeInterval == 0
	storer := initStorer(&zapLogger, config.restore, config.fileStoragePath, config.storeInterval)
	rout := router.Router(storer, shouldSyncWriteToFile, config.fileStoragePath, &zapLogger)

	zapLogger.Infow(
		"Starting server",
		"address", config.address,
		"dbAddress", config.dbAddress,
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

		metrics, err := storer.All()
		if err != nil {
			zapLogger.Errorw(
				"failed get metrics from storage before writing to file",
				"filePath", config.fileStoragePath,
				"err", err.Error())
			return
		}
		repo.RewriteFile(config.fileStoragePath, metrics, &zapLogger)
		defer cancel()
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

	if err := httpServer.ListenAndServe(); err != nil {
		zapLogger.Errorw(
			"Fail start server",
			"error", err,
		)
	}
}

func initStorer(logger *zap.SugaredLogger, restore bool, filePath string, storeInterval int) repo.MetricStorer {
	if storeInterval == 0 {
		return repo.NewFileMetricStorer(filePath, logger)
	}

	return repo.NewLocalMetricStorer(restore, filePath, logger)
}
