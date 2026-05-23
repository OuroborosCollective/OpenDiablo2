package d2player

import (
	"testing"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2config"
)

func TestGetIndexForValue(t *testing.T) {
	m := &EscapeMenu{}
	options := []string{"0.5", "1.0", "1.5"}

	tests := []struct {
		value    float64
		expected int
	}{
		{0.5, 0},
		{1.0, 1},
		{1.5, 2},
		{0.7, 5}, // default fallback for getIndexForValue as currently implemented
	}

	for _, tt := range tests {
		result := m.getIndexForValue(tt.value, options)
		if result != tt.expected {
			t.Errorf("getIndexForValue(%f) = %d; want %d", tt.value, result, tt.expected)
		}
	}
}
