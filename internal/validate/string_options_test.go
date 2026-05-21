package validate

import "testing"

func TestRange(t *testing.T) {
	var tests = []struct {
		wantErr error
		name    string
		input   string
	}{
		{
			wantErr: nil,
			name:    "Valid range",
			input:   "hello",
		},
		{
			wantErr: ErrMinLimit,
			name:    "Precedes minimum",
			input:   "he",
		},
		{
			wantErr: ErrMaxLimit,
			name:    "Exceeds maximum",
			input:   "hello, world!",
		},
	}

	opt := Range(3, 8)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := opt(tt.input)
			if tt.wantErr != actual {
				t.Fatalf("want %v, got %v", tt.wantErr, actual)
			}
		})
	}
}

func TestRange_Panic(t *testing.T) {
	var tests = []struct {
		name string
		fn   func() stringOption
	}{
		{
			name: "Minimum is 0",
			fn:   func() stringOption { return Range(0, 2) },
		},
		{
			name: "Maximum is 0",
			fn:   func() stringOption { return Range(2, 0) },
		},
		{
			name: "Minimum is greater than maximum",
			fn:   func() stringOption { return Range(7, 6) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if v := recover(); v == nil {
					t.Fatal("expected to panic")
				}
			}()
			tt.fn()
		})
	}
}

func TestMin(t *testing.T) {
	var tests = []struct {
		wantErr error
		name    string
		input   string
	}{
		{
			wantErr: nil,
			name:    "Exceeds minimum",
			input:   "hello",
		},
		{
			wantErr: ErrMinLimit,
			name:    "Precedes minimum",
			input:   "123",
		},
	}

	opt := Min(4)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := opt(tt.input)
			if tt.wantErr != actual {
				t.Fatalf("want %v, got %v", tt.wantErr, actual)
			}
		})
	}
}

func TestMin_Panic(t *testing.T) {
	defer func() {
		if v := recover(); v == nil {
			t.Fatal("expected to panic")
		}
	}()
	Min(0)
}

func TestMax(t *testing.T) {
	var tests = []struct {
		wantErr error
		name    string
		input   string
	}{
		{
			wantErr: ErrMaxLimit,
			name:    "Exceeds maximum",
			input:   "hello, world!",
		},
		{
			wantErr: nil,
			name:    "Precedes maximum",
			input:   "12345",
		},
	}

	opt := Max(6)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := opt(tt.input)
			if tt.wantErr != actual {
				t.Fatalf("want %v, got %v", tt.wantErr, actual)
			}
		})
	}
}

func TestMax_Panic(t *testing.T) {
	defer func() {
		if v := recover(); v == nil {
			t.Fatal("expected to panic")
		}
	}()
	Max(0)
}

func TestNotEmpty(t *testing.T) {
	var tests = []struct {
		wantErr error
		name    string
		input   string
	}{
		{
			wantErr: nil,
			name:    "Non-empty string",
			input:   "hello, world!",
		},
		{
			wantErr: ErrEmpty,
			name:    "Empty string",
			input:   "",
		},
		{
			wantErr: nil,
			name:    "String with only whitespace",
			input:   "    ",
		},
	}

	opt := NotEmpty()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := opt(tt.input)
			if tt.wantErr != actual {
				t.Fatalf("want %v, got %v", tt.wantErr, actual)
			}
		})
	}
}

func TestNoWhitespace(t *testing.T) {
	var tests = []struct {
		wantErr error
		name    string
		input   string
	}{
		{
			wantErr: nil,
			name:    "No whitespace",
			input:   "username123",
		},
		{
			wantErr: ErrContainsWhitespace,
			name:    "Contains whitespace",
			input:   "hello, world!",
		},
		{
			wantErr: nil,
			name:    "Empty string",
			input:   "",
		},
		{
			wantErr: ErrContainsWhitespace,
			name:    "Only whitespace",
			input:   "      ",
		},
	}

	opt := NoWhitespace()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := opt(tt.input)
			if tt.wantErr != actual {
				t.Fatalf("want %v, got %v", tt.wantErr, actual)
			}
		})
	}
}
