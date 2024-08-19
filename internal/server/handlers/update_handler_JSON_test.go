package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/mocks"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/router"
)

func TestUpdateHandlerJSON(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metricStorer := mocks.NewMockMetricStorer(ctrl)

	zapLogger := *logger.Sugar()
	server := httptest.NewServer(router.Router(metricStorer, &sql.DB{}, &zapLogger))
	defer server.Close()

	tests := []struct {
		name         string
		method       string
		url          string
		body         models.Metric
		expectedCode int
	}{
		{
			name:   "Valid update",
			method: "POST",
			url:    "/update/",
			body: models.Metric{
				ID:    "testMetric",
				MType: "gauge",
				Value: func() *float64 { f := float64(456.78); return &f }(),
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "Invalid counter metric update",
			method: "POST",
			url:    "/update/",
			body: models.Metric{
				ID:    "testMetric",
				MType: "counter",
				Value: func() *float64 { f := float64(456.78); return &f }(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "Invalid gauge metric update",
			method: "POST",
			url:    "/update/",
			body: models.Metric{
				ID:    "testMetric",
				MType: "gauge",
				Delta: func() *int64 { f := int64(456); return &f }(),
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if test.expectedCode == http.StatusOK {
				metricStorer.EXPECT().StoreSingle(gomock.Any(), test.body).Return(&test.body, nil)
			}

			jsonBody, err := json.Marshal(test.body)
			assert.NoError(t, err)

			req, err := http.NewRequest(test.method, server.URL+test.url, bytes.NewBuffer(jsonBody))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, test.expectedCode, resp.StatusCode, "Response code didn't match expected")
		})
	}
}
