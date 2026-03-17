package proc

import (
	"testing"
)

func TestParseHexPort(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		wantErr  bool
	}{
		{"0F02000A:1F90", 8080, false},
		{"00000000:0050", 80, false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseHexPort(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHexPort(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("parseHexPort(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}
