package d2config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Gamma != 1.0 {
		t.Errorf("expected default Gamma to be 1.0, got %f", cfg.Gamma)
	}

	if cfg.Contrast != 1.0 {
		t.Errorf("expected default Contrast to be 1.0, got %f", cfg.Contrast)
	}
}
