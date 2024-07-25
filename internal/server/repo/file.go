package repo

import (
	"encoding/json"
	"os"

	"github.com/VOTONO/go-metrics/internal/models"
	"go.uber.org/zap"
)

func Read(file string, logger *zap.SugaredLogger) (map[string]models.Metric, error) {
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

func Write(file string, metrics map[string]models.Metric, logger *zap.SugaredLogger) error {
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

func logError(logger *zap.SugaredLogger, message, file string, err error) {
	logger.Errorw(message, "file", file, "err", err)
}
