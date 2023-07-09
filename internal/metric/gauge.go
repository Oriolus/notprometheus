package metric

var _ Gauge = (*gauge)(nil)

type gauge struct {
	name  string
	value float64
}

func NewGauge(name string) Gauge {
	return &gauge{name: name}
}

func (g *gauge) Value() float64 {
	return g.value
}

func (g *gauge) Name() string {
	return g.name
}
