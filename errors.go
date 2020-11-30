package registry

import (
	"errors"
)

var (
	ErrDriverExists  = errors.New("driver has already registered")
	ErrUnknownDriver = errors.New("unknown registry driver")
)
