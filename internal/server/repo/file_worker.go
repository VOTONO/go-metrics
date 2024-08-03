package repo

import (
	"context"
	"encoding/json"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"go.uber.org/zap"
	"os"
	"time"
)

func ReadFile(file string, logger *zap.SugaredLogger) (map[string]models.Metric, error) {
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Infow("file not found", "file", file)
			return make(map[string]models.Metric), nil
		}
		logError(logger, "failed to open file", file, err)
		return nil, err
	}
	defer f.Close()

	var metrics map[string]models.Metric
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&metrics); err != nil {
		logError(logger, "failed to decode file", file, err)
		return nil, err
	}

	return metrics, nil
}

// RewriteFile all metrics in file
func RewriteFile(file string, metrics map[string]models.Metric, logger *zap.SugaredLogger) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logError(logger, "failed to open file", file, err)
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	if err := encoder.Encode(metrics); err != nil {
		logError(logger, "failed to encode metrics", file, err)
		return err
	}

	logger.Infow("successfully wrote metrics", "metrics", metrics, "file", file)
	return nil
}

func AddToFile(file string, metric models.Metric, logger *zap.SugaredLogger) (models.Metric, error) {
	metrics, err := ReadFile(file, logger)

	if err != nil {
		return models.Metric{}, err
	}

	updated, err := helpers.UpdateMetricInMap(&metrics, metric, logger)

	if err != nil {
		return models.Metric{}, err
	}

	return updated, nil
}

func StartWriting(ctx context.Context, storer MetricStorer, logger *zap.SugaredLogger, storeInterval int, filePath string) {
	if storeInterval <= 0 || filePath == "" {
		logger.Infow("skip periodical writing to file", "file", filePath, "interval", storeInterval)
		return
	}

	storeTicker := time.NewTicker(time.Duration(storeInterval) * time.Second)

	go func() {
		defer storeTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				metrics, err := storer.All()
				if err != nil {
					logError(logger, "failed get metrics from storage before writing to file", filePath, err)
					return
				}
				RewriteFile(filePath, metrics, logger)
				return
			case <-storeTicker.C:
				metrics, err := storer.All()
				if err != nil {
					logError(logger, "failed get metrics from storage before writing to file", filePath, err)
					return
				}
				RewriteFile(filePath, metrics, logger)
			}
		}
	}()
}

func logError(logger *zap.SugaredLogger, message, file string, err error) {
	logger.Errorw(message, "file", file, "err", err.Error())
}
