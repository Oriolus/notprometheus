package metric

var _ Gauge = (*gauge)(nil)

type gauge struct {
	name  string
	value float64
}

func NewGauge(name string, value float64) (Gauge, error) {
	if name == "" {
		return nil, EmptyNameError
	}

	return &gauge{name: name, value: value}, nil
}

func (g *gauge) Set(value float64) {
	g.value = value
}

func (g *gauge) Value() float64 {
	return g.value
}

func (g *gauge) Name() string {
	return g.name
}
