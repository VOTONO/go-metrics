package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/server/repo"
)

// AllValueHandler returns all metrics in HTML format.
func AllValueHandler(storer repo.MetricStorer, logger *zap.SugaredLogger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if req.URL.Path != "/" {
			http.Error(res, "Bad url", http.StatusNotFound)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
		defer cancel()

		metrics, err := fetchMetricsWithRetry(ctx, storer, 3, 1*time.Second)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		htmlContent, err := helpers.MetricsToHTML(metrics, logger)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		_, printErr := fmt.Fprintln(res, htmlContent)
		if printErr != nil {
			http.Error(res, printErr.Error(), http.StatusInternalServerError)
			return
		}
	}
}
