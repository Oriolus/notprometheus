package metric

import "errors"

var (
	ErrUnknownMetricType = errors.New("unknown metric type error")
	ErrEmptyName         = errors.New("empty name")
)
