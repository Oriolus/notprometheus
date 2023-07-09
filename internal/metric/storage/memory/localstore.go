package memory

import (
	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/metric/storage"
)

// MemStorage warning! thread unsafe
type MemStorage struct {
	gauges   map[string]metric.Gauge
	counters map[string]metric.Counter
}

func NewMemStorage() storage.Storage {
	return &MemStorage{
		gauges:   make(map[string]metric.Gauge),
		counters: make(map[string]metric.Counter),
	}
}

func (s *MemStorage) GetGauge(name string) (metric.Gauge, error) {
	g, ok := s.gauges[name]
	if !ok {
		return nil, storage.MetricNotFound
	}
	return g, nil
}

func (s *MemStorage) SetGauge(gauge metric.Gauge) error {
	if gauge == nil {
		return storage.ArgumentNilError
	}

	s.gauges[gauge.Name()] = gauge
	return nil
}

func (s *MemStorage) GetCounter(name string) (metric.Counter, error) {
	cnt, ok := s.counters[name]
	if !ok {
		return nil, storage.MetricNotFound
	}
	return cnt, nil
}

func (s *MemStorage) SetCounter(counter metric.Counter) error {
	if counter == nil {
		return storage.ArgumentNilError
	}

	s.counters[counter.Name()] = counter
	return storage.MetricNotFound
}
