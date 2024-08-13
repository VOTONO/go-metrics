package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/repo"
	"log"
	"net/http"
	"time"
)

func UpdateHandlerJSON(s repo.MetricStorer) http.HandlerFunc {
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

		stored, err := s.StoreSingle(ctx, metric)
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
