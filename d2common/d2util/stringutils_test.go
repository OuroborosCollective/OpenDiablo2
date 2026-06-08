package d2util

import "testing"

func TestAsterToEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"asterix", "*test", ""},
		{"no asterix", "test", "test"},
		{"empty", "", ""},
		{"space", " ", " "},
		{"asterix in middle", "te*st", "te*st"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsterToEmpty(tt.input); got != tt.expected {
				t.Errorf("AsterToEmpty(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestEmptyToZero(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", "0"},
		{"single space", " ", "0"},
		{"already zero", "0", "0"},
		{"positive number", "1", "1"},
		{"alpha string", "abc", "abc"},
		{"double space", "  ", "  "},
		{"leading space", " 1", " 1"},
		{"trailing space", "1 ", "1 "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EmptyToZero(tt.input); got != tt.expected {
				t.Errorf("EmptyToZero(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
