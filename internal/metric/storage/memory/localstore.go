package memory

import (
	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/metric/storage"
)

type MemStorage struct {
	gauges   map[string]metric.Gauge
	counters map[string]metric.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]metric.Gauge),
		counters: make(map[string]metric.Counter),
	}
}

func (s *MemStorage) GetGauge(name string) (metric.Gauge, error) {
	return nil, storage.MetricNotFound
}

func (s *MemStorage) SetGauge(gauge metric.Gauge) error {
	return storage.MetricNotFound
}

func (s *MemStorage) GetCounter(name string) (metric.Counter, error) {
	return nil, storage.MetricNotFound
}

func (s *MemStorage) SetCounter(counter metric.Counter) error {
	return storage.MetricNotFound
}
