package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VOTONO/go-metrics/internal/helpers"

	"github.com/VOTONO/go-metrics/internal/compressor"
	"github.com/VOTONO/go-metrics/internal/models"

	"go.uber.org/zap"
)

type MetricSender interface {
	Send(metrics map[string]models.Metric) error
	BuildSingleRequest(metric models.Metric) (*http.Request, error)
	BuildBatchRequest(metrics map[string]models.Metric) (*http.Request, error)
}

type MetricSenderImpl struct {
	Client  *http.Client
	Address string
	Logger  *zap.SugaredLogger
}

func New(client *http.Client, address string, logger *zap.SugaredLogger) *MetricSenderImpl {
	return &MetricSenderImpl{
		Client:  client,
		Address: address,
		Logger:  logger,
	}
}

// Send metrics to the server.
func (sender *MetricSenderImpl) Send(metrics map[string]models.Metric) error {
	err := sender.sendBatch(metrics)

	if err != nil {
		retryCount := 3
		retryPause := 1 * time.Second

		for i := 0; i < retryCount; i++ {
			time.Sleep(retryPause)
			err = sender.sendBatch(metrics)

			if err == nil {
				return nil
			}

			retryPause += 2
		}
	}

	return nil
}

func (sender *MetricSenderImpl) sendBatch(metrics map[string]models.Metric) error {
	req, err := sender.BuildBatchRequest(helpers.ConvertMapToSlice(metrics))
	if err != nil {
		return fmt.Errorf("failed to build batch request: %w", err)
	}

	resp, err := sender.Client.Do(req)
	if err != nil {
		sender.Logger.Errorw("Failed to send batch request", "error", err)
		return fmt.Errorf("error sending batch request for metrics: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			sender.Logger.Errorw("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		sender.Logger.Errorw("Received bad response", "status code", resp.StatusCode)
		return fmt.Errorf("batch request failed with status code %d", resp.StatusCode)
	}

	return nil
}

// BuildBatchRequest creates a compressed HTTP request for a batch of metrics.
func (sender *MetricSenderImpl) BuildBatchRequest(metrics []models.Metric) (*http.Request, error) {
	url := fmt.Sprintf("http://%s/updates/", sender.Address)

	body, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}

	compressedBody, err := compressor.GzipCompress(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(compressedBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	return req, nil
}

// NotFoundError is a custom error type for handling 404 status codes.
type NotFoundError struct {
	Message string
}
