package validate

import (
	"regexp"
)

type stringOption func(string) error

var (
	whitespaceRegex = regexp.MustCompile(`\s`)
)

// Range returns an option that checks if a string is within a range.
// Panics, if min/max is lesser than or equal to 0, or if min is greater
// than or equal to max.
func Range(min, max int) stringOption {
	if min <= 0 {
		panic("min limit must be greater than 0")
	}
	if max <= 0 {
		panic("max limit must be greater than 0")
	}
	if min > max {
		panic("min limit must be lesser than or equal to the max limit")
	}

	return func(s string) error {
		length := len(s)

		if length < min {
			return ErrMinLimit
		}

		if length > max {
			return ErrMaxLimit
		}

		return nil
	}
}

// Min returns an option that checks if the length of a string meets the
// minimum length.
func Min(min int) stringOption {
	if min <= 0 {
		panic("min limit must be greater than 0")
	}

	return func(s string) error {
		if len(s) < min {
			return ErrMinLimit
		}

		return nil
	}
}

// Max returns an option that checks if the length of a string meets the
// maximum length.
func Max(max int) stringOption {
	if max <= 0 {
		panic("max limit must be greater than 0")
	}

	return func(s string) error {
		if len(s) > max {
			return ErrMaxLimit
		}

		return nil
	}
}

// NotEmpty returns an option that checks if a string is non-empty.
func NotEmpty() stringOption {
	return func(s string) error {
		if len(s) == 0 {
			return ErrEmpty
		}

		return nil
	}
}

// NoWhitespace returns an option that checks if a string contains whitespace.
func NoWhitespace() stringOption {
	return func(s string) error {
		if whitespaceRegex.MatchString(s) {
			return ErrContainsWhitespace
		}

		return nil
	}
}
