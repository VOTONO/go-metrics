package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
)

func BatchUpdateHandler(storer repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metrics []models.Metric
		var buf bytes.Buffer

		_, readErr := buf.ReadFrom(req.Body)
		if readErr != nil {
			http.Error(res, readErr.Error(), http.StatusBadRequest)
			return
		}

		umarshalErr := json.Unmarshal(buf.Bytes(), &metrics)
		if umarshalErr != nil {
			http.Error(res, umarshalErr.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), 1000*time.Second)
		defer cancel()

		storeErr := storeMetricsWithRetry(ctx, storer, metrics, 3, 1*time.Second)
		if storeErr != nil {
			http.Error(res, "fail store metric", http.StatusInternalServerError)
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
	}
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

	return err
}
