package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/router"

	"github.com/stretchr/testify/assert"
)

type MockStorage struct{}

func (m MockStorage) Store(metric models.Metric) error {
	return nil
}

func (m MockStorage) Get(name string) (models.Metric, bool) {
	return models.Metric{Name: "foo", Type: "gauge", Value: "100"}, true
}

func (m MockStorage) All() map[string]models.Metric {
	return nil
}

func TestUpdateHandler(t *testing.T) {
	server := httptest.NewServer(router.Router(MockStorage{}))
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
		t.Run(test.method, func(t *testing.T) {
			req, err := http.NewRequest(test.method, server.URL+test.url, nil)
			assert.NoError(t, err)

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Response code didn't match expected")
		})
	}
}

func TestValueHandler(t *testing.T) {
	server := httptest.NewServer(router.Router(MockStorage{}))
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
		t.Run(test.method, func(t *testing.T) {
			req, err := http.NewRequest(test.method, server.URL+test.url, nil)
			assert.NoError(t, err)

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Response code didn't match expected")
		})
	}
}
