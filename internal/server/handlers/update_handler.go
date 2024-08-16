package handlers

import (
	"context"
	"errors"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"time"
)

func UpdateHandler(storer repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		metricType := chi.URLParam(req, "metricType")
		name := chi.URLParam(req, "metricName")
		value := chi.URLParam(req, "metricValue")

		newMetric, err := models.NewMetric(name, metricType, value)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
		defer cancel()

		_, err = storeMetricWithRetry(ctx, storer, newMetric, 3, 1*time.Second)
		if err != nil {
			http.Error(res, "fail store metric", http.StatusInternalServerError)
		}

		res.WriteHeader(http.StatusOK)
	}
}

// storeMetricsWithRetry handles the retry logic for storing a slice of metrics.
func storeMetricWithRetry(ctx context.Context, storer repo.MetricStorer, metric models.Metric, retryCount int, initialPause time.Duration) (*models.Metric, error) {
	retryPause := initialPause
	var err error

	for i := 0; i <= retryCount; i++ {
		storedMetric, err := storer.StoreSingle(ctx, metric)
		if err == nil {
			return storedMetric, nil // Success
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ConnectionException {
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
