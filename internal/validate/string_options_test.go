package validate

import "testing"

func TestRange(t *testing.T) {
	var tests = []struct {
		wantErr bool
		name    string
		input   string
	}{
		{
			wantErr: false,
			name:    "Valid range",
			input:   "hello",
		},
		{
			wantErr: true,
			name:    "Precedes min",
			input:   "he",
		},
		{
			wantErr: true,
			name:    "Exceeds max",
			input:   "hello, world!",
		},
		{
			wantErr: false,
			name:    "Contains whitespace but valid range",
			input:   " space ",
		},
		{
			wantErr: true,
			name:    "Contains whitespace and precedes min",
			input:   "   1          2      ",
		},
		{
			wantErr: true,
			name:    "Contains whitespace and exceeds max",
			input:   " hello world! ",
		},
	}

	opt := Range(3, 8)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := opt(tt.input)
			if !tt.wantErr && err != nil {
				t.Fatalf("expected error to be nil, got %v", err.Error())
			}
			if tt.wantErr && err == nil {
				t.Fatal("expected error to be non-nil")
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
			name: "Min is 0",
			fn:   func() stringOption { return Range(0, 2) },
		},
		{
			name: "Max is 0",
			fn:   func() stringOption { return Range(2, 0) },
		},
		{
			name: "Min is greater than max",
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
		wantErr bool
		name    string
		input   string
	}{
		{
			wantErr: false,
			name:    "Exceeds minimum",
			input:   "hello",
		},
		{
			wantErr: true,
			name:    "Precedes minimum",
			input:   "123",
		},
		{
			wantErr: false,
			name:    "Contains whitespace but exeeds minimum",
			input:   " hello ",
		},
		{
			wantErr: true,
			name:    "Contains whitespace but precedes minimum",
			input:   " 123 ",
		},
	}

	opt := Min(4)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := opt(tt.input)
			if !tt.wantErr && err != nil {
				t.Fatalf("expected error to be nil, got %v", err.Error())
			}
			if tt.wantErr && err == nil {
				t.Fatal("expected error to be non-nil")
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
		wantErr bool
		name    string
		input   string
	}{
		{
			wantErr: true,
			name:    "Exceeds maximum",
			input:   "hello, world!",
		},
		{
			wantErr: false,
			name:    "Precedes maximum",
			input:   "12345",
		},
		{
			wantErr: true,
			name:    "Contains whitespace but exeeds maximum",
			input:   " hello, world! ",
		},
		{
			wantErr: false,
			name:    "Contains whitespace but precedes maximum",
			input:   " 12345 ",
		},
	}

	opt := Max(6)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := opt(tt.input)
			if !tt.wantErr && err != nil {
				t.Fatalf("expected error to be nil, got %v", err.Error())
			}
			if tt.wantErr && err == nil {
				t.Fatal("expected error to be non-nil")
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
		wantErr bool
		name    string
		input   string
	}{
		{
			wantErr: false,
			name:    "Non-empty string",
			input:   "hello, world!",
		},
		{
			wantErr: true,
			name:    "Empty string",
			input:   "",
		},
		{
			wantErr: false,
			name:    "String with only whitespace",
			input:   "    ",
		},
	}

	opt := NotEmpty()

	for _, tt := range tests {
		err := opt(tt.input)
		if !tt.wantErr && err != nil {
			t.Fatalf("expected error to be nil, got %v", err.Error())
		}
		if tt.wantErr && err == nil {
			t.Fatal("expected error to be non-nil")
		}
	}
}
