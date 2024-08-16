package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"net/http"
	"time"
)

func ValueHandlerJSON(storer repo.MetricStorer) http.HandlerFunc {
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

		ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
		defer cancel()

		storedMetric, found, getErr := getMetriWithRetry(ctx, storer, metric.ID, 3, 1*time.Second)

		if getErr != nil {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !found {
			http.Error(res, "Metric not found", http.StatusNotFound)
			return
		}

		out, err := json.Marshal(storedMetric)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(out)
	}
}
