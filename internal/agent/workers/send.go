package workers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/agent/helpers"
	"github.com/VOTONO/go-metrics/internal/agent/semaphore"
	"github.com/VOTONO/go-metrics/internal/compressor"
	"github.com/VOTONO/go-metrics/internal/models"
)

type SendWorker struct {
	client       *http.Client
	logger       *zap.SugaredLogger
	ticker       *time.Ticker
	address      string
	stopChannel  chan struct{}
	inputChannel <-chan []models.Metric
	semaphore    *semaphore.Semaphore
	secretKey    string
	waitGroup    sync.WaitGroup
}

func NewSendWorker(
	client *http.Client,
	logger *zap.SugaredLogger,
	interval int,
	inputChannel <-chan []models.Metric,
	rateLimit int,
	address string,
	secretKey string) *SendWorker {
	return &SendWorker{
		client:       client,
		logger:       logger,
		ticker:       helpers.CreateTicker(interval),
		stopChannel:  make(chan struct{}),
		inputChannel: inputChannel,
		semaphore:    semaphore.NewSemaphore(rateLimit),
		address:      address,
		secretKey:    secretKey,
		waitGroup:    sync.WaitGroup{},
	}
}

func (w *SendWorker) Start() {
	w.logger.Info("starting send worker")
	for {
		select {
		case <-w.stopChannel:
			w.logger.Info("stopping send worker")
			w.ticker.Stop()
			w.waitGroup.Wait()
			return
		case metrics := <-w.inputChannel:
			go w.send(metrics)
		}
	}
}

func (w *SendWorker) Stop() {
	close(w.stopChannel)
}

// buildRequest creates a compressed HTTP request for a batch of metrics.
func (w *SendWorker) buildRequest(metrics []models.Metric) (*http.Request, error) {
	url := fmt.Sprintf("http://%s/updates/", w.address)

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

	if w.secretKey != "" {
		h := hmac.New(sha256.New, []byte(w.secretKey))
		h.Write(compressedBody)
		req.Header.Set("HashSHA256", hex.EncodeToString(h.Sum(nil)))
	}

	return req, nil
}

// send metrics to the server.
func (w *SendWorker) send(metrics []models.Metric) error {
	w.waitGroup.Add(1)
	defer w.waitGroup.Done()

	req, buildReqErr := w.buildRequest(metrics)
	if buildReqErr != nil {
		return fmt.Errorf("failed to build batch request: %s", buildReqErr.Error())
	}

	err := w.sendRequest(req)
	w.logger.Infow("sent batch", "count", len(metrics))
	if err != nil {
		retryCount := 3
		retryPause := 1 * time.Second

		for i := 0; i < retryCount; i++ {
			time.Sleep(retryPause)
			err = w.sendRequest(req)

			if err == nil {
				return nil
			}

			retryPause += 2
		}
	}

	return nil
}

func (w *SendWorker) sendRequest(req *http.Request) error {
	w.semaphore.Acquire()
	defer w.semaphore.Release()

	resp, err := w.client.Do(req)
	if err != nil {
		w.logger.Errorw("Failed to send batch request", "error", err)
		return fmt.Errorf("error sending batch request for metrics: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			w.logger.Errorw("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		w.logger.Errorw("Received bad response", "status code", resp.StatusCode)
		return fmt.Errorf("batch request failed with status code %d", resp.StatusCode)
	}

	return nil
}
