package logic

import (
	"math/rand"
	"runtime"

	"github.com/VOTONO/go-metrics/internal/constants"
	"github.com/VOTONO/go-metrics/internal/models"
)

type MetricReader interface {
	Read() map[string]models.Metric
}

type MetricReaderImpl struct {
	metrics map[string]models.Metric
	count   int64
}

func NewMetricReaderImpl() MetricReader {
	return &MetricReaderImpl{
		metrics: make(map[string]models.Metric),
		count:   0,
	}
}

func (m *MetricReaderImpl) Read() map[string]models.Metric {
	m.count++
	m.metrics[PollCount] = models.Metric{
		ID:    PollCount,
		MType: constants.Counter,
		Delta: &m.count,
	}
	random := rand.Float64()
	m.metrics[RandomValue] = models.Metric{
		ID:    RandomValue,
		MType: constants.Gauge,
		Value: &random,
	}

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	m.metrics[Alloc] = makeGaugeMetric(Alloc, float64(stats.Alloc))
	m.metrics[BuckHashSys] = makeGaugeMetric(BuckHashSys, float64(stats.BuckHashSys))
	m.metrics[Frees] = makeGaugeMetric(Frees, float64(stats.Frees))
	m.metrics[GCCPUFraction] = makeGaugeMetric(GCCPUFraction, stats.GCCPUFraction)
	m.metrics[GCSys] = makeGaugeMetric(GCSys, float64(stats.GCSys))
	m.metrics[HeapAlloc] = makeGaugeMetric(HeapAlloc, float64(stats.HeapAlloc))
	m.metrics[HeapIdle] = makeGaugeMetric(HeapIdle, float64(stats.HeapIdle))
	m.metrics[HeapInuse] = makeGaugeMetric(HeapInuse, float64(stats.HeapInuse))
	m.metrics[HeapObjects] = makeGaugeMetric(HeapObjects, float64(stats.HeapObjects))
	m.metrics[HeapReleased] = makeGaugeMetric(HeapReleased, float64(stats.HeapReleased))
	m.metrics[HeapSys] = makeGaugeMetric(HeapSys, float64(stats.HeapSys))
	m.metrics[LastGC] = makeGaugeMetric(LastGC, float64(stats.LastGC))
	m.metrics[Lookups] = makeGaugeMetric(Lookups, float64(stats.Lookups))
	m.metrics[MCacheInuse] = makeGaugeMetric(MCacheInuse, float64(stats.MCacheInuse))
	m.metrics[MCacheSys] = makeGaugeMetric(MCacheSys, float64(stats.MCacheSys))
	m.metrics[MSpanInuse] = makeGaugeMetric(MSpanInuse, float64(stats.MSpanInuse))
	m.metrics[MSpanSys] = makeGaugeMetric(MSpanSys, float64(stats.MSpanSys))
	m.metrics[Mallocs] = makeGaugeMetric(Mallocs, float64(stats.Mallocs))
	m.metrics[NextGC] = makeGaugeMetric(NextGC, float64(stats.NextGC))
	m.metrics[NumForcedGC] = makeGaugeMetric(NumForcedGC, float64(stats.NumForcedGC))
	m.metrics[NumGC] = makeGaugeMetric(NumGC, float64(stats.NumGC))
	m.metrics[OtherSys] = makeGaugeMetric(OtherSys, float64(stats.OtherSys))
	m.metrics[PauseTotalNs] = makeGaugeMetric(PauseTotalNs, float64(stats.PauseTotalNs))
	m.metrics[StackInuse] = makeGaugeMetric(StackInuse, float64(stats.StackInuse))
	m.metrics[StackSys] = makeGaugeMetric(StackSys, float64(stats.StackSys))
	m.metrics[Sys] = makeGaugeMetric(Sys, float64(stats.Sys))
	m.metrics[TotalAlloc] = makeGaugeMetric(TotalAlloc, float64(stats.TotalAlloc))

	return m.metrics
}

func makeGaugeMetric(id string, value float64) models.Metric {
	return models.Metric{
		ID:    id,
		MType: constants.Gauge,
		Value: &value,
	}
}
