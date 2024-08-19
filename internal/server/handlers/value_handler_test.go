package handlers_test

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/mocks"
	"github.com/VOTONO/go-metrics/internal/server/handlers/utils"
	"github.com/VOTONO/go-metrics/internal/server/router"
)

func TestValueHandler(t *testing.T) {
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
		expectedCode int
	}{
		{
			name:         "Valid get",
			method:       "GET",
			url:          fmt.Sprintf("/value/%v/%v", utils.ValidGaugeMetric.MType, utils.ValidGaugeMetric.ID),
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

			if test.expectedCode == http.StatusOK {
				metricStorer.EXPECT().Get(gomock.Any(), utils.ValidGaugeMetric.ID).Return(utils.ValidGaugeMetric, true, nil)
			}

			req, err := http.NewRequest(test.method, server.URL+test.url, nil)
			assert.NoError(t, err, "Error create HTTP request")

			resp, err := server.Client().Do(req)
			if err != nil {
				assert.NoError(t, err)
			}
			defer resp.Body.Close()

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Response code didn't match expected")

			if test.expectedCode == http.StatusOK {
				bodyBytes, err := io.ReadAll(resp.Body)
				assert.NoError(t, err, "Error reading response body")
				bodyString := strings.TrimSpace(string(bodyBytes))
				expectedString := fmt.Sprintf("%v", *utils.ValidGaugeMetric.Value)
				assert.Equal(t, expectedString, bodyString, "Response body didn't match expected string")
			}
		})
	}
}
