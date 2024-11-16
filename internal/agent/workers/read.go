package workers

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/agent/helpers"
	"github.com/VOTONO/go-metrics/internal/constants"
	"github.com/VOTONO/go-metrics/internal/models"
)

// ReadWorker periodically read runtimes metrics and sendWithRetry it to result channel.
type ReadWorker struct {
	ResultChannel chan []models.Metric
	logger        *zap.SugaredLogger
	ticker        *time.Ticker
	stopChannel   chan struct{}
	count         int64
}

func NewReadWorker(logger *zap.SugaredLogger, interval int) ReadWorker {
	readResultChannel := make(chan []models.Metric, 1)
	return ReadWorker{
		ResultChannel: readResultChannel,
		logger:        logger,
		ticker:        helpers.CreateTicker(interval),
		stopChannel:   make(chan struct{}),
		count:         0,
	}
}

// Start read metrics.
func (w *ReadWorker) Start() {
	w.logger.Infow("starting readRuntimeMetrics worker")
	for {
		select {
		case <-w.stopChannel:
			w.logger.Infow("stopping readRuntimeMetrics worker")
			w.ticker.Stop()
			close(w.ResultChannel)
			return
		case <-w.ticker.C:
			w.logger.Infow("readRuntimeMetrics metrics")
			w.ResultChannel <- w.readRuntimeMetrics()
			w.ResultChannel <- w.readMemoryMetrics()
		}
	}
}

// Stop work and close channels.
func (w *ReadWorker) Stop() {
	close(w.stopChannel)
}

// readMemoryMetrics reads metrics from mem package.
func (w *ReadWorker) readMemoryMetrics() []models.Metric {
	virtualMemory, _ := mem.VirtualMemory()
	var metrics []models.Metric

	totalMemory := float64(virtualMemory.Total)
	freeMemory := float64(virtualMemory.Free)
	usedMemory := float64(virtualMemory.UsedPercent)

	metrics = append(metrics, models.Metric{ID: "TotalMemory", MType: constants.Gauge, Value: &totalMemory})
	metrics = append(metrics, models.Metric{ID: "FreeMemory", MType: constants.Gauge, Value: &freeMemory})
	metrics = append(metrics, models.Metric{ID: "UsedMemory", MType: constants.Gauge, Value: &usedMemory})

	return metrics
}

// readRuntimeMetrics reads metrics from runtime package.
func (w *ReadWorker) readRuntimeMetrics() []models.Metric {
	var metrics []models.Metric

	w.count++
	random := rand.Float64()

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	metrics = append(metrics, models.Metric{ID: "PollCount", MType: "counter", Delta: &w.count})
	metrics = append(metrics, models.Metric{ID: "RandomValue", MType: constants.Gauge, Value: float64Ptr(random)})

	metrics = append(metrics, models.Metric{ID: "Alloc", MType: constants.Gauge, Value: float64Ptr(float64(stats.Alloc))})
	metrics = append(metrics, models.Metric{ID: "BuckHashSys", MType: constants.Gauge, Value: float64Ptr(float64(stats.BuckHashSys))})
	metrics = append(metrics, models.Metric{ID: "Frees", MType: constants.Gauge, Value: float64Ptr(float64(stats.Frees))})
	metrics = append(metrics, models.Metric{ID: "GCCPUFraction", MType: constants.Gauge, Value: float64Ptr(stats.GCCPUFraction)})
	metrics = append(metrics, models.Metric{ID: "GCSys", MType: constants.Gauge, Value: float64Ptr(float64(stats.GCSys))})
	metrics = append(metrics, models.Metric{ID: "HeapAlloc", MType: constants.Gauge, Value: float64Ptr(float64(stats.HeapAlloc))})
	metrics = append(metrics, models.Metric{ID: "HeapIdle", MType: constants.Gauge, Value: float64Ptr(float64(stats.HeapIdle))})
	metrics = append(metrics, models.Metric{ID: "HeapInuse", MType: constants.Gauge, Value: float64Ptr(float64(stats.HeapInuse))})
	metrics = append(metrics, models.Metric{ID: "HeapObjects", MType: constants.Gauge, Value: float64Ptr(float64(stats.HeapObjects))})
	metrics = append(metrics, models.Metric{ID: "HeapReleased", MType: constants.Gauge, Value: float64Ptr(float64(stats.HeapReleased))})
	metrics = append(metrics, models.Metric{ID: "HeapSys", MType: constants.Gauge, Value: float64Ptr(float64(stats.HeapSys))})
	metrics = append(metrics, models.Metric{ID: "LastGC", MType: constants.Gauge, Value: float64Ptr(float64(stats.LastGC))})
	metrics = append(metrics, models.Metric{ID: "Lookups", MType: constants.Gauge, Value: float64Ptr(float64(stats.Lookups))})
	metrics = append(metrics, models.Metric{ID: "MCacheInuse", MType: constants.Gauge, Value: float64Ptr(float64(stats.MCacheInuse))})
	metrics = append(metrics, models.Metric{ID: "MCacheSys", MType: constants.Gauge, Value: float64Ptr(float64(stats.MCacheSys))})
	metrics = append(metrics, models.Metric{ID: "MSpanInuse", MType: constants.Gauge, Value: float64Ptr(float64(stats.MSpanInuse))})
	metrics = append(metrics, models.Metric{ID: "MSpanSys", MType: constants.Gauge, Value: float64Ptr(float64(stats.MSpanSys))})
	metrics = append(metrics, models.Metric{ID: "Mallocs", MType: constants.Gauge, Value: float64Ptr(float64(stats.Mallocs))})
	metrics = append(metrics, models.Metric{ID: "NextGC", MType: constants.Gauge, Value: float64Ptr(float64(stats.NextGC))})
	metrics = append(metrics, models.Metric{ID: "NumForcedGC", MType: constants.Gauge, Value: float64Ptr(float64(stats.NumForcedGC))})
	metrics = append(metrics, models.Metric{ID: "NumGC", MType: constants.Gauge, Value: float64Ptr(float64(stats.NumGC))})
	metrics = append(metrics, models.Metric{ID: "OtherSys", MType: constants.Gauge, Value: float64Ptr(float64(stats.OtherSys))})
	metrics = append(metrics, models.Metric{ID: "PauseTotalNs", MType: constants.Gauge, Value: float64Ptr(float64(stats.PauseTotalNs))})
	metrics = append(metrics, models.Metric{ID: "StackInuse", MType: constants.Gauge, Value: float64Ptr(float64(stats.StackInuse))})
	metrics = append(metrics, models.Metric{ID: "StackSys", MType: constants.Gauge, Value: float64Ptr(float64(stats.StackSys))})
	metrics = append(metrics, models.Metric{ID: "Sys", MType: constants.Gauge, Value: float64Ptr(float64(stats.Sys))})
	metrics = append(metrics, models.Metric{ID: "TotalAlloc", MType: constants.Gauge, Value: float64Ptr(float64(stats.TotalAlloc))})

	return metrics
}

// Helper function to create a pointer to float64
func float64Ptr(v float64) *float64 {
	return &v
}
