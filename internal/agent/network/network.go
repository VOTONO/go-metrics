package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VOTONO/go-metrics/internal/models"

	"go.uber.org/zap"
)

type Network interface {
	Send(metrics map[string]models.Metric) error
	BuildRequest(metric models.Metric) (*http.Request, error)
}

// NetworkImpl implements MetricsSender using HTTP.
type NetworkImpl struct {
	Client  *http.Client
	Address string
	Logger  zap.SugaredLogger
}

func New(client *http.Client, address string, logger zap.SugaredLogger) *NetworkImpl {
	return &NetworkImpl{
		Client:  client,
		Address: address,
		Logger:  logger,
	}
}

// Send metrics to the server.
func (sender *NetworkImpl) Send(metrics map[string]models.Metric) error {
	for _, metric := range metrics {
		req, err := sender.BuildRequest(metric)
		if err != nil {
			sender.Logger.Errorw(
				"Fail build request",
				"metric", metric,
			)
			return fmt.Errorf("error creating HTTP request for metric %s: %v", metric.ID, err)
		}

		resp, err := sender.Client.Do(req)
		if err != nil {
			sender.Logger.Errorw(
				"Fail send request",
				"metric ", metric,
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

// Build an HTTP request for a metric.
func (sender *NetworkImpl) BuildRequest(metric models.Metric) (*http.Request, error) {
	url := fmt.Sprintf("http://%s/update", sender.Address)

	body, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
