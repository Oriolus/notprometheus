package storage

import (
	"errors"

	"github.com/oriolus/notprometheus/internal/metric"
)

var (
	MetricNotFoundError = errors.New("metric not found")
)

type Storage interface {
	GetGauge(name string) (metric.Gauge, error)
	SetGauge(gauge metric.Gauge) error

	GetCounter(name string) (metric.Counter, error)
	SetCounter(counter metric.Counter) error
}
