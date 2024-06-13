package handlers

import (
	"fmt"
	"net/http"

	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
)

func UpdateHandler(s storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		metricType := chi.URLParam(req, "metricType")
		name := chi.URLParam(req, "metricName")
		value := chi.URLParam(req, "metricValue")

		fmt.Printf("Received metricType: %s, metricName: %s, metricValue: %s\n", metricType, name, value)

		metric, err := models.New(name, metricType, value)

		if err != nil {
			if err.Error() == "bad name" {
				http.Error(res, "Invalide metric name", http.StatusNotFound)
			} else {
				http.Error(res, "Invalide metric", http.StatusBadRequest)
			}
		}

		err = s.Store(metric)
		if err != nil {
			http.Error(res, "fail store metric", http.StatusInternalServerError)
		}
		res.WriteHeader(http.StatusOK)
	}
}

func ValueHandler(s storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		name := chi.URLParam(req, "metricName")

		if name == "" {
			http.Error(res, "Invalide metric name", http.StatusNotFound)
		}

		value, found := s.Get(name)

		if !found {
			http.Error(res, "Metric not found", http.StatusNotFound)
			return
		}
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte(fmt.Sprintf("%v", value)))
	}
}

func AllValueHandler(memStorage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if req.URL.Path != "/" {
			http.Error(res, "Bad url", http.StatusNotFound)
			return
		}

		metrics := memStorage.All()

		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		fmt.Fprintln(res, "<html><body><h1>Metrics</h1><table border='1'><tr><th>Metric</th><th>Value</th></tr>")
		for key, value := range metrics {
			fmt.Fprintf(res, "<tr><td>%s</td><td>%v</td></tr>", key, value)
		}
		fmt.Fprintln(res, "</table></body></html>")
	}
}
