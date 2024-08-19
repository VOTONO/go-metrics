package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
)

func ValueHandler(storer repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		name := chi.URLParam(req, "metricName")

		if name == "" {
			http.Error(res, "Invalide metric name", http.StatusNotFound)
		}

		ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
		defer cancel()

		metric, found, getErr := getMetricWithRetry(ctx, storer, name, 3, 1*time.Second)

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
		_, writeErr := res.Write([]byte(fmt.Sprintf("%v", value)))
		if writeErr != nil {
			http.Error(res, writeErr.Error(), http.StatusInternalServerError)
			return
		}
	}
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

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ConnectionException {
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
