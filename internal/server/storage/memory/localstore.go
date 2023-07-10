package memory

import (
	"github.com/oriolus/notprometheus/internal/metric"
	"github.com/oriolus/notprometheus/internal/server/storage"
)

var _ storage.Storage = (*MemStorage)(nil)

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
		return nil, storage.MetricNotFoundError
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

func (s *MemStorage) AllGauges() []metric.Gauge {
	gs := make([]metric.Gauge, 0, len(s.gauges))
	for _, g := range s.gauges {
		gs = append(gs, g)
	}
	return gs
}

func (s *MemStorage) GetCounter(name string) (metric.Counter, error) {
	cnt, ok := s.counters[name]
	if !ok {
		return nil, storage.MetricNotFoundError
	}
	return cnt, nil
}

func (s *MemStorage) SetCounter(counter metric.Counter) error {
	if counter == nil {
		return storage.ArgumentNilError
	}

	s.counters[counter.Name()] = counter
	return nil
}

func (s *MemStorage) AllCounters() []metric.Counter {
	cs := make([]metric.Counter, 0, len(s.counters))
	for _, c := range s.counters {
		cs = append(cs, c)
	}
	return cs
}
