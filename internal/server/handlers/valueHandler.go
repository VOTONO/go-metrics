package handlers

import (
	"fmt"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func ValueHandler(s repo.MetricStorer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		name := chi.URLParam(req, "metricName")

		if name == "" {
			http.Error(res, "Invalide metric name", http.StatusNotFound)
		}

		metric, found := s.Get(name)

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
