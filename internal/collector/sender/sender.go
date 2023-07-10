package sender

import (
	"github.com/oriolus/notprometheus/internal/metric"
)

type MetricSender interface {
	UpdateGauge(gauge metric.Gauge) error
	UpdateCounter(counter metric.Counter) error
}
