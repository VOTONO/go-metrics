package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/server/repo"
)

// ValueHandler retrieve metric name from URLParams and return value.
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
		_, writeErr := fmt.Fprintf(res, "%v", value)
		if writeErr != nil {
			http.Error(res, writeErr.Error(), http.StatusInternalServerError)
			return
		}
	}
}
