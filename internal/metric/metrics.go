package metric

type Type string

const (
	TypeNone    = Type("none")
	TypeGauge   = Type("gauge")
	TypeCounter = Type("counter")
)

func GetMetricType(typ string) (Type, error) {
	if typ == TypeGauge {
		return TypeGauge, nil
	}
	if typ == TypeCounter {
		return TypeCounter, nil
	}

	return TypeNone, UnknownMetricTypeError
}

type Gauge interface {
	Set(value float64)
	Value() float64
	Name() string
}

type Counter interface {
	Inc()
	Value() int64
	Name() string
}
