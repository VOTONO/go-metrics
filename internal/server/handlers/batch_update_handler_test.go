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
	"github.com/VOTONO/go-metrics/internal/server/handlers/utils"
	"github.com/VOTONO/go-metrics/internal/server/router"
)

func TestBatchUpdateHandler(t *testing.T) {
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
		body         []models.Metric
		expectedCode int
	}{
		{
			name:         "Valid update",
			method:       "POST",
			url:          "/updates/",
			body:         []models.Metric{utils.ValidCounterMetric, utils.ValidGaugeMetric},
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if test.expectedCode == http.StatusOK {
				metricStorer.EXPECT().StoreSlice(gomock.Any(), test.body).Return(
					nil)
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
