package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
)

func AllValueHandler(storer repo.MetricStorer, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if req.URL.Path != "/" {
			http.Error(res, "Bad url", http.StatusNotFound)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
		defer cancel()

		metrics, err := fetchMetricsWithRetry(ctx, storer, 3, 1*time.Second)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		htmlContent, err := helpers.MetricsToHTML(metrics, logger)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		_, printErr := fmt.Fprintln(res, htmlContent)
		if printErr != nil {
			http.Error(res, printErr.Error(), http.StatusInternalServerError)
			return
		}
	}
}

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

	return nil, err
}
