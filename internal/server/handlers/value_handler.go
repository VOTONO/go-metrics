package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"time"
)

func ValueHandler(storer repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		name := chi.URLParam(req, "metricName")

		if name == "" {
			http.Error(res, "Invalide metric name", http.StatusNotFound)
		}

		ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
		defer cancel()

		metric, found, getErr := getMetriWithRetry(ctx, storer, name, 3, 1*time.Second)

		if getErr != nil {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !found {
			http.Error(res, "Metric not found", http.StatusNotFound)
			return
		}

		value, err := helpers.ExtractValue(metric)

		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusInternalServerError)
		}

		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte(fmt.Sprintf("%v", value)))
	}
}

// storeMetricsWithRetry handles the retry logic for storing a slice of metrics.
func getMetriWithRetry(ctx context.Context, storer repo.MetricStorer, id string, retryCount int, initialPause time.Duration) (models.Metric, bool, error) {
	retryPause := initialPause
	var err error

	for i := 0; i <= retryCount; i++ {
		metric, found, getErr := storer.Get(ctx, id)
		if getErr == nil {
			return metric, found, nil // Success
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

	return models.Metric{}, false, err
}
