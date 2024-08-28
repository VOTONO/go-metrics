package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
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
