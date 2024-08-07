package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VOTONO/go-metrics/internal/compressor"
	"github.com/VOTONO/go-metrics/internal/models"

	"go.uber.org/zap"
)

type MetricSender interface {
	Send(metrics map[string]models.Metric) error
	BuildRequest(metric models.Metric) (*http.Request, error)
}

// MetricSenderImpl implements MetricsSender using HTTP.
type MetricSenderImpl struct {
	Client  *http.Client
	Address string
	Logger  zap.SugaredLogger
}

func New(client *http.Client, address string, logger zap.SugaredLogger) *MetricSenderImpl {
	return &MetricSenderImpl{
		Client:  client,
		Address: address,
		Logger:  logger,
	}
}

// Send metrics to the server.
func (sender *MetricSenderImpl) Send(metrics map[string]models.Metric) error {
	for _, metric := range metrics {
		req, err := sender.BuildRequest(metric)
		if err != nil {
			sender.Logger.Errorw(
				"Fail build request",
				"metric", metric,
				"error", err,
			)
			return fmt.Errorf("error creating HTTP request for metric %s: %v", metric.ID, err)
		}

		resp, err := sender.Client.Do(req)
		if err != nil {
			sender.Logger.Errorw(
				"Fail send request",
				"metric ", metric,
				"error", err,
			)
			return fmt.Errorf("error sending HTTP request for metric %s: %v", metric.ID, err)
		}

		if resp.StatusCode != http.StatusOK {
			sender.Logger.Errorw(
				"Bad response",
				"code ", resp.StatusCode,
			)
			return fmt.Errorf("request for metric %s failed with status code %d", metric.ID, resp.StatusCode)
		}

		resp.Body.Close()
	}
	return nil
}

// BuildRequest  for a metric.
func (sender *MetricSenderImpl) BuildRequest(metric models.Metric) (*http.Request, error) {
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
