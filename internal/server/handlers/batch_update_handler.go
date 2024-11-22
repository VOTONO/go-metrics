package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
)

// BatchUpdateHandler receives array of metrics in JSON format and updates storage with them.
func BatchUpdateHandler(storer repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metrics []models.Metric
		var buf bytes.Buffer

		_, readErr := buf.ReadFrom(req.Body)
		if readErr != nil {
			http.Error(res, readErr.Error(), http.StatusBadRequest)
			return
		}

		unmarshalErr := json.Unmarshal(buf.Bytes(), &metrics)
		if unmarshalErr != nil {
			http.Error(res, unmarshalErr.Error(), http.StatusBadRequest)
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
