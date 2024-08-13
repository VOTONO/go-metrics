package handlers_test

import (
	"database/sql"
	"fmt"
	"github.com/VOTONO/go-metrics/internal/mocks"
	"github.com/VOTONO/go-metrics/internal/models"
	"github.com/VOTONO/go-metrics/internal/server/handlers/utils"
	"github.com/VOTONO/go-metrics/internal/server/router"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
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
		metric       models.Metric
		expectedCode int
	}{
		{
			name:         "Valid gauge metric",
			method:       "POST",
			url:          fmt.Sprintf("/update/%v/%v/%v", utils.ValidGaugeMetric.MType, utils.ValidGaugeMetric.ID, *utils.ValidGaugeMetric.Value),
			metric:       utils.ValidGaugeMetric,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Valid counter metric",
			method:       "POST",
			url:          fmt.Sprintf("/update/%v/%v/%v", utils.ValidCounterMetric.MType, utils.ValidCounterMetric.ID, *utils.ValidCounterMetric.Delta),
			metric:       utils.ValidCounterMetric,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid metric type",
			method:       "POST",
			url:          "/update/invalid/testInvalid/123",
			metric:       models.Metric{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid gauge value",
			method:       "POST",
			url:          "/update/gauge/testGauge/abc",
			metric:       models.Metric{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid counter value",
			method:       "POST",
			url:          "/update/counter/testCounter/123.45",
			metric:       models.Metric{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid URL format",
			method:       "POST",
			url:          "/update/gauge/123",
			metric:       models.Metric{},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if test.expectedCode == http.StatusOK {
				metricStorer.EXPECT().StoreSingle(gomock.Any(), test.metric).Return(&test.metric, nil)
			}

			req, err := http.NewRequest(test.method, server.URL+test.url, nil)
			assert.NoError(t, err)

			resp, err := server.Client().Do(req)
			if err != nil {
				assert.NoError(t, err)
			}
			defer resp.Body.Close()

			assert.NoError(t, err, "Error making HTTP request")
			assert.Equal(t, test.expectedCode, resp.StatusCode, "Response code didn't match expected")
		})
	}
}
