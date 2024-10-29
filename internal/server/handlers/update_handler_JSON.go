package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
)

// UpdateHandlerJSON retrieve metric from body in JSON format and updates it in storage.
func UpdateHandlerJSON(storer repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var metric models.Metric
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		valid := helpers.ValidateMetric(metric)
		if !valid {
			http.Error(res, "invalid metric", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
		defer cancel()

		stored, storeErr := storeMetricWithRetry(ctx, storer, metric, 3, 1*time.Second)
		if storeErr != nil {
			http.Error(res, storeErr.Error(), http.StatusInternalServerError)
			return
		}

		out, marshalErr := json.Marshal(stored)
		if marshalErr != nil {
			http.Error(res, marshalErr.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)

		_, writeErr := res.Write(out)
		if writeErr != nil {
			http.Error(res, writeErr.Error(), http.StatusInternalServerError)
			return
		}
	}
}
