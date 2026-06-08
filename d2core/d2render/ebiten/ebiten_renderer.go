package ebiten

import (
	"errors"
	"image"
	"image/color"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2config"
)

const (
	screenWidth       = 800
	screenHeight      = 600
	defaultSaturation = 1.0
	defaultBrightness = 1.0
	defaultSkewX      = 0.0
	defaultSkewY      = 0.0
	defaultScaleX     = 1.0
	defaultScaleY     = 1.0
	defaultGamma      = 1.0
	defaultContrast   = 1.0
)

const gammaShaderSource = `
package main

var Gamma float32
var Contrast float32

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	clr := imageSrc0At(texCoord)
	if clr.a == 0.0 {
		return clr
	}

	rgb := clr.rgb
	rgb = (rgb - 0.5) * Contrast + 0.5
	rgb = clamp(rgb, 0.0, 1.0)
	rgb = pow(rgb, vec3(1.0 / Gamma))

	return vec4(rgb, clr.a)
}
`

type renderCallback = func(surface d2interface.Surface) error

type updateCallback = func() error

// static check that we implement our renderer interface
var _ d2interface.Renderer = &Renderer{}

// Renderer is an implementation of a renderer
type Renderer struct {
	updateCallback
	renderCallback
	*d2util.GlyphPrinter
	lastRenderError error
	gamma           float64
	contrast        float64
	offscreen       *ebiten.Image
	shader          *ebiten.Shader
	shaderFailed    bool
}

// Update calls the game's logical update function (the `Advance` method)
func (r *Renderer) Update() error {
	if r.updateCallback == nil {
		return errors.New("no update callback defined for ebiten renderer")
	}

	return r.updateCallback()
}

const drawError = "no render callback defined for ebiten renderer"

// Draw updates the screen with the given *ebiten.Image
func (r *Renderer) Draw(screen *ebiten.Image) {
	r.lastRenderError = nil

	if r.renderCallback == nil {
		r.lastRenderError = errors.New(drawError)
		return
	}

	// If gamma and contrast are default, and we don't have a shader, just draw normally
	if r.gamma == defaultGamma && r.contrast == defaultContrast {
		r.lastRenderError = r.renderCallback(createEbitenSurface(r, screen))
		return
	}

	if r.shaderFailed {
		r.lastRenderError = r.renderCallback(createEbitenSurface(r, screen))
		return
	}

	// Setup offscreen buffer if needed
	w, h := screen.Size()
	if r.offscreen == nil {
		r.offscreen = ebiten.NewImage(w, h)
	} else {
		ow, oh := r.offscreen.Size()
		if ow != w || oh != h {
			r.offscreen.Dispose()
			r.offscreen = ebiten.NewImage(w, h)
		}
	}

	r.offscreen.Clear()
	r.lastRenderError = r.renderCallback(createEbitenSurface(r, r.offscreen))
	if r.lastRenderError != nil {
		return
	}

	// Compile shader if needed
	if r.shader == nil {
		var err error
		r.shader, err = ebiten.NewShader([]byte(gammaShaderSource))
		if err != nil {
			r.shaderFailed = true
			// Fallback to normal draw for this frame
			screen.DrawImage(r.offscreen, nil)
			return
		}
	}

	// Draw offscreen to screen with shader
	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = r.offscreen
	op.Uniforms = map[string]interface{}{
		"Gamma":    float32(r.gamma),
		"Contrast": float32(r.contrast),
	}
	screen.DrawRectShader(w, h, r.shader, op)
}

// Layout returns the renderer screen width and height
func (r *Renderer) Layout(_, _ int) (width, height int) {
	return screenWidth, screenHeight
}

