package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VOTONO/go-metrics/internal/models"
	fileworker "github.com/VOTONO/go-metrics/internal/server/fileWorker"
	"github.com/VOTONO/go-metrics/internal/server/router"
	"github.com/VOTONO/go-metrics/internal/server/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	zapLogger := *logger.Sugar()
	config := getConfig()

	db, err := sql.Open("pgx", config.dbAdress)
	if err != nil {
		zapLogger.Errorw(
			"Fail open db",
			"address", config.dbAdress,
		)
	}
	defer db.Close()

	var initialMetrics map[string]models.Metric

	if config.restore {
		restoredMetrics, err := fileworker.Read(config.fileStoragePath, &zapLogger)
		initialMetrics = restoredMetrics
		if err != nil {
			zapLogger.Errorw(
				"Fail read metrics from file",
				"path", config.fileStoragePath,
			)
		}
	}
	shouldSyncWriteToFile := config.storeInterval == 0
	stor := storage.New(initialMetrics, db, zapLogger)
	rout := router.Router(stor, shouldSyncWriteToFile, config.fileStoragePath, &zapLogger)

	zapLogger.Infow(
		"Starting server",
		"address", config.address,
		"dbAdress", config.dbAdress,
		"fileStoragePath", config.fileStoragePath,
		"storeInterval", config.storeInterval,
		"restore", config.restore,
	)

	go func() {
		er := http.ListenAndServe(config.address, rout)
		if er != nil {
			zapLogger.Errorw(
				"Fail start server",
				"error", er,
			)
		}
	}()

	if !shouldSyncWriteToFile && config.fileStoragePath != "" {
		storeTicker := time.NewTicker(time.Duration(config.storeInterval) * time.Second)
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		for {
			select {
			case <-stop:
				metrics := stor.All()
				for _, metric := range metrics {
					err := fileworker.Write(config.fileStoragePath, metric, &zapLogger)

					if err != nil {
						zapLogger.Errorw(
							"Fail write metrics to file",
							"path", config.fileStoragePath,
							"error", err,
						)
					}
				}
				storeTicker.Stop()
				return
			case <-storeTicker.C:
				metrics := stor.All()
				for _, metric := range metrics {
					err := fileworker.Write(config.fileStoragePath, metric, &zapLogger)
					if err != nil {
						zapLogger.Errorw(
							"Fail write metrics to file",
							"path", config.fileStoragePath,
							"error", err)
					}
				}
			}
		}
	}
}
