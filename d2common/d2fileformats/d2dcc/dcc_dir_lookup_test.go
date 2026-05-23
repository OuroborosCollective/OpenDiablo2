package d2dcc

import (
	"testing"
)

func TestDir64ToDcc(t *testing.T) {
	tests := []struct {
		name          string
		direction     int
		numDirections int
		want          int
	}{
		{"4 directions - 0", 0, 4, 0},
		{"4 directions - 8", 8, 4, 1},
		{"4 directions - 24", 24, 4, 2},
		{"4 directions - 40", 40, 4, 3},
		{"4 directions - 63", 63, 4, 0},

		{"8 directions - 0", 0, 8, 4},
		{"8 directions - 4", 4, 8, 0},
		{"8 directions - 12", 12, 8, 5},
		{"8 directions - 20", 20, 8, 1},
		{"8 directions - 28", 28, 8, 6},
		{"8 directions - 36", 36, 8, 2},
		{"8 directions - 44", 44, 8, 7},
		{"8 directions - 52", 52, 8, 3},
		{"8 directions - 60", 60, 8, 4},

		{"16 directions - 0", 0, 16, 4},
		{"16 directions - 2", 2, 16, 8},
		{"16 directions - 6", 6, 16, 0},
		{"16 directions - 10", 10, 16, 9},
		{"16 directions - 14", 14, 16, 5},

		{"32 directions - 0", 0, 32, 4},
		{"32 directions - 1", 1, 32, 16},
		{"32 directions - 3", 3, 32, 8},
		{"32 directions - 5", 5, 32, 17},
		{"32 directions - 7", 7, 32, 0},

		{"64 directions - 0", 0, 64, 4},
		{"64 directions - 1", 1, 64, 32},
		{"64 directions - 2", 2, 64, 16},
		{"64 directions - 3", 3, 64, 33},
		{"64 directions - 4", 4, 64, 8},
		{"64 directions - 63", 63, 64, 63},

		{"Invalid numDirections", 0, 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dir64ToDcc(tt.direction, tt.numDirections); got != tt.want {
				t.Errorf("Dir64ToDcc(%d, %d) = %v, want %v", tt.direction, tt.numDirections, got, tt.want)
			}
		})
	}
}

func TestDir64ToDcc_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Dir64ToDcc did not panic on out of bounds direction")
		}
	}()

	// This should panic as it accesses dir4[64]
	Dir64ToDcc(64, 4)
}

func TestDir64ToDcc_PanicNegative(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Dir64ToDcc did not panic on negative direction")
		}
	}()

	// This should panic as it accesses dir4[-1]
	Dir64ToDcc(-1, 4)
}
