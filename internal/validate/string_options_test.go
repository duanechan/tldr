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
			input:   "hello",
		},
	}

	opt := Range(3, 8)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := opt(tt.input); !tt.wantErr && err != nil {
				t.Fatalf("expected error to be nil, got %v", err.Error())
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
