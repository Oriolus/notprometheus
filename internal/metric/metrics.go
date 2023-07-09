package metric

import (
	"errors"
)

type Type string

var (
	UnknownMetricTypeError = errors.New("unknown metric type error")
	WrongMetricName        = errors.New("wrong metric name")
)

const (
	TypeNone    = Type("none")
	TypeGauge   = Type("gauge")
	TypeCounter = Type("counter")
)

func GetMetricType(typ string) (Type, error) {
	innerType := Type(typ)
	if innerType == TypeGauge {
		return TypeGauge, nil
	}
	if innerType == TypeCounter {
		return TypeCounter, nil
	}

	return TypeNone, UnknownMetricTypeError
}

type Gauge interface {
	Value() float64
	Name() string
}

type Counter interface {
	Value() int64
	Name() string
}
