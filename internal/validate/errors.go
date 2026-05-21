package validate

import "errors"

var (
	ErrMaxLimit           = errors.New("input exceeds maximum limit")
	ErrMinLimit           = errors.New("input precedes minimum limit")
	ErrEmpty              = errors.New("input is empty")
	ErrContainsWhitespace = errors.New("input contains whitespace")
)
