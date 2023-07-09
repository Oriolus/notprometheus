package metric

var _ Counter = (*counter)(nil)

type counter struct {
	name  string
	value int64
}

func NewCounter(name string) (Counter, error) {
	if name == "" {
		return nil, EmptyNameError
	}

	return &counter{name: name, value: 1}, nil
}

func (c *counter) Inc() {
	c.value++
}

func (c *counter) Value() int64 {
	return c.value
}

func (c *counter) Name() string {
	return c.name
}
