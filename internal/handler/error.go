package handler

import "errors"

var (
	ErrArgumentNil = errors.New("argument nil error")
)

var (
	errInvalidCounterReq = errors.New("req.Delta cannot be nil for counter")
	errInvalidGaugeReq   = errors.New("req.Value cannot be nil for gauge")
)
