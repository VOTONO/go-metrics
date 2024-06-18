package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VOTONO/go-metrics/internal/server/handlers"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type MockMemStorage struct{}

func (m MockMemStorage) GetAll() map[string]interface{} {
	return nil
}

func (m MockMemStorage) Get(name string) interface{} {
	return "100"
}

func (m MockMemStorage) Increment(name string, value interface{}) error {
	return nil
}

func (m MockMemStorage) Replace(name string, value interface{}) error {
	return nil
}

func TestUpdateHandler(t *testing.T) {
	// handler := http.HandlerFunc(handlers.UpdateHandler)
	server := httptest.NewServer(handlers.UpdateHandler(MockMemStorage{}))
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
			url:          "/update/gauge/testGauge",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {

			req := resty.New().R()
			req.Method = test.method
			req.URL = server.URL + test.url

			resp, err := req.Send()

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, test.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestValueHandler(t *testing.T) {
	// handler := http.HandlerFunc(handlers.UpdateHandler)
	server := httptest.NewServer(handlers.ValueHandler(MockMemStorage{}))
	defer server.Close()

	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
		expectedBody string
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

			req := resty.New().R()
			req.Method = test.method
			req.URL = server.URL + test.url

			resp, err := req.Send()

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, test.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}
