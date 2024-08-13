package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"net/http"

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
	Client                  *http.Client
	Address                 string
	Logger                  *zap.SugaredLogger
	shouldSendBatchRequests bool
}

func New(client *http.Client, address string, logger *zap.SugaredLogger) *MetricSenderImpl {
	return &MetricSenderImpl{
		Client:                  client,
		Address:                 address,
		Logger:                  logger,
		shouldSendBatchRequests: true,
	}
}

// Send metrics to the server.
func (sender *MetricSenderImpl) Send(metrics map[string]models.Metric) error {
	if sender.shouldSendBatchRequests {
		if err := sender.sendBatch(metrics); err != nil {
			if _, ok := err.(*NotFoundError); ok {
				sender.shouldSendBatchRequests = false
			} else {
				return err
			}
		} else {
			return nil
		}
	}

	for _, metric := range metrics {
		if err := sender.sendSingle(metric); err != nil {
			return err
		}
	}

	return nil
}

func (sender *MetricSenderImpl) sendBatch(metrics map[string]models.Metric) error {
	req, err := sender.BuildBatchRequest(helpers.ConvertMapToSlice(metrics))
	if err != nil {
		return err
	}

	resp, err := sender.Client.Do(req)
	if err != nil {
		sender.Logger.Errorw("Fail send batch request", "error", err.Error())
		return fmt.Errorf("error sending batch request for metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{Message: "batch request returned 404"}
	}

	if resp.StatusCode != http.StatusOK {
		sender.Logger.Errorw("Bad response", "code", resp.StatusCode)
		return fmt.Errorf("batch request failed with status code %d", resp.StatusCode)
	}

	return nil
}

func (sender *MetricSenderImpl) sendSingle(metric models.Metric) error {
	req, err := sender.BuildSingleRequest(metric)
	if err != nil {
		sender.Logger.Errorw("Fail build single request", "metric", metric, "error", err.Error())
		return fmt.Errorf("error creating  single request for metric %s: %v", metric.ID, err)
	}

	resp, err := sender.Client.Do(req)
	if err != nil {
		sender.Logger.Errorw("Fail send request", "metric", metric, "error", err)
		return fmt.Errorf("error sending single request for metric %s: %v", metric.ID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sender.Logger.Errorw("Bad response", "code", resp.StatusCode)
		return fmt.Errorf("request for metric %s failed with status code %d", metric.ID, resp.StatusCode)
	}

	return nil
}

// BuildSingleRequest creates a compressed HTTP request for a single metric.
func (sender *MetricSenderImpl) BuildSingleRequest(metric models.Metric) (*http.Request, error) {
	url := fmt.Sprintf("http://%s/update/", sender.Address)

	body, err := json.Marshal(metric)
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

func (e *NotFoundError) Error() string {
	return e.Message
}
