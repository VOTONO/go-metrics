package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/router"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

type MockStorage struct{}

func (m MockStorage) Store(metric models.Metric) (*models.Metric, error) {
	return nil, nil
}

func (m MockStorage) Get(name string) (models.Metric, bool) {
	value := float64(100)
	return models.Metric{ID: "testMetric", MType: "gauge", Value: &value}, true
}

func (m MockStorage) All() map[string]models.Metric {
	return nil
}

func TestUpdateHandlerJSON(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugar := *logger.Sugar()
	server := httptest.NewServer(router.Router(MockStorage{}, sugar))
	defer server.Close()

	tests := []struct {
		name         string
		method       string
		url          string
		body         interface{}
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

func TestValueHandlerJSON(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugar := *logger.Sugar()
	server := httptest.NewServer(router.Router(MockStorage{}, sugar))
	defer server.Close()

	tests := []struct {
		name         string
		method       string
		url          string
		body         interface{}
		expectedCode int
	}{
		{
			name:   "Valid get",
			method: "GET",
			url:    "/value/",
			body: models.Metric{
				ID:    "testMetric",
				MType: "gauge",
				Value: func() *float64 { f := float64(100); return &f }(),
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

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

func TestUpdateHandler(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugar := *logger.Sugar()
	server := httptest.NewServer(router.Router(MockStorage{}, sugar))
	defer server.Close()

	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
	}{
		{
			name:         "Valid gauge metric",
			method:       "POST",
			url:          "/update/gauge/testGauge/123.45",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Valid counter metric",
			method:       "POST",
			url:          "/update/counter/testCounter/123",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid metric type",
			method:       "POST",
			url:          "/update/invalid/testInvalid/123",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid gauge value",
			method:       "POST",
			url:          "/update/gauge/testGauge/abc",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid counter value",
			method:       "POST",
			url:          "/update/counter/testCounter/123.45",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid URL format",
			method:       "POST",
			url:          "/update/gauge/123",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, server.URL+test.url, nil)
			assert.NoError(t, err)

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Response code didn't match expected")
		})
	}
}

func TestValueHandler(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugar := *logger.Sugar()
	server := httptest.NewServer(router.Router(MockStorage{}, sugar))
	defer server.Close()

	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
	}{
		{
			name:         "Valid get",
			method:       "GET",
			url:          "/value/gauge/testGauge",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid method",
			method:       "POST",
			url:          "/value/gauge/testGauge",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "Invalid url",
			method:       "GET",
			url:          "/value/gauge",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, server.URL+test.url, nil)
			assert.NoError(t, err)

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Response code didn't match expected")
		})
	}
}
