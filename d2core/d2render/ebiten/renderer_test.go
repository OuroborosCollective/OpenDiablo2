package ebiten

import (
	"testing"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2config"
	"github.com/stretchr/testify/assert"
)

func TestSetGamma(t *testing.T) {
	r, err := CreateRenderer(nil)
	assert.NoError(t, err)

	// Default values
	assert.Equal(t, defaultGamma, r.gamma)

	// Update gamma
	newGamma := 1.5
	r.SetGamma(newGamma)
	assert.Equal(t, newGamma, r.gamma)

	// Invalid gamma (should be ignored)
	r.SetGamma(-1.0)
	assert.Equal(t, newGamma, r.gamma)
}

func TestSetContrast(t *testing.T) {
	r, err := CreateRenderer(nil)
	assert.NoError(t, err)

	// Default values
	assert.Equal(t, defaultContrast, r.contrast)

	// Update contrast
	newContrast := 1.2
	r.SetContrast(newContrast)
	assert.Equal(t, newContrast, r.contrast)

	// Invalid contrast (should be ignored)
	r.SetContrast(-0.5)
	assert.Equal(t, newContrast, r.contrast)
}

func TestCreateRendererWithConfig(t *testing.T) {
	cfg := &d2config.Configuration{
		Gamma:    2.2,
		Contrast: 1.5,
	}

	r, err := CreateRenderer(cfg)
	assert.NoError(t, err)

	assert.Equal(t, 2.2, r.gamma)
	assert.Equal(t, 1.5, r.contrast)
}
