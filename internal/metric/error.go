package metric

import "errors"

var (
	UnknownMetricTypeError = errors.New("unknown metric type error")
	EmptyNameError         = errors.New("empty name")
)
