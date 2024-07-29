package handlers

import (
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func UpdateHandler(s repo.MetricStorer, zap *zap.SugaredLogger, shouldSyncWriteToFile bool, filePath string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		metricType := chi.URLParam(req, "metricType")
		name := chi.URLParam(req, "metricName")
		value := chi.URLParam(req, "metricValue")

		metric, err := models.NewMetric(name, metricType, value)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = s.Store(metric)
		if err != nil {
			http.Error(res, "fail store metric", http.StatusInternalServerError)
		}

		res.WriteHeader(http.StatusOK)
	}
}
