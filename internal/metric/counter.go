package metric

var _ Counter = (*counter)(nil)

type counter struct {
	name  string
	value int64
}

func NewCounterWithValue(name string, value int64) (Counter, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	return &counter{name: name, value: value}, nil
}

func NewCounter(name string) (Counter, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	return &counter{name: name, value: 1}, nil
}

func (c *counter) Add(value int64) {
	c.value += value
}

func (c *counter) Value() int64 {
	return c.value
}

func (c *counter) Name() string {
	return c.name
}

func (c *counter) Inc() {
	c.value++
}
