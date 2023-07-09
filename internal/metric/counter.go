package metric

var _ Counter = (*counter)(nil)

type counter struct {
	name  string
	value int64
}

func NewCounter(name string) Counter {
	return &counter{name: name}
}

func (c *counter) Value() int64 {
	return c.value
}

func (c *counter) Name() string {
	return c.name
}
