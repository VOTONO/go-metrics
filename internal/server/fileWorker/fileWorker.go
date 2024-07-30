package fileworker

import (
	"bufio"
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

	metrics := make(map[string]models.Metric)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var metric models.Metric
		if err := json.Unmarshal(scanner.Bytes(), &metric); err != nil {
			logError(logger, "failed to parse file", file, err)
			return nil, err
		}
		metrics[metric.ID] = metric
	}

	if err := scanner.Err(); err != nil {
		logError(logger, "failed to scan file", file, err)
		return nil, err
	}

	return metrics, nil
}

func Write(file string, metric models.Metric, logger *zap.SugaredLogger) error {
	metrics, err := Read(file, logger)
	if err != nil {
		logError(logger, "failed to read file", file, err)
		return err
	}

	logger.Infow("metrics before writing", "metrics", metrics)

	metrics[metric.ID] = metric

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logError(logger, "failed to open file", file, err)
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, m := range metrics {
		data, err := json.Marshal(m)
		if err != nil {
			logError(logger, "failed to marshal metric", file, err)
			return err
		}
		if _, err := writer.Write(data); err != nil {
			logError(logger, "failed to write metric", file, err)
			return err
		}
		if _, err := writer.WriteString("\n"); err != nil {
			logError(logger, "failed to write metric newline", file, err)
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		logError(logger, "failed to flush writer", file, err)
		return err
	}

	logger.Infow("successfully write", "metrics", metrics, "file", file)

	return nil
}

func logError(logger *zap.SugaredLogger, message, file string, err error) {
	logger.Errorw(message, "file", file, "err", err)
}
