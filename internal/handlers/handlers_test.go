package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VOTONO/go-metrics/internal/handlers"
)

type MockMemStorage struct{}

func (m *MockMemStorage) Increment(name string, value interface{}) error { return nil }
func (m *MockMemStorage) Replace(name string, value interface{}) error   { return nil }
func (m *MockMemStorage) Get(name string) interface{}                    { return nil }

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid gauge metric",
			method:       "POST",
			url:          "/update/gauge/testGauge/123.45",
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "Valid counter metric",
			method:       "POST",
			url:          "/update/counter/testCounter/123",
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "Invalid metric type",
			method:       "POST",
			url:          "/update/invalid/testInvalid/123",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid metric type\n",
		},
		{
			name:         "Invalid request method",
			method:       "GET",
			url:          "/update/gauge/testGauge/123.45",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "Only POST requests are allowed!\n",
		},
		{
			name:         "Invalid gauge value",
			method:       "POST",
			url:          "/update/gauge/testGauge/abc",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid value for gauge, must be float64\n",
		},
		{
			name:         "Invalid counter value",
			method:       "POST",
			url:          "/update/counter/testCounter/123.45",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid value for counter, must be int64\n",
		},
		{
			name:         "Invalid URL format",
			method:       "POST",
			url:          "/update/gauge/testGauge",
			expectedCode: http.StatusNotFound,
			expectedBody: "Invalid request format\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := handlers.UpdateHandler(&MockMemStorage{})

			req, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedCode {
				t.Errorf("Wrong status code: got %v want %v",
					status, test.expectedCode)
			}

			if body := rr.Body.String(); body != test.expectedBody {
				t.Errorf("Unexpected body: got %v want %v",
					body, test.expectedBody)
			}
		})
	}
}
