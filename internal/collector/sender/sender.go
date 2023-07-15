package sender

import (
	"errors"

	"github.com/oriolus/notprometheus/internal/metric"
)

var (
	ErrStringIsEmpty = errors.New("string argument is empty")
)

type MetricSender interface {
	UpdateGauge(gauge metric.Gauge) error
	UpdateCounter(counter metric.Counter) error
}
