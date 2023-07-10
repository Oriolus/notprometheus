package collector

import "errors"

var (
	ArgumentNilError   = errors.New("argument is nil")
	StringIsEmptyError = errors.New("string argument is empty")
)
