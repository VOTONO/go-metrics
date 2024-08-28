package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VOTONO/go-metrics/internal/constants"
)

func TestMemStatsMonitor_Read(t *testing.T) {
	monitor := NewMetricReaderImpl()

	metrics := monitor.Read()

	assert.NotNil(t, metrics)
	assert.Equal(t, constants.Counter, metrics[PollCount].MType)
	assert.Equal(t, constants.Gauge, metrics[RandomValue].MType)

	// Ensure MemStats metrics are present
	memStatsFields := []string{
		Alloc, BuckHashSys, Frees, GCCPUFraction, GCSys, HeapAlloc, HeapIdle,
		HeapInuse, HeapObjects, HeapReleased, HeapSys, LastGC, Lookups, MCacheInuse,
		MCacheSys, MSpanInuse, MSpanSys, Mallocs, NextGC, NumForcedGC, NumGC,
		OtherSys, PauseTotalNs, StackInuse, StackSys, Sys, TotalAlloc,
	}

	for _, field := range memStatsFields {
		t.Run(field, func(t *testing.T) {
			assert.Equal(t, constants.Gauge, metrics[field].MType)
		})
	}
}
