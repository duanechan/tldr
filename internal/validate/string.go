package validate

// String validates a string, s, using the given validation options.
func String(s string, opts ...stringOption) (string, []error) {
	var errs []error
	for _, opt := range opts {
		if err := opt(s); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return "", errs
	}

	return s, errs
}
