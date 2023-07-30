package sender

import (
	"errors"
	"time"

	"github.com/oriolus/notprometheus/internal/metric"
)

const ClientTimeout = 300 * time.Millisecond

var (
	ErrStringIsEmpty = errors.New("string argument is empty")
)

type MetricSender interface {
	UpdateGauge(gauge metric.Gauge) error
	UpdateCounter(counter metric.Counter) error
}
