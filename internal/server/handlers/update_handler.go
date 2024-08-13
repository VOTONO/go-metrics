package handlers

import (
	"context"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func UpdateHandler(s repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		metricType := chi.URLParam(req, "metricType")
		name := chi.URLParam(req, "metricName")
		value := chi.URLParam(req, "metricValue")

		metric, err := models.NewMetric(name, metricType, value)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
		defer cancel()

		_, err = s.StoreSingle(ctx, metric)
		if err != nil {
			http.Error(res, "fail store metric", http.StatusInternalServerError)
		}

		res.WriteHeader(http.StatusOK)
	}
}
