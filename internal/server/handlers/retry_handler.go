package handlers

import (
	"context"
	"time"

	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
)

// fetchMetricsWithRetry handles the retry logic for fetching metrics.
func fetchMetricsWithRetry(ctx context.Context, storer repo.MetricStorer, retryCount int, initialPause time.Duration) (map[string]models.Metric, error) {
	var metrics map[string]models.Metric
	var err error
	retryPause := initialPause

	for i := 0; i <= retryCount; i++ {
		metrics, err = storer.All(ctx)
		if err == nil {
			return metrics, nil
		}

		shouldRetry := helpers.DecideShouldRetryAfterError(err)
		if shouldRetry {
			if i < retryCount {
				time.Sleep(retryPause)
				retryPause += 2
			}
		} else {
			break
		}
	}

	return nil, err
}

// storeMetricsWithRetry handles the retry logic for storing a slice of metrics.
func storeMetricsWithRetry(ctx context.Context, storer repo.MetricStorer, metrics []models.Metric, retryCount int, initialPause time.Duration) error {
	retryPause := initialPause
	var err error

	for i := 0; i <= retryCount; i++ {
		err = storer.StoreSlice(ctx, metrics)
		if err == nil {
			return nil // Success
		}

		shouldRetry := helpers.DecideShouldRetryAfterError(err)
		if shouldRetry {
			if i < retryCount {
				time.Sleep(retryPause)
				retryPause *= 2
			}
		} else {
			break
		}
	}

	return err
}

// storeMetricsWithRetry handles the retry logic for storing a slice of metrics.
func storeMetricWithRetry(ctx context.Context, storer repo.MetricStorer, metric models.Metric, retryCount int, initialPause time.Duration) (*models.Metric, error) {
	retryPause := initialPause
	var err error

	for i := 0; i <= retryCount; i++ {
		var storedMetric *models.Metric
		storedMetric, err = storer.StoreSingle(ctx, metric)
		if err == nil {
			return storedMetric, nil // Success
		}

		shouldRetry := helpers.DecideShouldRetryAfterError(err)
		if shouldRetry {
			if i < retryCount {
				time.Sleep(retryPause)
				retryPause *= 2
			}
		} else {
			break
		}
	}

	return nil, err
}

func getMetricWithRetry(ctx context.Context, storer repo.MetricStorer, id string, retryCount int, initialPause time.Duration) (models.Metric, bool, error) {
	retryPause := initialPause
	var err error

	for i := 0; i <= retryCount; i++ {
		var metric models.Metric
		var found bool
		metric, found, err = storer.Get(ctx, id)

		if err == nil {
			return metric, found, nil
		}

		shouldRetry := helpers.DecideShouldRetryAfterError(err)
		if shouldRetry {
			if i < retryCount {
				time.Sleep(retryPause)
				retryPause += 2
			}
		} else {
			break
		}
	}

	return models.Metric{}, false, err
}
