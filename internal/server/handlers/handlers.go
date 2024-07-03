package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func UpdateHandlerJSON(s storage.Storage) http.HandlerFunc {
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

		valid := validateMetric(metric)
		if !valid {
			http.Error(res, "invalid metric", http.StatusBadRequest)
			return
		}

		stored, err := s.Store(metric)
		if err != nil {
			http.Error(res, "fail store metric", http.StatusInternalServerError)
		}

		out, err := json.Marshal(stored)
		if err != nil {
			log.Fatal(err)
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(out)
	}
}

func ValueHandlerJSON(s storage.Storage) http.HandlerFunc {
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

		value, found := s.Get(metric.ID)

		if !found {
			http.Error(res, "Metric not found", http.StatusNotFound)
			return
		}

		out, err := json.Marshal(value)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(out)
	}
}

func UpdateHandler(s storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		metricType := chi.URLParam(req, "metricType")
		name := chi.URLParam(req, "metricName")
		value := chi.URLParam(req, "metricValue")

		var metric models.Metric
		var err error

		switch metricType {
		case "gauge":
			var v float64
			v, err = strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(res, "Invalid metric value", http.StatusBadRequest)
				return
			}
			metric = models.Metric{
				ID:    name,
				MType: metricType,
				Value: &v,
			}
		case "counter":
			var v int64
			v, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				http.Error(res, "Invalid metric delta", http.StatusBadRequest)
				return
			}
			metric = models.Metric{
				ID:    name,
				MType: metricType,
				Delta: &v,
			}
		default:
			http.Error(res, "Invalid metric type", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusBadRequest)
			return
		}

		_, err = s.Store(metric)
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

		metric, found := s.Get(name)

		if !found {
			http.Error(res, "Metric not found", http.StatusNotFound)
			return
		}

		value, err := extractValue(metric)

		if err != nil {
			http.Error(res, "Invalid metric value", http.StatusInternalServerError)
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
		for key, metric := range metrics {
			value, err := extractValue(metric)

			if err != nil {
				http.Error(res, "Invalid metric value", http.StatusInternalServerError)
			}

			fmt.Fprintf(res, "<tr><td>%s</td><td>%v</td></tr>", key, value)
		}
		fmt.Fprintln(res, "</table></body></html>")
	}
}

func extractValue(m models.Metric) (string, error) {
	var value string

	switch m.MType {
	case "gauge":
		if m.Value != nil {
			value = strconv.FormatFloat(*m.Value, 'f', -1, 64)
		} else {
			return "", fmt.Errorf("metric value not found")
		}
	case "counter":
		if m.Delta != nil {
			value = strconv.FormatInt(*m.Delta, 10)
		} else {
			return "", fmt.Errorf("metric delta not found")
		}
	default:
		return "", fmt.Errorf("unknown metric type")
	}
	return value, nil
}

func validateMetric(m models.Metric) bool {
	switch m.MType {
	case "gauge":
		if m.Value == nil {
			return false
		}
		return true
	case "counter":
		if m.Delta == nil {
			return false
		}
		return true
	default:
		return false
	}
}
