package handlers_test

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/constants"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/mocks"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/router"
)

func TestAllValueHandler(t *testing.T) {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metricStorer := mocks.NewMockMetricStorer(ctrl)

	zapLogger := *logger.Sugar()
	server := httptest.NewServer(router.Router(metricStorer, &sql.DB{}, &zapLogger, ""))
	defer server.Close()

	tests := []struct {
		name         string
		method       string
		url          string
		metrics      map[string]models.Metric
		expectedCode int
	}{
		{
			name:   "Valid GET All Value",
			method: http.MethodGet,
			url:    "/",
			metrics: map[string]models.Metric{
				"metric1": {ID: "metric1", MType: constants.Gauge, Value: func() *float64 { f := 123.45; return &f }()},
				"metric2": {ID: "metric2", MType: constants.Counter, Delta: func() *int64 { i := int64(678); return &i }()},
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectedCode == http.StatusOK {
				metricStorer.EXPECT().All(gomock.Any()).Return(test.metrics, nil)
			}

			req, err := http.NewRequest(test.method, server.URL+test.url, nil)
			if err != nil {
				assert.NoError(t, err)
			}

			resp, err := server.Client().Do(req)
			if err != nil {
				assert.NoError(t, err)
			}
			defer resp.Body.Close()

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Response code didn't match expected")

			if test.expectedCode == http.StatusOK {
				bodyBytes, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				bodyString := strings.TrimSpace(string(bodyBytes))

				expectedHTML, er := helpers.MetricsToHTML(test.metrics, &zapLogger)

				assert.NoError(t, er)

				assert.Equal(t, expectedHTML, bodyString, "Response body didn't match expected HTML")
			}
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Response code didn't match expected")
		})
	}
}
