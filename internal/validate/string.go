package validate

// String validates a string, s, using the given validation options.
func String(s string, opts ...stringOption) (string, error) {
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return "", err
		}
	}
	return s, nil
}
