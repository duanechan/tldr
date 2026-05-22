package validate

import (
	"testing"
)

func TestString_NoOptions(t *testing.T) {
	var tests = []struct {
		name  string
		input string
	}{
		{
			name:  "String with length",
			input: "hello, world!",
		},
		{
			name:  "Empty string",
			input: "",
		},
		{
			name:  "String with only whitespace",
			input: "     ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, errs := String(tt.input)
			if errs != nil {
				t.Fatal("expected errors to be nil, got non-empty errors")
			}

			if tt.input != actual {
				t.Fatalf("want %s, got %s", tt.input, actual)
			}
		})
	}
}

func TestString_WithOptions(t *testing.T) {
	var tests = []struct {
		wantErr  bool
		name     string
		input    string
		expected string
		fn       func(string) (string, []error)
	}{
		{
			wantErr:  false,
			name:     "String within range",
			input:    "Hello!",
			expected: "Hello!",
			fn: func(s string) (string, []error) {
				return String(s, Range(3, 8))
			},
		},
		{
			wantErr:  true,
			name:     "String exceeding range",
			input:    "Hello, world!",
			expected: "",
			fn: func(s string) (string, []error) {
				return String(s, Range(3, 8))
			},
		},
		{
			wantErr:  false,
			name:     "String within minimum limit",
			input:    "yes",
			expected: "yes",
			fn: func(s string) (string, []error) {
				return String(s, Min(3))
			},
		},
		{
			wantErr:  true,
			name:     "String preceding minimum limit",
			input:    "no",
			expected: "",
			fn: func(s string) (string, []error) {
				return String(s, Min(3))
			},
		},
		{
			wantErr:  false,
			name:     "String within maximum limit",
			input:    "Hello!",
			expected: "Hello!",
			fn: func(s string) (string, []error) {
				return String(s, Max(6))
			},
		},
		{
			wantErr:  true,
			name:     "String exceeding maximum limit",
			input:    "Hello, world!",
			expected: "",
			fn: func(s string) (string, []error) {
				return String(s, Max(6))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.fn(tt.input)
			if !tt.wantErr && err != nil {
				t.Fatal("expected errors to be nil, got non-empty errors")
			}

			if tt.wantErr && err == nil {
				t.Fatal("expected error to be non-nil")
			}

			if tt.expected != actual {
				t.Fatalf("want %v, got %v", tt.expected, actual)
			}
		})
	}
}