// CreateRenderer creates an ebiten renderer instance
func CreateRenderer(cfg *d2config.Configuration) (*Renderer, error) {
	result := &Renderer{
		GlyphPrinter: d2util.NewDebugPrinter(),
		gamma:        defaultGamma,
		contrast:     defaultContrast,
	}

	if cfg != nil {
		config := cfg

		ebiten.SetCursorMode(ebiten.CursorModeHidden)
		ebiten.SetFullscreen(config.FullScreen)
		ebiten.SetRunnableOnUnfocused(config.RunInBackground)
		ebiten.SetVsyncEnabled(config.VsyncEnabled)
		ebiten.SetMaxTPS(config.TicksPerSecond)

		if config.Gamma > 0 {
			result.gamma = config.Gamma
		}
		if config.Contrast > 0 {
			result.contrast = config.Contrast
		}
	}

	return result, nil
}

// GetRendererName returns the name of the renderer
func (*Renderer) GetRendererName() string {
	return "Ebiten"
}

// SetWindowIcon sets the icon for the window, visible in the chrome of the window
func (*Renderer) SetWindowIcon(fileName string) {
	_, iconImage, err := ebitenutil.NewImageFromFile(fileName)
	if err == nil {
		ebiten.SetWindowIcon([]image.Image{iconImage})
	}
}

// IsDrawingSkipped returns a bool for whether or not the drawing has been skipped
func (r *Renderer) IsDrawingSkipped() bool {
	return r.lastRenderError != nil
}

// Run initializes the renderer
func (r *Renderer) Run(f renderCallback, u updateCallback, width, height int, title string) error {
	r.renderCallback = f
	r.updateCallback = u

	ebiten.SetWindowTitle(title)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(width, height)

	return ebiten.RunGame(r)
}

// CreateSurface creates a renderer surface from an existing surface
func (r *Renderer) CreateSurface(surface d2interface.Surface) (d2interface.Surface, error) {
	img := surface.(*ebitenSurface).image
	sfcState := surfaceState{
		filter:     ebiten.FilterNearest,
		effect:     d2enum.DrawEffectNone,
		saturation: defaultSaturation,
		brightness: defaultBrightness,
		skewX:      defaultSkewX,
		skewY:      defaultSkewY,
		scaleX:     defaultScaleX,
		scaleY:     defaultScaleY,
	}
	result := createEbitenSurface(r, img, sfcState)

	return result, nil
}

// NewSurface creates a new surface
func (r *Renderer) NewSurface(width, height int) d2interface.Surface {
	img := ebiten.NewImage(width, height)

	return createEbitenSurface(r, img)
}

// IsFullScreen returns a boolean for whether or not the renderer is currently set to fullscreen
func (r *Renderer) IsFullScreen() bool {
	return ebiten.IsFullscreen()
}

// SetFullScreen sets the renderer to fullscreen, given a boolean
func (r *Renderer) SetFullScreen(fullScreen bool) {
	ebiten.SetFullscreen(fullScreen)
}

// SetVSyncEnabled enables vsync, given a boolean
func (r *Renderer) SetVSyncEnabled(vsync bool) {
	ebiten.SetVsyncEnabled(vsync)
}

// GetVSyncEnabled returns a boolean for whether or not vsync is enabled
func (r *Renderer) GetVSyncEnabled() bool {
	return ebiten.IsVsyncEnabled()
}

// GetCursorPos returns the current cursor position x,y coordinates
func (r *Renderer) GetCursorPos() (x, y int) {
	return ebiten.CursorPosition()
}

// CurrentFPS returns the current frames per second of the renderer
func (r *Renderer) CurrentFPS() float64 {
	return ebiten.CurrentFPS()
}

// ShowPanicScreen shows a panic message in a forever loop
func (r *Renderer) ShowPanicScreen(message string) {
	errorScreen := CreatePanicScreen(message)

	err := ebiten.RunGame(errorScreen)
	if err != nil {
		panic(err)
	}
}

// SetGamma sets the gamma for the renderer
func (r *Renderer) SetGamma(gamma float64) {
	if gamma > 0 {
		r.gamma = gamma
	}
}

// SetContrast sets the contrast for the renderer
func (r *Renderer) SetContrast(contrast float64) {
	if contrast > 0 {
		r.contrast = contrast
	}
}
